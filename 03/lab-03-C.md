# Laboratorio 03-C: ***Volúmenes***
 
En este laboratorio repasaremos cómo usar volúmenes para conseguir persistencia en el sistema de archivos del contenedor. También pueden ser usados para pasar configuraciones en tiempo de inicio del contenedor.

Requisitos:

1. Una máquina virtual con Ubuntu 20.04 LTS a la que poder hacer ssh o escritorio remoto.
2. Tener instalado ***Docker-ce*** como runtime de contenedor en dicha máquina virtual.


## Ejercicio 1: ***Publicar aplicación en el contenedor***

Vamos a crear una aplicación web elemental en un archivo html. Presentaremos este archivo al contenedor a través de un volumen. Cuando actualicemos la aplicación, el contenedor verá el archivo actualizado.

Cambiamos al directorio de trabajo
```
cd ~/k8s_desarrolladores/03/volumen
```

En el directorio existen dos archivos de configuración de nginx: ***nginx.conf*** y ***global.conf***. Comprobemos el contenido de cada uno.
```
nano nginx.conf
```

Este archivo configura a ***nginx*** para que se inicie solo con llamar a su ejecutable, sin necesidad de establecer parámetros (***-g daemon off;***).

Una vez estudiado, lo cerramos sin modificar.

En siguiente archivo ***global.conf*** tiene una configuración importante para entender el laboratorio. Lo editamos con el siguiente comando:
```
nano global.conf
```

* *Línea 5*: que tiene este valor: ***/var/www/html/website;*** indica el directorio donde vamos a incluir el archivo ***index.html*** de la aplicación. Esto es importante porque posteriormente usaremos un volumen que va a hacer referencia a esa ruta.

El archivo ***Dockerfile*** instala nginx y copia los archivos de configuración a la imagen de contenedor. Editamos el archivo:
```
nano Dockerfile
```

El contenido del archivo es el siguiente:

* *Línea 1*: Usamos una imagen base de ***Ubuntu 16.04***.
* *Línea 2*: Se actualiza el ***repositorio de paquetes*** y se ***instala nginx***.
* *Línea 3*: Creamos la carpeta de la aplicación en la ruta ***/var/www/html/website***. Será dentro de esa carpeta ***website*** donde coloquemos el archivo de la aplicación ***index.html*** a través de un volumen.
* *Líneas 4 y 5*: Se copian los archivos de configuración ***global.conf*** y ***nginx.conf*** desde el directorio de trabajo al sistema de archivos de la imagen de contenedor. Cuando se inicie nginx, leerá estos archivos.
* *Línea 6*: Se abre el puerto ***80*** en el contenedor.

Cerramos el arhivo sin modificarlo y procedemos a compilar la imagen con el siguiente comando:
```
sudo docker image build -t antsala/nginx .
```

Listamos las imágenes.
```
sudo docker image ls
```

En el directorio de trabajo, existe una carpeta llamada ***website***, que contiene el archivo ***index.html*** con la aplicación web. Este archivo por comodidad en la explicación es muy elemental. Lo editamos con el siguiente comando.
```
cp ./website/index_original.html ./website/index.html
nano ./website/index.html
```

Una vez estudiada la aplicación, cerramos sin modificar el archivo.

Aquí viene lo importante. Vamos a ejecutar un contenedor que ***monte un volumen***. En la sintaxis para montar el volumen, debemos especificar dos rutas:

1. Ruta del host: Es la carpeta donde está la aplicación, el archivo ***index.html***. Si el directorio actual es ***~/k8s_desarrolladores/03/volumen***, la ruta de la aplicación es ***$PWD/website***.
2. Ruta en el contenedor: Es la carpeta donde se realizará el montaje en el contenedor. En este caso, necesitamos que sea la carpeta en la que nginx espera encontrar los archivos de la aplicación, es decir, ***/var/www/html/website***.
3. El tercer argumento del volumen es ***ro*** o ***rw***, e indica si el contenedor puede modificar los archivos del host.

El comando que inicia el contenedor con el volumen es el siguiente. Nota: observar el al final ***docker run*** se ejecuta la aplicación ***nginx***, sin argumentos.
```
sudo docker container run -d -p 8080:80 --name website -v $PWD/website:/var/www/html/website:ro antsala/nginx nginx
```

Comprobamos que nginx "puede ver" al archivo 'index.html' a través del volumen:
```
curl localhost:8080
```

Editamos el archivo ***index.html*** e insertamos algún mensaje.
```
nano ./website/index.html
```

Añadimos, debajo de la línea 6, el siguiente código:
```
<h2>Acabo de modificar la aplicación desde el host</h2>
```

Guardamos y salimos.

Al volver a realizar una request al servidor, veremos la aplicación actualizada.
```
curl localhost:8080
```

Eliminamos los contenedores.
```
sudo docker container rm -f `sudo docker container ls -a -q`
```
