# Laboratorio 55-C: ***MongoDB con StatefulSet y Sidecar***
 
En este laboratorio aprenderemos a usar el objeto ***StatefulSet*** junto con contenedores ayudantes o ***sidecars***.

Vamos desplegar topología replicada de ***MongoDB***, usando un contenedor ***Sidecar***.

Requisitos:

1. Una máquina virtual con Ubuntu 20.04 LTS a la que poder hacer ssh o escritorio remoto.
2. Cluster ***Minikube*** iniciado.


Como ya hemos aprendido en este curso, hemos aportado la capacidad de autoescalado al ***ReplicaSet***, pero en este caso usaremos otra técnica que hace uso de un ***Sidecar***.

El ***Sidecar*** monitorizará los pods y, si se levanta una nueva réplica de MongoDB, el Sidecar la añadirá al ReplicaSet. Así mismo, si retiramos pods, el Sidecar lo tendrá en cuenta.

Otra cuestión a aprender es conseguir que nuestros despliegues en Kubernetes sean seguros, así que tendremos en cuenta las buenas prácticas a usar con MongoDB. Protegeremos las comunicaciones entre los nodos MongoDB por medio de una clave, crearemos un usuario y habilitaremos la autenticación.

## Ejercicio 1: ***Crear secretos y script de inicio***

La documentación de ***MongoDB*** nos dice que la imagen tiene la posibilidad de configurar scripts y ejecutarlos la primera vez que se inicia el contenedor. Utilizaremos esta técnica para configurar usuarios (y passwords)

Cambiamos al directorio de trabajo:
```
cd ~/k8s_desarrolladores/55
```

Una técnica muy interesante el crear un script y almacenarlo en un ***configmap***. El script leerá información confidencial que se almacena en un ***secreto***.

Estudiemos el archivo del secreto.
```
code lab-55-C-mongo-secret.yaml
```

Las líneas más destacables son:

* *Línea 2*: Indica que estamos creando un ***secreto***.
* *Línea 5*: Con nombre ***mongo-secret***.
* *Líneas 7 y 8*: Codificación en ***Base64*** de los passwords para los usuarios ***root*** y ***mi_database_user***. El password de ***root*** es ***Pa55w.rdRoot***. El password de ***mi_database_user*** es ***Pa55w.rdUser***.

Salimos sin modificar nada. Aplicaremos más tarde.

A continuación declaramos un ***ConfigMap*** que almacenará el script. Este script leerá el nombre del usuario administrador y su password desde variables de entorno, cargadas a su vez, desde el ***secreto***. Posteriormente crea una base de datos y un usuario.

Echamos un vistazo al siguiente archivo:
```
code lab-55-C-mongo-init-configmap.yaml
```

Las líneas más importantes son:

* *Línea 2*: Es un ***ConfigMap***.
* *Línea 4*: Su nombre es ***mongo-init-script***.
* *Línea 6*: El script se llama ***mongo-user-sh***.
* *Líneas 7-12*: El script crea un usuario administrador cuyo nombre y password se toma respectivamente de las variables de entorno ***MONGO_INITDB_ROOT_USERNAME*** y ***MONGO_INITDB_ROOT_PASSWORD***, que serán inicializadas desde el secreto. Posteriormente crea (si no existe) una base de datos con nombre ***mi_database***. Luego da de alta un usuario llamado ***mi_database_user*** con el password definido en la variable de entorno ***SECOND_USER_DB_PASSWORD***. A este usuario se le da el permiso ***RW*** en la base de datos que se acaba de crear.

Cerramos sin modificar nada. Aplicamos más tarde.

## Ejercicio 2: ***Proteger las comunicaciones de MongoDB***

***MongoDB*** puede (y debe) cifrar las comunicaciones entre sus instancias. Para ello simplemente debemos proporcionar una clave a modo de vector de inicialización para el algoritmo de cifrado. Esta clave la almacenamos en un ***ConfigMap***.

Editamos el archivo:
```
code lab-55-C-mongo-key.yaml
```

Las líneas más importantes son:

* *Línea 2*: Es un ***ConfigMap***.
* *Línea 4*: El nombre del ConfigMap es ***mongo-key***.
* *Línea 5 y 6*:  Nombre de la clave y su valor.

Salimos sin modificar nada. Aplicamos más tarde.

## Ejercicio 3: ***Cambiar permisos al script***

Cuando Mongo se ejecute necesitaremos proporcionar permisos para la ***mongo-key***. Debido a que los valores almacenador por los ***ConfigMaps*** se cargarán en el sistema de archivos del contenedor Mongo, debemos tener cuidado con los permisos de archivo.

El contenedor Mongo se inicia con el usuario ***mongodb*** y no con ***root***, así que tenemos que cambiar los permisos y luego iniciar Mongo. Para conseguirlo cargamos temporalmente el ConfigMap ***mongo-key*** en una carpeta temporal, que luego copiaremos en otra ruta mejor. 

La razón de hacer esto es que el contenedor carga el ***ConfigMap*** como un enlace simbólico en el sistema de archivos y por ello no nos permitirá cambiar ni el propietario ni los permisos del archivo. Para conseguir esto nos apoyamos en otro script.

Abrimos el archivo:
```
code lab-55-C-mongo-script-permissions.yaml
```

Las líneas más importantes son:

* *Línea 2*: Es un ***ConfigMap***.
* *Línea 4*: El nombre es ***mongo-script-permissions***.
* *Línea 6*: ***Nombre*** del script.
* *Línea 7*: Se copia el directorio temporal a la ruta ***/var/lib/mongoKey***, se cambia propietario y se asigna permiso de lectura al archivo de la clave.

Cerramos sin cambiar nada. Aplicamos más tarde.


## Ejercicio 4: ***Creación de un espacio de nombres***

Vamos a crear un ***NameSpace*** para tener organizados los objetos de la aplicación.

Editamos el archivo:
```
code lab-55-C-mongo-namespace.yaml
```

El archivo es muy sencillo y permite ver como se crea el espacio de nombres ***mongodb-repl-system***

Salimos sin cambiar nada. Aplicamos más tarde.

## Ejercicio 5: ***Creación del servicio y cuenta de servicio***

Procedemos a crear un servicio interno de tipo ***Headless***. Estudiamos el siguiente código:
```
code lab-55-C-mongo-service.yaml
```

Las líneas más importantes son:

* *Línea 2*: Es un servicio.
* *Línea 11*: ***ClusterIP: none*** indica que es un servicio ***Headless*** (Sin dirección IP).

Salimos sin modificar. Aplicamos más tarde

Procedemos a crear una ***Service Account*** y un ***Cluster Role Binding*** ya que los sidecars necesitan permisos para observar al pod. Solo es necesario aportar los permisos ***watch*** y ***list***, pero en este ejemplo asignamos más para que podamos entender cómo se asignan.

Editamos el siguiente archivo:
```
code lab-55-C-mongo-service-account-rbac.yaml
```

El contenido más interesante es el siguiente:

* *Líneas 1-5*: Define una ***Service Account***.
* *Línea 4*: El nombre de la cuenta es ***mongo-account***.
* *Líneas 7-23*: Se crea un ***Rol*** en el cluster.
* *Línea 10*: El rol se llama ***mongo-role***.
* *Líneas 12-14*: Se asignan todos los permisos para los ***ConfigMaps***.
* *Líneas 15-17*: Se asignan los permisos ***list*** y ***watch*** en el contexto de los ***Deployments***.
* *Líneas 18-20*: Se asignan todos los permisos en el contexto de ***Services***.
* *Líneas 21-23*: Se asignan los permisos ***get***, ***list*** y ***watch*** en los ***Pods***.
* *Líneas 25-36*: Se asigna el rol y la cuenta de servicio al cluster.
* *Línea 28*: El ***ClusterRoleBinding*** se llama ***mongo_role_binding***.
* *Líneas 30-32*: Se asigna la cuenta de servicio al cluster en el espacio de nombres ***mongodb-repl-system***.
* *Líneas 33-36*: Se signa el rol al cluster en el mismo espacio de nombres.

Salimos sin modificar nada. Aplicamos más tarde.

## Ejercicio 6: ***Creación del StatefulSet de MongoDB***

Procedemos a crear el ***StatefulSet*** con dos contenedores. El primero es el de MongoDB, mientras que el segundo será el ***Sidecar***.

Editamos el archivo:
```
code lab-55-C-mongo-statefulset.yaml
```

Las líneas más importantes:

* *Línea 2*: Se declara un ***StatefulSet***.
* *Línea 6*: ***podManagementPolicy: Parallel*** indica al controlador del StatefulSet que inicie o elimine todos los pods a la vez (en paralelo). En consecuencia, no espera a que un pod se inicie o se elimine para procesar el siguiente. La otra opción es ***OrderedReady***, en la que no se procesa un nuevo pod hasta que el anterior está en ***Ready*** o ***Terminated***. 
* *Línea 7*: En principio levantamos un solo pod, el maestro de MongoDB.
* *Línea 11*: Conectamos el StatefulSet al servicio ***mongo***.
* *Línea 17*: Usamos la cuenta de servicio llamada ***mongo-account***.
* *Líneas 20-49*: Declaramos el ***contenedor principal*** del pod. Pueden observarse todas las propiedades.
* *Líneas 51-69*: Declaramos un contenedor helper o ***sidecar***. 
* *Líneas 70-82*: Declaración de los volúmenes que montarán los contenedores.
* *Líneas 85-93*: PVCs 

Salimos sin guardar. Aplicamos más tarde.

## Ejercicio 7: ***Creación del objeto Kustomization***

En Kubernetes existe una clase de objeto muy útil llamada ***Kustomization***. Su principal utilidad es agrupar los archivos YAML que forman una aplicación o proyecto y crear un espacio de nombres para aislarlo del resto.

Editamos el archivo y lo consultamos:
```
code kustomization.yaml
```

Las líneas más interesantes son:

* *Línea 2*: Indicamos que es un objeto ***Kustomization***.
* *Líneas 5-12*: Aquí indicamos todos los archivos YAML que forman el proyecto.
* *Línea 5*: Creamos un espacio de nombres para el proyecto.
* *Línea 13*: Establece el espacio de nombres ***mongodb-repl-system*** para la creación de los objetos que pueden ser creados en el ámbito de un espacio de nombres.

A continuación debemos ejecutar el siguiente comando para generar un único archivo YAML que contendrá todos los objetos.
```
kubectl kustomize . > mongo-backend.yaml
```

Editamos el archivo:
```
code mongo-backend.yaml
```

Observar como aparecen todos los objetos en el archivo, y como ciertas líneas (10, 56, 70, 79, 92, 101, 110 y 123 declaran el espacio de nombres en el que deseamos que se creen los objetos.

Ahora solo queda aplicar:
```
kubectl apply -f mongo-backend.yaml 
```

Comprobamos que todos los objetos se ha iniciado correctamente:
```
kubectl get all --namespace mongodb-repl-system
```

Comprobar que el maestro de MongoDB (mongo-0) se ha iniciado:
```
kubectl describe pod mongo-0 --namespace mongodb-repl-system
```

Podemos comprobar el escalado del StatefulSet de la siguiente forma:
```
kubectl scale --replicas=3 statefulset mongo --namespace mongodb-repl-system
```

Comprobamos:
```
kubectl get all --namespace mongodb-repl-system
```

Limpiamos:
```
kubectl delete -f mongo-backend.yaml 
```

Listamos los PVCs y los PVs para comprobar que se han eliminado.
```
kubectl get pvc 
```
