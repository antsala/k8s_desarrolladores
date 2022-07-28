# Laboratorio 65-A: ***Endpoints externos***
 
En este laboratorio aprenderemos a conectar pods de Kubernetes con servicios que corren ***fuera del cluster*** de Kubernetes.

Requisitos:

1. Una máquina virtual con Ubuntu 20.04 LTS a la que poder hacer ssh o escritorio remoto.
2. Cluster ***Minikube*** iniciado.

Durante el curso hemos hablado mucho sobre los servicios en Kubernetes. Vimos que pueden ser ***internos***, ***externos***, e incluso no tener IP como los ***Headless***. Todos ellos se caracterizan porque envían tráfico a una serie de pods que son balanceados por el servicio.

También hemos hablado de la dificultad existente en Kubernetes para hacer que una applicación con estado funcione bien en un contexto de alta disponibilidad, y los objetos de Kubernetes implicados en ello (***PV***, ***PVC***, ***StatefulSet***) así como la posibilidad de usar ***Operators***.

En ciertos escenarios, se prefiere ***no implementar*** el servicio en el cluster de Kubernetes. Ejemplos típicos son las bases de datos, que se suelen desplegar en máquinas virtuales. La propia tecnología de la base de datos ofrece la alta disponibilidad.

En este laboratorio aprenderemos a indicar a los pods que el servicio a consultar está fuera del cluster de Kubernetes. Para ello haremos uso de un nuevo objeto llamado ***Endpoint***.

## Ejercicio 1: ***Creación del endpoint***

Lo primero que debemos conocer es la ***IP*** del servicio externo al cluster. En este ejemplo será la aplicación ***hello_container***  y estará en la única máquina que tenemos, pero funcionará igualmente.

Cambiamos al directorio de trabajo:
```
cd ~/k8s_desarrolladores/65
```

Tomamos la IP del nodo de ***Minikube*** y la anotamos:
```
IP=`minikube ip`
```

La visualizamos y la copiamos al portapapeles.
```
echo $IP
```

A continuación editamos el siguiente archivo:
```
code lab-65-A-external-web-endpoint.yaml
```

Las líneas más interesantes son:

* *Línea 2*: Declaramos un objeto ***Endpoints***.
* *Línea 4*: Su nombre es ***external-web***. Esta etiqueta debe coincidir con la del servicio interno que crearemos luego.
* *Línea 9*: IMPORTANTE!!!! Pegar aquí la IP capturada .
* *Línea 12*: ***Puerto*** al que se redirigirá el tráfico al servicio externo.

Aplicamos:
```
kubectl apply -f lab-65-A-external-web-endpoint.yaml
```

Comprobamos:
```
kubectl get endpoints external-web 
```

La salida mostrará algo similar a esto: (Nota: Observar la IP del servicio externo al cluster)
```
NAME           ENDPOINTS         AGE
external-web   192.168.1.38:80   44s
```

Aun es necesario dar otro paso, ya que los pods del cluster necesitan conectar con un ***servicio*** (no lo pueden resolver el nombre del endpoint)

## Ejercicio 2: ***Creación del servicio interno***

Procedemos a crear el servicio interno y asociarlo con el endpoint. 

Abrimos el siguiente archivo:
```
code lab-65-A-external-web-service.yaml 
```

Las líneas más importantes son:

* *Línea 4*. Esta etiqueta ***debe coincidir*** con la declarada en el Endpoint.
* *Líneas 9 y 10: se crea la regla de ***natting***.

Aplicamos:
```
kubectl apply -f lab-65-A-external-web-service.yaml 
```

Observamos el resultado:
```
kubectl describe service external-web
```

Comprobar en la salida como ***Endpoints*** es ***<IP del servicio externo al cluster>:80***, que realmente no es un pod, sino el objeto ***Endpoint*** que hemos creado antes.

## Ejercicio 3: ***Prueba del endpoint***


Ahora necesitamos ***simular*** ser el servicio externo al cluster. Para ello vamos a usar directamente el ***runtime***, es decir, Kubernetes no interviene para nada. Usaremos ***podman*** para levantar un contenedor de la aplicación ***hello_container*** sobre la propia IP del nodo.
```
sudo podman run -d --name mi_servicio_web_externo -p 80:8080 docker.io/antsala/hello_container
```

Ya tenemos nuestro servidor web funcionando. Se ha publicado en el puerto 80. El siguiente paso es levantar un pod en Kubernetes que, usando el nombre del servicio interno ***external-web*** pueda alcanzar a ***hello_container***.

Editamos el siguiente archivo y lo estudiamos:
```
code lab-65-A-test-external-web.yaml
```

Las líneas más interesantes son:

* *Líneas 17-21*: Se levanta un contenedor de ubuntu que realiza un simple bucle infinito para que no se detenga el contenedor. El contenedor se llamará ***test-container***.

Aplicamos:
```
kubectl apply -f lab-65-A-test-external-web.yaml
```

Comprobamos:
```
kubectl get pods
```

Ahora necesitamos entrar en el contenedor del pod. Ajustar el nombre del pod según convenga: (Nota: Esto debe devolver un prompt dentro del contenedor)
```
kubectl exec -i -t test-external-service-96c77b685-n2smg test-container -- /bin/bash
```

Ejecutar los comandos siguientes dentro del contenedor:
```
apt-get update -y 
apt-get install -y curl
```

Ahora probamos si desde dentro del contenedor se puede alcanzar el servicio externo a.l cluster
```
curl external-web
```

Como resultado debemos ver la salida de la aplicación ***hello_container***.

Salimos del contenedor.
```
exit
```

Borramos los recursos.
```
kubectl delete -f lab-65-A-external-web-endpoint.yaml
kubectl delete -f lab-65-A-external-web-service.yaml 
kubectl delete -f lab-65-A-test-external-web.yaml
```
