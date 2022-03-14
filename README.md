# k8s_desarrolladores

# Para poder realizar los laboratorios es necesario clonar el repo desde github. 
# Para ello realizamos lo siguiente:

# Abrir una consola de comandos e instalar git en la máquina virtual

sudo apt-get -y update

sudo apt-get -y install git


# Clonar el repositorio con los laboratorios del curso.

cd ~

git clone https://github.com/antsala/k8s_desarrolladores.git

cd ~/k8s_desarrolladores

ls -l


# Instalar Visual Studio Code

sudo snap install code --classic

# Carpeta 03

# Laboratorio 03-A: "Creación de contenedores con Docker"
# Archivo: lab-03-A.txt
#
# Ejercicios:
#   1. Instalación de 'Docker'.
#   2. Primeros contenedores con 'Docker',
#   3. Imágenes con 'Docker'.


# Laboratorio 03-B: "Construir imágenes desde Dockerfile"
# Archivo: lab-03-B.txt
#
# Ejercicios:
#   1. Creación de imagen desde Dockerfile.
#   2. Publicación de puertos en el host.
#   3. 'ENTRYPOINT' en 'Dockerfile'.
#   4. 'ENTRYPOINT' y 'CMD' en 'Dockerfile'.


# Laboratorio 03-C: "Volúmenes"
# Archivo: lab-03-C.txt
#
# Ejercicios:
#   1. Publicar aplicación en el contenedor.


# Laboratorio 03-D: "Frontend-Backend"
# Archivo: lab-03-D.txt
#
# Ejercicios:
#   1. Creación del Frontend.
#   2. Creación del Backend.
#   3. Creación de una red.
#   4. Recreación del Backend conectado a la nueva red.
#   5. Despliegue de la versión de Frontend que conecta con Backend.


# Laboratorio 03-E: "Microservicios"
# Archivo: lab-03-E.txt
#
# Ejercicios:
#   1. Instalación de Go.



# Carpeta 20

# Laboratorio 20-A: "Instalación de Minikube"
# Archivo: lab-20-A.txt
#
# Ejercicios:
#   1. Instalación de 'Podman'.
#   2. Instalación de 'Minikube'.
#   3. Instalación de 'kubectl',
#   4. Modificar 'sudoers' e instalar 'uidmap',
#   5. Iniciar Minikube,


# Laboratorio 20-B: "Comandos básicos de kubectl"
# Archivo: lab-20-B.txt
#
# Ejercicios:
#   1. Primera toma de contacto con 'kubectl',
#   2. Crear un deployment con 'kubectl'.
#   3. El 'ReplicaSet'.
#   4. Editar un deployment con 'kubectl'.
#   5. 'Rollout undo' del deployment con 'kubectl'.
#   6. Describir un objeto con 'kubectl'.
#   7. Ver la salida estándar del contenedor con 'kubectl'.
#   8. Ejecutar comandos en el contenedor con 'kubectl'.
#   9. Eliminar objetos del cluster con 'kubectl'.
#   10. Aplicar un archivo YAML con 'kubectl'.