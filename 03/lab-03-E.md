# Laboratorio 03-E: ***Microservicios***
 
En este laboratorio crearemos una aplicación el lenguaje ***Go***. Esta app la compilaremos en el mismo momento de generar la imagen de contenedor. La app levanta un servidor web que, si nos conectamos, nos responde indicando el nombre del host y las direcciones IPs que tiene asignada. Posteriormente crearemos un servicio en ***Swarm*** para demostrar en balanceo.

Requisitos:

1. Una máquina virtual con Ubuntu 20.04 LTS a la que poder hacer ssh o escritorio remoto.
2. Tener instalado ***Docker-ce*** como runtime de contenedor en dicha máquina virtual.

## Ejercicio 1: ***Instalación de Go*** 

Nos ponemos en ***HOME***.
```
cd ~
```

Descargamos el compilador de Go, lo descoprimimos.
```
curl -O https://dl.google.com/go/go1.15.5.linux-amd64.tar.gz
tar xvf go1.15.5.linux-amd64.tar.gz
```

Cambiamos propietario.
```
sudo chown -R root:root ./go
```

Lo cambiamos de carpeta.
```
sudo mv go /usr/local
```

Editamos ***.profile*** para agregar variables en entorno al final del archivo.
```
nano .profile
```

Agregar estas dos líneas al final del archivo.
```
export GOPATH=$HOME/k8s_desarrolladores/03/work
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
```

Guardar y salir. Cargamos el profile para que se actualicen las varibles de entorno.
```
source .profile
```

## Ejercicio 2: ***Compilar una app en Go***

Editamos el archivo ***hello.go*** que usaremos para probar el compilador.
```
nano ~/k8s_desarrolladores/03/work/src/hello/hello.go
```

Como se puede comprobar es un código sencillo que muestra el mensaje ***Hola Mundo!!!***

Procedemos a compilarlo, indicando el nombre de la carpeta. Se procederá a compilar todos los archivos ***.go*** que existan (en este caso solo está ***hello.go***)
```
go install hello
```

Comprobamos si ha compilado escribiendo el nombre del ejecutable (***hello***)
```
hello
```

Para determinar la carpeta en la que Go coloca el ejecutable:
```
which hello
```

Acabamos de comprobar que el compilador de Go funciona. Procedemos a crear una app que haga algo más apropiado. Se encuentra en el archivo ***~/k8s_desarrolladores/03/work/src/hellocontainer/helloContainer.go***. La editamos.
```
nano ~/k8s_desarrolladores/03/work/src/helloContainer/helloContainer.go
```

El programa hace lo siguiente:

* *Líneas 1-9*: Declara módulos a usar.
* *Línea 11*: Función ***main***. Es el punto de entrada a la aplicación.
* *Líneas 12-15*: Inicializa la variable ***port*** al valor ***8080***. Si se pasa el valor del puerto como variable de entorno al iniciar la app, el código puede leer dicha variable y actualizar el puerto.
* *Línea 17*: Instancia un objeto servidor web.
* *Línea 18*: Asocia función de respuesta a evento. Si hay una request a la URI ***/***, el código llamará a la función ***helloContainer*** para procesar dicha request.
* *Línea 20*: Se muestra mensaje por la salida estándar.
* *Línea 21 y 22*: Se inicia el servidor web. Cualquier error es guardado en el objeto ***err*** y almacenado en el log.
* *Línea 25*: Inicio de la función de respuesta a evento.
* *Líneas 26-30*: Se disponen diversos mensajes, entre ellos el nombre del host, que serán mostrados en la response.
* *Líneas 32-39*: Se leen todas las IPs asignadas al servidor que se mostrarán en la response.

Salimos sin modificar y procedemos a compilar el programa
```
go install helloContainer
```

Lanzamos la aplicación para probarla
```
helloContainer
```

Abrimos otra consola, porque en la anterior se está ejecutando el servidor, y hacemos una request al puerto ***8080***.
```
curl localhost:8080
```

La salida será similar a la siguiente, donde se muestra el nombre del host y las direcciones IPs asignadas.
```
Hola Mundo!
Version: 1.0.0
Hostname: ubu
Dirección IP: 192.168.1.45
Dirección IP: 172.17.0.1
Dirección IP: 172.18.0.1
```

Cerramos la última consola y, paramos el servidor con ***CTRL+C***.


## Ejercicio 3: ***Contenerizar la app Go***

Vamos a contenerizar la app. Para ello nos colocamos en el directorio de contexto con el siguiente comando.
```
cd ~/k8s_desarrolladores/03/helloContainerCtx
```

Copiamos la carpeta ***helloContainer*** (que contiene ***helloContainer.go***) a este directorio de contexto.
```
cp -r ~/k8s_desarrolladores/03/work/src/helloContainer/ .
```

En el directorio de contexto tenemos un archivo Dockerfile que pasamos a detallar.
```
nano ~/k8s_desarrolladores/03/helloContainerCtx/Dockerfile
```

Este Dockerfile va a hacer dos cosas:

1.Compilará el código fuente (***helloContainer.go***) y generará el ejecutable.
2.Tomará el ejecutable generado y lo almacenará en la imagen de contenedor.

* *Línea 1*: Crea un contenedor basado en la imagen ***golang:1.11-alpine***, que es una imagen que tiene instalada el compilador de Go.
* *Línea 2*: Añade los archivos del directorio de contexto, en este caso la carpeta ***helloContainer*** (Que a su vez contiene el archivo ***helloContainer.go***) en la carpeta ***/go/src*** de la imagen de contenedor.
* *Línea 3*: Compila nuestra aplicación en Go. El ejecutable se almacena en la ruta ***/go/bin/helloContainer***.

En este punto, la imagen de contenedor, a la que podemos llamar ***PREVIA***, tiene el ejecutable de nuestra aplicación.

El archivo Dockerfile contínua, desde otra imagen base.

* *Línea 5*: Se crea un contenedor con la imagen ***alpine***.
* *Línea 6*: Esto es lo importante: Se establece como contexto la imagen ***PREVIA***, copia de ella el ejecutable de la aplicación (***/go/bin/hellocontainer***) y lo pone en el directorio actual (***.***) de la ***NUEVA IMAGEN***.
* *Línea 7*: Por último. se ejecuta el programa ***helloContainer***, que está en el directorio actual (***.***)


Con esta técnica se puede compilar el código fuente y crear la imagen de contenedor de la aplicación. Es ideal para usar en los pipelines de CI/CD.

Creamos la imagen (que a su vez compilará el código fuente)
```
sudo docker build -t antsala/hello_container .
```

Comprobamos que la imagen se ha creado correctamente.
```
sudo docker image ls
```

La salida será similar a esta. Nota: Se han omitido el resto de imágenes previas)
```
REPOSITORY                       TAG           IMAGE ID       CREATED          SIZE
antsala/hello_container          latest        b66a10e42c57   52 seconds ago   12.2MB
```

Subimos la imagen a DockerHub.
```
sudo docker login
```
```
sudo docker push antsala/hello_container
sudo docker logout
```

Lanzamos un contenedor para probar nuestra aplicación contenerizada.
```
sudo docker container run --name helloContainer -d -p 85:8080 antsala/hello_container
```

Probamos
```
curl localhost:85
```

La salida será similar a esto.(Nota: Observar como el nombre del host es el id del contenedor)
```
Hola Mundo!
Version: 1.0.0
Hostname: 93a224a3858d
Dirección IP: 172.17.0.2
```

Una vez comprobado, eliminamos el contenedor.
```
sudo docker container rm -f helloContainer
```


## Ejercicio 4: ***Desplegar servicio en Swarm***


En este ejercicio vamos a desplegar un servicio basado en la imagen de contenedor que hemos creado. La práctica va a servirnos para introducir los ***archivos de manifiesto*** que definen el servicio (y que serán la clave de Kubernetes). También podremos probar el balanceo, característico de los micro servicios.

En primer lugar debemos desplegar el cluster de Docker (Swarm). Trabajaremos con un único nodo, pero esto no es importante.
```
sudo docker swarm init
```

Como resultado, un mensaje nos informará que el cluster está levantado, así como el procedimiento a seguir para añadir más nodos al cluster.

Para poder crear archivos YAML de manifiesto, necesitamos que esté instalada la herramienta ***Docker-compose***. La instalamos con el siguiente comando:
```
sudo apt -y install docker-compose
```

Ahora procedemos a estudiar un archivo de manifiesto de ejemplo. Este se encuentra en la carpeta *** '~/k8s_desarrolladores/03/helloContainerSvc***, así que entramos en ella.
```
cd ~/k8s_desarrolladores/03/helloContainerSvc
```

En esta carpeta tenemos el archivo ***Docker-compose.yaml***. Lo editamos para estudiarlo:
```
nano ./Docker-compose.yaml
```

* *Línea 18 y 19*: Crea una red, llamada ***webnet*** para uso exclusivo de los contenedores que van a crearse.
* *Línea 2*: Define los servicios que se implementarán. En este caso uno solo, el servicio ***web***, que se declara desde la línea 4 hasta la 17.
* *Línea 5*: Indicamos la imagen que usarán los contenedores.
* *Línea 7*: Deseamos 5 réplicas (5 contenedores)
* *Líneas 10-11*: Limitamos cada contenedor a usar el 10% de la CPU total y 50 MB de RAM.
* *Líneas 12-13*: Los contenedores se reiniciarán si la aplicación falla.
* *Línea 15*: Regla de nateo para acceder al servicio. Conectaremos contra la IP de cualquier nodo del cluster, al puerto ***4000***. El balanceador irá repartiendo el tráfico entre los 5 contenedores, que responden en el puerto ***8080***.
* *Línea 17*: Conectamos los contenedores a la red ***webnet***.

Salimos sin modificar y procedemos a desplegar este servicio en el cluster. Lo llamaremos ***helloContainerStack***.
```
sudo docker stack deploy -c Docker-compose.yaml helloContainerStack
```

Swarm responderá que ha creado la red y el servicio. Para comprobarlo ejecutamos el siguiente comando:
```
sudo docker stack ls
```

La salida es como esta:
```
NAME                  SERVICES   ORCHESTRATOR
helloContainerStack   1          Swarm
```

Para ver los contenedores que ha levantado, usamos este comando:
```
sudo docker stack ps helloContainerStack
```

La salida será similar a la siguiente, en la que podemos ver los 5 contenedores corriendo.
```
ID             NAME                        IMAGE                            NODE      DESIRED STATE   CURRENT STATE           ERROR     PORTS
a69nbga2tm61   helloContainerStack_web.1   antsala/hello_container:latest   ubu       Running         Running 2 minutes ago             
mkze4dpg85rz   helloContainerStack_web.2   antsala/hello_container:latest   ubu       Running         Running 2 minutes ago             
jfx4eweg2kv1   helloContainerStack_web.3   antsala/hello_container:latest   ubu       Running         Running 2 minutes ago             
wqo1cvz212z1   helloContainerStack_web.4   antsala/hello_container:latest   ubu       Running         Running 2 minutes ago             
ypqzhk0cp9x1   helloContainerStack_web.5   antsala/hello_container:latest   ubu       Running         Running 2 minutes ago 
```

Lo que realmente nos interesa demostrar es el funcionamiento del balanceo. Lo conseguimos repitiendo el siguiente comando, con el que verificamos que van contestando los diferentes contenedores que forman en servicio.
```
curl localhost:4000
```

Eliminamos la aplicación. 
```
sudo docker stack rm helloContainerStack
```


Borramos los contenedores (por si quedara alguno perdido)
```
sudo docker container rm -f `sudo docker container ls -a -q`
```

