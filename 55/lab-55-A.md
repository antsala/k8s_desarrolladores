# Laboratorio 55-A: ***Usar volúmenes en la aplicaciones***
 
En este laboratorio aprenderemos a usar Volúmenes, ConfigMaps y Secrets. Vamos desplegar una app con un frontend y backend de redis (un master y dos réplicas)

Requisitos:

1. Una máquina virtual con Ubuntu 20.04 LTS a la que poder hacer ssh o escritorio remoto.
2. Tener un cluster con con plugin de ***Storage Class*** configurado, por ejemplo AKS (también vale Minikube)


## Ejercicio 1: ***Instalar aplicaciones con estado en el cluster***

(Nota: Si tienes un error, antes de reintentarlo, borra los PVCs y los PVs)

Como ya hemos visto, Helm es el administrador de paquetes de Kubernetes. Permite desplegar, actualizar, y administrar las aplicaciones de Kubernetes. Para ello, se escriben los denominados ***charts***. 

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

Vamos a instalar ***WordPress*** por medio de Helm. En primer lugar añadimos el repositorio que contiene los charts de Helm:
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

## Ejercicio 2: ***Persistent Volume Claims***.

Las ***PVCs*** se utilizan para abstraer el almacenamiento del proveedor (cloud) subyacente. El chart de instalación de ***WordPress*** depende del chart de instalación de ***MariaDB*** para la instalación de su base de datos. A diferencia de las aplicaciones sin estado, como los frontends que hemos usado, ***MariaDB*** requiere una cuidadosa gestión del almacenamiento. Para hacer que K8s administre cargas con estado, se define un objeto específico llamado ***StatefulSet***. 

Un ***Statefulset*** es como un Deployment con capacidades adicionales que se asocian con pods individuales. Esto significa que k8s se asegurará que el ***pod y su almacenamiento se mantengan juntos***. Otra curiosidad es que los StatefulSets nombran los pods con ***números***, en lugar de un id aleatorio. Comprobamos el estado. Observar cómo se ha nombrado al pod de MariaDB (***pod/myakswp-mariadb-0***)
```
kubectl get pods
```

Otra diferencia es como se administra la eliminación del pod. Cuando el pod de un Deployment se elimina, K8s lo replanificará en cualquier nodo, mientras que cuando se borra un pod de un StatefulSet, K8s lo planificará solamente en el nodo que estaba corriendo. Solo cambiará la ubicación del pod si el nodo donde corría se quita del cluster o no está disponible.

Normalmente, desearemos conectar almacenamiento a un StatefulSet. Para ello se requiere un ***PersistentVolume (PV)***. Este volumen persistente puede ser respaldado por diversos mecanismos, como bloques, blobs, EBS, iSCSI, NFS, ...

Los StatefulSets requieren, o un volumen pre-aprovisionado, o un volumen aprovisionado dinámicamente mediante una ***PersistentVolumeClaim (PVC)***. Una PVC permite al usuario solicitar almacenamiento de forma dinámica, lo que resultará en la creación de un ***PersistentVolume (PV)***.

En este ejemplo de WordPress, se está usando una PVC. La PVC proporciona una abstracción sobre los mecanismos de almacenamiento subyacentes.

Observemos lo que hizo el chart de Helm de MariaDB con el siguiente comando:
```
kubectl get statefulset -o yaml > mariadbss.yaml
nano mariadbss.yaml
```

Las líneas más relevantes son:

* *Línea 4*: Se declara un ***StatefulSet***.
* *Líneas 119-124*: Monta el volumen definido como ***data*** en el path ***/bitnami/mariadb***.  
* *Lineas 140-157*: Estas líneas declaran el PVC. Concretamente...
* *Línea 149*: Le asigna el nombre ***data***, que será reutilizado en la línea 119 anterior.
* *Línea 152*: Establece el modo de acceso ***ReadWriteOnce***, que creará un almacenamiento de bloque, que en Azure es un disco. También tenemos los modos ***ReadOnlyMany*** y ***ReadWriteMany***. Como su nombre sugiere, un volumen ***ReadWriteOnce*** solo puede ser conectado a un único pod, mientras que un ***ReadOnlyMany*** o un ***ReadWriteMany*** pueden ser conectados a diferentes pods a la vez. Estos dos últimos requieren un mecanismo de almacenamiento subyacente del tipo ***Azure Files*** o ***Azure Blob***.
* *Línea 155*: Define el tamaño del disco.

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

