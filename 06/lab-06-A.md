# Laboratorio 06-A: ***Instalar podman***
 
En este laboratorio instalaremos el runtime de contenedor ***podman*** que nos servirá para presentar la abstracción del pod, objeto que se usará en Kubernetes.

Requisitos:

1. Una máquina virtual con Ubuntu 20.04 LTS a la que poder hacer ssh o escritorio remoto.


## Ejercicio 1: ****Desinstalación de Docker***

Si se han realizado los ejercicios del módulo anterior, Docker estará instalado. Para poder realizar una instalación de ***podman*** correctamente, debemos desinstalar docker.

Ejecutamos los siguientes comandos para desinstalar ***Docker***:
```
sudo apt-get purge -y docker-ce
```
```
sudo apt-get purge -y docker-engine docker docker.io docker-ce
sudo apt-get autoremove -y --purge docker-engine docker docker.io docker-ce
sudo rm -rf /var/lib/docker /etc/Docker
sudo groupdel docker
sudo rm -rf /var/run/docker.sock
```

Comprobamos que Docker ya no está con los siguientes comandos:
```
sudo systemctl status docker
```
```
sudo docker --version
```

## Ejercicio 2: ***Instalación de Podman***

Necesitamos un runtime de contenedor. Vamos a instalar ***podman***, de RedHat, que a diferencia de Docker permite trabajar con la abstracción del ***POD***. Lo que aprendamos aquí sobre los pods, será de directa aplicación en Kubernetes.

Aseguramos que 'curl' está instalado.
```
sudo apt-get -y update
```
```
sudo apt-get -y install curl
```

Para Ubuntu 20.04 no se puede instalar podman directamente desde los repositorios oficiales. Para conseguirlo, añadimos a la lista de repos, el de podman (RedHat).
```
echo "deb https://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/xUbuntu_20.04/ /" | sudo tee /etc/apt/sources.list.d/devel:kubic:libcontainers:stable.list
curl -L "https://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/xUbuntu_20.04/Release.key" | sudo apt-key add -
```

Volvemos a actualizar el repositorio de paquetes.
```
sudo apt-get -y update
```

Instalamos podman.
```
sudo apt-get -y install podman
```

Comprobamos que podman se ha instalado correctamente. 
```
podman --version
```

Probamos a lanzar un contenedor y ejecutar un comando en él. Como resultado aparecerá información sobre la versión del sistema operativo del contenedor.
```
podman run --rm docker.io/library/ubuntu:latest cat /etc/lsb-release
```

Con esto finalizamos la instalación del runtime Podman.

