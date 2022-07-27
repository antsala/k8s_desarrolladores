# Laboratorio 55-A: ***Usar volúmenes en la aplicaciones***
 
En este laboratorio aprenderemos a usar Volúmenes, ConfigMaps y Secrets. Vamos desplegar una app con un frontend y backend de redis (un master y dos réplicas)

Requisitos:

1. Una máquina virtual con Ubuntu 20.04 LTS a la que poder hacer ssh o escritorio remoto.
2. Tener un cluster con con plugin de ***Storage Class*** configurado, por ejemplo AKS.


## Ejercicio 1: ***Despliegue del servidor Redis***

Entramos en el directorio del laboratorio
```
cd ~/k8s_desarrolladores/55
```

Estudiamos el YAML:
```
code lab-55-A-redis-master-deployment.yaml
```

Desplegamos el backend del maestro de redis:
```
kubectl apply -f lab-55-A-redis-master-deployment.yaml
```

Examinamos el Despliegue. (Nota: No continuar hasta que se haya desplegado completamente)
```
kubectl get all
```

Comprobar los detalles del deployment:
```
kubectl describe deployment/redis-master
```

Deseamos pasar configuraciones al servidor Redis así que borramos el deployment actual desde la línea de comandos:
```
kubectl delete -f lab-55-A-redis-master-deployment.yaml
```

Comprobar que se elimina:
```
kubectl get all
```

Ahora vamos a hacer lo mismo pero usando un ***ConfigMap***. Hay dos formas de crear un ConfigMap: Desde un archivo de texto o desde un archivo yaml.

## Ejercicio 2: ***Creación de un ConfigMap desde un archivo***


Creamos un archivo, llamado ***lab-55-A-redis-config*** y le ponemos estas dos líneas:
   maxmemory 2mb
   maxmemory-policy allkeys-lru

Realmente no hace falta hacerlo porque en el directorio de ejemplos ya existe ese archivo. Así que nos limitamos a abrirlo para comprobarlo:
```
code lab-55-A-redis-config
```

***allkeys-lru*** --> lru = Less Recently Used

Para crear el ConfigMap, ejecutamos el siguiente comando:
```
kubectl create configmap example-redis-config --from-file=lab-55-A-redis-config
```

Ahora listamos el configmap:
```
kubectl get configmaps
```

Observar que hay una clave llamada ***redis-config*** y una lista con los valores de la configuración:
```
kubectl describe configmap/example-redis-config
```

Practicamos la segunda forma de crear el configmap, es decir, desde un archivo YAML.


## Ejercicio 3: ***Creación de un ConfigMap desde archivo YAML***

Para ello borramos el objeto recién creado:
```
kubectl delete configmap/example-redis-config
```

En este caso partiríamos de un archivo YAML, que coincidirá exactamente con el que se creó anteriormente. Lo visualizamos. (Nota: ***|-*** se utiliza para definir un string en varias líneas. Elimina el salto de linea al final y los espacios en blanco si los hubiera)
```
code lab-55-A-example-redis-config.yaml
```

Volvemos a crear el configmap, pero esta vez directamente desde el YAML:
```
kubectl create -f lab-55-A-example-redis-config.yaml
```

Y lo describimos:
```
kubectl describe configmap/example-redis-config
```

K8s tiene una opción muy útil, la ***-o*** (***--output***), que puede ser usada para obtener la salida de cualquier objeto presente en el cluster en formato YAML o JSON. Así se puede generar el archivo yaml del objeto. Lo prodriamos redireccionar a un archivo si lo vemos necesario:
```
kubectl get configmap/example-redis-config --output yaml 
```

Ahora vamos a usar un ConfigMap para pasar información al contenedor en tiempo de ejecución. Vamos a abrir el archivo modificado:
```
code lab-55-A-redis-master-deployment_modified.yaml
```

Observar lo siguiente:

* *Líneas 24-26*: Lanza el contenedor de redis pasándole información sobre el archivo de configuración que debe utilizar: ***/redis-master/redis.conf***.
* *Líneas 27-29*: Pasa al contenedor una variable de entorno: ***MASTER=true***.
* *Líneas 30-32*: Monta un volumen llamado ***config*** en la ruta ***/redis-master*** del contenedor.
* *Líneas 39-45*: Almacena el valor de la clave ***redis-config*** del configmap ***example-redis-config*** en el archivo ***redis.conf*** en el volumen ***config***. Cuando el contenedor arranca, lee el archivo ***/redis-master/redis.conf***, que contiene los valores del configmap que se creó.

Desplegamos esta versión actualizada para usar configmap:
```
kubectl create -f lab-55-A-redis-master-deployment_modified.yaml
```

Miramos los pods:
```
kubectl get pods
```

Ejecutamos un comando dentro del pod para comprobar si realmente ha leído los valores de configuración. (Nota: Cambiar el id del pod por el apropiado. ***--*** indica que lo que viene después es el comando que ejecutan los contenedores del pod. En este caso, el pod tiene un único contenedor)
```
kubectl exec -it redis-master-<Poner aquí el id apropiado> -- redis-cli
```

El comando anterior debe abrir una conexión con REDIS. Ejecutamos los siguientes comandos para verificar que hay ***2MB*** y la politica es ***LRU*** (Borra las claves que se usan menos Less Recently Used cuando le haga falta memoria):
```
CONFIG GET maxmemory
```
```
CONFIG GET maxmemory-policy
```

Salimos con ***exit***.

Ahora vamos a desplegar un servicio para los pods del deployment redis-master. Echamos un vistazo al YAML:
```
code lab-55-A-redis-master-service.yaml
```

Desplegamos:
```
kubectl apply -f lab-55-A-redis-master-service.yaml
```

Comprobamos el despliegue del servicio:
```
kubectl get service
```

Verificar que el tipo de servicio es ***ClusterIP***, por lo que solo se puede acceder a él desde dentro del cluster (desde otros pods) y no desde el exterior del cluster. Un servicio también introduce un nombre DNS para dicho servicio, en la forma ***<nombreServicio>.<espacioDeNombres>svc.cluster.local***. Puesto que estamos usando el espacio de nombres ***default***, la DNS del servicio ***redis-master*** es ***redis-master.default.svc.cluster.local***.

Para ver esto funcionando, nos metemos en el pod para ver si hay resolución DNS. Lo vamos a hacer con ***ping*** ya que ***nslookup*** no está instalado en el pod. No va a haber respuesta de ping, solo nos interesa la resolución DNS.

Listamos pods y nos quedamos con su id:
```
kubectl get pods
```

```
kubectl exec -it redis-master-<Poner aquí el ID apropiado> -- ping redis-master
```

Comprobar que el registro A ****redis-master.default.svc.cluster.local***. se resuelve a la ***IP de cluster*** anterior. CTRL+C para salir.

Convolución. Se completa el dominio si no se especifica completamente:
```
kubectl exec -it redis-master-<Poner aquí el ID apropiado> -- ping redis-master
```
```
kubectl exec -it redis-master-<Poner aquí el ID apropiado> -- ping redis-master.default
```
```
kubectl exec -it redis-master-<Poner aquí el ID apropiado> -- ping redis-master.default.svc.cluster.local
```

Ahora vamos a desplegar las réplicas (2) de redis, que se sincronizarán desde redis-master. Estudiamos el YAML:
```
code lab-55-A-redis-replica-deployment.yaml
```

* *Línea 13*: Desplegamos 2 pods (2 réplicas)
* *Línea 23*: Se usa la imagen ***follower***, que se conectará a una máquina llamada ***redis-server***. (Realmente es un servicio)
* *Línea 29-30*: Le decimos que la resolución se hace por DNS. De esta forma, los followers se conectarán a la IP asociada a ***redis-server***, es decir a ***redis-server.default.svc.cluster.local***, que es la IP del pod de redis-server. (Aun nos queda por crear el servicio)

Aplicamos:
```
kubectl apply -f lab-55-A-redis-replica-deployment.yaml
```

Comprobamos que arrancan los pods:
```
kubectl get all
```

Para que el frontend (aún por desplegar) pueda contactar con las réplicas (además del master de redis), es necesario exponerlas mediante un servicio, que nos dará la respectiva ClusterIP:
```
kubectl apply -f lab-55-A-redis-replica-service.yaml
```

Comprobamos el despliegue:
```
kubectl get all
```

Ahora desplegamos en frontend.
```
code lab-55-A-frontend-deployment.yaml
```

Lo más importante en el archivo de despliegue anterior es:

* *Línea 12*: El número de réplicas del frontend son 3 (***3 pods***)
* *Línea 10-11*:  Las etiquetas son ***app: guestbook***, ***tier: frontend***.
* *Línea 21*: Se está usando la imagen de contenedor ***gb-frontend:v4***.

en K8s existe tres formas de exponer un servicio:

* *ClusterIP*: Por defecto. Se crea una IP para el servicio y k8s redirige el tráfico al nodo apropiado. Al ser una IP interna, el servicio no puede ser accedido desde Internet.
* *NodePort*: El servicio puede ser accedido desde fuera del cluster, conectándose a la IP (y puerto del nodo.
* *LoadBalancer*: Se creará un balanceador externo, con una IP pública. Se balancea entre los pods.

Desplegamos el frontend:
```
kubectl apply -f lab-55-A-frontend-deployment.yaml
```

Comprobamos:
```
kubectl get pods
```

Vamos a crear un servicio de tipo ***LoadBalancer*** para el frontend.
```
code lab-55-A-frontend-service.yaml
```

Aplicamos:
```
kubectl create -f lab-55-A-frontend-service.yaml
```

Comprobamos los servicios. La IP Externa tarda un rato, hasta que se cree la regla en el balanceador:
```
kubectl get service
```

Tomar nota de la IP External del frontend y conectarse con un navegador.


Limpiamos recursos del cluster:
```
kubectl delete deployment frontend redis-master redis-replica
kubectl delete service frontend redis-master redis-replica
```

Comprobamos que solo queda el servicio de Kubernetes
```
kubectl get all
```

## Ejercicio 4: ***Instalar aplicaciones con estado en el cluster***

(Nota: Si tienes un error, antes de reintentarlo, borra los PVCs y los PVs)

Como ya hemos visto, Helm es el administrador de paquetes de Kubernetes. Permite desplegar, actualizar, y administrar las aplicaciones de Kubernetes. Para ello, se escriben los denominados ***charts***'***. 

Instalación de Helm. (Nota: instalar Helm en Ubuntu solo si no estuviera instalado)
```
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3
chmod 700 get_helm.sh
```
```
./get_helm.sh
```
```
source ~/.profile
```

Vamos a instalar WordPress por medio de Helm. En primer lugar añadimos el repositorio que contiene los charts de Helm:
```
helm repo add bitnami https://charts.bitnami.com/bitnami
```

Ahora instalamos WordPress desde Helm:
```
helm install myakswp bitnami/wordpress
```

Tarda un rato en desplegarse. comprobarlo:
```
kubectl get all
```

Mientras se despliega, vamos a ver las ***Persistent Volume Claims***.

## Ejercicio 5: ***Persistent Volume Claims***.

Las ***PVCs*** se utilizan para abstraer el almacenamiento del proveedor (cloud) subyacente. El chart de instalación de ***WordPress*** depende del chart de instalación de ***MariaDB*** para la instalación de su base de datos. A diferencia de las aplicaciones sin estado, como los frontends que hemos usado, ***MariaDB*** requiere una cuidadosa gestión del almacenamiento. Para hacer que K8s administre cargas con estado, se define un objeto específico llamado ***StatefulSet***. 

Un ***Statefulset*** es como un Deployment con capacidades adicionales que se asocian con pods individuales. Esto significa que k8s se asegurará que el ***pod y su almacenamiento se mantengan juntos***. Otra curiosidad es que los StatefulSets nombran los pods con ***números***, en lugar de un id aleatorio. Comprobamos el estado. Observar cómo se ha nombrado al pod de MariaDB (***pod/myakswp-mariadb-0***)
```
kubectl get pods
```

Otra diferencia es como se administra la eliminación del pod. Cuando el pod de un Deployment se elimina, K8s lo replanificará en cualquier nodo, mientras que cuando se borra un pod de un StatefulSet, K8s lo planificará solamente en el nodo que estaba corriendo. Solo cambiará la ubicación del pod si el nodo donde corría se quita del cluster o no está disponible.

Normalmente, desearemos conectar almacenamiento a un StatefulSet. Para ello se requiere un ***PersistentVolume (PV)***. Este volumen persistente puede ser respaldado por diversos mecanismos, como bloques, blobs, EBS, iSCSI, NFS, ...

Los StatefulSets requieren, o un volumen pre-aprovisionado, o un volumen aprovisionado dinámicamente mediante una ***PersistentVolumeClaim (PVC)***. Una PVC permite al usuario solicitar almacenamiento de forma dinámica, lo que resultará en la creación de un ***PersistentVolume (PV)***.

En este ejemplo de WordPress, se está usando una PVC. La PVC proporciona una abstracción sobre los mecanismos de almacenamiento subyacentes.

Observemos lo que hizo el charts de Helm de MariaDB con el siguiente comando:
```
kubectl get statefulset -o yaml > mariadbss.yaml
code mariadbss.yaml
```

Las líneas más relevantes son:

* *Línea 4*: Se declara un ***StatefulSet***.
* *Líneas 119-124*: Monta el volumen definido como ***data*** en el path ***/bitnami/mariadb***.  
* *Lineas 140-157*: Estas líneas declaran el PVC. Concretamente...
* *Línea 149*: Le asigna el nombre ***data***, que será reutilizado en la línea 119 anterior.
* *Línea 152*: Establece el modo de acceso ***ReadWriteOnce***, que creará un almacenamiento de bloque, que en Azure es un disco. También tenemos los modos ***ReadOnlyMany*** y ***ReadWriteMany***. Como su nombre sugiere, un volumen ***ReadWriteOnce*** solo puede ser conectado a un único pod, mientras que un ***ReadOnlyMany*** o un ***ReadWriteMany*** pueden ser conectados a diferentes pods a la vez. Estos dos últimos requieren un mecanismo de almacenamiento subyacente del tipo ***Azure Files*** o ***Azure Blob***.
* *Línea 156*: Define el tamaño del disco.

En resumen, K8s crea de forma dinámica (a través de la clase de almacenamiento) y conecta un volumen de 8 GiB a este pod. El aprovisionador de almacenamiento dinámico es de tipo ***disco***. Estos aprovisionadores (clases de almacenamiento) se configuraron al crear el cluster.


Para ver las clases de almacenamiento disponibles en el cluster, ejecutamos el siguiente comando. (Nota: CSI = Container Storage Interface)
```
kubectl get storageclass
```

Para mostrar los detalles del PVC:
```
kubectl get pvc
```

El nombre del PVC puede buscarse en los recursos de Azure (discos).

El concepto del PVC abstrae las cuestiones específicas del proveedor cloud. Esto permite que la misma plantilla de Helm funcione en ***Azure***, ***AWS*** o ***GCP (Google Cloud Platform)***. 

En AWS se utilizará un ***Elastic Block Storage (EBS)***, mientras que en GCP será un ***Persistent Disk***, y en Azure un ***Disco***, un ***blob*** o un ***Azure File***.

Los PVCs se pueden crear directamente sin la necesidad de usar Helm. Para comprobar el despliegue, usamos el comando:
```
helm ls
```

Información de estado:
```
helm status myakswp
```

Para ver los objetos de K8s que ha creado Helm:
```
kubectl get all
```

Para ver el estado del despliegue del chart:
```
helm status myakswp
```

Para desinstalar el chart:
```
helm delete myakswp
```

Comprobar que se eliminan los objetos de k8s:
```
kubectl get all
```

Comprobar que los PVs y las PVCs no se eliminan. La eliminación del pod no implica su eliminación.
```
kubectl get pvc
```
```
kubectl get pv
```

Borramos el PVC, que a su vez borrará el PV.
```
kubectl delete pvc --all
```

