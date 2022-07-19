# Laboratorio 20-A: ***Instalación de Minikube***
 
En este laboratorio instalará Minikube en una máquina virtual Ubuntu 20.04 LTS

Requisitos:

1. Una máquina virtual con Ubuntu 20.04 LTS a la que poder hacer ssh o escritorio remoto.
2. Tener instalado el runtime de ***podman***. (ver lab-06-A.md, Ejercicios 1 y 2)


## Ejercicio 1: ***Instalación de Minikube***

Ya tenemos el runtime de contenedor podman instalado. Ahora prodecemos a descargar el binario de Minikube desde la web de Google.
```
wget https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
```

Comprobamos que se ha descargado listando el directorio. Debe aparecer un archivo con nombre ***minikube-linux-amd64***
```
ls -l 
```

Vamos a mover el archivo descargado a un directorio del sistema. En el mismo comando aprovechamos y le cambiamos el nombre a ***minikube***
```
sudo mv minikube-linux-amd64 /usr/local/bin/minikube
```

Comprobamos que ha quedado en su sitio. Observar que no tiene permiso de ejecución.
```
sudo ls -l /usr/local/bin/minikube
```

Añadimos el permiso de ejecución al binario de Minikube.
```
sudo chmod +x /usr/local/bin/minikube
```

Comprobamos que ya es ejecutable. El propietario del archivo sigue siendo el usuario con el que estás logado. No hace falta cambiarlo.
```
sudo ls -l /usr/local/bin/minikube
```

Comprobamos que podemos ejecutarlo
```
minikube version
```

Con esto termina la instalación del binario de Minikube.


## Ejercicio 2: ***Instalación de kubectl***


***kubectl*** es la herramienta de línea de comando para interactuar con el ***API-Server (Control Plane)***. Esta herramienta también hay que descargarla e instalarla.

Procedemos a descargar la última versión estable.
```
curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl
```

Comprobamos. Debe haberse descargado un archivo llamado ***kubectl***.
```
ls -l
```

Procedemos a moverlo a una ubicación más apropiada y asignarle permisos de ejecución.
```
sudo mv kubectl /usr/local/bin/
chmod +x /usr/local/bin/kubectl
```

Comprobamos
```
ls -l /usr/local/bin/kubectl
```

## Ejercicio 3: ***Modificar sudoers e instalar uidmap

Ahora probamos la herramienta ***kubectl*** mostrando su versión. Aparecerá un error que indica que no puede conectar con ***localhost:8080***. Este comportamiento ***es correcto*** porque aun no hemos iniciado ***Minikube***, solo hemos descargado su binario, así que ***kubectl*** no puede todavía contactar con el ***API Server***. 
```
kubectl version -o yaml
```

El siguiente paso sería iniciar Minikube, pero daría error. Esto es así porque el driver podman no se puede ejecutar con permisos de root. Para resolverlo, es necesario anadir una línea al fichero ***sudoers***.

Editar ***sudoers*** con ***visudo*** de la siguiente forma:
```
sudo visudo
```

Desplazarse hasta el final del archivo ***sudoers*** y añadir una nueva línea debajo de ***#includedir /etc/sudoers.d***. Esta línea debe tener el siguiente contenido. Nota: Sustituye ***<usuario>*** por tu usuario de Linux.
```
<usuario> ALL=(ALL) NOPASSWD: /usr/bin/podman
```

Guardar con CTRL+X, Y, ENTER.


Debemos asegurar que la distribución contiene el paquete ***uidmap***. Este paquete es necesario para poder trabajar con podman en modo ***rootless*** cuando está instalado Minikube.
```
sudo apt-get -y install uidmap
```

## Ejercicio 4: ***Iniciar Minikube***

Si ejecutamos el comando ***minikube***, se mostrará la ayuda. 
```
minikube
```

Observar los comandos básicos:
```
start          Inicia un cluster local de Kubernetes.
status         Muestra el estado del cluster.
stop           Detiene la ejecución del cluster de Kubernetes.
delete         Elimina el cluster.
dashboard      Levanta la interfaz gráfica de administración.
```

Para iniciar Minikube no es necesario indicar el runtime de contenedor que queremos usar. Esto solo sería necesario si usamos uno diferente a Docker, que es precísamente nuestro caso. Así que para iniciar Minikube con el driver (runtime) podman, debemos ejecutar el siguiente comando.
```
minikube start --driver=podman 
```

Todos los servicios del ***Control Plane***, ***Kubelet*** y ***k-proxy*** se ejecutan en un único contenedor (en este caso de podman). Podemos ver el contenedor, que se llama ***minikube*** con el siguiente comando.
```
sudo podman container ls
```

Una vez finalizado el arranque de Minikube, probamos que ***kubectl**** puede comunicarse con el ***Control Plane***. Aparecerá información sobre el único nodo de nuestro cluster de Kubernetes.
```
kubectl get nodes
```

El comando ***minikube status*** mostrará información de estado de ejecución de Minikube. Ejecutar el siguiente comando y observar el resultado.
```
minikube status
```

Aunque ***kubectl*** será la principal herramienta para administrar el cluster, también podemos levantar una interfaz gráfica. Para ello ejecutamos el siguiente comando. Nota: Se abrirá automáticamente el navegador y conectará con la URL: 'http://127.0.0.1:40743/api/v1/namespaces/kubernetes-dashboard/services/http:kubernetes-dashboard:/proxy/'
```
minikube dashboard
```

La consola ejecuta el dashboard. Para finalizar dicha ejecución pulsamos ***CTRL+C***.

Minikube realmente no se instala en el equipo, simplemente se inicia. Cuando finalicemos nuestra sesión de trabajo debemos detener el contenedor que nos ofrece los servicios del cluster. Esto lo hacemos para liberar recursos de hardware de nuestro equipo.
```
minikube stop
```

Realmente, lo que hemos hecho es detener el contenedor de podman que ofrece los servicios del cluster. Un posterior ***minikube start*** volverá a iniciar el contenedor y, con él, los servicios de cluster.

Si en algún momento se produjera un error al iniciar el contenedor de podman con los servicios del cluster, debemos eliminarlo, con el comando ***minikube delete*** y proceder a un nuevo start, que volverá a descargar la imagen de contenedor desde el repositorio.

Si aun así, Minikube no se inicia, podemos arrancarlo en modo de depuración (DEBUG). Para ello debemos eliminar primero el contenedor y volver a iniciarlo en modo de depuración con los siguientes comandos y observar dónde se produce el problema.
```
minikube delete
minikube start --driver=podman --v=7 --alsologtostderr
```
