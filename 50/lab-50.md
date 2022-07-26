# Laboratorio 50: ***Instalar aplicaciones usando HELM***
 
En este laboratorio aprenderemos a usar Helm.

Requisitos:

1. Una máquina virtual con Ubuntu 20.04 LTS a la que poder hacer ssh o escritorio remoto.
2. Tener instalado ***Minikube***.

## Ejercicio 1: ***Instalar Helm***

Entramos en el directorio del laboratorio
```
cd ~/k8s_desarrolladores/50
```

Helm es el administrador de paquetes de Kubernetes. Permite desplegar, actualizar, y administrar las aplicaciones de Kubernetes. Para ello, se escribe algo llamado ***charts***. Puedes pensar en los charts de Helm como archivos YAML parametrizados.

Vamos a instalar Helm:
```
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3
```
```
chmod 700 get_helm.sh
./get_helm.sh
source ~/.profile
```
Comprobamos que está instalado:
```
helm version
```

Es conveniente familizarse con los comandos. Echamos un vistazo a la ayuda:
```
helm --help
```

## Ejercicio 2: ***Instalar WordPress con Helm***


Vamos a instalar ***WordPress*** por medio de Helm. El repositorio oficial de Helm no es el único registro en el que se pueden encontrar charts. De hecho, los fabricantes más importantes también publican los charts para poder desplegar sus aplicaciones en Kubernetes.

Bitnami ofrece muchas aplicaciones, con sus respectivos charts. Consultemos su web en (https://bitnami.com/stacks/helm)

En este ejercicio instalaremos WordPress: https://bitnami.com/stack/wordpress/helm. En primer lugar añadimos el repositorio que contiene los charts de Helm:
```
helm repo add bitnami https://charts.bitnami.com/bitnami
```

La salida muestra el siguiente mensajes:
```
"bitnami" has been added to your repositories
```

Para listar los repositorios de Helm que tenemos añadidos, usamos el comando:
```
helm repo list
```

Y la salida mostrará lo siguiente:
```
NAME   	URL                               
bitnami	https://charts.bitnami.com/bitnami
```

Ahora instalamos WordPress desde Helm:
```
helm install my-minikube-wp bitnami/wordpress
```

La salida es extensa y es muy importante leerla. Será así: (Nota: Tomemos unos minutos para leerla detenidamente)
```
NAME: my-minikube-wp
LAST DEPLOYED: Sun Mar 20 18:26:31 2022
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
CHART NAME: wordpress
CHART VERSION: 13.1.4
APP VERSION: 5.9.2

** Please be patient while the chart is being deployed **

Your WordPress site can be accessed through the following DNS name from within your cluster:

    my-minikube-wp-wordpress.default.svc.cluster.local (port 80)

To access your WordPress site from outside the cluster follow the steps below:

1. Get the WordPress URL by running these commands:

  NOTE: It may take a few minutes for the LoadBalancer IP to be available.
        Watch the status with: 'kubectl get svc --namespace default -w my-minikube-wp-wordpress'
 
   export SERVICE_IP=$(kubectl get svc --namespace default my-minikube-wp-wordpress --include "{{ range (index .status.loadBalancer.ingress 0) }}{{ . }}{{ end }}")
   echo "WordPress URL: http://$SERVICE_IP/"
   echo "WordPress Admin URL: http://$SERVICE_IP/admin"
 
2. Open a browser and access WordPress using the obtained URL.
 
3. Login with the following credentials below to see your blog:

  echo Username: user
  echo Password: $(kubectl get secret --namespace default my-minikube-wp-wordpress -o jsonpath="{.data.wordpress-password}" | base64 --decode)
```

Repasemos todo lo que ha ocurrido. Estos son los objetos que el chart ha creado en el cluster:
```
kubectl get all
```

En la salida se muestra que se han creado: 2 pods (Wordpress y su base de datos), 2 servicios, uno interno para MariaDB y el otro externo (LoadBalancer). 1 Deployment para el pod del frontend con su replicaset y, por último 1 statefulset (aun por explicar) para mariaDB.
```
pod/my-minikube-wp-mariadb-0                    1/1     Running   0          2m20s
pod/my-minikube-wp-wordpress-8599667678-g2wh8   1/1     Running   0          2m20s
 
NAME                               TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
service/kubernetes                 ClusterIP      10.96.0.1        <none>        443/TCP                      47h
service/my-minikube-wp-mariadb     ClusterIP      10.104.185.233   <none>        3306/TCP                     2m20s
service/my-minikube-wp-wordpress   LoadBalancer   10.101.12.249    <pending>     80:31040/TCP,443:30653/TCP   2m20s

NAME                                       READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/my-minikube-wp-wordpress   1/1     1            1           2m20s

NAME                                                  DESIRED   CURRENT   READY   AGE
replicaset.apps/my-minikube-wp-wordpress-8599667678   1         1         1       2m20s

NAME                                      READY   AGE
statefulset.apps/my-minikube-wp-mariadb   1/1     2m20s
```

Si observamos, el servicio externo ***my-minikube-wp-wordpress*** aun no tiene EXTERNAL-IP. En una terminal diferente ejecutamos:
```
minikube tunnel
```

En la terminal principal mostramos los servicios para comprobar la EXTERNAL-IP asignada.
```
kubectl get service my-minikube-wp-wordpress
```

La salida será similar a la siguiente:
```
NAME                       TYPE           CLUSTER-IP      EXTERNAL-IP     PORT(S)                      AGE
my-minikube-wp-wordpress   LoadBalancer   10.101.12.249   10.101.12.249   80:31040/TCP,443:30653/TCP   9m6s
```

Copiamos la IP externa y, nos conectamos con el navegador. Debe mostrarse la página principal de WordPress.


## Ejercicio 3: ***Desinstalar WordPress con Helm***

Si hemos realizado despliegues en el cluster con Helm, debemos serguir usándolo para las actualizaciones o las desinstalaciones, como en este ejercicio. Primero listamos los stacks que se han desplegado con Helm.
```
helm list
```

La salida es similar a esta:
```
NAME          	NAMESPACE	REVISION	UPDATED                                	STATUS  	CHART           	APP VERSION
my-minikube-wp	default  	1       	2022-03-20 18:26:31.839328282 +0100 CET	deployed	wordpress-13.1.4	5.9.2 
```

Para proceder a la desinstalación ejecutamos este comando:
```
helm uninstall my-minikube-wp
```

Damos unos segundos y comprobamos que se han borrado los objetos del cluster:
```
kubectl get all
```

Después de aprender a usar los volúmenes persistentes (PV) y las claims de volúmenes persistentes (PVC) revisitaremos Helm para ver como se producen las actualizaciones del chart y cómo podemos resolver problemas. Helm tiene una curva de aprendizaje muy plana y en consecuencia es muy asequible. En esta URL tienes la documentación del producto para aprender a crear tus propios charts: https://v2.helm.sh/docs/developing_charts/


