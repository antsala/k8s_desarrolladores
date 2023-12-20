# Laboratorio 03-D: ***Frontend-Backend***
 
En este laboratorio crearemos una aplicación de dos capas. Un Frontend con un servidor web que guardará datos en un Backend de Redis.

En la parte final del laboratorio profundizaremos en la Redis de Docker.

Requisitos:

1. Una máquina virtual con Ubuntu 20.04 LTS a la que poder hacer ssh o escritorio remoto.
2. Tener instalado ***Docker-ce*** como runtime de contenedor en dicha máquina virtual.

## Ejercicio 1: ***Creación del Frontend***

Para mayor aclaración, ver ***diapositiva 18*** (Frontend-Backend 1)

Creamos una applicación de Frontend con un servidor web. Para este ejemplo usaremos ***Sinatra***, que es un framework de aplicación web escrito en Ruby. Contiene una biblioteca de aplicación y un DSL sencillo. 

Un DSL es un ***Domain Specific Language*** que permite crear apps de forma declarativa.

Cambiamos al directorio de trabajo. En él hay otro directorio llamado ***sinatra***, donde se encuentra el Dockerfile.

```
cd ~/k8s_desarrolladores/03/frontend-backend/sinatra
```
```
ls -l
```
```
nano Dockerfile
```

El Dockerfile hace lo siguiente:

* *Línea 1*: Usamos la imagen base de ***Ubuntu 16.04***.
* *Línea 2*: Descargamos e instalamos los paquetes: ***ruby***, ***ruby-dev***, ***build-essential*** y las herramientas de desarrollo de redis (redis-tools)
* *Linea 3*: Una vez instalado ruby, procedemos a instalar el framework ***sinatra*** y  el ***cliente de redis***.
* *Línea 4*: Creamos el directorio ***/opt/webapp*** en el sistema de archivos de la imagen de contenedor. Posteriormente usaremos un volumen para publicar la aplicación.
* *Línea 5*: Sinatra levanta un servidor web que escucha en el puerto ***4567***.
* *Línea 6*: ***/opt/webapp/bin/webapp*** es el lanzador de la aplicación.

Compilamos el Dockerfile
```
sudo docker image build -t antsala/sinatra .
```

Comprobamos que se ha creado la imagen.
```
sudo docker image ls
```

Debe salir algo parecido a esto:
```
REPOSITORY                       TAG       IMAGE ID       CREATED             SIZE
antsala/sinatra                  latest    215d44b4b057   2 seconds ago       392MB
antsala/nginx                    latest    47e37f7e59ed   About an hour ago   222MB
antsala/web_estatica_entry_cmd   latest    958b8110c445   2 hours ago         222MB
antsala/web_estatica_entry       latest    c0887b707b1c   3 hours ago         222MB
antsala/web_estatica             latest    b1813072ba8d   3 hours ago         222MB
antsala/apache2                  latest    9f6e1d4e26ca   7 hours ago         220MB
ubuntu                           latest    2b4cba85892a   10 days ago         72.8MB
nginx                            latest    c919045c4c2b   12 days ago         142MB
ubuntu                           16.04     b6f507652425   6 months ago        135MB
```

Subimos la imagen al repositorio.
```
sudo docker login
```
```
sudo docker push antsala/sinatra
```
```
sudo docker logout
```


En la carpeta ***sinatra*** tenemos la carpeta de la aplicación, llamada ***webapp***. Entramos en ella.
```
cd webapp
```

Si listamos el contenido, vemos otras dos carpetas, cada una con un archivo:
```
ls -l
```

* ***./bin/webapp*** es un script que tiene el iniciador de la aplicación.
* ***./lib/app.rb*** es la aplicación en sí.

Procedemos a estudiar los archivos.
```
nano ./bin/webapp
```

El código que puede observarse es el lanzador de la aplicación. Realmente ejecuta el archivo ***app.rb*** que está en la carpeta ***lib***. Cerramos sin guardar.

Vamos a darle permiso de ejecución al script, de lo contrario no funcionaría.
```
chmod +x ./bin/webapp
```

Procedemos a consultar ***app.rb***.
```
nano ./lib/app.rb
```

***app.rb*** hace lo siguiente:

* *Líneas 1-3*: Carga las dependencias.
* *Línea 7*: El servidor web se enlaza con todas las IPs. 
* *Líneas 9-11*: Si llega una request de tipo ***GET*** a la URI ***/***, entonces la response que se devuelve muestra el mensaje ***Aplicación de prueba Sinatra***.
* *Líneas 13-15*: Si llega una request de tipo ***POST*** a la URI ***/json/***, entonces la response devuelve los parámetros que se han recibido por la request, convertidos a JSON.

Para mayor aclaración, ver ***diapositiva 19*** (Frontend-Backend 2)

Por el momento este Frontend no hace otra cosa. Procedemos a crear un contenedor para probarlo.
```
sudo docker run -d -p 4567 --name webapp -v $PWD:/opt/webapp antsala/sinatra
```

Verificamos que el contenedor se ha iniciado.
```
sudo docker container ls
```

La salida será algo parecido a esto:
```
CONTAINER ID   IMAGE             COMMAND                  CREATED          STATUS          PORTS                                         NAMES
4a7f7fb60528   antsala/sinatra   "/opt/webapp/bin/web…"   27 seconds ago   Up 26 seconds   0.0.0.0:49155->4567/tcp, :::49155->4567/tcp   webapp
```
```
puerto_externo=<Poner aquí el puerto externo>
```

Comprobamos los logs.
```
sudo docker container logs webapp
```

La salida debe indicar que el contenedor ha levantado el servidor web, algo así:
```
[2022-03-13 20:02:43] INFO  WEBrick 1.3.1
[2022-03-13 20:02:43] INFO  ruby 2.3.1 (2016-04-26) [x86_64-linux-gnu]
== Sinatra (v2.2.0) has taken the stage on 4567 for development with backup from WEBrick
[2022-03-13 20:02:43] INFO  WEBrick::HTTPServer#start: pid=1 port=4567
```

La primera prueba será mandar una request (GET) a la URI ***/***
```
curl localhost:$puerto_externo
```

La respuesta mostrará el mensaje ***<h1>Aplicación de prueba Sinatra</h1>***

Ahora vamos a construir una request que mande parámetros. Si en la response vemos esos parámetros convertidos a JSON, entonces el Frontend estará funcionando bien.
```
curl -i -H 'Accept: application/json' -d 'nombre=Antonio&apellidos=Salazar Gravan&telefono=666123321' http://localhost:$puerto_externo/json/
```

La respuesta debe ser así: ***{"nombre":"Antonio","apellidos":"Salazar Gravan","telefono":"666123321"}***

Como funciona bien, borramos el contenedor.
```
sudo docker container rm -f webapp
```

## Ejercicio 2: ***Creación del Backend***

Procedemos a desplegar el Backend, que será un contenedor Redis que almacenará los parámetros recibidos por el Frontend. Para mayor aclaración, ver ***diapositiva 20*** (Frontend-Backend 3)

La carpeta ***redis*** contiene los archivos del Backend. En primer lugar consultamos el Dockerfile que construirá el contenedor.
```
cd ~/k8s_desarrolladores/03/frontend-backend/sinatra/redis
```
```
nano ./Dockerfile
```

La imagen se construirá de la siguiente forma:

* *Línea 1*: Imagen base Ubuntu 16.04
* *Línea 2*: Actualizamos paquetes del repositorio e instalamos el servidor redis y las herramientas.
* *Línea 3*: Se abre el puerto ***6379***.
* *Línea 4*: Se inicia el servidor Redis.

Compilamos la imagen del Backend.
```
sudo docker image build -t antsala/redis .
```
Listamos las imágenes.
```
sudo docker image ls
```

La salida será parecida a esta:
```
REPOSITORY                       TAG       IMAGE ID       CREATED          SIZE
antsala/redis                    latest    2acc22b39230   25 seconds ago   168MB
antsala/sinatra                  latest    215d44b4b057   13 hours ago     392MB
antsala/nginx                    latest    47e37f7e59ed   14 hours ago     222MB
antsala/web_estatica_entry_cmd   latest    958b8110c445   15 hours ago     222MB
antsala/web_estatica_entry       latest    c0887b707b1c   16 hours ago     222MB
antsala/web_estatica             latest    b1813072ba8d   16 hours ago     222MB
antsala/apache2                  latest    9f6e1d4e26ca   20 hours ago     220MB
ubuntu                           latest    2b4cba85892a   10 days ago      72.8MB
nginx                            latest    c919045c4c2b   12 days ago      142MB
ubuntu                           16.04     b6f507652425   6 months ago     135MB
```

Subimos la imagen al repositorio.
```
sudo docker login
```
```
sudo docker push antsala/redis
```
```
sudo docker logout
```

Procedemos a levantar un contenedor de Redis para probar.
```
sudo docker run -d -p 6379 --name db antsala/redis
```

Mostramos los contenedores en ejecución:
```
sudo docker container ls
```

La salida es la siguiente:
```
CONTAINER ID   IMAGE           COMMAND                  CREATED          STATUS          PORTS                                         NAMES
3346b5dabb18   antsala/redis   "/usr/bin/redis-serv…"   41 seconds ago   Up 40 seconds   0.0.0.0:49157->6379/tcp, :::49157->6379/tcp   db
```

## Ejercicio 3: ***Creación de una red***

El servidor Redis está en ejecución, pero aparece un problema de seguridad muy importante. Observemos la columna ***PORTS***. La regla de nateo ***49157->6379*** quiere decir que desde el localhost se puede acceder a la base de datos Redis. La ***diapositiva 21*** (Frontend-Backend 4) lo visualiza: Un usuario podría alcanzar a Redis gracias a esta regla de nateo, cosa que NUNCA debe poder hacerse. El acceso al Backend debe realizarse siempre desde el Frontend.

Para poder controlar este comportamiento, debemos crear siempre redes nuevas. Esta red conectará Frontend y Backend, pero no crearemos reglas de nateo para el backend. La ***diapositiva 22*** (Frontend-Backend 5) lo visualiza. Por lo tanto, el contenedor ***db*** que está corriendo debe ser eliminado para crear uno nuevo conectado a la nueva red.

# Eliminamos el contenedor
```
sudo docker container rm -f db
```

Creamos una red nueva, a la que asignaremos el nombre ***app***
```
sudo docker network create app
```

Listamos las redes para ver que se ha creado correctamente.
```
sudo docker network ls
```

La salida es la siguiente:
```
NETWORK ID     NAME      DRIVER    SCOPE
8a88acbd8802   app       bridge    local
e1619c41bce7   bridge    bridge    local
aeace24a7889   host      host      local
7f9cc20838df   none      null      local
```

## Ejercicio 4: ***Recreación del Backend conectado a la nueva red***

Procedemos a desplegar de nuevo el contenedor de Redis, pero esta vez conectado a la red ***app***. Para ello hacemos uso del parámetro ***--net app*** al crear el contenedor.
```
sudo docker run -d --name db --net app antsala/redis
```

Listamos los contenedores en ejecución:
```
sudo docker container ls
```

La salida es similiar a la siguiente. (Nota: Observar cómo en este caso no aparece la regla de nateo)
```
CONTAINER ID   IMAGE           COMMAND                  CREATED         STATUS         PORTS      NAMES
33ecdec64238   antsala/redis   "/usr/bin/redis-serv…"   4 seconds ago   Up 4 seconds   6379/tcp   db
```

## Ejercicio 5: ***Despliegue de la versión de Frontend que conecta con Backend***

Cuando creamos redes nuevas conseguimos tres efectos beneficiosos:

1. Garantizamos el ***aislamiento*** entre las aplicaciones. Puesto que cada aplicación debería tener sus propias redes con sus contenedores conectados a las mismas y, no como no hay enrutamiento entre las redes, se garantiza que los contenedores de una app no pueden interactuar a nivel de IP con los contenedores de otra.
2. Creando una nueva red podemos controlar cuándo se crea la regla de nateo y, en consecuencia, no podrán ser accedidos desde el exterior.
3. Cuando se crea una red, se pueden resolver los contenedores por su nombre. Esto es muy beneficioso ya que permite no tener que usar direcciones IPs en las conexiones (Las IPs de los contenedores podrían cambiar entre reinicios)

Comprobamos los contenedores conectado a la red ***app***
```
sudo docker network inspect app
```

A continuación se muestra la salida del comando anterior, donde se puede apreciar cómo el contenedor ***db*** está conectado a dicha red.
```
"Containers": {
            "33ecdec64238ee7afe4fc7b25432a148fcf2d9e428f0353974b8a0a5f766b8c4": {
                "Name": "db",
                "EndpointID": "07d0e3c507498e6cea3cd8420294ee6c213fede4645170bd28e67c026b3e07c6",
                "MacAddress": "02:42:ac:12:00:02",
                "IPv4Address": "172.18.0.2/16",
                "IPv6Address": ""
           }
       },
```

Procedemos a deslplegar la nueva versión del Frontend, esta vez conectando con el contenedor ***db*** (Redis). Cambiamos al directorio donde se encuentra la nueva versión.
```
cd ~/k8s_desarrolladores/03/frontend-backend/sinatra/webapp_redis
```

Vuelven a aparecer las dos carpetas de la versión anterior (***bin*** y ***lib***). En la carpeta ***bin*** tenemos el script ***webapp*** que lanza la app. Éste no ha cambiado, pero debemos darle permiso de ejecución.
```
chmod +x ./bin/webapp
```

En la carpeta ***lib**** está el archivo ***app.rb*** con la nueva versión del Frontend. Procedemos a editarlo para consultarlo.
```
nano ./lib/app.rb
```

El contenido del archivo es similar a la versión anterior, pero contiene cambios significativos:

* *Línea 8*: Se instancia el objeto ***redis*** que gestionará la comunicación con la base de datos. Observemos cómo se resuelve el host por el nombre del contenedor (***db***). Esto es muy importante.
* *Líneas 12-14*: Si hay una request de tipo GET a la URI */*, se devuelve el mensaje ***Aplicación Sinatra de ejemplo conectada a Redis***.
* *Líneas 21-24*: Si hay una request de tipo POST que envíe parámetros a la URI ***/json/***, se almacenan en la base de datos  por medio de la llamada ***redis.set***. A continuación se devuelve en la response los parámetros convertidos a JSON.
* *Líneas 16-19*: Si hay una request de tipo GET a la URI ***/json***, se hace una lectura de la base de datos para recuperar los parámetros almacenados. A continuación se devuelve en la response los parámetros convertidos a JSON.

Otro aspecto notable a tener en cuenta es que no debemos volver a compilar la imagen del contenedor de Frontend, ya que los archivos de la aplicación se los pasaremos por medio del volumen. Así que solo queda volver a levantar el contenedor, desde el directorio ***webapp_redis*** y conectarlo a la red ***app***. En este caso sí le creamos la regla de nateo para que podamos conectar con el servidor web desde el exterior. La ***diapositiva 23*** (Frontend-Backend 6) resume todos estos conceptos.
```
sudo docker run -d -p 8080:4567 --name webapp_redis --net app -v $PWD:/opt/webapp antsala/sinatra
```

Comprobamos los contenedores en ejecución. 
```
sudo docker container ls
```

La salida debe ser parecida a esta:
```
CONTAINER ID   IMAGE             COMMAND                  CREATED          STATUS          PORTS                                       NAMES
1a18ac8b594d   antsala/sinatra   "/opt/webapp/bin/web…"   38 seconds ago   Up 36 seconds   0.0.0.0:8080->4567/tcp, :::8080->4567/tcp   webapp_redis
33ecdec64238   antsala/redis     "/usr/bin/redis-serv…"   39 minutes ago   Up 39 minutes   6379/tcp                                    db
```

Ahora probamos. En primer lugar una GET al directorio raíz, debe devolver un mensaje.
```
curl localhost:8080
```

Enviamos una request con POST a ***/json/*** con los campos de un formulario. Esta vez se guardará en redis, además de devolverse en formato JSON.
```
curl -i -H 'Accept: application/json' -d 'nombre=Antonio&apellidos=Salazar Gravan&telefono=666123321' http://localhost:8080/json/
```

Por último, una GET a ***/json*** que provocará una lectura de Redis para leer los parámetros del formulario, que serán devueltos en la response en formato JSON.
```
curl http://localhost:8080/json
```

Borramos los contenedores:
```
sudo docker container rm -f `sudo docker container ls -a -q`
```
