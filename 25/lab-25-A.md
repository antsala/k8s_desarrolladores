# Laboratorio 25-A: ***Despliegue de archivos YAML***
 
En este laboratorio aprenderemos a desplegar archivos YAML.

Requisitos:

1. Una máquina virtual con Ubuntu 20.04 LTS a la que poder hacer ssh o escritorio remoto.
2. Tener instalado el ***runtime de podman***. (ver lab-06-A.md, Ejercicio 1 y 2)


## Ejercicio 1:  ***Aplicar un deployment desde archivo YAML***

Como hemos aprendido, se puede usar ***kubectl*** para crear objetos en el cluster de Kubernetes, pero es necesario escribir todos los parámetros y opciones en la línea de comandos. Esto resulta poco práctico. En realidad lo que se debe hacer es ***describir los objetos*** por medio de ***archivos de manifiesto***, en formato ***YAML***. Posteriormente usaremos el comando ***kubectl apply -f archivo.yaml*** para crear o actualizar dichos objetos en el cluster.

En la carpeta de los laboratorios del curso tenemos el archivo ***lab-20-B-nginx-deployment.yaml***, que contiene la estructura básica de un deployment. Procedemos a abrir el archivo con el editor ***VSC*** si se dispone de interfaz gráfica, sino, otro a elección.
```
cd ~/k8s_desarrolladores/25
code lab-25-A-nginx-deployment.yaml
```

La sintaxis y su interpretación para este archivo YAML es la siguiente:

* *Línea 1*: ***Versión*** del lenguaje de manifiesto a usar.
* *Linea 2*: Indicamos que vamos a crear un objeto de tipo ***deployment***.
* *Líneas 3-6*: Establece el nombre del deployment ***nginx-deployment*** y asigna una etiqueta ***app: my-app*** que podrá ser usada para asociar este deployment a otros objetos de k8s.
* *Línea 12*: Empieza la definición del pod.
* *Líneas 13-15*: Se asigna una etiqueta para identificar el pod (***app: my-app***)
* *Línea 17*: Commienza la definición de los contenedores que contendrá el pod.
* *Líneas 17-21*: Indicamos que existe un contenedor que se llama ***nginx***, basado en la imagen ***nginx:1.16*** y que abre el puerto ***80***.
* *Línea 7*: Comienza la definición de la especificación de los pods del deployment.
* *Líneas 9-11*: Asocia el pod con un conjunto de contenedores, definidos a partir de la línea 17. Se utiliza la etiqueta ***app: my-app*** para establecer la asociación.
* *Línea 8*: Indica cuántas instancias de pod levantará el deployment (ReplicaSet).

Una vez entendido el manifiesto, lo cerramos sin modificar. 

Si ***Minikube*** no estuviera iniciado, lo arrancamos en el siguiente comando:
```
minikube start
```

Ahora vamos a proceder a aplicar el manifiesto por medio de ***kubectl***.
```
kubectl apply -f lab-25-A-nginx-deployment.yaml
```

Procedemos a comprobar que se han creado el deployment, replicaset y pod.
```
kubectl get all
```

La salida del comando anterior debe ser parecida a esta:
```
NAME                                    READY   STATUS    RESTARTS   AGE
pod/nginx-deployment-6ff5b4564f-gdnst   1/1     Running   0          74s

NAME                 TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
service/kubernetes   ClusterIP   10.96.0.1    <none>        443/TCP   8h

NAME                               READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/nginx-deployment   1/1     1            1           74s

NAME                                          DESIRED   CURRENT   READY   AGE
replicaset.apps/nginx-deployment-6ff5b4564f   1         1         1       74s
```

La gran ventaja de usar los archivos YAML es que la modificación de los objetos del cluster es muy simple. Para demostrarlo vamos a editar el archivo YAML y cambiar el número de replicas del pod de 1 a 3.

Editar el archivo y cambiar ***replicas: 1** por ***replicas: 3***. Salir y guardar el cambio.
```
nano lab-25-A-nginx-deployment.yaml
```

Ahora volvemos a aplicar el archivo YAML. Observar cómo nos dice que el deployment ha sido configurado.
```
kubectl apply -f lab-25-A-nginx-deployment.yaml
```

Comprobamos cuantas réplicas tiene el deployment
```
kubectl get deployment nginx-deployment
```

La salida debe ser similar a esta:
```
NAME               READY   UP-TO-DATE   AVAILABLE   AGE
nginx-deployment   3/3     3            3           7m26s
```

Comprobamos también los pods y los replicasets.
```
kubectl get pods
```
```
kubectl get replicasets
```

## Ejercicio 2:  ***Aplicar un servicio desde archivo YAML***


Como ya sabemos, un servicio de Kubernetes permite que el tráfico llegue a los contenedores de los pods. Editemos el siguiente archivo:
```
code lab-25-A-nginx-service.yaml
```

La sintaxis y su interpretación para este archivo YAML es la siguiente:

* *Línea 1*: Declara la versión de la sintaxis.
* *Línea 2*: Declara el ***tipo de objeto*** que se creará, en este caso un servicio.
* *Líneas 3 y 4*: Asigna, mediante una ***etiqueta***, un ***nombre*** al servicio.
* *Línea 5*: Empieza la especificación del servicio.
* *Líneas 6 y 7*: ***Asocian*** en servicio con un deployment y definición de pod que tengan asignada la etiqueta ***app: myapp***.
* *Línea 8*: Comienza la ***especificación de puertos*** en el servicio.
* *Línea 9*: Establece el protocolo de transporte (TCP).
* *Línea 10*: ***Puerto externo*** del servicio.
* *Línea 11*: ***Puerto del contenedor*** al que se enviará el tráfico (8080)

Salimos sin cambiar nada y procedemos a crear este objeto en el cluster.
```
kubectl apply -f lab-25-A-nginx-service.yaml
```
Comprobamos que el servicio se ha creado:
```
kubectl get service nginx-service
```

La salida será similar a la siguiente.
```
NAME            TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)   AGE
nginx-service   ClusterIP   10.108.103.253   <none>        80/TCP    28s
```

Hemos creado un servicio ***INTERNO***. Lo sabemos porque ***no existe una IP externa***, tal y como indica el valor ***<none>*** en la columna ***EXTERNAL-IP***.

Un servicio interno, por definición, no puede ser accedido desde fuera del cluster, es decir, los usuarios de la aplicación no podrán acceder a los contenedores de nginx a través de este servicio.

En consecuencia, los servicios internos solo pueden ser usados por otros objetos dentro del cluster. Kubernetes expresa esto indicando que el tipo de servicio es ***ClusterIP***, en la columna ***TYPE*** y, como puede observarse en la columna ***CLUSTER-IP***, se ha asignado una IP (10.108.103.253) que solo es accesible dentro del cluster.

En el siguiente ejercicio aprenderemos a crear ***servicios externos***, que Kubernetes llama ***LoadBalancer***. El nombre no es especialmente acertado, ya que el servicio de tipo ***ClusterIP*** también balancea el tráfico.

Procedemos a estudiar la salida de la opción ***describe*** sobre el servicio, con este comando:
```
kubectl describe service nginx-service
```

La salida ofrece mucha información. y será similar a la siguiente:
```
Name:              nginx-service
Namespace:         default
Labels:            <none>
Annotations:       <none>
Selector:          app=my-app
Type:              ClusterIP
IP Family Policy:  SingleStack
IP Families:       IPv4
IP:                10.108.103.253
IPs:               10.108.103.253
Port:              <unset>  80/TCP
TargetPort:        8080/TCP
Endpoints:         172.17.0.3:8080,172.17.0.4:8080,172.17.0.5:8080
Session Affinity:  None
Events:            <none>
```

De la salida anterior destacamos la siguiente información:

* *Selector: app=my-app*. Muestra la pareja clave/valor que deberá tener los pods del deployment para la asociación.
* *Type: ClusterIP*. El servicio es de tipo ***ClusterIP*** (Explicado anteriormente)
* *IPs: 10.108.103.253*. Dirección IP (privada e interna) asignada al servicio.
* *Port: <unset>  80/TCP*. Puerto externo del servicio (80). ***<unset>*** indica que no hay controlador Ingress asociado a este servicio. Los controladores Ingress serán explicados posteriormente.
* *TargetPort: 8080/TCP*. Puerto de los contenedores. El tráfico se mandará a este puerto en el contenedor.
* *Endpoints: 172.17.0.3:8080,172.17.0.4:8080,172.17.0.5:8080*. Direcciones IPs (y puerto) de los PODs. Como el despliegue tiene ***replicas: 3***, se han creado 3 pods y, en consecuencia, el servicio tiene que repartir el tráfico entre los tres pods.


## Ejercicio 3:  ***Obtener el estado de Kubernetes***

***etcd*** en Kubernetes almacena en todo momento el estado de todos los objetos que viven en el cluster. Como ya vimos en la parte teórica, Kubernetes intenta que el estado actual de los objetos coincida con el que deseamos y declaramos en los archivos yaml.

Es muy conveniente poder consultar este estado, y para ello hacemos uso de ***kubectl***, de la siguiente forma:
```
kubectl get deployment nginx-deployment -o yaml > nginx-deployment-status.yaml
```

Hemos generado un archivo con el estado. Lo abrimos con el editor.
```
code nginx-deployment-status.yaml
```

Estudiar detenidamente el archivo, prestando especial atención a la seccción ***status***


## Ejercicio 4:  ***Crear un servicio de tipo LoadBalancer desde archivo YAML***

En este ejercicio vamos a crear un servicio de tipo ***LoadBalancer***, de forma que podamos usarlo desde fuera del cluster. Además, crearemos un deployment que use la imagen ***antsala/hello_container***, para poder verificar el balanceo en Kubernetes.

En primer lugar eliminamos el servicio y el deployment anterior. Aprovechamos la potencia de los archivos YAML para retirar los objetos del cluster.
```
kubectl delete -f lab-25-A-nginx-service.yaml
kubectl delete -f lab-25-A-nginx-deployment.yaml
```

Comprobamos que solo quede el servicio de Kubernetes:
```
kubectl get all
```

La salida debe ser similar a esta:
```
NAME                 TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
service/kubernetes   ClusterIP   10.96.0.1    <none>        443/TCP   3h
```

Editamos el archivo ***lab-25-A-helloContainer-deployment.yaml***:
```
code lab-25-A-helloContainer-deployment.yaml
```

El despliegue es muy parecido al anterior. Las líneas más destacables son:

* *Línea 4*: El nombre del deployment es ***helloContainerDeployment***.
* *Línea 8*: Queremos 3 pods.
* *Líneas 10 y 11*: Se utiliza el selector ***app: helloContainer***.
* *Línea 19*: Se usará la imagen ***antsala/hello_container*** para los contenedores.

Salimos sin guardar y aplicamos el objeto.
```
kubectl apply -f lab-25-A-helloContainer-deployment.yaml
```

Comprobamos que se han desplegado los tres pods.
```
kubectl get deployment hello-container-deployment
```

La salida será similar a esta:
```
NAME                         READY   UP-TO-DATE   AVAILABLE   AGE
hello-container-deployment   3/3     3            3           46s
```

Aseguramos que los pods están en running:
```
kubectl get pods
```

La salida será como esta:
```
NAME                                         READY   STATUS    RESTARTS   AGE
hello-container-deployment-566d999d9-4m9kk   1/1     Running   0          96s
hello-container-deployment-566d999d9-7296q   1/1     Running   0          96s
hello-container-deployment-566d999d9-qfvb9   1/1     Running   0          96s
```

Como el pod tiene un único contenedor y ***STATUS*** pone ***Running***, podemos garantizar que los contenedores se han iniciado correctamente.

Ahora procedemos a desplegar el servicio. Editamos el archivo ***lab-25-A-helloContainer-service.yaml***:
```
code lab-25-A-helloContainer-service.yaml
```

Las líneas más importantes son:

* *Línea 3*: El nombre del servicio es ***hello-container-service***.
* *Línea 6*: Aquí tenemos el cambio fundamental. ***type: LoadBalancer*** hace que Kubernetes genere para el servicio una ***IP Externa*** (que puede ser pública o privada). De esta forma podemos contactar con el servicio desde fuera del cluster, por ejemplo desde el ***localhost***.
* *Líneas 7 y 8*: Se utiliza el selector ***app: helloContainer*** para asociar este servicio con el deployment y definición de pod apropiadas.

Aplicamos el servicio:
```
kubectl apply -f lab-25-A-helloContainer-service.yaml
```

Comprobamos que el servicio está funcionando.
```
kubectl get service hello-container-service
```

La salida es la siguiente:
```
NAME                      TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
hello-container-service   LoadBalancer   10.106.173.212   <pending>     80:32758/TCP   25s
```

Como se puede ver en la columna ***TYPE***, el servicio es un ***LoadBalancer***. 

En la columna ***EXTERNAL-IP*** aparece el valor ***<pending>'***. Esto es así porque minikube no sabe qué dirección IP externa debe asignar al servicio. En un proveedor cloud tendremos inmediatamente la IP externa, ya que estos proveedores automatizan la asignación de esta IP. Solo en el caso de Minikube debemos hacer lo siguiente: 

En una terminal diferente, ejecutamos el siguiente comando:
```
minikube tunnel
```

Cuando se inicie el túnel, se asignará la IP externa al servicio ***hello-container-service***. Si volvemos a mostrar información del servicio, la veremos:
```
kubectl get service hello-container-service
```

En este caso la salida es:
```
NAME                      TYPE           CLUSTER-IP       EXTERNAL-IP      PORT(S)        AGE
hello-container-service   LoadBalancer   10.106.173.212   10.106.173.212   80:32758/TCP   6m4s
```

En la columna ***EXTERNAL-IP*** podemos ver la IP asignada (10.106.173.212)
```
ip_externa=<poner aquí EXTERNAL-IP>
```

Ya solo queda probar el balanceo en Kubernetes. Repetir varias veces el siguiente comando:
```
curl $ip_externa:80
```

Cerramos la terminal del túnel de Minikube.

Limpiamos recursos.
```
kubectl delete -f lab-25-A-helloContainer-service.yaml
kubectl delete -f lab-25-A-helloContainer-deployment.yaml
```

Comprobamos que solo quede el servicio de Kubernetes:
```
kubectl get all
```

