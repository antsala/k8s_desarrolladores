# Laboratorio 03-B: ***Construir imágenes desde Dockerfile***
 
En este laboratorio repasaremos cómo crear imágenes desde Dockerfile, integrando la compilación del  código fuente en tiempo de creación de la imagen.

Requisitos:

1. Una máquina virtual con Ubuntu 20.04 LTS a la que poder hacer ssh o escritorio remoto.
2. Tener instalado ***Docker-ce*** como runtime de contenedo en dicha máquina virtual.


## Ejercicio 1: ***Creación de imagen desde Dockerfile***

Cambiamos al directorio del laboratorio:
```
cd ~/k8s_desarrolladores/03/web_estatica
```

En este directorio tenemos un archivo de ejemplo llamado ***Dockerfile***, que editamos para estudiarlo
```
code Dockerfile
```

Este archivo Dockerfile es muy sencillo y realiza lo siguiente:

* *Línea 2*: Utiliza como imagen base ***Ubuntu 16.04***.
* *Línea 3*: Actualiza repositorio de paquetes e instala ***nginx***.
* *Línea 4*: Crea el archivo ***/var/www/html/index.html*** y le añade el contenido *Hola,               estoy dentro de tu contenedor*. De ahí el nombre de web estática.
* *Línea 5*:  Expone el puerto ***80*** en el contenedor.

Procedemos a crear una imagen basada en este Dockerfile, así que cerramos el editor y, en la consola escribimos el siguiente comando:
```
sudo docker image build -t antsala/web_estatica .
```

Cuando finalice la compilación de la imagen, listamos imágenes.
```
sudo docker image ls
```

## Ejercicio 2: ***Publicación de puertos en el host***

Para poder acceder al contenedor, es necesario publicar un puerto en el host. Tenemos dos métodos:

1. Dejar que Docker asigne un puerto aleatorio (***32768-61000***)
2. Elegir nosotros el puerto externo que más nos guste.

El siguiente comando abrirá un puerto ***aleatorio*** en el host y lo conectará con el puerto ***80*** del contenedor.
```
sudo docker run -d -p 80 --name web_estatica antsala/web_estatica nginx -g "daemon off;"
```

Comprobamos que el contenedor está corriendo:
```
sudo docker container ls
```

La salida del comando anterior tendrá la siguiente forma:
```
CONTAINER ID   IMAGE                  COMMAND                  CREATED          STATUS          PORTS                                     NAMES
aeac57545f04   antsala/web_estatica   "nginx -g 'daemon of…"   32 seconds ago   Up 30 seconds   0.0.0.0:49153->80/tcp, :::49153->80/tcp   web_estatica
```
Observar la columna ***PORTS***. Aparecen los nateos para IPv4 e IPv6. En ambos casos, el puerto externo, en este ejemplo es el ***49153*** (Nota: pegar el comando y poner el  puerto correcto)
```
puerto=poner_aquí_el_puerto_externo
```
```
curl localhost:$puerto
```

Tras la ejecución de ***curl***, veremos la response con el mensaje *Hola, estoy dentro de tu contenedor*.

Para poder elegir el puerto externo, usamos la sintaxis ***-p puerto_externo:puerto_contenedor***. Por ejemplo, si queremos publicar la web estática en el puerto ***8080*** del host, el contenedor debe crearse con el siguiente comando:
```
sudo docker run -d -p 8080:80 --name web_estatica_8080 antsala/web_estatica nginx -g "daemon off;"
```

Al igual que antes procedemos a listar los contenedores y probar la conexión con ***curl***
```
sudo docker container ls
```
```
curl localhost:8080
```

Debe funcionar correctamente.


## Ejercicio 3: ***ENTRYPOINT en Dockerfile***

En el Dockerfile suele ponerse el comando ***ENTRYPOINT*** o el comando ***CMD***, que sirven para indicar el programa o aplicación que se ejecutará al iniciar el contenedor. De esta forma, no será necesario indicar el ejecutable en la línea del ***docker run***.

Procedamos a abrir el archivo ***Dockerfile_entry*** en el directorio de trabajo y le echamos un vistazo.
```
code Dockerfile_entry
```

Cerramos el editor y procedemos a compilar la imagen a partir de este Dockerfile:
```
sudo docker image build -t antsala/web_estatica_entry . -f Dockerfile_entry
```

Comprobamos las imágenes:
```
sudo docker image ls
```

La salida del comando debe ser similar a esta:
```
REPOSITORY                   TAG       IMAGE ID       CREATED          SIZE
antsala/web_estatica_entry   latest    c0887b707b1c   58 seconds ago   222MB
antsala/web_estatica         latest    b1813072ba8d   23 minutes ago   222MB
antsala/apache2              latest    9f6e1d4e26ca   5 hours ago      220MB
ubuntu                       latest    2b4cba85892a   9 days ago       72.8MB
nginx                        latest    c919045c4c2b   12 days ago      142MB
ubuntu                       16.04     b6f507652425   6 months ago     135MB
```

Procedemos a crear un nuevo contenedor publicando el puerto ***80*** del host. Nótese que ya no es necesario indicar el ejecutable en la línea de ***docker run***.
```
sudo docker run -d -p 80:80 --name web_estatica_80_entry antsala/web_estatica_entry
```

Procedemos a listar los contenedores y probar la conexión con ***curl***
```
sudo docker container ls
```
```
curl localhost:80
```

Subimos la imagen ***antsala/web_estatica*** a DockerHub. Nos autenticamos en DockerHub.
```
sudo docker login
```

Subimos la imagen.
```
sudo docker image push antsala/web_estatica
```
```
sudo docker logout
```

Eliminamos los contenedores:
```
sudo docker container rm -f `sudo docker container ls -a -q`
```

## Ejercicio 4: ***ENTRYPOINT y CMD en Dockerfile***

En muchos Dockerfiles, aparecen conjuntamente los comandos ***ENTRYPOINT*** y ***CMD***. Cuando se da esta coincidencia se debe interpretar de la siguiente forma: 

1. ***ENTRYPOINT*** se utiliza para indicar cual es la aplicación o programa que queremos ejecutar en   el contenededor. 
2. ***CMD*** establece los parámetros por defecto que se usarán al llamar a ese ejecutable. 

Docker empleará los parámetros por defecto si en la línea de docker run ***no se especifican parámetros***.

Veámoslo con un ejemplo. En el directorio de trabajo tenemos el archivo ***Dockerfile_entry_cmd***. Lo editamos:
```
code Dockerfile_entry_cmd
```

Observemos los cambios:

* *Línea 6*: Ahora ***ENTRYPOINT*** solo establece el ejecutable. No indica ningún parámetro.
* *Línea 7*: ***CMD*** establece ***-h*** como parámetro por defecto, que mostrará la ayuda de nginx.

Cerramos el editor y procedemos a compilar una nueva imagen.
```
sudo docker image build -t antsala/web_estatica_entry_cmd . -f Dockerfile_entry_cmd
```

Comprobamos las imágenes:
```
sudo docker image ls
```

La salida debe ser similar a esta:
```
REPOSITORY                       TAG       IMAGE ID       CREATED          SIZE
antsala/web_estatica_entry_cmd   latest    958b8110c445   37 seconds ago   222MB
antsala/web_estatica_entry       latest    c0887b707b1c   25 minutes ago   222MB
antsala/web_estatica             latest    b1813072ba8d   47 minutes ago   222MB
antsala/apache2                  latest    9f6e1d4e26ca   5 hours ago      220MB
ubuntu                           latest    2b4cba85892a   9 days ago       72.8MB
nginx                            latest    c919045c4c2b   12 days ago      142MB
ubuntu                           16.04     b6f507652425   6 months ago     135MB
```

Creamos un contenedor con la nueva imagen. Observar que no ponemos ningún parámetro al final de la línea de ***docker run***.
```
sudo docker run -d -p 80:80 --name web_estatica_80_entry_cmd antsala/web_estatica_entry_cmd
```

Lo que ha ocurrido es que Docker ha levantado el contenedor con la instrucción ***/usr/sbin/nginx -h***, cuyo resultado es mostrar la ayuda del comando nginx por la salida estándar. Lo comprobamos con el siguiente comando:
```
sudo docker container logs web_estatica_80_entry_cmd
```

La salida mostrará lo siguiente:
```
nginx version: nginx/1.10.3 (Ubuntu)
Usage: nginx [-?hvVtTq] [-s signal] [-c filename] [-p prefix] [-g directives]

Options:
  -?,-h         : this help
  -v            : show version and exit
  -V            : show version and configure options then exit
  -t            : test configuration and exit
  -T            : test configuration, dump it and exit
  -q            : suppress non-error messages during configuration testing
  -s signal     : send signal to a master process: stop, quit, reopen, reload
  -p prefix     : set prefix path (default: /usr/share/nginx/)
  -c filename   : set configuration file (default: /etc/nginx/nginx.conf)
  -g directives : set global directives out of configuration file
```

Como podemos ver, la utilidad de ***CMD*** es la de proporcionar unos parámetros por defecto que se usarán si la línea de ***docker run*** no pone ninguno.

De la misma forma, si ponemos parámetros en ***docker run***, éstos tendrán prioridad y serán usados en lugar de los parámetros establecidos en ***CMD***.

Vamos a lanzar otro contenedor que inicie nginx. Observar los parámetros que se ponen al final de la línea de ***docker run***.
```
sudo docker run -d -p 90:80 --name web_estatica_90_entry_cmd antsala/web_estatica_entry_cmd "-g daemon off;"
```

Esta vez se ha iniciado nginx en lugar de mostrar la ayuda. Lo comprobamos así:
```
curl localhost:90
```

El propósito de usar ***CMD como parámetro por defecto*** es conseguir lanzar los contenedores con una configuración deseada en la mayoría de los casos, mientras que también podemos cambiar la configuración de la aplicación pasando nuevos parámetros sin necesidad de recompilar la imagen de contenedor.

Más adelante aprenderemos otras formas de conseguir pasar configuraciones a la aplicación sin tener que recurrir a esta técnica. Usaremos ***variables de entorno***, ***archivos de variables de entorno***, ***ConfigMaps*** y ***Secrets*** de Kubernetes.

Existen muchas otras instrucciones para el Dockerfile: ***WORKDIR***, ***ENV***, ***USER***, ***VOLUME***, ***ADD***, ***COPY***, ***HEALTHCHECK***, ***ONBUILD***, etc, que iremos explicando conforme vayan siendo necesitadas.

Eliminamos los contenedores:
```
sudo docker container rm -f `sudo docker container ls -a -q`
```















