# Laboratorio 45-B: ***Horizontal POD Autoscaler (HPA) en AWS***
 
En este laboratorio aprenderemos a usar el HPA.

Requisitos:

1. Una máquina virtual con Ubuntu 20.04 LTS a la que poder hacer ssh o escritorio remoto.
2. ***TENER UNA SUBSCRIPCIÓN DE AWS PARA DESPLEGAR EKS***.
3. Desplegar el cluster EKS en AWS (Ver lab-01.md)


## Ejercicio 1: ***Desplegar la aplicación de ejemplo***

Entramos en el directorio del laboratorio
```
cd ~/k8s_desarrolladores/45
```

Hay dos dimensiones para realizar el escalado de una aplicación en Kubernetes. La primera dimensión del escalado consiste en determinar el número de pods, mientras que la segunda es el número de nodos del cluster.

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

Se usa el editor ***vim***. Pulsar ***i***. Cambiar en el código el tipo de ***ClusterIP*** a ***LoadBalancer***. Luego ***:wq!***.

Comprobamos que el servicio se ha actualizado a ***LoadBalancer***.
```
kubectl get service frontend
```

Salida debe ser así:
```
NAME       TYPE           CLUSTER-IP       EXTERNAL-IP                                                               PORT(S)        AGE
frontend   LoadBalancer   10.100.157.230   aa554514004624ce9b896ff5781c9995-2125204655.eu-west-1.elb.amazonaws.com   80:30689/TCP   101s
```

Si 'EXTERNAL-IP' se quedan en ***Pending*** esperar un poco porque AWS aún no ha asignado IP. (NOTA: no continuar hasta tener la ***EXTERNAL-IP***)

Podemos conectar con un navegador a: 'http://aa554514004624ce9b896ff5781c9995-2125204655.eu-west-1.elb.amazonaws.com' para ver la app funcionando en EKS.


## Ejercicio 2: ***Escalar el frontend de GuestBook***

Para que funcione HPA en AWS, debemos instalar el servidor de métricas, de esta forma:
```
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

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
NAME                             READY   STATUS    RESTARTS   AGE     IP               NODE                                           NOMINATED NODE   READINESS GATES
frontend-7c8cb4c59f-22h46        1/1     Running   0          6m33s   192.168.50.134   ip-192-168-40-174.eu-west-1.compute.internal   <none>           <none>
frontend-7c8cb4c59f-685nj        1/1     Running   0          7s      192.168.32.166   ip-192-168-40-174.eu-west-1.compute.internal   <none>           <none>
frontend-7c8cb4c59f-cdxkz        1/1     Running   0          6m34s   192.168.14.70    ip-192-168-24-240.eu-west-1.compute.internal   <none>           <none>
frontend-7c8cb4c59f-nflmn        1/1     Running   0          7s      192.168.91.201   ip-192-168-91-163.eu-west-1.compute.internal   <none>           <none>
frontend-7c8cb4c59f-p5cch        1/1     Running   0          7s      192.168.23.78    ip-192-168-24-240.eu-west-1.compute.internal   <none>           <none>
frontend-7c8cb4c59f-qgdxw        1/1     Running   0          6m34s   192.168.71.87    ip-192-168-91-163.eu-west-1.compute.internal   <none>           <none>
redis-master-f46ff57fd-xzfkr     1/1     Running   0          6m34s   192.168.44.220   ip-192-168-40-174.eu-west-1.compute.internal   <none>           <none>
redis-replica-5bc7bcc9c4-j6f6p   1/1     Running   0          6m34s   192.168.94.225   ip-192-168-91-163.eu-west-1.compute.internal   <none>           <none>
redis-replica-5bc7bcc9c4-wsnvk   1/1     Running   0          6m34s   192.168.11.181   ip-192-168-24-240.eu-west-1.compute.internal   <none>           <none>
```

Comprobar que hay 6 instancias de pod del frontend y que están repartidas entre los tres nodos workers. (NOTA: Las columnas ***NOMINATED NODE*** y ***READINESS GATES*** serán tratadas en un laboratorio posterior)

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
NAME                             READY   STATUS    RESTARTS   AGE
frontend-7c8cb4c59f-qgdxw        1/1     Running   0          8m2s
redis-master-f46ff57fd-xzfkr     1/1     Running   0          8m2s
redis-replica-5bc7bcc9c4-j6f6p   1/1     Running   0          8m2s
redis-replica-5bc7bcc9c4-wsnvk   1/1     Running   0          8m2s
```

Abrimos el archivo de ejemplo para configurar el Horizontal POD Autoscaler 'lab-45-hpa.yaml'
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

En la columna ***TARGETS*** aparece el valor ***<unknown>*** porque el HPA aun no se ha puesto en funcionamiento y no tiene estadísticas que mostrar.

Ahora vamos a ver el autoescalado en acción. Creamos una nueva terminal y ejecutamos el comando.
```
kubectl get pods -w
```

La idea es que vayamos viendo en tiempo real cómo se van creando los pods. En la terminal original vamos a instalar un programa llamado ***hey*** cuya finalidad es meter carga en el servidor web, pero para ello, antes instalamos ***go***.
```
sudo apt-get update
```
```
sudo apt install -y golang-go
```
```
export GOPATH=~/go
export PATH=$GOPATH/bin:$PATH
go get -u github.com/rakyll/hey
```

Copiamos la EXTERNAL-IP del servicio de frontend.
```
kubectl get service frontend
```

La salida será parecida a esto:
```
NAME       TYPE           CLUSTER-IP       EXTERNAL-IP                                                               PORT(S)        AGE
frontend   LoadBalancer   10.100.157.230   aa554514004624ce9b896ff5781c9995-2125204655.eu-west-1.elb.amazonaws.com   80:30689/TCP   10m
```

Metemos carga (20 millones de Request) Usar la IP Externa del servicio de frontend.
```
hey -z 20m http://<Poner aquí la IP Externa del Servicio>
```

Esperar unos minutos (tener paciencia) y ver en la segunda terminal como se escalan los pods. 

Hacer CTRL+C en ***hey*** y esperar un tiempo para ver como se quitan los pods. Tarda un buen rato (al menos 10 minutos) Lo podemos ver con detalle describiendo el HPA.

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

