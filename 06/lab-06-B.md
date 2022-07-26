# Laboratorio 06-B: ***Frontend-Backend con POD***
 
En este laboratorio crearemos una aplicación de dos capas. Para ello implementaremos un pod con dos contenedores. El contenedor de ***Frontend*** tendrá una imagen de ***phpMyAdmin***, mientras que el contenedor de ***Backend*** usará una imagen de ***mySQL***. La información de usuario y contraseña se inyectará a los contenedores en tiempo de ejecución, mediante un archivo de variables de entorno.

Requisitos:

1. Una máquina virtual con Ubuntu 20.04 LTS a la que poder hacer ssh o escritorio remoto.
2. Tener instalado el runtime de contenedor ***podman***, o lo que es lo mismo, haber realizado el ejercico anterior ***lab-06-A.md***

## Ejercicio 1: ***Descargar imágenes de contenedor***

Nos ubicamos en el directorio de trabajo
```
cd ~/k8s_desarrolladores/06
```

Este laboratorio desplegará un pod con dos contenedores. El Frontend con ***phpMyAdmin*** y el Backend con ***mySQL***. En primer lugar comprobamos que no existen contenedores en ejecución con podman.
```
podman container ls -a
```

La salida mostrará que no hay contenedores.
```
CONTAINER ID  IMAGE       COMMAND     CREATED     STATUS      PORTS       NAMES
```

Procedemos a descargar las imagenes mencionadas. (Nota: elegir DockerHub como registro)
```
podman image pull phpmyadmin
podman image pull mysql:8
```

Comprobamos que las imágenes se han descargado:
```
podman image ls
```

La salida será similar a esto:
```
REPOSITORY                    TAG         IMAGE ID      CREATED     SIZE
docker.io/library/phpmyadmin  latest      e53ae368a8e5  3 days ago  520 MB
docker.io/library/mysql       8           826efd84393b  5 days ago  526 MB
```

## Ejercicio 2: ***Archivo con variables de entorno***

El Frontend de phpMySQL necesitará credenciales para conectar con mySQL. Al mismo tiempo, a mySQL debemos proporcionarle las ***credenciales de administrador***. Todo ello se debe hacer mediante ***variables de entorno*** y está documentado en el repositorio de las imágenes.

Para ***mySQL***, la documentación está en: (https://hub.docker.com/_/mysql). En el apartado ***Environment Variables***, podemos leer... (Nota: ponemos fragmentos de la documentación)
```
...
When you start the mysql image, you can adjust the configuration of the MySQL instance by passing one or more environment variables on the docker run command line. 
...
MYSQL_ROOT_PASSWORD
This variable is mandatory and specifies the password that will be set for the MySQL root superuser account...

MYSQL_DATABASE
This variable is optional and allows you to specify the name of a database to be created on image startup...

MYSQL_USER, MYSQL_PASSWORD
These variables are optional, used in conjunction to create a new user and to set that user's password. This user will be granted superuser permissions (see above) for the database specified by the MYSQL_DATABASE variable. 
Both variables are required for a user to be created.
```

Para ***phpMyAdmin***, la documentación está en (https://hub.docker.com/_/phpmyadmin). En el apartado ***Environment variables summary*** podemos leer... (Nota: ponemos fragmentos de la documentación)
```
...
Variables that can be read from a file using _FILE

PMA_HOST (Esta variable será necesaria para permitir la comunicación con el contenedor dentro del pod) ...
MYSQL_ROOT_PASSWORD
...
MYSQL_PASSWORD
...
```

Así que para facilitar la manipulación de estos secretos, creamos un archivo que contenga estas variables, con la excepción de ***PMA_HOST*** que la inyectaremos manualmente. Editemos el archivo con estas variables con el siguiente comando:
```
code lab-06-B-env-vars
```

## Ejercicio 3: ***Creación del pod***

Procedemos a crear el pod, al que llamaremos ***mi_inventario***:
```
podman pod create --name mi_inventario
```

Al crear un pod también se crea el contenedor de infraestructura, que administra los contenedores que vamos a incluir posteriormente en el pod. La ***diapositiva 9*** (Frontend-Backend con POD 1) aclara el escenario.

Listamos los pods:
```
podman pod ls
```

La salida será similar a esta:
```
POD ID        NAME           STATUS      CREATED             INFRA ID      # OF CONTAINERS
b9323ddee5de  mi_inventario  Created     About a minute ago  43004cb9a3a3  1
```

Deseamos que ***phpMyAdmin*** sea alcanzable desde el exterior a través del puerto ***8085***, por lo que lo más lógico sería pensar que al iniciar el contenedor, debemos crear una regla de nateo con el parámetro ***-p 8085:80***, pero esto no se hace así.

Al usar pods, las reglas de nateo se deben ***aplicar en el pod***, no en los contenedores. Esto es así porque en el modo ***rootless***  (no usamos ***sudo***) los contenedores ***no adquieren dirección IP***. Cuando el pod reciba el tráfico lo reenviará a los contenedores. La ***diapositiva 10*** (Frontend-Backend con POD 2) aclara el escenario.

Además, estas reglas se deben poner en el momento de crear el pod, y no a posteriori. La forma más efectiva será destruir el pod y recrearlo con la regla de nateo.

Por último, los contenedores dentro del pod se pueden comunicar sin problemas, a pesar de no tener direcciones IPs, no les hace falta.
```
podman pod rm mi_inventario
```

Procedemos a crear el pod, pero esta vez incluímos la regla de nateo.
```
podman pod create --name mi_inventario -p 8085:80
```

Listamos los pods:
```
podman pod ls
```
La salida será similar a esta: (Nota: Se puede ver que existe el contenedor de infraestructura porque el número de contenedores es ***1***)
```
POD ID        NAME           STATUS      CREATED         INFRA ID      # OF CONTAINERS
eea96c70c998  mi_inventario  Created     5 seconds ago  c30a65ef86ae  1
```

Creamos el contenedor del Backend. Se tienen que cumplir dos consideraciones:

1. Debe pertenecer al pod ***mi_inventario***.
2. Debe incluir el ***archivo de variables de entorno***.
```
podman container run -d --pod mi_inventario --name mi_inventario_db --env-file ./lab-06-B-env-vars mysql:8
```

Ya solo nos queda levantar el contenedor del Frontend, que debe cumplir los siguientes requisitos:

1. Debe pertenecer al pod ***mi_inventario***.
2. Debe incluir el ***archivo de variables de entorno***.
3. Debemos inyectarle la variable de entorno ***PMA_HOST=127.0.0.1*** para que la regla de nateo funcione sobre este contenedor.
```
podman container run -d --pod mi_inventario --name mi_inventario_phpmyadmin --env-file ./lab-06-B-env-vars -e PMA_HOST=127.0.0.1 phpmyadmin
```

Podemos comprobar que el runtime ***no ha asignado*** una dirección IP a ninguno de los contenedores:
```
podman container inspect mi_inventario_phpmyadmin --format '{{.NetworkSettings.IPAddress}}'
podman container inspect mi_inventario_db --format '{{.NetworkSettings.IPAddress}}'
```

Listamos los contenedores en ejecución: (Nota: el parámetro ***--pod*** hace que se muestre información del pod al que pertenece el contenedor)
```
podman container ls --pod
```

La salida será similar a esta:
```
CONTAINER ID  IMAGE                                COMMAND               CREATED         STATUS             PORTS                 NAMES                     POD ID        PODNAME
c30a65ef86ae  k8s.gcr.io/pause:3.5                                       About a minute ago  Up 18 seconds ago  0.0.0.0:8085->80/tcp  eea96c70c998-infra        eea96c70c998  mi_inventario
3b779247b0e7  docker.io/library/mysql:8            mysqld                17 seconds ago      Up 18 seconds ago  0.0.0.0:8085->80/tcp  mi_inventario_db          eea96c70c998  mi_inventario
f6b7cfe407c9  docker.io/library/phpmyadmin:latest  apache2-foregroun...  8 seconds ago       Up 8 seconds ago   0.0.0.0:8085->80/tcp  mi_inventario_phpmyadmin  eea96c70c998  mi_inventario
```

Podemos observar que hay tres reglas de nateo, una para el contenedor de infraestructura, otra para el Frontend y la última para el Backend. Debería resultar extraño como se permite tráfico http a los contenedores de infraestructura y de la base de datos, pero no tiene ninguna implicación puesto que la regla sea válida se debe inyectar la variable de entorno ***PMA_HOST=127.0.0.1*** al contenedor que deseemos que reciba tráfico.

Para probar que todo funciona debemos conectar con un navegador a ***http://localhost:8085***. La autenticación usará el usuario ***root*** con la contraseña ***secreto***.

Cuando phpMyAdmin conecte, mostrará la base de datos ***nominas***, que podremos administrar perfectamente.


## Ejercicio 4: ***Eliminación del pod***

Para eliminar el pod, primero debemos parar los contenedores con el siguiente comando:
```
podman pod stop mi_inventario
```

Comprobamos que los contenedores se han detenido.
```
podman container ls -a --pod
```

La salida debe ser similar a esta:
```
CONTAINER ID  IMAGE                                COMMAND               CREATED         STATUS                    PORTS                 NAMES                     POD ID        PODNAME
c30a65ef86ae  k8s.gcr.io/pause:3.5                                       12 minutes ago  Exited (0) 8 seconds ago  0.0.0.0:8085->80/tcp  eea96c70c998-infra        eea96c70c998  mi_inventario
3b779247b0e7  docker.io/library/mysql:8            mysqld                10 minutes ago  Exited (0) 7 seconds ago  0.0.0.0:8085->80/tcp  mi_inventario_db          eea96c70c998  mi_inventario
6748b4cd893a  docker.io/library/phpmyadmin:latest  apache2-foregroun...  5 minutes ago   Exited (0) 7 seconds ago  0.0.0.0:8085->80/tcp  mi_inventario_phpmyadmin  eea96c70c998  mi_inventario
```

Ahora procedemos a eliminar el pod (y sus contenedores)
```
podman pod rm mi_inventario
```

Borramos las imágenes.
```
podman image rm -f `podman image ls -q`
```

