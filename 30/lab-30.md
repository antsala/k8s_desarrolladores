# Laboratorio 30: ***Espacios de Nombres***
 
En este laboratorio aprenderemos a usar los espacios de nombres (***Namespaces***)

Requisitos:

1. Una máquina virtual con Ubuntu 20.04 LTS a la que poder hacer ssh o escritorio remoto.
2. Tener instalado el ***runtime de podman***. (ver lab-06-A.md, Ejercicio 1 y 2)

## Ejercicio 1:  ***Creación de un espacio de nombres***

Entramos en el directorio del laboratorio
```
cd ~/k8s_desarrolladores/30
```

Listamos los espacios de nombre que existen el cualquier implementación del cluster de Kubernetes.
```
kubectl get namespaces
```

La salida será similar a esta:
```
NAME              STATUS   AGE
default           Active   26h
kube-node-lease   Active   26h
kube-public       Active   26h
kube-system       Active   26h
```

Vamos a volver a desplegar la aplicación ***helloContainer***, pero en este caso usaremos un espacio de nombres. Creamos el espacio de nombres ***hello-container-namespace***.
```
kubectl create namespace hello-container-namespace
```

Comprobamos que se ha creado.
```
kubectl get namespace hello-container-namespace
```

La salida será parecida a esta:
```
NAME                        STATUS   AGE
hello-container-namespace   Active   52s
```

Podemos ver las características del espacio de nombre creado, describiéndolo:
```
kubectl describe namespace hello-container-namespace
```

La salida es:
```
Name:         hello-container-namespace
Labels:       kubernetes.io/metadata.name=hello-container-namespace
Annotations:  <none>
Status:       Active

No resource quota.
 
No LimitRange resource.
```

Podemos confirmar que el espacio de nombres se crea sin cuotas de recursos. Si queremos aprender cómo funcionan las cuotas en K8s, procedemos a leer este artículo: (https://kubernetes.io/docs/concepts/policy/resource-quotas/)

Para los límites, este artículo lo explica detalladamente: (https://kubernetes.io/docs/concepts/policy/limit-range/)


El siguiente paso es filtrar los objetos para el espacio de nombres. Para ello hacemos uso del parámetro ***--namespace***.
```
kubectl get all --namespace hello-container-namespace
```

No debe existir ningún objeto.
```
No resources found in hello-container-namespace namespace.
```

## Ejercicio 2:  ***Aplicar un archivo YAML en un espacio de nombres***


Si el archivo YAML no contiene indicación alguna sobre el espacio de nombres en el que se deben crear los objetos, podemos especificarlo al aplicar el YAML. Los archivos ***lab-30-helloContainer-deployment.yaml*** y ***lab-30-helloContainer-service.yaml*** no tienen ninguna propiedad que indique el espacio de nombres en el que deseamos crear los objetos. Podemos comprobarlo abriendo los archivos y buscando la propiedad ***namespace***, que no estará.
```
code lab-30-helloContainer-deployment.yaml
```
```
code lab-30-helloContainer-service.yaml
```

Salimos del editor sin modificar los archivos.

Si procedieramos a realizar un ***apply***, los objetos se crearían en el espacio de nombres por defecto, que se llama ***default***. En este ejercicio, pretendemos que se depliegue en el espacio de nombres que hemos creado (***hello-container-namespace***)
```
kubectl apply -f lab-30-helloContainer-deployment.yaml --namespace hello-container-namespace
```

Comprobamos que se ha desplegado el deployment en el espacio de nombres 'hello-container-namespace'
```
kubectl get deployments --namespace hello-container-namespace
```

La salida debe ser parecida a esta y debe mostrar el deployment en nuestro espacio de nombres:
```
NAME                         READY   UP-TO-DATE   AVAILABLE   AGE
hello-container-deployment   3/3     3            3           74s
```

Procedemos a aplicar el archivo del servicio de forma análoga:
```
kubectl apply -f lab-30-helloContainer-service.yaml --namespace hello-container-namespace
```

Listamos los servicios en el espacio de nombres.
```
kubectl get services --namespace hello-container-namespace
```

El resultado debe ser el siguiente:
```
NAME                      TYPE           CLUSTER-IP     EXTERNAL-IP   PORT(S)        AGE
hello-container-service   LoadBalancer   10.97.21.188   <pending>     80:30000/TCP   40s
```

La ***EXTERNAL-IP*** está pendiente. Necesitamos ejecutar ***minikube tunnel***. Abrimos otra terminal y ejecutamos:
```
minikube tunnel
```

Listamos los servicios en el espacio de nombres.
```
kubectl get services --namespace hello-container-namespace
```

El resultado debe ser el siguiente:
```
NAME                      TYPE           CLUSTER-IP     EXTERNAL-IP    PORT(S)        AGE
hello-container-service   LoadBalancer   10.111.92.10   10.111.92.10   80:30883/TCP   93s
```
```
external_ip=<Copiar aquí la EXTERNAL-IP asignada>
```

Probamos que funciona conectacto a la EXTERNAL-IP.
```
curl $external_ip:80
```

Procedemos a borrar el deployment y el servicio.
```
kubectl delete -f lab-30-helloContainer-deployment.yaml --namespace hello-container-namespace
kubectl delete -f lab-30-helloContainer-service.yaml --namespace hello-container-namespace
```

Comprobamos que el espacio de nombres está vacío.
```
kubectl get all --namespace hello-container-namespace
```

## Ejercicio 2:  ***Predeterminar el espacio de nombres***

Podemos indicar a Kubernetes que deseamos trabajar siempre con el mismo espacio de nombres y, de esta forma no tener que especificarlo en la línea de comandos.
```
kubectl config set-context --current --namespace hello-container-namespace
```

Mostramos todos los objetos:
```
kubectl get all
```

Observar como se refiere al espacio de nombres ***hello-container-namespace***.


Volvemos a aplicar los despliegues, esta vez sin indicar el espacio de nombres:
```
kubectl apply -f lab-30-helloContainer-deployment.yaml
kubectl apply -f lab-30-helloContainer-service.yaml
```

Listamos los servicios en el espacio de nombres.
```
kubectl get services --namespace hello-container-namespace
```

El resultado debe ser el siguiente:
```
NAME                      TYPE           CLUSTER-IP       EXTERNAL-IP      PORT(S)        AGE
hello-container-service   LoadBalancer   10.106.177.149   10.106.177.149   80:30963/TCP   11m
```
```
external_ip=<Copiar aquí la EXTERNAL-IP asignada>
```

Probamos que funciona conectando a la ***EXTERNAL-IP***.
```
curl $external_ip:80
```

Procedemos a borrar el deployment y el servicio.
```
kubectl delete -f lab-30-helloContainer-deployment.yaml 
kubectl delete -f lab-30-helloContainer-service.yaml 
```

Comprobamos que el espacio de nombres está vacío.
```
kubectl get all 
```

Volvemos a predeterminar ***default*** como espacio de nombres.
```
kubectl config set-context --current --namespace default
```

Comprobamos
```
kubectl get all 
```

## Ejercicio 4:  ***Predeterminar el espacio de nombres en el archivo YAML***

Sin duda esta es la forma más cómoda, porque en el propio archivo YAML indicamos el espacio de nombres en el que nos interesa que se creen los objetos. El archivo ***lab-30-helloContainer-deployment-namespace.yaml*** ha sido modificado para establecer el espacio de nombres para el deployment.

Editamos el archivo:
```
code lab-30-helloContainer-deployment-namespace.yaml
```

* *Línea 5*: Se indica el espacio de nombres en el que deseamos crear el deployment.

Lo mismo ocurre para ***lab-30-helloContainer-service-namespace.yaml***. Lo Editamos
```
code lab-30-helloContainer-service-namespace.yaml
```

* *Línea 5*: Se indica el espacio de nombres en el que deseamos crear el servicio.

Solo queda aplicar los archivos:
```
kubectl apply -f lab-30-helloContainer-deployment-namespace.yaml
kubectl apply -f lab-30-helloContainer-service-namespace.yaml
```

Comprobamos
```
kubectl get all --namespace hello-container-namespace
```

Listamos los servicios en el espacio de nombres.
```
kubectl get services --namespace hello-container-namespace
```

El resultado debe ser el siguiente:
```
NAME                      TYPE           CLUSTER-IP       EXTERNAL-IP      PORT(S)        AGE
hello-container-service   LoadBalancer   10.103.39.1      10.103.39.1      80:30579/TCP   89s
```
```
external_ip=<Copiar aquí la EXTERNAL-IP asignada>
```

Probamos que funciona conectandoo a la ***EXTERNAL-IP***.
```
curl $external_ip:80
```

Procedemos a borrar el deployment y el servicio.
```
kubectl delete -f lab-30-helloContainer-deployment.yaml --namespace hello-container-namespace
kubectl delete -f lab-30-helloContainer-service.yaml --namespace hello-container-namespace
```

Comprobamos que el espacio de nombres está vacío.
```
kubectl get all --namespace hello-container-namespace
```

La salida debe indicar:
```
No resources found in hello-container-namespace namespace.
```

Limpiamos los recursos:

Cerrar la terminal de ***minikube tunnel***.

Eliminamos el espacio de nombres creado:
```
kubectl delete namespace hello-container-namespace
```

Comprobamos:
```
kubectl get namespaces
```

Solo deben quedar los espacios de nombres originales
```
NAME              STATUS   AGE
default           Active   56m
kube-node-lease   Active   56m
kube-public       Active   56m
kube-system       Active   56m
```

Predeterminamos el espacio de nombres ***default***
```
kubectl config set-context --current --namespace default
```

Comprobamos los objetos presentes en el cluster.
```
kubectl get all 
```

Solo debe quedar el servicio de Kubernetes.
```
NAME                 TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
service/kubernetes   ClusterIP   10.96.0.1    <none>        443/TCP   53m
```

