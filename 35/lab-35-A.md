# Laboratorio 35-A: ***Ingress***
 
En este laboratorio aprenderemos a usar los objetos Ingress.

Requisitos:

1. Una máquina virtual con Ubuntu 20.04 LTS a la que poder hacer ssh o escritorio remoto.
2. Tener instalado el ***runtime de podman***. (ver lab-06-A.md, Ejercicio 1 y 2)


## Ejercicio 1:  ***Despliegue de las aplicaciones***

Entramos en el directorio del laboratorio
```
cd ~/k8s_desarrolladores/35
```

Procedemos a redesplegar las aplicaciones ***hellocontainer*** y ***mongodb*** de los ejercicios anteriores, pero con un cambio ***NOTABLE***: Puesto que vamos a usar un ***controlador Ingress*** para reenviar el tráfico a los servicios del Frontend de las aplicaciones, éstos servicios deben convertirse a servicios ***INTERNOS***.

Para el Frontend de ***helloContainer*** no hay cambios. Editamos el archivo para recordar qué contiene.
```
code lab-35-A-helloContainer-deployment.yaml
```

* *línea 21*: Expresa que los contenedores estarán dando servicio en el puerto ***8080***.

Cerramos el archivo sin modificarlo.

***helloContainer*** definía un servicio EXTERNO para exponer los pods. Ahora ese servicio debe ser reconvertido a ***INTERNO***. Editamos el archivo ***lab-35-A-helloContainer-service.yaml***.
```
code lab-35-A-helloContainer-service.yaml
```

El cambio más notable es que ***NO APARECE*** el parámetro ***type: LoadBalancer*** en la especificación del servicio. Esto lo convierte en un ***SERVICIO INTERNO***. Además, en la línea 4, se cambia el nombre del servicio a ***hello-container-internal-service***. Por otro lado, el servicio escucha en el puerto ***4000*** (Línea 10) y reenvía el tráfico al puerto ***8080*** de los pods (Línea 11)

Guardamos sin salir.

Aplicamos ambos archivos.
```
kubectl apply -f lab-35-A-helloContainer-deployment.yaml
kubectl apply -f lab-35-A-helloContainer-service.yaml
```

Comprobamos que se hayan creado los objetos:
```
kubectl get all
```

Miramos cómo ha quedado el servicio interno.
```
kubectl get service hello-container-internal-service
```

La salida del comando anterior será similar a esta:
```
NAME                               TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
hello-container-internal-service   ClusterIP   10.101.242.181   <none>        4000/TCP   8m15s
```

El tipo ***ClusterIP*** y la ausencia de ***EXTERNAL-IP*** identifican al servicio como ***INTERNO***.

Procedemos ahora al despliegue de la segunda aplicación. Esta contiene un Frontend, un Backend, un secreto y un configmap. El objeto ***secret*** y el ***configmap*** no tienen cambios. Aplicamos directamente sus archivos YAML.
```
kubectl apply -f lab-35-A-mongodb-secret.yaml
kubectl apply -f lab-35-A-mongodb-configmap.yaml
```

Comprobamos que se ha creado el secreto.
```
kubectl get secret mongodb-secret
```

La salida debe ser parecida a esta:
```
NAME             TYPE     DATA   AGE
mongodb-secret   Opaque   2      40s
```

Hacemos lo propio para el configmap.
```
kubectl get configmap mongodb-configmap
```

Y la salida es:
```
NAME                DATA   AGE
mongodb-configmap   1      12h
```

Procedemos a desplegar el Frontend y el Backend.

El Backend no sufre cambios, editamos el archivo para recordarlo:
```
code lab-35-A-mongodb.yaml
```

Las líneas 42 y 43 define el puerto en el que escuchará sel servicio de backend (***27017***) y el tráfico será reenviado al puerto ***27017*** del pod de ***mongodb server***. Salimos sin modificar nada.


Por último, para el Frontend si hay cambios. Editamos el archivo ***lab-35-A-mongo-express.yaml***.
```
code lab-35-A-mongo-express.yaml
```

El cambio más notable es que ***NO APARECE*** el parámetro ***type: LoadBalancer*** en la especificación del servicio. Esto lo convierte en un ***SERVICIO INTERNO***. Además, en la línea 42, se cambia el nombre del servicio a ***mongo-express-internal-service***. Por otro lado, el servicio escucha en el puerto ***5000*** (Línea 48) y reenvía el tráfico al puerto ***8081*** del pod ***mongo-express*** (Línea 49).

Salimos sin guardar. Aplicamos los archivos YAML.
```
kubectl apply -f lab-35-A-mongodb.yaml
kubectl apply -f lab-35-A-mongo-express.yaml
```

Verificamos los servicios asociados a la aplicación.
```
kubectl get services mongodb-service mongo-express-internal-service
```

La salida es:
```
NAME                             TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)     AGE
mongodb-service                  ClusterIP   10.107.141.33   <none>        27017/TCP   4m25s
mongo-express-internal-service   ClusterIP   10.104.152.51   <none>        5000/TCP    3m44s
```

Ambos servicios son internos (***ClusterIP***). A falta de crear el objeto ***Ingress***, el cluster debe contener lo siguiente:
```
kubectl get all
```

La salida indica que hay ***5 pods***, ***4 servicios***, ***3 deployments*** sus respectivos (3) ***replicasets***.
```
NAME                                             READY   STATUS    RESTARTS   AGE
pod/hello-container-deployment-566d999d9-4wpdw   1/1     Running   0          30m
pod/hello-container-deployment-566d999d9-9nf5z   1/1     Running   0          30m
pod/hello-container-deployment-566d999d9-xh5cz   1/1     Running   0          30m
pod/mongo-express-deployment-68c4748bd6-jwcrx    1/1     Running   0          5m13s
pod/mongodb-deployment-7bb6c6c4c7-mp8jq          1/1     Running   0          5m54s

NAME                                       TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)     AGE
service/hello-container-internal-service   ClusterIP   10.101.242.181   <none>        4000/TCP    30m
service/kubernetes                         ClusterIP   10.96.0.1        <none>        443/TCP     22h
service/mongo-express-internal-service     ClusterIP   10.104.152.51    <none>        5000/TCP    5m13s
service/mongodb-service                    ClusterIP   10.107.141.33    <none>        27017/TCP   5m54s

NAME                                         READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/hello-container-deployment   3/3     3            3           30m
deployment.apps/mongo-express-deployment     1/1     1            1           5m13s
deployment.apps/mongodb-deployment           1/1     1            1           5m54s

NAME                                                   DESIRED   CURRENT   READY   AGE
replicaset.apps/hello-container-deployment-566d999d9   3         3         3       30m
replicaset.apps/mongo-express-deployment-68c4748bd6    1         1         1       5m13s
replicaset.apps/mongodb-deployment-7bb6c6c4c7          1         1         1       5m54s
```


## Ejercicio 2:  ***Creación del objeto Ingress***

En este momento tenemos las aplicaciones funcionando, pero no son accesibles a ser todos los servicios INTERNOS. Vamos a agregar un objeto ***Ingress*** que definirá las reglas que aplicará el ***Controlador Ingress*** (aun por crear). Para enviar el tráfico externo a cada servicio de Frontend. Crearemos dos dominios: ***www.hellocontainer.com*** y ***www.mongoexpress.com***. En función de la URL que escriba el usuario en su navegador el controlador ingress redirigirá el tráfico al servicio interno de Frontend apropiado.

Editamos el archivo ***lab-35-A-ingress.yaml***.
```
code lab-35-A-ingress.yaml
```

Este YAML hace lo siguiente:

* *Línea 1*: Versión de la API.
* *Línea 2*: Declara el tipo de objeto como ***Ingress***.
* *Línea 4*: Este objeto se llamará ***my-ingress***.
* *Línea 6*: Empieza la definición de las reglas que se van a aplicar al tráfico.
* *Línea 7*: Se aplicará si la URL contiene el dominio dns ***www.hellocontainer.com***.
* *Línea 8*: Comienza la definición para las reglas ***http***.
* *Líneas 9-16*: Regla para reenviar tráfico a la aplicación ***helloContainer***.
* *Línea 10 y 11*: Si la URI comienza por ***/***.
* *Línea 13-16*: Se reenvia el tráfico al servicio ***hello-container-internal-service*** puerto ***4000***.
* *Líneas 17-26*: Regla para reenviar tráfico a la aplicación ***mongo-express***.
* *Línea 17*: Se aplicará si la URL contiene el dominio dns ***www.mongoexpress.com***.
* *Línea 20*: Si la URI comienza por ***/***.
* *Línea 23-26*: Se reenvia el tráfico al servicio ***mongo-express-internal-service*** puerto ***5000***.

En resumen, el flujo de redirecciones será:
Si usuario conecta a ***http://www.hellocontainer.com/*** --> hello-container-internal-service:4000 --> pod:8080
Si usuario conecta a ***http://www.mongoexpress.com/*** --> mongo-express-internal-service:5000 --> pods:8081

Salimos sin modificar nada. Creamos el objeto ***Ingress***.
```
kubectl apply -f lab-35-A-ingress.yaml
```

Comprobamos que se ha creado.
```
kubectl get ingress my-ingress
```

La salida es:
```
NAME         CLASS    HOSTS                                         ADDRESS     PORTS   AGE
my-ingress   nginx    www.hellocontainer.com,www.mongoexpress.com   localhost   80      78m
```

Pedimos más información:
```
kubectl describe ingress my-ingress
```

La salida mostrará algo como esto: (Nota: observar los ***path*** y los ***backends***)
```
Name:             my-ingress
Labels:           <none>
Namespace:        default
Address:          localhost
Default backend:  default-http-backend:80 (<error: endpoints "default-http-backend" not found>)
Rules:
  Host                    Path  Backends
  ----                    ----  --------
  www.hellocontainer.com  
                          /   hello-container-internal-service:4000 (172.17.0.3:8080,172.17.0.5:8080,172.17.0.6:8080)
  www.mongoexpress.com    
                          /   mongo-express-internal-service:5000 (172.17.0.8:8081)
Annotations:              <none>
Events:
  Type    Reason  Age                 From                      Message
  ----    ------  ----                ----                      -------
  Normal  Sync    11m (x16 over 79m)  nginx-ingress-controller  Scheduled for sync
```

## Ejercicio 3:  ***Instalar el Controlador Ingress***

Minikube ya viene con un controlador ingress instalado, de forma que solo debemos habilitarlo:
```
minikube addons enable ingress
```

La salida mostrará lo siguiente:
```
    ▪ Using image k8s.gcr.io/ingress-nginx/kube-webhook-certgen:v1.1.1
    ▪ Using image k8s.gcr.io/ingress-nginx/controller:v1.1.0
    ▪ Using image k8s.gcr.io/ingress-nginx/kube-webhook-certgen:v1.1.1
    Verifying ingress addon...
    The 'ingress' addon is enabled
```

El ***addon*** crea un espacio de nombres (llamado ***ingress-nginx***) y una serie de objetos dentro para poder funcionar. Listamos los objetos:
```
kubectl get all --namespace ingress-nginx
```

La salida es similar a la siguiente. Crea 3 pods, 2 servicios, 1 deployment, 1 replicaset y 2 jobs.

Los jobs son pods que se instancian, realizan su cometido y cuando finalizan, se destruyen.
```
NAME                                            READY   STATUS      RESTARTS   AGE
pod/ingress-nginx-admission-create-pq7zs        0/1     Completed   0          89m
pod/ingress-nginx-admission-patch-l86xq         0/1     Completed   0          89m
pod/ingress-nginx-controller-6d5f55986b-tdvcf   1/1     Running     0          89m

NAME                                         TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)                      AGE
service/ingress-nginx-controller             NodePort    10.96.44.103   <none>        80:32496/TCP,443:30467/TCP   89m
service/ingress-nginx-controller-admission   ClusterIP   10.101.14.77   <none>        443/TCP                      89m
 
NAME                                       READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/ingress-nginx-controller   1/1     1            1           89m

NAME                                                  DESIRED   CURRENT   READY   AGE
replicaset.apps/ingress-nginx-controller-6d5f55986b   1         1         1       89m

NAME                                       COMPLETIONS   DURATION   AGE
job.batch/ingress-nginx-admission-create   1/1           5s         89m
job.batch/ingress-nginx-admission-patch    1/1           5s         89m
```

De todos los objetos anteriores, el pod del controlador de ingress es el más importante. En esta captura es ***pod/ingress-nginx-controller-6d5f55986b-tdvcf***. Debe estar siempre en ***Running***. Mirando sus ***logs*** o con el comando ***describe*** podemos ver como entra el tráfico en el cluster.
```
kubectl logs log <Poner aquí el nombre del pod del controlador> --namespace ingress-nginx 
```

## Ejercicio 4:  ***Configurar el Registro de Recurso de DNS***

En el arhchivo YAML del objeto ingress hemos añadido dos URLs que deben entrar al cluster: ***www.hellocontainer.com*** y ***wwww.mongoexpress.com***. Debemos hacer que se resuelvan a la IP del entrypoint del cluster.

Tomamos las IP del nodo de minikube
```
minikube_ip=$(minikube ip)
echo $minikube_ip
```

Copiar la IP en el portapapeles. A falta de usar un servidor DNS para crear registros ***A***, editamos el archivo ***hosts***.
```
sudo nano /etc/hosts
```

En la sección IPv4, creamos un par de líneas nuevas en la que debe aparecer la IP copiada anteriormente y los dominios ***www.hellocontainer.com*** y ***wwww.mongoexpress.com***. Debería ser algo parecido a esto:
```
127.0.0.1       localhost
127.0.1.1       ubu
192.168.49.2    www.hellocontainer.com       # <---- Esta es la línea que hemos añadido.
192.168.49.2    www.mongoexpress.com         # <---- Esta también.
```

Salimos y guardamos. Comprobamos la resolución DNS:
```
nslookup www.hellocontainer.com
```
```
nslookup www.mongoexpress.com
```

La salida para ambos será como esta, donde aparece la IP de Minikube.
```
midominio.com has address 192.168.49.2
```

Ya solo queda probar las aplicaciones: Para ***hellocontainer*** cogemos un navegador y nos conectamos a ***http://www.hellocontainer.com/*** y para ***mongoexpress*** a ***http://www.mongoexpress.com***. (Nota: Para ***hellocontainer*** se puede verificar el balanceo a los pods)

En el siguiente enlace está toda la documentación de Ingress, en la que podemos ver, todos los detalles, por ejemplo la reescritura de la URL. Un artículo sin duda muy importante que tenemos que leer: https://kubernetes.io/docs/concepts/services-networking/ingress/


## Ejercicio 5:  Configurar el ***Default Backend*** 

En nuestro ejemplo, las reglas de ingress reenvían el tráfico al directorio raíz de las dos apps, de los dominios anteriores, pero ¿qué pasa si se pone una URL para la que no hay establecida una propiedad ***path***? Pues el resultado es que no se puede encontrar un Backend apropiado, es decir, no hay ningún servicio interno al que redirigir el tráfico.  Kubernetes ofrece para estos casos el ***Default Backend***, al que será redirigido todo el tráfico que no pueda verificar ninguna de las reglas creadas. Por ejemplo, si intentamos conecta a un path que no esta definido en las reglas, obtendremos un error:
```
curl http://www.mongoexpress.com/path_no_existente
```

La salida es la siguiente:
```
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>Error</title>
</head>
<body>
<pre>Cannot GET /path_no_existente</pre>
</body>
</html>
```

La respuesta es una página de error genérica que podemos mejorar. Ejecutemos el siguiente comando:
```
kubectl describe ingress my-ingress
```

La salida es la siguiente 
```
Name:             my-ingress
Labels:           <none>
Namespace:        default
Address:          localhost
Default backend:  default-http-backend:80 (<error: endpoints "default-http-backend" not found>)
Rules:
...
```

Nos interesa la línea ***Default backend:  default-http-backend:80 (<error: endpoints "default-http-backend" not found>)***, que especifica el servicio (y puerto) al que se redirigirá al navegador del usuario si el ingress no verifica ninguna de las reglas. El servicio que será llamado tiene el nombre ***default-http-backend*** (puerto ***80***), y nos está diciendo que ese servicio no tiene endpoints (pods).

Lo que debemos hacer es crear un servidor web que muestre la página de error deseada. Para ello creamos un deployment que levantará un pod de la imagen ***antsala/page_not_found***. Esta imagen levanta un servidor web que muestra un mensaje amigable de página no encontrada (error 404).

Abrimos el archivo.
```
code lab-35-A-page-404-deployment.yaml
```

Es el típico deployment con un pod. Cerramos sin modificar y aplicamos.
```
kubectl apply -f lab-35-A-page-404-deployment.yaml
```

Comprobamos que el pod se inicia.
```
kubectl get deployment page-404-deployment
```

Nos resta por crear un servicio, al que llamaremos ***page-404-internal-service***, que escucha en el puerto ***6666*** y que reenviará el tráfico al puerto ***80*** del pod seleccionado por la etiqueta ***app: page-404***.

Abrimos el archivo ***lab-35-A-page-404-internal-service.yaml*** para verificarlo.
```
code lab-35-A-page-404-internal-service.yaml
```

Cerramos sin cambiar nada. Creamos el servicio.
```
kubectl apply -f lab-35-A-page-404-internal-service.yaml
```

Comprobamos que está funcionando.
```
kubectl get service page-404-internal-service
```

La salida mostrará algo similar a esto:
```
NAME                        TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)   AGE
page-404-internal-service   ClusterIP   10.109.167.175   <none>        6666/TCP   33s
```

Ahora debemos modificar el objeto ingress para que configure el endpoint para el Default Backend. El archivo ***lab-35-A-ingress-with-default-backend.yaml*** ya tiene estas modificaciones. Lo editamos para estudiarlo:
```
code lab-35-A-ingress-with-default-backend.yaml
```


El ***Default backend*** responde en el path ***/***, es decir, si nos conectamos a ***nombre_de_dominio_dns/***, saltará la página de error personalizada. 

Observa los cambios en este YAML:

* *Línea 6 y 7*: Definen anotaciones, indicando que el controlador ***Ingress*** es de la clase ***Nginx***, y que se producirá una reescritura de la URL, de forma que si dicha URL no existe, se reenviará el tráfico al default backend en el path ***/***.
* *Línea 18 y 19 *: Nuestra aplicación responderá en el path ***/app*** . Observa el uso de expresiones regulares.

Para la reescritura de URL te recomendamos que leas este artículo del controlador ingress de Nginx: https://kubernetes.github.io/ingress-nginx/user-guide/ingress-path-matching/


Lo aplicamos:
```
kubectl apply -f lab-35-A-ingress-with-default-backend.yaml
```

Mostramos información:
```
kubectl describe ingress my-ingress
```

En este caso la salida muestra:
```
Name:             my-ingress
Labels:           <none>
Namespace:        default
Address:          localhost
Default backend:  page-404-internal-service:6666 (172.17.0.9:80)
....
```

Como se puede observar, el Default backend está actualizado para que se llame al servicio ***page-404-internal-service:6666***, que a su vez redirige el tráfico al pod ***172.17.0.9:80***, que mostrará la página error personalizada.

Prueba a conectar a la siguiente URL. Debes ver la aplicación.
```
http://www.hellocontainer.com/app
```

Si conectas a esta otra, debes ver la página de error personalizada.
```
http://www.hellocontainer.com/otra_url
```

Limpiamos los recursos.
```
kubectl delete -f lab-35-A-helloContainer-deployment.yaml
kubectl delete -f lab-35-A-helloContainer-service.yaml
kubectl delete -f lab-35-A-mongodb.yaml
kubectl delete -f lab-35-A-mongo-express.yaml
kubectl delete -f lab-35-A-mongodb-secret.yaml 
kubectl delete -f lab-35-A-mongodb-configmap.yaml
kubectl delete -f lab-35-A-ingress.yaml
kubectl delete -f lab-35-A-page-404-deployment.yaml 
kubectl delete -f lab-35-A-page-404-internal-service.yaml 
```

Comprobamos los objetos presentes en el cluster.
```
kubectl get all 
```

Solo debe quedar el servicio de Kubernetes.
```
NAME                 TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
service/kubernetes   ClusterIP   10.96.0.1    <none>        443/TCP   53m
```
