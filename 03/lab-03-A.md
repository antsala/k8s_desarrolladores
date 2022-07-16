# Laboratorio 03-A: ***Creación de contenedores con Docker***

En este laboratorio repasaremos los conceptos más importantes de Docker y crearemos una imagen de contenedor para su uso posterior.

Requisitos:

1. Una máquina virtual con Ubuntu 20.04 LTS a la que poder hacer ssh o escritorio remoto.


## Ejercicio 1: Instalación de ***Docker*** 

Procedemos a instalar Docker (Community Edition) en la máquina. Este procedimiento es para Ubuntu 20.04.
```
sudo apt -y update
```
```
sudo apt -y install apt-transport-https ca-certificates curl software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu focal stable"
sudo apt update
sudo apt -y install docker-ce 
```

Comprobamos que se ha instalado correctamente.
```
sudo docker --version
sudo systemctl status docker
```

## Ejercicio 2: Primeros contenedores con ***Docker***

Creamos un contenedor y entramos dentro (contenedor interactivo ***-i -t***)
```
sudo docker run --name mi_primer_contenedor -i -t ubuntu /bin/bash
```

Al iniciarse el contenedor, estaremos dentro. Observar el prompt, será parecido a este
***root@141be1940196:/#***. Somos ***root*** en una maquina que tiene por nombre ***141be1940196***, que es un id aleatorio asignado al contenedor.

Dentro del contenedor procedemos a realizar cuantas acciones estimemos oportunas. Para salir del contenedor escribimos:
```
exit
```

Al salir de un contenedor interactivo, se detiene el programa que estaba ejecutando (***/bin/bash***), por lo que el contenedor se detiene. Podemos comprobarlo listando los contenedores:
```
sudo docker container ls -a 
```

La salida del comando anterior será similar a esto:
```
CONTAINER ID   IMAGE     COMMAND       CREATED         STATUS                    PORTS     NAMES
141be1940196   ubuntu    "/bin/bash"   3 minutes ago   Exited (0) 1 second ago             mi_primer_contenedor 
```
Observar como la columna ***STATUS*** indica que el contenedor ha finalizado su ejecución.

Podemos iniciar un contenedor detenido con el comando ***start***, de esta forma:
```
sudo docker container start mi_primer_contenedor
```
y comprobar que está en ejecución listando los contenedores activos:
```
sudo docker container ls
```

La salida del comando será así:
```
CONTAINER ID   IMAGE     COMMAND       CREATED         STATUS          PORTS     NAMES
141be1940196   ubuntu    /bin/bash     6 minutes ago   Up 38 seconds             mi_primer_contenedor
```

Comprobar que en la columna ***STATUS*** se indica que el contenedor está corriendo (***UP***)

Si necesitamos realizar cambios en un contenedor, podemos "conectarnos" a él con el comando ***attach***
```
sudo docker container attach mi_primer_contenedor
```

Volvemos a escribir ***exit*** para salir.
```
exit
```

Lo que nos lleva a la etapa anterior donde el contenedor se ha detenido. Podemos comprobarlo listando los contenedores y ver como la columna ***STATUS*** indica ***Exited***
```
sudo docker container ls -a
```

Para que un contenedor siga en ejecución de forma desatendida (***dettached***), es necesario indicar el parámetro ***-d*** y hacer que ejecute un programa que no finalice, como un servidor web, de base de datos, o en este ejemplo, un bucle infinito.
```
sudo docker run --name mi_daemon -d ubuntu /bin/sh -c "while true; do echo Hello World; sleep 1; done"
```

Si listamos los contenedores, tendremos uno corriendo (***mi_daemon***) y otro detenido (***mi_primer_contenedor***)
```
sudo docker container ls -a 
```
La salida del comando mostrará algo similiar a esto:
```
CONTAINER ID   IMAGE     COMMAND                  CREATED              STATUS                     PORTS     NAMES
316c46bb5f47   ubuntu    "/bin/sh -c 'while t…"   About a minute ago   Up About a minute                    mi_daemon
141be1940196   ubuntu    "/bin/bash"              15 minutes ago       Exited (0) 6 minutes ago             mi_primer_contenedor
```

Podemos ver la salida estándar de cualquier contenedor que esté en ejecución con el comando ***logs***, por ejemplo:
```
sudo docker container logs mi_daemon
```

La salida será similar a esta:
```
Hello World
Hello World
Hello World
Hello World
Hello World
antonio@ubu:~$ 
```

Levantemos otro contenedor similar al que está en ejecución, y que llamaremos ***mi_daemon_2*** con el siguiente comando:
```
sudo docker run --name mi_daemon_2 -d ubuntu /bin/sh -c "while true; do echo Hello World; sleep 1; done"
```

Comprobemos que tenemos dos contenedores corriendo:
```
sudo docker container ls -a 
```

La salida será similar a esto:
```
CONTAINER ID   IMAGE     COMMAND                  CREATED          STATUS                      PORTS     NAMES
209e2796365f   ubuntu    "/bin/sh -c 'while t…"   41 seconds ago   Up 41 seconds                         mi_daemon_2
316c46bb5f47   ubuntu    "/bin/sh -c 'while t…"   5 minutes ago    Up 5 minutes                          mi_daemon
141be1940196   ubuntu    "/bin/bash"              20 minutes ago   Exited (0) 11 minutes ago             mi_primer_contenedor
```

Un aspecto muy importante a controlar es el consumo de recursos que realizan los contenedores. El comando ***stats*** está ahí para eso. Ejecutemos el siguiente comando:
```
sudo docker container stats mi_daemon mi_daemon_2
```

La salida muestra el uso de CPU, Memoria, Red, E/S y número de procesos para los contenedores que hemos indicado.
```
CONTAINER ID   NAME          CPU %     MEM USAGE / LIMIT   MEM %     NET I/O       BLOCK I/O    PIDS
316c46bb5f47   mi_daemon     0.12%     684KiB / 3.834GiB   0.02%     3.86kB / 0B   176kB / 0B   2
209e2796365f   mi_daemon_2   0.12%     500KiB / 3.834GiB   0.01%     3.28kB / 0B   0B / 0B      2
```
Salimos con ***CTRL+C***

En otras ocasiones necesitaremos ejecutar un comando en el contenedor, pero no deseamos hacerlo de forma interactiva. Para ello hacemos uso de ***exec***, por ejemplo, ver el contenido de un archivo, ejecutamos:
```
sudo docker container exec mi_daemon cat /etc/lsb-release
```

Que mostrará la siguiente salida:
```
DISTRIB_ID=Ubuntu
DISTRIB_RELEASE=20.04
DISTRIB_CODENAME=focal
DISTRIB_DESCRIPTION="Ubuntu 20.04.4 LTS"
```

Para detener un contenedor en ejecución, usamos el comando ***stop***. Por ejemplo:
```
sudo docker container stop mi_daemon
```

Comprobamos que ***mi_daemon*** se ha detenido.
```
sudo docker container ls -a 
```

Lo verificamos en la salida del comando. El códido de salida es ***127***, que significa que hemos "matado" al contenedor.
```
CONTAINER ID   IMAGE     COMMAND                  CREATED          STATUS                        PORTS     NAMES
209e2796365f   ubuntu    "/bin/sh -c 'while t…"   10 minutes ago   Up 10 minutes                           mi_daemon_2
316c46bb5f47   ubuntu    "/bin/sh -c 'while t…"   15 minutes ago   Exited (137) 25 seconds ago             mi_daemon
141be1940196   ubuntu    "/bin/bash"              30 minutes ago   Exited (0) 21 minutes ago               mi_primer_contenedor
```

La inspección de contenedores es una herramienta fundamental, que nos muestra mucha información sobre el estado de ejecución del contenedor. La salida es un objeto JSON con la mencionada información.
```
sudo docker container inspect mi_daemon_2
```

Generalmente estaremos interesado en conocer cierta información del contenedor, como su estado de ejecución o dirección IP. Hacemos uso de las plantilla del lenguaje ***Go*** para filtrar la salida del JSON anterior. Por ejemplo, para mostrar el estado de ejecución de nuestros contnenedores ponemos:
```
sudo docker container inspect --format='{{.State.Running}}' mi_daemon
```
```
sudo docker container inspect --format='{{.State.Running}}' mi_daemon_2
```

Y para conocer sus direcciones IP los siguiente comandos: 
(***Nota***: Observar como los contenedores detenidos no tienen dirección IP)
```
sudo docker container inspect --format='{{.NetworkSettings.IPAddress}}' mi_daemon
sudo docker container inspect --format='{{.NetworkSettings.IPAddress}}' mi_daemon_2
```

En cuanto a la eliminación de contenedores, en principio no se pueden eliminar contenedores en ejecución. El siguiente comando dará un error:
```
sudo docker container rm mi_daemon_2
```

El mensaje de error es similar a este:
```
Error response from daemon: You cannot remove a running container 209e2796365f182942531bcef47c9ed7c3605d390c8d26bcfa0ab6248a1221a6. 
Stop the container before attempting removal or force remove
```

El propio mensaje indica la forma de proceder: O bien, detenemos el contenedor antes de su eliminación, o usamos el modificador ***-f*** (-***-force***) para ello. Borremos de nuevo el contenedor.
```
sudo docker container rm -f mi_daemon_2
```

Para borrar el resto de contenedores, que están detenidos, procedemos así:
```
sudo docker container rm mi_daemon
```
```
sudo docker container rm mi_primer_contenedor
```

Listamos los contenedores para verificar que no existe ninguno en este momento.
```
sudo docker container ls -a 
```

## Ejercicio 3: Imágenes con ***Docker***

Cuando lanzamos un contenedor, Docker descarga la imagen. Podemos ver las imágenes descargadas con el siguiente comando:
```
sudo docker image ls
```

La salida debería ser similar a la siguiente:
```
REPOSITORY   TAG       IMAGE ID       CREATED      SIZE
ubuntu       latest    2b4cba85892a   9 days ago   72.8MB
```

Para descargar imágenes con antelación, usamos el comando ***pull***:
```
sudo docker image pull nginx
```

Si volvemos a listar las imágenes, veremos la de ***nginx***.
```
sudo docker image ls
```

Esta es la salida:
```
REPOSITORY   TAG       IMAGE ID       CREATED       SIZE
ubuntu       latest    2b4cba85892a   9 days ago    72.8MB
nginx        latest    c919045c4c2b   11 days ago   142MB
```

En los despliegues es muy importante hacer referencia una imagen concreta. Más adelante hablaremos del etiquetado de imagen detalladamente, pero por ahora supongamos que queremos ejecutar un contenedor de forma interactiva, basado en la versión 16.04 de Ubuntu:
```
sudo docker container run -i -t --name mi_ubuntu_16_04 ubuntu:16.04 /bin/bash
```

Comprobamos versión y salimos del contenedor.
```
cat /etc/lsb-release
exit
```

En breve estaremos creando imágenes. El siguiente paso será subirlas a un repositorio como ECR, ACR, Quay, etc. En este caso la subimos a ***DockerHub***. Primero debemos hacer un login.
```
sudo docker login
```

Observar el mensaje que aparecer:
```
WARNING! Your password will be stored unencrypted in /root/.docker/config.json.
Configure a credential helper to remove this warning. See
https://docs.docker.com/engine/reference/commandline/login/#credentials-store
```

Dice que las credenciales se almacenan de forma NO ENCRIPTADA. Concretamente en base64. Por lo que si alguien ajeno tiene acceso al archivo ***/root/.docker/config.json*** puede obtener la credencial en texto plano simplemente usando la herramienta ***base64***.

Por eso es muy recomendable hacer un ***logout*** cuando hayamos terminado de subir las imágenes:
```
sudo docker logout
```

Procedemos a crear una imagen desde un contenedor. Para ello creamos un nuevo contenedor interactivo basado en Ubuntu al que llamaremos ***mi_servidor_web***.
```
sudo docker container run --name mi_servidor_web -i -t ubuntu /bin/bash
```

Dentro del contenedor procedemos a instalar el servidor web.
```
apt-get -yqq update
apt-get -y install apache2
service apache2 start
exit
```

Al ejecutar ***exit*** hemos salido del contenedor y éste se detiene. Lo comprobamos.
```
sudo docker container ls -a 
```

En la salida comprobamos que ***mi_servidor_web*** está detenido.
```
CONTAINER ID   IMAGE          COMMAND       CREATED          STATUS                      PORTS     NAMES
2253067a8f35   ubuntu         "/bin/bash"   2 minutes ago    Exited (0) 7 seconds ago              mi_servidor_web
2b4fc5d72193   ubuntu:16.04   "/bin/bash"   16 minutes ago   Exited (0) 15 minutes ago             mi_ubuntu_16_04
```

Ahora procedemos a crear una imagen a partir del contenedor. Podemos usar el id del contenedor, ***2253067a8f35*** en este ejemplo, o su nombre ***mi_servidor_web***, para seleccionarlo.

En la última parte del comando ponemos el nombre de la imagen. En este ejemplo será ***antsala/apache2***, porque ya la vamos preparando para subirla al registro y ***antsala*** es el nombre de un repositorio de imágenes. ***apache2*** es el nombre de la imagen.
```
sudo docker container commit mi_servidor_web antsala/apache2
```

Listamos las imágenes para ver la nueva:
```
sudo docker image ls
```

Y en la salida podremos verla:
```
REPOSITORY        TAG       IMAGE ID       CREATED          SIZE
antsala/apache2   latest    9f6e1d4e26ca   32 seconds ago   220MB
ubuntu            latest    2b4cba85892a   9 days ago       72.8MB
nginx             latest    c919045c4c2b   11 days ago      142MB
ubuntu            16.04     b6f507652425   6 months ago     135MB
```

Para finalizar este laboratorio, eliminamos TODOS contenedores:
```
sudo docker container rm -f `sudo docker container ls -a -q`
```