# Laboratorio 03-F: ***Frontend y Backend por medio de archivo YAML***
 
En este laboratorio crearemos volvemos  desplegar la aplicación ***SINATRA/REDIS*** pero esta vez de forma declarativa por medio de un archivo YAML.

Requisitos:

1. Una máquina virtual con Ubuntu 20.04 LTS a la que poder hacer ssh o escritorio remoto.
2. Tener instalado ***Docker-ce*** como runtime de contenedor en dicha máquina virtual.
3. Haber desplegado el cluster ***Swarm***.

## Ejercicio 1: ***Desplegar por medio de archivo YAML*** 

Cambiamos al directorio de trabajo.
```
cd ~/k8s_desarrolladores/03/frontend-backend/sinatra/webapp_redis
```

```
chmod +x ./bin/webapp
```

Editamos el archivo de manifiesto.
```
code Docker-compose.yaml
```

* *Línea 35 y 36*: Se crea una red llamada ***app*** para que el frontend y el backend se comuniquen.
* *Líneas 3-19*: Se define el servicio de frontend ***Sinatra***, donde destacamos:
* *Línea 15*: El puerto externo lo ponemos a ***8000***.
* *Línea 17*: Se monta un volumen para que el contenedor acceda a la aplicación.
* *Línea 18 y 19: Conectamos al contenedor a la red ***app***.

Muy importante!!!

En el ejercicio anterior te habrás dado cuenta que las instancias reciben un nombre especial cuando se despliegan en el cluster. Tienen la forma siguiente: ***NombreDelStack_NombreDelServicio.NumeroDeInstancia.IDContenedor***. De esta forma, el contenedor de redis se llamaría, por ejemplo así: ***webappRedisStack_db.1.r8d70eaexrpevi4scdpfdfldg***.
(Nota: Puedes verificarlo listando los contenedores cuando hayas desplegado el stack)         

Esto introduce una dificultad importante, porque si recuerdas, el código del frontend Sinatra, se conectaba al backend por medio del nombre del contenedor, que en el ejercicio anterior era ***db***, y si bien podríamos saber la parte fija del nombre del contenedor, resutará del todo imposible "adivinar" cual será el id del contenedor.

Para solucionar este problema, en Swarm un contenedor no se conecta a otro directamente (como hicimos en el ejercicio ***lab-03-D***), sino que el contenedor deberá conectarse a un SERVICIO y este, el servicio, redirigirá el tráfico a los contenedores apropiados.

Pues bien, para no modificar la aplicación Sinatra, que conectaba a "algo" llamado ***db***, lo que vamos a hacer es nombrar al servicio de backend (Redis) precísamente con ese nombre (***db***).

Seguimos con el archivo YAML.

* *Líneas 21 a 33*: Se define el servicio de backend de Redis, que como hemos aclarado debe llamarse obligatoriamente ***db***.
* *Línees 32 y 33*: Conectamos los contenedores del servicio Redis a la red ***App***.

Desplegamos el Stack
```
sudo docker stack deploy -c Docker-compose.yaml webappRedisStack
```

Comprobamos que se ha desplegado correctamente.
```
sudo docker stack ps webappRedisStack
```

Como curiosidad, listamos los contenedores.
```
sudo docker container ls -a
```

Solo resta probar la aplicación.

Ahora probamos. En primer lugar una GET al directorio raíz, debe devolver un mensaje.
```
curl localhost:8000
```

Enviamos una request con POST a ***/json/*** con los campos de un formulario. Esta vez se guardará en redis, además de devolverse en formato JSON.
```
curl -i -H 'Accept: application/json' -d 'nombre=Antonio&apellidos=Salazar Gravan&telefono=666123321' http://localhost:8000/json/
```

Por último, una GET a ***/json*** que provocará una lectura de Redis para leer los parámetros del formulario, que serán devueltos en la response en formato JSON.
```
curl http://localhost:8000/json
```

Nuestra aplicación está formada por dos partes, donde el Frontend NO TIENE ESTADO y el Backend SÍ lo tiene. 

Para las aplicaciones SIN estado, puedes cambiar el número de replicas sin problemas, pero para las que lo tienen, esto no se puede hacer. Kubernetes ofrece un tipo de objeto especial para resolver esta limitación, llamado ***StatefulSet***, pero Docker no lo tiene.

Como práctica adicional te animo a que escales a 5 réplicas el ***FRONTEND*** y pruebes la aplicación. Podrás comprobar que funciona perfectamente.

Eliminamos la aplicación. 
```
sudo docker stack rm webappRedisStack
```

Destruimos el cluster.
```
sudo docker swarm leave --force
```

Borramos los contenedores y la imágenes.
```
sudo docker container rm -f `sudo docker container ls -a -q`
sudo docker image rm -f `sudo docker image ls -q`
```

