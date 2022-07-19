# Laboratorio 20-B: ***Comandos básicos de kubectl***
 
En este laboratorio aprenderemos a usar ***kubectl*** e interactuaremos con el cluster de Kubernetes.

Requisitos:

1. Una máquina virtual con Ubuntu 20.04 LTS a la que poder hacer ssh o escritorio remoto.
2. La VM debe tener descargado el binario de ***Minikube*** y la herramienta ***kubectl***. Es decir, haber realizado el laboratorio anterior (lab-20-A.md)


## Ejercicio 1: ***Primera toma de contacto con kubectl***

Iniciamos Minikube.
```
minikube start --driver=podman
```

En primer lugar verificamos que tenemos un solo nodo en el cluster. El único nodo deberá esta ejecutando el ***Control Plane***.
```
kubectl get nodes
```

Consultamos los pods que está ejecutando el cluster. No deben haber ninguno.
```
kubectl get pods
```

Comprobamos los servicios que están corriendo en el cluster. De haber uno, llamado ***Kubernetes***, que es precisamente el endpoint al que se conecta la herramienta ***kubectl***. En consecuencia, nunca debemos eliminar este servicio.
```
kubectl get services
```

Vamos a proceder a crear objetos en Kubernetes. Para ello hacemos uso del comando ***kubectl create***. La ayuda muestra todos los tipos (kind) de objetos que se pueden crear.
```
kubectl create --help
```

Estos son los objetos que iremos creando a lo largo del curso, ya sea con ***kubectl create*** o, más habitualmente, en archivos de ***manifiesto*** o declarativos con ***sintaxis YAML***.
```
clusterrole         Crea un rol en el cluster.
clusterrolebinding  Asocia un rol de cluster a un objeto/usuario del cluster
configmap           Crea un config map desde un archivo, directorio o valor literal.
cronjob             Programa la ejecucion de un pod.
deployment          Crea deployments.
ingress             Crea un controlador ingress para asociar URIs con los servicios.
job                 Crea un job.
namespace           Crea un espacio de nombres para conseguir aislar objetos entre aplicaciones.
poddisruptionbudget Determina el número mínimos de pods en ejecución durante las operaciones de mantenimiento del cluster.
priorityclass       Crea una clase de prioridad, lo que permite desahuciar (evict) pods con prioridad baja si no hay recursos en el cluster.
quota               Crea cuotas (cpu, memoria, GPU, storage, ...) para los pods 
role                Crea un rol.
rolebinding         Asocia el rol.
secret              Crea un secreto en el cluster.
service             Crea un servicio.
serviceaccount      Crea una cuenta de servicio.
```

## Ejercicio 2: ***Crear un deployment con kubectl***

El pod es la unidad más pequeña de computación en Kubernetes, pero en la práctica los pods se crearán al definir otro objeto, el ***deployment***, que creará los pods.

Para crear un deployment, debemos indicar como mínimo, la imagen que se usará en su contenedor (por ahora uno solo)
```
kubectl create deployment nginx-deployment --image=nginx:latest
```

Para ver los deployments, ejecutamos el siguiente comando.
```
kubectl get deployments
```

La salida del comando debe ser similar a esta:
```
NAME               READY   UP-TO-DATE   AVAILABLE   AGE
nginx-deployment   1/1     1            1           111s
```

El significado de las columnas es el siguiente:

* *NAME*: Muestra el nombre del deployment, en este caso 'nginx-deployment'
* *READY*: Indica cuántos pods (del total) están ejecutándose. En ese caso, el deployment solo tiene un pod, y éste se está ejecutando. 
* *UP-TO-DATE* : Cuando se actualice el deployment, pod ejemplo si cambiamos la imagen de los contenedores de sus pods. Esta columna indicará cuantos pods están actualizados.
* *AVAILABLE*: Indica el número de pods están pasando las pruebas ***readiness***. Esto se verá más adelante en el curso, pero si quieres verlo por adelantado, consulta esto: (https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/)
* *AGE*: Tiempo desde que se creó el deployment.

Los deployments crean los pods. Para ver los pods ejecutamos el siguiente comando.
```
kubectl get pods
```

La salida del comando anterior debe ser similar a esta:
```
NAME                                READY   STATUS    RESTARTS   AGE
nginx-deployment-7fd6754bf7-29ft2   1/1     Running   0          13m
```
El significado y el valor de cada columna es el siguiente:

* *NAME*: Nombre del pod. En este ejemplo ***nginx-deployment-7fd6754bf7-29ft2***. El nombre del pod es una concatenación del nombre del deployment, el nombre del objeto ReplicaSet y un identificador          para el pod. Es decir: ***nginx-deployment*** es el nombre del deployment que creó el pod.***7fd6754bf7*** hace referencia a un objeto ***ReplicaSet***, asociado al deployment. El objeto ReplicaSet es quien está controlando, en todo momento, que el número de pods que forma el deployment sea el deseado. Si algún pod cae, el ReplicaSet creará otro nuevo. En este ejemplo, hemos creado un deployment sin indicar cuántos pods queremos (réplicas), así que solo se instancia un único pod. ***29ft2*** es el identificador del pod dentro del ReplicaSet. Cada pod en el ReplicaSet, tendrá un valor único en este campo.
* *READY*: Indican cúantos CONTENEDORES están corriendo en el pod.
* *STATUS*: Indica el estado de ejecución del pod. Si al menos un contenedor se inicia, pondrá ***Running***.
* *RESTARTS*: Indica cuántos reinicios se han producido en el pod. Si es ***>0***, suele indicar que algo no va bien y, en consecuencia, el orquestador está reiniciando el pod.
* *AGE*: Indica el tiempo transcurrido desde que se creo el pod.

## Ejercicio 3: ***El ReplicaSet***

Como venimos indicando, al crear un deployment, se crea de forma automática otro objeto, de tipo ReplicaSet. Para ver los ReplicaSets del cluster, ejecutamos el siguiente comando:
```
kubectl get replicasets
```

La salida del comando anterior será más o menos así:
```
NAME                          DESIRED   CURRENT   READY   AGE
nginx-deployment-7fd6754bf7   1         1         1       28m
```

Y el significado de las columnas, el siguiente:

* *NAME*: Es el nombre del ReplicaSet, en este ejemplo: ***nginx-deployment-7fd6754bf7***. Este nombre está formado por la concatenación del nombre del deployment ***nginx-deployment*** que creó este ReplicaSet, y un identificador único del ReplicaSet: ***7fd6754bf7***. Este id único también aparece en los nombres de los pods que administra el ReplicaSet.
* *DESIRED*: Como decíamos, la finalidad del ReplicaSet es asegurar que el número de pods declarado en el deployment coincide con el número de los que están instanciados en el cluster. Por lo tanto, esta columna muestra el número de pods que "deseamos" tener corriendo.
* *CURRENT*: Indica cuántos pods ha podido poner en funcionamiento el ReplicaSet. Lo habitual es que el estado deseado, en cuanto a número de pods, coincida con el estado actual del Cluster. Esto no será así en el caso de que el cluster se quede sin recursos y no pueda levantar nuevos pods.
* *READY*: Indica cuántos pods se han iniciado.
* *AGE*: La antiguedad del ReplicaSet.

No debemos tocar el objeto ReplicaSet. Cuando creamos un deployment, éste creará el ReplicaSet. Así mismo, cuando eliminemos el deployment, se eliminará sus ReplicaSets y pods automáticamente.

## Ejercicio 4: ***Editar un deployment con kubectl***


En Kubernetes, es el objeto deployment quien gestiona los pods (nosotros no debemos modificar los pods. Por ejemplo, si deseamos cambiar la versión de la imagen que ejecutan los contenedores de los pods, debemos hacerlo desde el deployment, y no editando los pods que se han creado en el cluster.

En consecuencia, para modificar el deployment, ejecutamos el siguiente comando:
```
kubectl edit deployment nginx-deployment
```

Se abre el editor con la configuración del deployment que hemos creado. La sintaxis y las propiedades que aparecen en este archivo serán tratadas posteriormente en este curso. Por ahora localicemos las única configuración que hemos aportado a la hora de crear el deployment: La imagen de contenedor que ser usará.

El editor usado es ***vim***, que habrá que saber usar. (Pulsar ***R*** y sustituir caracteres hasta quese pulse ***ESC***)

Localizar ***spec/template/spec/containers/image***. Debe aparecer el valor ***nginx:latest***.

Vamos a cambiar la versión de nginx a la ***1.16***, así que procedemos a cambiar ***latest*** por ***1.16***.
```
deberá quedar así:
     spec:
       containers:
       - image: nginx:1.16
```

Guardamos los cambios.

Comprobamos los pods.
```
kubectl get pods
```

La salida del comando anterior será similar a esta:
```
NAME                                READY   STATUS    RESTARTS   AGE
nginx-deployment-6bbc464978-tsrb5   1/1     Running   0          116s
```
Se ha creado un nuevo pod (con la versión 1.16 de nginx en su contenedor) y está en ejecución. Lo importante es darse cuenta que el ***ReplicaSet ha cambiado***. Lo podemos ver en el nombre del pod,
***nginx-deployment-6bbc464978-tsrb5***. El id del ReplicaSet es ahora ***6bbc464978***, que es diferente al que tenía el pod anterior (***7fd6754bf7***).

Este comportamiento ***ES NORMAL*** en Kubernetes. Si se modifica el cambio el en deployment va más allá de cambiar el número de réplicas, por ejemplo, cambiando la imagen de los contenedores, K8s creará un nuevo ReplicaSet asociado a ese deployment. Lo podemos ver así:
```
kubectl get replicasets
```

Que mostrará la siguiente salida.
```
NAME                          DESIRED   CURRENT   READY   AGE
nginx-deployment-6bbc464978   1         1         1       7m26s
nginx-deployment-7fd6754bf7   0         0         0       68m
```

Los ReplicaSets tienen en la parte izquierda de su nombre el nombre del deployment que lo ha creado. Observar cómo el número de pods del nuevo ReplicaSet es el deseado a la vez que el del antiguo se ha puesto a cero. En consecuencia, ahora existe un pod con la versión ***1.16*** de nginx y se ha eliminado el pod con la version ***latest***.


## Ejercicio 5: ***Rollout undo' del deployment con kubectl***

Kubernetes no borra los ReplicaSets. La razón es que siguen ahí por si queremos volver a la versión anterior. Esto se conoce como ***Rollout undo***.

Este comando lista el historial del deployment (sus ReplicaSets):
```
kubectl rollout history deployment nginx-deployment
```

La salida del comando muestra que hay dos revisiones.
```
REVISION  CHANGE-CAUSE
1         <none>
2         <none>
```

La actual es la ***2***. Si queremos volver a la anterior, la ***1***, hacemos:
```
kubectl rollout undo deployment nginx-deployment --to-revision=1
```

Si mostramos los ReplicaSets, podemos comprobar que se ha recuperado el pod anterior y eliminado el actual. Ahora se está ejecutando ***nginx:latest***
```
kubectl get replicasets
```

La salida del comando es:
```
NAME                          DESIRED   CURRENT   READY   AGE
nginx-deployment-6bbc464978   0         0         0       16m
nginx-deployment-7fd6754bf7   1         1         1       77m
```

El id del ReplicaSet activo vuelve a ser ***7fd6754bf7***, y si listamos los pods:
```
kubectl get pods
```

La salida del comando indica que se ha reutilizado el ReplicaSet original y el nuevo pod es gestionado por éste.
```
NAME                                READY   STATUS    RESTARTS   AGE
nginx-deployment-7fd6754bf7-8xhc9   1/1     Running   0          2m52s
```

Se puede comprobar la reversión del despliegue, editando su configuración y comprobando que la versión de nginx vuelve a ser ***latest***.
```
kubectl edit deployment nginx-deployment
```

## Ejercicio 6: ***Describir un objeto con kubectl***

Procedemos a crear otro ejemplo. En este caso un deployment de ***MongoDB***.
```
kubectl create deployment mongo-deployment --image=mongo
```

La salida será similar a la siguiente, en la que podemos ver el deployment ***mongo-deployment***.
```
NAME                                READY   STATUS    RESTARTS   AGE
mongo-deployment-7994f64674-6pnrq   1/1     Running   0          20s
nginx-deployment-7fd6754bf7-8xhc9   1/1     Running   0          4h14m
```

Listamos los pods del deployment. Si queremos filtrar por los pods de un deploymento en concreto usamos el parámetro ***-l app=mongo-deployment***.
```
kubectl get pods -l app=mongo-deployment
```

La salida del comando será similar a la siguiente, mostrando información del único pod del despliegue.
```
NAME                                READY   STATUS    RESTARTS   AGE
mongo-deployment-7994f64674-6pnrq   1/1     Running   0          7m40s
```

Es muy importante saber qué esta ocurriendo dentro de un pod, sobre todo cuando hay problemas y no se inicia correctamente. La opción ***describe*** aplicada a un pod muestra mucha información sobre el pod, incluso el histórico de eventos que se han producido. Ejecutemos el siguiente comando y estudiemos detenidamente su salida. (Nota: cambiar el nombre del pod por el correcto)
```
mi_pod=<poner aquí el nombre del pod>
```
```
kubectl describe pod $mi_pod
```

La salida (parte final de la misma) será similar a la siguiente. Se pueden comprobar los eventos que se han producido dentro del pod.
```
Events:
Type    Reason     Age   From               Message
----    ------     ----  ----               -------
Normal  Scheduled  13m   default-scheduler  Successfully assigned default/mongo-deployment-7994f64674-6pnrq to minikube
Normal  Pulling    13m   kubelet            Pulling image "mongo"
Normal  Pulled     12m   kubelet            Successfully pulled image "mongo" in 15.375152034s
Normal  Created    12m   kubelet            Created container mongo
Normal  Started    12m   kubelet            Started container mongo
```

## Ejercicio 7: ***Ver la salida estándar del contenedor con kubectl***

Otra opción que nos interesa conocer es ***logs***, que muestra la salida estándar de los contenedores que se están ejecutando dentro del pod. Por ejemplo, el siguiente comando mostrará la salida del contenedor ***MongoDB***
```
kubectl logs $mi_pod
```

## Ejercicio 8: ***Ejecutar comandos en el contenedor con kubectl***

En algunas ocasiones necesitaremos ejecutar comandos en el contenedor del pod o, incluso abrir una shell en él. El siguiente comando hace precisamente eso. Los parámetros ***-it*** permiten enviar al contenedor MongoDB nuestra entrada estándar al mismo tiempo que vemos la salida estándar del contenedor. ***--*** es simplemente un ***separador***. Depués de él, ponemos el comando que queremos ejecutar, p.e. ***/bin/bash***
```
kubectl exec -it $mi_pod -- /bin/bash
```

## Ejercicio 9: ***Eliminar objetos del cluster con kubectl***


Como resultado del comando anterior, estaremos dentro de una shell en el contenedor. Probar a ejecutar cualquier comando. Para salir del contenedor, escribimos ***exit*** y pulsamos ***Enter***.

Procedemos a borrar los dos deployments que tenemos corriendo. Al hacerlo se borrarán los ReplicaSets asociados así como los pods (y sus contenedores)
```
kubectl get deployments
kubectl delete deployment mongo-deployment
kubectl delete deployment nginx-deployment
```

Podemos verificar que no hay objetos en el cluster (excepto el servicio Kubernetes) con el siguiente comando.
```
kubectl get all
```

