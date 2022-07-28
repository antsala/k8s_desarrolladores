# Laboratorio 60: ***Despliegue de Prometeus con Operator***.
 
En este laboratorio aprenderemos desplegar una app que hace uso de ***Operator***.

Requisitos:

1. Una máquina virtual con Ubuntu 20.04 LTS a la que poder hacer ssh o escritorio remoto.
2. Cluster ***Minikube*** iniciado.

La monitorización del cluster es crítica. ***Prometheus*** es una solución de monitorización para Kubernetes muy usada. Su arquitectura tiene tres capas:

La capa de ***captura de métrica***, que recopìla información del cluster de Kubernetes. Una capa de ***almacenamiento***, que ingiere y almacena en una base de datos toda esta información. Un ***frontend*** con el que puede interactuar otras app, como Grafana, para consultar los datos.

El servidor de Prometheus también puede emitir alertas a otros servicios.

La cuestión es que para desplegar un servidor Prometheus podemos actuar de dos formas:

1. Crear nosotros mismos los archivos de configuración YAML para agregar los respectivos objetos al cluster (STS, CM, Secret, Deployment, ...) y aplicarlos en el orden correcto.
2. Usar un ***Operator*** que se encargará de desplegar y mantener todos los componentes de Prometheus.

Parece claro que la primera solución es muy ineficiente, compleja y requiere mantenimiento manual. Para proceder con la segunda, debemos localizar un ***operador*** para Prometheus que suele encontrarse en la forma de un ***chart*** de Helm.

## Ejercicio 1: ***Instalación del Operator de Prometheus***

Cambiamos al directorio de trabajo:
```
cd ~/k8s_desarrolladores/60
```

En este link: https://prometheus-community.github.io/helm-charts/ tenemos la documentación del chart de Prometheus. Si la leemos atentamente, veremos que se instala el ***Operator***, con las ventajas que este ofrece.

Instalamos Helm (solo si no estuviera presente):
```
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3
chmod 700 get_helm.sh
```
```
./get_helm.sh
source ~/.profile
```

Comprobamos que está instalado:
```
helm version
```

Añadimos los repos de Helm para Prometheus.
```
helm repo add prometheus-community https://charts.helm.sh/stable
helm repo add stable https://kubernetes-charts.storage.googleapis.com/
helm repo update
```

Instalamos el Chart:
```
helm install prometheus prometheus-community/kube-prometheus-stack
```

Una vez que el instalador nos indique que se han instalado los componentes, ejecutamos este comando:
```
kubectl --namespace default get pods -l "release=prometheus"
```

La salida mostrará algo parecido a esto: (Nota: Observar el pod con el nombre ***prometheus-operator***)
```
NAME                                                  READY   STATUS    RESTARTS   AGE
prometheus-kube-prometheus-operator-b7fdd56db-qq66n   1/1     Running   0          97s
prometheus-kube-state-metrics-94f76f559-6t9xd         1/1     Running   0          97s
prometheus-prometheus-node-exporter-9lskx             1/1     Running   0          97s
```

Si deseamos ver todos los objetos relacionados con Prometheus, hacemos:
```
kubectl --namespace default get all -l "release=prometheus"
```

La salida será como esta:(Nota: Observar todos los objetos con ***operator*** en el nombre)
```
NAME                                                      READY   STATUS    RESTARTS   AGE
pod/prometheus-kube-prometheus-operator-b7fdd56db-qq66n   1/1     Running   0          57m
pod/prometheus-kube-state-metrics-94f76f559-6t9xd         1/1     Running   0          57m
pod/prometheus-prometheus-node-exporter-9lskx             1/1     Running   0          57m
 
NAME                                              TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
service/prometheus-kube-prometheus-alertmanager   ClusterIP   10.106.131.89    <none>        9093/TCP   57m
service/prometheus-kube-prometheus-operator       ClusterIP   10.103.184.123   <none>        443/TCP    57m
service/prometheus-kube-prometheus-prometheus     ClusterIP   10.104.159.183   <none>        9090/TCP   57m
service/prometheus-kube-state-metrics             ClusterIP   10.99.57.98      <none>        8080/TCP   57m
service/prometheus-prometheus-node-exporter       ClusterIP   10.107.117.31    <none>        9100/TCP   57m

NAME                                                 DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR   AGE
daemonset.apps/prometheus-prometheus-node-exporter   1         1         1       1            1           <none>          57m

NAME                                                  READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/prometheus-kube-prometheus-operator   1/1     1            1           57m
deployment.apps/prometheus-kube-state-metrics         1/1     1            1           57m

NAME                                                            DESIRED   CURRENT   READY   AGE
replicaset.apps/prometheus-kube-prometheus-operator-b7fdd56db   1         1         1       57m
replicaset.apps/prometheus-kube-state-metrics-94f76f559         1         1         1       57m

NAME                                                                    READY   AGE
statefulset.apps/alertmanager-prometheus-kube-prometheus-alertmanager   1/1     57m
statefulset.apps/prometheus-prometheus-kube-prometheus-prometheus       1/1     57m
```

El ***daemonset*** es un componente de Kubernetes que se ejecutará en cada uno de los nodos workers.

Helm crea un montón de objetos de configuración: (Nota: observar los configmaps)
```
kubectl get configmap
```

Helm crea un montón de Secretos: (Nota: observar los secretos)
```
kubectl get secrets
```

Si deseamos estudiar los objetos individualmente, por ejemplo, uno de los servicios:
```
kubectl get service prometheus-kube-prometheus-operator -o yaml > operator.yaml 
```

Lo editamos y lo estudiamos:
```
code operator.yaml
```

Limpiamos recursos:
```
helm delete prometheus
```

Comprobamos:
```
kubectl get all
```
