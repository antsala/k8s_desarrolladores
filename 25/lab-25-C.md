# Laboratorio 25-C: ***Backend de Redis con un master y dos réplicas***
 
En este laboratorio aprenderemos a desplegar una app con in Frontend y un  Backend de Redis con un master y dos réplicas (master y dos slaves)

Requisitos:

1. Una máquina virtual con Ubuntu 20.04 LTS a la que poder hacer ssh o escritorio remoto.
2. Tener instalado el runtime de podman. (ver lab-06-A.md, Ejercicio 1 y 2)

## Ejercicio 1:  ***Despliegue del maestro de Redis***

Vamos desplegar una app con un frontend y backend de redis (un master y dos réplicas)

Cambiamos al directorio de trabajo.
```
cd ~/k8s_desarrolladores/25
```

Desplegamos el backend del maestro de redis. Abrimos el achivo ***lab-25-redis-master-deployment.yaml***.
```
code lab-25-C-redis-master-deployment.yaml
```

El contenido a destacar de este archivo lo comentamos a continuación:

* *Línea 2*: Es un ***deployment***.
* *Línea 4*: Asignamos el nombre de ***redis-master-deployment***.
* *Línea 6*: Se define la etiqueta ***app: redis***.
* *Líneas 8-12*: Importante: Hasta el momento hemos usado una sola etiqueta para asociar el deployment a la plantilla de pod. En este ejemplo se usan varias etiquetas: ***app: redis***, ***role: master*** y ***tier: backend***. Esto va a darnos más libertad a la hora de asociar servicios/deployments/pods, a la vez que añade más información al archivo YAML que servirá para entender mejor el propósito del objeto que se está creando.

Para que se produzca la asociación entre objetos, se deben verificar TODAS las etiquetas.

* *Líneas 17-19*: En la especificación de la plantilla del pod, volvemos a poner las mismas etiquetas.
* *Líneas 22*:    El contenedor ser llamará ***master***.
* *Línea 23*:     y estará basado en la imagen es ***k8s.gcr.io/redis:e2e***. 
* *Líneas 24-30*: Se especifica la contención de recursos para el contenedor. Si el servidor va sobrado de recursos, el contenedor podrá usar más recursos de los que se declaran en ***requests***, pero en ningún caso más de lo que aparece en ***limits***. La CPU se puede expresar en tanto por uno o, como en este caso, usando la unidad ***milis*** (m). 1000 milis equivalen al 100% de la CPU disponible. Es muy conveniente leer el siguiente artículo: (https://kubernetes.io/es/docs/concepts/configuration/manage-resources-containers/)

Aplicamos el deployment:
```
kubectl apply -f lab-25-C-redis-master-deployment.yaml
```

Examinamos el deployment:
```
kubectl get deployment redis-master-deployment
```

La salida será parecida a esta:
```
NAME                      READY   UP-TO-DATE   AVAILABLE   AGE
redis-master-deployment   1/1     1            1           8s
```

Miramos el detalle:(Nota: Leer detenidamente la salida y observar la información de escalado del ReplicaSet)
```
kubectl describe deployment/redis-master-deployment
```

La implementación ha sido una prueba, borramos el deployment desde la línea de comandos.
```
kubectl delete deployment/redis-master-deployment
```

Comprobar que se elimina y que solo quede el servicio de Kubernetes:
```
kubectl get all
```

Ahora vamos a hacer lo mismo pero usando un ***ConfigMap***. Hay dos formas de crear un ConfigMap: Desde un archivo de texto o desde un archivo yaml.


## Ejercicio 2:  ***Creación de ConfigMap desde un archivo***

Creamos un archivo, llamado ***lab-25-C-redis-config*** y le ponemos estas dos líneas:
```
maxmemory 2mb
maxmemory-policy allkeys-lru
```

Realmente no hace falta hacerlo porque en el directorio de ejemplos ya existe ese archivo. Así que nos limitamos a abrirlo para comprobarlo.
```
code lab-25-C-redis-config
```

* *allkeys-lru'*: lru = Less Recently Used

Para crear el ConfigMap, ejecutamos el siguiente comando.
```
kubectl create configmap redis-configmap-from-file --from-file=lab-25-C-redis-config
```

y ahora listamos el configmap.
```
kubectl get configmaps redis-configmap-from-file
```

La salida debe mostrar el objeto creado:
```
NAME                        DATA   AGE
redis-configmap-from-file   1      6s
```

Observamos contenido. Nótese que hay una clave llamada ***redis-config*** y una lista con los valores de la configuración:
```
kubectl describe configmap/redis-configmap-from-file
```

Esta es la salida:
```
Name:         redis-configmap-from-file
Namespace:    default
Labels:       <none>
Annotations:  <none>

Data
====
redis-config:
----
maxmemory 2mb
maxmemory-policy allkeys-lru

BinaryData
====

Events:  <none>
```

Existe otra forma de crear el objetos ***ConfigMap***, mediante un archivo YAML. Para ello borramos el objeto recién creado:
```
kubectl delete configmap/redis-configmap-from-file
```

En este caso partiríamos desde un archivo YAML, que tenemos disponible en la carpeta del laboratorio. Abrimos el archivo 'redis-config.yaml' 
```
code lab-25-C-redis-config.yaml
```

El contenido más destacado del archivo es el siguiente:

* *Línea 2*: Comienza el contenido de datos del archivo.
* *Línea 3*: ***redis-config*** será el nombre la clave. Observar el uso de ***|-*** que se utiliza para definir un string en varias líneas. Elimina el salto de linea al final y los espacios en blanco si los hubiera.
* *Línea 6*: El tipo de objeto que se creará es un ***ConfigMap***. Este parámetro se puede poner en cualquier línea siempre que este sangrado a la izquierda, pero lo más habitual es ponerlo en las primeras líneas del archivo.
* *Línea 8*: ***redis-config-from-yaml*** es el nombre del ***ConfigMap***.
* *Línea 9*: Se pone el espacio de nombres en el que se creará este objeto. Si lo queremos en ***default*** no es necesario especificarlo.

Volvemos a crear el configmap, pero esta vez directamente desde el archivo YAML:
```
kubectl create -f lab-25-C-redis-config.yaml
```

Y lo describimos:
```
kubectl describe configmap/redis-config-from-yaml
```

La salida es idéntica al caso anterior.

K8s tiene una opción muy útil, la ***-o*** (***--output***), que puede ser usada para obtener la salida de un objeto presente en ***etcd*** en formato YAML o JSON. Así se puede tener el archivo yaml del objeto. Lo prodriamos redireccionar a un archivo si lo vemos necesario:
```
kubectl get configmap/redis-config-from-yaml --output yaml 
```

La salida (por pantalla): (Nota: Observar cómo K8s agrega sus propios parámetros)
```
apiVersion: v1
  data:
   redis-config: |-
    maxmemory 2mb
    maxmemory-policy allkeys-lru
kind: ConfigMap
metadata:
  creationTimestamp: "2022-03-17T18:34:31Z"
  name: redis-config-from-yaml
  namespace: default
  resourceVersion: "13147"
  uid: d22709fa-c5f6-47a2-9e5b-c9f088a85419
```

Ahora vamos a usar un ***ConfigMap*** para pasar información al contenedor en tiempo de ejecución. Vamos a abrir el archivo modificado:
```
code lab-25-C-redis-master-deployment-modified.yaml
```

Observar lo siguiente:

* *Líneas 24-26*: Lanza el contenedor de redis pasándole información sobre el archivo de configuración que debe utilizar: ***/redis-master/redis.conf***.
* *Líneas 27-29*: Pasa al contenedor una variable de entorno: ***MASTER=true***.
* *Líneas 30-32*: Monta un volumen llamado ***config*** en la ruta ***/redis-master*** del contenedor.
* *Líneas 39-45*: Almacena el valor de la clave ***redis-config*** del configmap ***redis-config-from-yaml*** en el archivo ***redis.conf*** en el volumen ***config***.

Cuando el contenedor arranca, lee el archivo ***/redis-master/redis.conf***, que contiene los valores del configmap que se creó.

Desplegamos esta versión actualizada para usar configmap.
```
kubectl create -f lab-25-C-redis-master-deployment-modified.yaml
```

Miramos los pods
```
kubectl get pods
```

La salida debe mostrar algo parecido a esto, indicando que el pod está en ejecución.
```
NAME                                       READY   STATUS    RESTARTS   AGE
redis-master-deployment-754ccc67d4-ctp9v   1/1     Running   0          7s
```

Ejecutamos un comando dentro del pod para comprobar si realmente ha leido los valores de configuración. Cambiar el nombre del pod. ***--*** indica que lo que viene después es el comando que ejecutan los contenedores del pod. En este caso, el pod tiene un único contenedor.
```
kubectl exec -it <Poner aquí nombre del pod> -- redis-cli
```

El comando anterior debe abrir una conexión con REDIS. Ejecutamos los siguientes comandos para verificar que hay ***2MB*** y la politica es ***LRU***: Borra las claves que se usan menos Less Recently Used cuando le haga falta memoria.
```
CONFIG GET maxmemory
CONFIG GET maxmemory-policy
```

Salimos con ***exit***.


Ahora vamos a desplegar un servicio INTERNO para los pods del deployment redis-master.Abrimos el archivo ***redis-master-service.yaml***.
```
code lab-25-C-redis-master-service.yaml
```

Estudiamos el contenido del archivo.

* *Línea 2*: El objeto es un servicio.
* *Línea 4*: Con nombre ***redis-master***.
* *Líneas 5-8*:   El servicio crea las etiquetas ***app: redis***, ***role: master*** y ***tier: backend***.
* *Línea 12*:     El servicio atiende en el puerto ***6379***...
* *Línea 13*:     y reenvía el tráfico a los pod al puerto ***6379***.
* *Líneas 13-16*: El servicio enviará tráfico a los pods que tengan definidas las etiquetas ***app: redis***, ***role: master*** y ***tier: backend***.


Creamos el objeto:
```
kubectl apply -f lab-25-C-redis-master-service.yaml
```

Comprobamos el despliegue del servicio:
```
kubectl get service redis-master
```

La salida del comando anterior es:
```
NAME                TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
redis-master        ClusterIP   10.100.15.241   <none>        6379/TCP   32s
```
Comprobar que el tipo de servicio es ***ClusterIP***, por lo que solo se puede acceder a él desde dentro del cluster (desde otros pods) y no desde el exterior del cluster.

Un servicio también introduce un ***nombre DNS*** para dicho servicio, en la forma: ***<nombreServicio>.<espacioDeNombres>svc.cluster.local***. Puesto que estamos usando el espacio de nombres ***default***, la DNS del servicio ***redis-master*** es: ***redis-master.default.svc.cluster.local***.

Para ver esto funcionando, nos metemos en el pod para ver si hay resolución DNS. Lo vamos a hacer con ***ping*** ya que ***nslookup*** no está instalado en el pod. No va a haber respuesta de ping, solo nos interesa la resolución DNS.

Listamos pods y nos quedamos con su nombre:
```
kubectl get pods
```

La salida se parecerá a esta:
```
NAME                                       READY   STATUS    RESTARTS   AGE
redis-master-deployment-754ccc67d4-ctp9v   1/1     Running   0          13h
```
```
kubectl exec -it <Poner aquí el nombre del pod> -- ping redis-master
```

Comprobar que el registro A ***redis-master.default.svc.cluster.local***. se resuelve a la ***CLUSTER-IP anterior***. CTRL+C para salir.

Convolución. Se completa el dominio si no se especifica completamente.
```
kubectl exec -it <Poner aquí el nombre del pod> -- ping redis-master
```
```
kubectl exec -it <Poner aquí el nombre del pod> -- ping redis-master.default
```
```
kubectl exec -it <Poner aquí el nombre del pod> -- ping redis-master-.internal-service.default.svc.cluster.local
```

## Ejercicio 3:  ***Despliegue de las réplicas de Redis***

Ahora vamos a desplegar las réplicas (2) de redis, que se sincronizarán desde redis-master. Para ello estudiamos el archivo ***lab-25-C-redis-replica-deployment.yaml***.
```
code lab-25-C-redis-replica-deployment.yaml
```

Las líneas más importantes del archivo son:

* *Línea 2*: Indicamos que es un deployment.
* *Línea 4*: Su nombre es ***redis-replica-deployment***.
* *Líneas 8-12*:  El deployment se asociará con una plantilla de pod que tenga definidas las etiquetas: ***app: redis***, ***role: replica*** y ***tier: backend***.
* *Línea 13*: Se instanciarán dos pods.
* *Línea 21*: Comienza la definición del contenedor.
* *Línea 22*: El nombre del contenedor será ***replica***.
* *Línea 23*: Estará basado en la imagen de réplica de redis ***gcr.io/google_samples/gb-redis-follower:v1***.
* *Línea 25*: Se le asignará el 10% de la CPU (100 milis)
* *Línea 26*: y 100 Mebibits de memoria. (Leer este artículo: https://es.wikipedia.org/wiki/Mebibit)
* *Línea 28-30*: Se crea la variable de entorno ***GET_HOST_FROM*** con el valor ***dns***. Cuando el contenedor se inicie se conectará a una máquina (su master) llamada ***redis-server***, es decir a ***redis-server-internal-service.default.svc.cluster.local***, que es la IP del servicio de ***redis-server***.

Aplicamos el deployment:
```
kubectl apply -f lab-25-C-redis-replica-deployment.yaml
```

Comprobamos que arranca el deployment:
```
kubectl get deployment redis-replica-deployment
```

La salida será similar a esta:
```
NAME                       READY   UP-TO-DATE   AVAILABLE   AGE
redis-replica-deployment   2/2     2            2           41s
```

Para que el frontend (aún por desplegar) pueda contactar con las réplicas (además del master de redis), es necesario exponerlas mediante un servicio, que nos dará la respectiva ClusterIP.

Editamos el archivo:
```
code lab-25-C-redis-replica-service.yaml
```

Las líneas más importantes son:

* *Línea 2*: Se crea un servicio interno.
* *Línea 4*: Con nombre ***redis-replica-internal-service***.
* *Línea 11*: Escuchará en el puerto ***6379***...
* *Línea 12*: y reenviará al puerto ***6379*** de los pods.
* *Líneas 13-16*: Se asociará con los pods que tengan definidas las etiquetas: ***app: redis***, ***role: replica***, ***tier: backend***.

Aplicamos el servicio:
```
kubectl apply -f lab-25-C-redis-replica-service.yaml
```

Comprobamos el despliegue del servicio:
```
kubectl get service redis-replica-internal-service
```

La salida será similar a esta:
```
NAME                             TYPE        CLUSTER-IP    EXTERNAL-IP   PORT(S)    AGE
redis-replica-internal-service   ClusterIP   10.109.5.10   <none>        6379/TCP   105s
```

## Ejercicio 4: ***Despliegue del Frontend***


Ahora desplegamos en frontend, que se contectará al servidor máster de Redis o a las réplicas en función del tipo de consulta que se necesite hacer.

Editamos la definición del deployment:
```
code lab-25-C-frontend-deployment.yaml
```

Lo más importante en el archivo de despliegue es:

* *Línea 9-11*: Este deployment se asociará con la plantilla de pod que tiene definidas las etiquetas ***app: guestbook*** y ***tier: frontend***.
* *Línea 12*: El número de réplicas del frontend son 3 (3 pods)
* *Línea 21*: Se está usando la imagen de contenedor ***gb-frontend:v4***.
* *Líneas 27 y 28*: La variable de entorno ***GET_HOST_FROM*** sirve para determinar dónde está el máster de Redis. 
                 
                   El valor es 'env'.

                   El código del Frontend en este punto es:

                   $host = 'redis-master';
                   if (getenv('GET_HOSTS_FROM') == 'env') {
                       $host = getenv('REDIS_MASTER_SERVICE_HOST');
                   }

                   Es decir, que el servidor máster de Redis será el que indique la variable del sistema ***REDIS_MASTER_SERVICE_HOST***. Esta variable la crea Kubernetes cuando se despliega el master.

* *Líneas 29 y 30*: Se inicializa la variable de entorno ***REDIS_SLAVE_SERVICE_HOST*** al valor ***redis-replica-internal-service***, para que los contenedores de Frontend puedan contactar con la réplicas.

Desplegamos el frontend.
```
kubectl apply -f lab-25-C-frontend-deployment.yaml
```

Comprobamos
```
kubectl get deployment frontend
```

La salida será parecida a esta:
```
NAME       READY   UP-TO-DATE   AVAILABLE   AGE
frontend   3/3     3            3           26s
```

Lo vemos con más detalle:
```
kubectl describe deployment frontend
```

Estudiar la salida.

Comprobamos los pods.
```
kubectl get pods
```

La salida será similar a la siguiente. Comprobar que todos los pods están en ***Running***.
```
NAME                                        READY   STATUS    RESTARTS   AGE
frontend-69bff8766c-r5mv9                   1/1     Running   0          4m10s
frontend-69bff8766c-thqh7                   1/1     Running   0          4m10s
frontend-69bff8766c-x67rg                   1/1     Running   0          4m10s
redis-master-deployment-754ccc67d4-ctp9v    1/1     Running   0          14h
redis-replica-deployment-6bf68ddfbd-dvj62   1/1     Running   0          40m
redis-replica-deployment-6bf68ddfbd-zzzdv   1/1     Running   0          40m
```

## Ejercicio 5: ***Despliegue del balanceador para el Frontend***


Vamos a crear un servicio de tipo ***LoadBalancer*** para el frontend.

En K8s existe tres formas de exponer un servicio:

* *ClusterIP*: Por defecto. Se crea una IP para el servicio y k8s redirige el tráfico al nodo apropiado. Al ser una IP privada, el servicio no puede ser accedido desde fuera del cluster.
* *NodePort*: El servicio puede ser accedido desde fuera del cluster, conectándose a la IP (y puerto del nodo.
* *LoadBalancer*: Se creará un balanceador externo, con una IP pública. Se balancea entre los pods.

Editamos el archivo:
```
code lab-25-C-frontend-service.yaml
```

Las líneas más importantes son:

* *Línea 2*: Define un servicio.
* *Línea 9*: Es de tipo ***LoadBalancer***.
* *Línea 11*: Define el puerto en el que escuchará el servicio (***80***)
* *Línea 12*: Puerto del pod al que se reenviará el tráfico (***80***)
* *Líneas 13-15*: Se enviará tráfico a los pods que declaren las etiquetas ***app: guestbook***, ***tier: frontend***.

Aplicamos el servicio:
```
kubectl create -f lab-25-C-frontend-service.yaml
```

Comprobamos los servicios.
```
kubectl get service frontend-load-balancer
```

La salida es:
```
NAME                     TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
frontend-load-balancer   LoadBalancer   10.102.188.44   <pending>     80:30087/TCP   50s
```

El servicio es de tipo ***load balancer*** y la EXTERNAL-IP está en ***Pending***. Cuando se nos asigne la EXTERNAL-IP, debemos conectar con el navegador al puerto ***80***. De este forma, el tráfico iria así:

Navegador --> EXTERNAL-IP:80  --> frontend_load_balancer --> endpoint_pod:80

En otra terminal ejecutamos:
```
minikube tunnel
```

Tomar nota de la IP External del frontend y conectarse con un navegador.

Limpiamos recursos del cluster.
```
kubectl delete deployment frontend redis-master-deployment redis-replica-deployment
kubectl delete service frontend-load-balancer redis-master redis-replica
```

Comprobamos que solo queda el servicio de Kubernetes
```
kubectl get all
```
