# Laboratorio 45-A: ***Horizontal POD Autoscaler (HPA) en AZURE***
 
Este laboratorio aprenderemos a usar el HPA.

Requisitos:

1. Una máquina virtual con Ubuntu 20.04 LTS a la que poder hacer ssh o escritorio remoto.
2. ***TENER UNA SUBSCRIPCIÓN DE AZURE PARA DESPLEGAR AKS***
3. Desplegar el cluster AKS en Azure (Ver lab-00.md)

## Ejercicio 1: ***Desplegar la aplicación de ejemplo***

Entramos en el directorio del laboratorio
```
cd ~/k8s_desarrolladores/45
```

Hay dos dimensiones para realizar el escalado de una aplicación en Kubernetes. La primera dimensión del escalado consiste en determinar el ***número de pods***, mientras que la segunda es el ***número de nodos*** del cluster.

El escalado de nodos es una característica nativa en los clusters administrados del cloud, que dejaremos para más adelante. En ese laboratorio aprenderemos a usar HPA, que desplegará las instancias de POD en los nodos que el cluster tenga disponibles (uno solo para Minikube)


Desplegamos la aplicación Guestbook de ejemplo de Google, que es la misma del laboratorio 25-C. Editamos el archivo para comprobarlo:
```
code lab-45-guestbook-all-in-one.yaml
```

Creamos todos los objetos en el cluster:
```
kubectl create -f lab-45-guestbook-all-in-one.yaml
```

Comprobamos que se han iniciado todos los objetos.
```
kubectl get all
```

Vamos a cambiar el servicio de frontend de ClusterIP a ***LoadBalancer***. Para ello editamos el servicio directamente.
```
kubectl edit service frontend
```

Se usa el editor ***vim***. Pulsar ***i*** para insertar, ***x*** para borrar carácter. Cambiar en el código el tipo de ***ClusterIP*** a ***LoadBalancer***. Luego ***:wq!***

Comprobamos que el servicio se ha actualizado a ***LoadBalancer***
```
kubectl get service frontend
```

Salida debe ser así:
```
NAME       TYPE           CLUSTER-IP       EXTERNAL-IP      PORT(S)        AGE
frontend   LoadBalancer   10.107.119.101   10.107.119.101   80:31325/TCP   5m11s
```

Si EXTERNAL-IP se quedan en ***Pending*** esperar un poco porque Azure aún no ha asignado IP. (NOTA: no continuar hasta tener la 'EXTERNAL-IP')

## Ejercicio 2: ***Escalar el frontend de GuestBook***

El frontend tiene actualmente 3 réplicas. Podemos escalarlo dinámicamente con el siguiente comando.
```
kubectl scale deployment frontend --replicas=6
```

Para ver dónde están corriendo los 6 pods...
```
kubectl get pods -o wide
```

La salida del comando anterior será similar a esta:
```
NAME                             READY   STATUS    RESTARTS   AGE     IP            NODE                                NOMINATED NODE   READINESS GATES
frontend-7c8cb4c59f-68tqs        1/1     Running   0          17s     10.244.1.8    aks-nodepool1-19313883-vmss000001   <none>           <none>
frontend-7c8cb4c59f-6sbpw        1/1     Running   0          17s     10.244.1.7    aks-nodepool1-19313883-vmss000001   <none>           <none>
frontend-7c8cb4c59f-7sp76        1/1     Running   0          4m44s   10.244.0.10   aks-nodepool1-19313883-vmss000000   <none>           <none>
frontend-7c8cb4c59f-c9tjr        1/1     Running   0          4m44s   10.244.1.5    aks-nodepool1-19313883-vmss000001   <none>           <none>
frontend-7c8cb4c59f-lgdcr        1/1     Running   0          17s     10.244.0.11   aks-nodepool1-19313883-vmss000000   <none>           <none>
frontend-7c8cb4c59f-pgbj8        1/1     Running   0          4m44s   10.244.1.6    aks-nodepool1-19313883-vmss000001   <none>           <none>
redis-master-f46ff57fd-qqk4h     1/1     Running   0          4m44s   10.244.1.3    aks-nodepool1-19313883-vmss000001   <none>           <none>
redis-replica-5bc7bcc9c4-nlhxv   1/1     Running   0          4m44s   10.244.1.4    aks-nodepool1-19313883-vmss000001   <none>           <none>
redis-replica-5bc7bcc9c4-wj2ws   1/1     Running   0          4m44s   10.244.0.9    aks-nodepool1-19313883-vmss000000   <none>           <none>
```

Comprobar que hay 6 instancias de pod del frontend y que están repartidas entre los dos nodos workers. (NOTA: Las columnas ***NOMINATED NODE*** y ***READINESS GATES*** serán tratadas en un laboratorio posterior)

En K8s se puede configurar el autoescalado de los pods usando el objeto HPA. HPA monitoriza las métricas a intervalos regulares, basándose en las reglas que definimos. Por ejemplo, que se añadan pods adicionales si el uso de CPU es superior al 50% y que se quiten si es inferior al 10%.

Para probar HPA, primero reducimos el escalado de forma manual.
```
kubectl scale deployment frontend --replicas=1
```

Comprobar que solo hay un pod en el frontend
```
kubectl get pods
```

La salida debe ser similar a esta:
```
frontend-7c8cb4c59f-7sp76        1/1     Running   0          6m19s
redis-master-f46ff57fd-qqk4h     1/1     Running   0          6m19s
redis-replica-5bc7bcc9c4-nlhxv   1/1     Running   0          6m19s
redis-replica-5bc7bcc9c4-wj2ws   1/1     Running   0          6m19s
```

Abrimos el archivo de ejemplo para configurar el Horizontal POD Autoscaler ***lab-45-hpa.yaml***.
```
code lab-45-hpa.yaml
```

Las líneas más importantes son:

* *Línea 2*: Se define el tipo de objeto a crear como ***HorizontalPodAutoscaler***.
* *Línea 6-9*: Lo asocia con el deployment que se auto escalará.
* *Línea 10-11*: Configuramos número mínimo y máximo de pods para el autoescalado.
* *Línea 12*: Definimos el porcentaje de CPU para que se dispare el autoescalado.

Salimos sin modificar.

Realizamos el despliegue del HPA.
```
kubectl create -f lab-45-hpa.yaml
```

Comprobamos el objeto de autoescalado.
```
kubectl get hpa frontend-scaler
```

La salida será similar a esta:
```
NAME              REFERENCE             TARGETS         MINPODS   MAXPODS   REPLICAS   AGE
frontend-scaler   Deployment/frontend   <unknown>/50%   1         10        1          65s
```

En la columna ***TARGETS*** aparece el valor ***<unknown>*** porque el HPA aun no se ha puesto en funcionamiento y no tiene estadísticas que mostrar. Ahora vamos a ver el autoescalado en acción. Creamos una nueva terminal y ejecutamos el comando.
```
kubectl get pods -w
```

La idea es que vayamos viendo en tiempo real cómo se van creando los pods. En la terminal original vamos a instalar un programa llamado ***hey*** cuya finalidad es generar carga en el servidor web, pero para ello, antes instalamos ***go***.
```
sudo apt-get update
```
```
sudo apt install -y hey
```

Copiamos la EXTERNAL-IP del servicio de frontend.
```
kubectl get service frontend
```

La salida será parecida a esto:
```
NAME       TYPE           CLUSTER-IP     EXTERNAL-IP     PORT(S)        AGE
frontend   LoadBalancer   10.0.159.221   20.126.186.17   80:31450/TCP   9m39s
```

Metemos carga (20 millones de Request) Usar la IP Externa del servicio de frontend.
```
hey -z 20m http://<Poner aquí la IP Externa del Servicio>
```

Esperar unos minutos (tener paciencia) y ver en la segunda terminal como se escalan los pods.

Hacer ***CTRL+C*** en ***hey*** y esperar un tiempo para ver como se quitan los pods. Tarda un buen rato (al menos 10 minutos) Lo podemos ver con detalle describiendo el HPA.

En una tercera terminal:
```
kubectl describe hpa frontend-scaler
```

Parar (CTRL+C) ***hey*** y comprobar como se hace el desescalado (Tarda un buen rato)

Limpiamos los recursos:
```
kubectl delete -f lab-45-hpa.yaml
kubectl delete -f lab-45-guestbook-all-in-one.yaml
```

Comprobamos
```
kubectl get all
```

Solo debe quedar el servicio de Kubernetes.
```
kubectl get all
```

## Ejercicio 3: ***Autoescalado de nodos***

Se puede hacer con la GUI, pero por comodidad lo hacemos con ***az aks***. Primero obtenemos el nombre del nodepool del cluster:
```
az aks nodepool list --resource-group myaks-rg --cluster-name myaks | grep name
```

Tomar el nombre del NopePool de la salida anterior. Vamos a poner el número de nodos del NodePool a 1:
```
az aks nodepool scale \
    --name <Poner aquí en nombre del NodePool>  \
    --node-count 1 \
    --resource-group myaks-rg \
    --cluster-name myaks
```

Mostramos los nodos.
```
kubectl get nodes
```

La salida debe ser como esta: (Nota: Solo debe quedar un único nodo)
```
NAME                                STATUS   ROLES   AGE   VERSION
aks-nodepool1-19313883-vmss000000   Ready    agent   38m   v1.21.9
```

Desplegamos la app en el único nodo que tenemos:
```
kubectl create -f lab-45-guestbook-all-in-one.yaml
```

Confirmamos que se ha desplegado correctamente:
```
kubectl get all
```

La salida debe ser similar a esta:
```
NAME                                 READY   STATUS    RESTARTS   AGE
pod/frontend-7c8cb4c59f-74bpf        1/1     Running   0          11s
pod/frontend-7c8cb4c59f-csx4w        1/1     Running   0          11s
pod/frontend-7c8cb4c59f-lb2mh        1/1     Running   0          11s
pod/redis-master-f46ff57fd-kmpm5     1/1     Running   0          11s
pod/redis-replica-5bc7bcc9c4-45vbr   1/1     Running   0          11s
pod/redis-replica-5bc7bcc9c4-6kxrf   1/1     Running   0          11s
 
NAME                    TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)    AGE
service/frontend        ClusterIP   10.0.56.255    <none>        80/TCP     11s
service/kubernetes      ClusterIP   10.0.0.1       <none>        443/TCP    42m
service/redis-master    ClusterIP   10.0.250.238   <none>        6379/TCP   11s
service/redis-replica   ClusterIP   10.0.114.105   <none>        6379/TCP   11s
 
NAME                            READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/frontend        3/3     3            3           11s
deployment.apps/redis-master    1/1     1            1           11s
deployment.apps/redis-replica   2/2     2            2           11s
 
NAME                                       DESIRED   CURRENT   READY   AGE
replicaset.apps/frontend-7c8cb4c59f        3         3         3       11s
replicaset.apps/redis-master-f46ff57fd     1         1         1       11s
replicaset.apps/redis-replica-5bc7bcc9c4   2         2         2       11s
```

El autoescalado del cluster determina el número de pods que no pueden ser planificados debido a recursos insuficientes (falta CPU o de RAM en los nodos). Para demostrar el autoescalado, forzaremos el deployment de forma que no se puedan poner a correr todos los pods. 

Para forzar que el cluster se quede sin recursos, modificamos el deployment ***redis-replica*** subiéndolo a ***5 pods*** de Redis (Debería ser suficiente para agotar la CPU)
```
kubectl scale deployment redis-replica --replicas 5
```

Comprobamos que hay pods que no se pueden planificar porque no hay recursos. Importante: Mostrarán el estado ***Pending***.
```
kubectl get pods
```

La salida será parecida a esta: (Nota: Comprobar que hay dos réplicas en ***Pending***)
```
NAME                             READY   STATUS    RESTARTS   AGE
frontend-7c8cb4c59f-74bpf        1/1     Running   0          3m
frontend-7c8cb4c59f-csx4w        1/1     Running   0          3m
frontend-7c8cb4c59f-lb2mh        1/1     Running   0          3m
redis-master-f46ff57fd-kmpm5     1/1     Running   0          3m
redis-replica-5bc7bcc9c4-45vbr   1/1     Running   0          3m
redis-replica-5bc7bcc9c4-6kxrf   1/1     Running   0          3m
redis-replica-5bc7bcc9c4-9dlvk   0/1     Pending   0          32s
redis-replica-5bc7bcc9c4-dvsqf   1/1     Running   0          32s
redis-replica-5bc7bcc9c4-fhx2v   0/1     Pending   0          32s
```

Para ver que no hay CPU, ejecutar este comando:
```
kubectl describe pod <Poner el nombre de un pod que esté en estado 'pending'>
```

La salida (recortada) mostrará algo como lo siguiente: (Nota: Localizar el mensaje ***Insufficient cpu***)
```
Events:
   Type     Reason            Age                  From               Message
   ----     ------            ----                 ----               -------
    Warning  FailedScheduling  11s (x4 over 2m25s)  default-scheduler  0/1 nodes are available: 1 Insufficient cpu.
```

Ahora habilitamos el autoescalado del cluster:
```
az aks nodepool update \
    --name nodepool1 \
    --enable-cluster-autoscaler \
    --resource-group myaks-rg \
    --cluster-name myaks \
    --min-count 1 \
    --max-count 2
```

Hay que esperar unos minutos (aprox. 5) a que se despliegue el nuevo nodo y se replanifiquen los pods. Comprobarlo con:
```
kubectl get pods
```
```
kubectl get nodes
```

Limpiamos los recursos:
```
kubectl delete -f lab-45-guestbook-all-in-one.yaml
```

Comprobamos que se han eliminado:
```
kubectl get all
```

Quitamos autoescalado, ponemos 2 nodos (o destruimos el cluster: Ver lab-00.md Ejercicio 3):
```
az aks nodepool update \
    --name nodepool1 \
    --disable-cluster-autoscaler \
    --resource-group myaks-rg \
    --cluster-name myaks 
```

Ponemos dos nodos:
```
az aks nodepool scale \
    --name nodepool1 \
    --node-count 2 \
    --resource-group myaks-rg \
    --cluster-name myaks
```

Volvemos a poner el contexto en Minikube:
```
kubectl config use-context minikube
```
