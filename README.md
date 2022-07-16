# k8s_desarrolladores

Para poder realizar los laboratorios es necesario clonar el repo desde github.Para ello realizamos abrimos una consola de comandos e instalamos git en la máquina virtual.

```
sudo apt-get -y update
sudo apt-get -y install git
```


# Clonar el repositorio con los laboratorios del curso.

cd ~

git clone https://github.com/antsala/k8s_desarrolladores.git

cd ~/k8s_desarrolladores

ls -l


# Instalar Visual Studio Code

sudo snap install code --classic


# Carpeta 00 (AZURE)

# Laboratorio 00: "Herramientas administración Azure"
# Archivo: lab-00.txt
#
# Ejercicios: 
#   1. Instalación de Azure CLI.
#   2. Creación de AKS desde Azure CLI.
#   3. Eliminación de AKS desde Azure CLI.


# Carpeta 01 (AWS)

# Laboratorio 01: "Herramientas administración AWS"
# Archivo: lab-01.txt
#
# Ejercicios: 
#   1. Instalación de AWS CLI.
#   2. Configuración de la credencial AWS para la cli.
#   3. Instalación y configuración de 'eksctl'.
#   4. Creación de EKS desde AWS CLI.
#   5. Eliminación de EKS desde AWS CLI.



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


# Laboratorio 03-E: "Micro servicios"
# Archivo: lab-03-E.txt
#
# Ejercicios:
#   1. Instalación de Go.
#   2. Compilar una app en Go.
#   3. Contenerizar la app Go.
#   4. Desplegar servicio en Swarm.



# Carpeta 06

# Laboratorio 06-A: "Instalar Podman"
# Archivo: lab-20-A.txt
#
# Ejercicios:
#   1.  Desinstalación de Docker.
#   2.  Instalación de Podman.


# Laboratorio 06-B: "Frontend-Backend con POD"
# Archivo: lab-20-B.txt
#
# Ejercicios:
#   1.  Descargar imágenes de contenedor.
#   2.  Archivo con variables de entorno.
#   3.  Creación del pod.
#   4.  Eliminación del pod.



# Carpeta 20

# Laboratorio 20-A: "Instalación de Minikube"
# Archivo: lab-20-A.txt
#
# Ejercicios:
#   1. Instalación de 'Minikube'.
#   2. Instalación de 'kubectl'.
#   3. Modificar 'sudoers' e instalar 'uidmap'.
#   4. Iniciar Minikube.


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



# Carpeta 25

# Laboratorio 25-A: "Despliegue de archivos YAML"
# Archivo: lab-25-A.txt
#
# Ejercicios:
#   1. Aplicar un deployment desde archivo YAML.
#   2. Aplicar un servicio desde archivo YAML.
#   3. Obtener el estado de Kubernetes.
#   4. Crear un servicio de tipo 'LoadBalancer' desde archivo YAML.


# Laboratorio 25-B: "Despliegue de MongoDB"
# Archivo: lab-25-B.txt
#
# Ejercicios:
#   1. Descripción del sistema.
#   2. Crear la base de datos MongoDB.
#   3. Crear un secreto en Kubernetes.
#   4. Aplicar el deployment de MongoDB.
#   5. Crear el deployment 'Mongo Express'.
#   6. Aplicar el deployment de 'Mongo Express'.


# Laboratorio 25-C: "Backend de Redis con un master y dos réplicas"
# Archivo: lab-25-C.txt
#
# Ejercicios:
#   1. Despliegue del maestro de Redis.
#   2. Creación de ConfigMap desde un archivo.
#   3. Despliegue de las réplicas de Redis.
#   4. Despliegue del Frontend.
#   5. Despliegue del balanceador para el Frontend.


# Carpeta 30

# Laboratorio 30: "Espacios de Nombres"
# Archivo: lab-30.txt
#
# Ejercicios:
#   1. Creación de un espacio de nombres.
#   2. Aplicar un archivo YAML en un espacio de nombres.
#   3. Predeterminar el espacio de nombres.
#   4. Predeterminar el espacio de nombres en el archivo YAML.


# Carpeta 35

# Laboratorio 35-A: "Ingress"
# Archivo: lab-35-A.txt
#
# Ejercicios:
#   1. Despliegue de 'helloContainer' y 'mongodb'.
#   2. Creación del objeto Ingress.
#   3. Instalar el Controlador Ingress.
#   4. Configurar el Registro de Recurso de DNS.
#   5. Configurar el 'Default Backend'.


# Laboratorio 35-B: "Asegurar el Ingress con TLS"
# Archivo: lab-35-B.txt
#
# Ejercicios:
#   1. Configuración de un Gateway de aplicación de Azure como Ingress de K8s.
#   2. Añadir una regla de entrada (Ingress) a la aplicación.
#   3. Instalación de 'cert-manager'.
#   4. Instalación del emisor de certificador (issuer).
#   5. Crear el certificado TLS y asegurar la Ingress.
#   6. Cambiar al entorno de producción de Let´s Encrypt.


# Carpeta 45

# Laboratorio 45-A: "Horizontal POD Autoscaler (HPA) en Azure y autoescalado de nodos."
# Archivo: lab-45-A.txt
#
# Ejercicios:
#   1. Desplegar la aplicación de ejemplo.
#   2. Escalar el frontend de GuestBook.
#   3. Autoescalado de nodos.


# Laboratorio 45-B: "Horizontal POD Autoscaler (HPA) en AWS."
# Archivo: lab-45-B.txt
#
# Ejercicios:
#   1. Desplegar la aplicación de ejemplo.
#   2. Escalar el frontend de GuestBook.


# Carpeta 50

# Laboratorio 50: "Instalar aplicaciones usando HELM."
# Archivo: lab-50.txt
#
# Ejercicios:
#   1. Instalar Helm.
#   2. Instalar WordPress con Helm.
#   3. Desinstalar WordPress con Helm.


# Carpeta 55

# Laboratorio 55-A: "Usar volúmenes en la aplicaciones."
# Archivo: lab-55-A.txt
#
# Ejercicios:
#   1. Despliegue del servidor Redis.
#   2. Creación de un ConfigMap desde un archivo.
#   3. Creación de un ConfigMap desde archivo YAML.
#   4. Instalar aplicaciones con estado en el cluster.
#   5. Persistent Volume Claims.


# Laboratorio 55-B: "MySQL replicado con StatefulSet."
# Archivo: lab-55-B.txt
#
# Ejercicios:
#   1. Crear el ConfigMap.
#   2. Creación de los servicios.  
#   3. Creación del StatefulSet.
#   4. Enviar tráfico desde el Frontend.


# Laboratorio 55-C: "MongoDB con StatefulSet y Sidecar"
# Archivo: lab-55-B.txt
#
# Ejercicios:
#   1. Crear secretos y script de inicio.
#   2. Proteger las comunicaciones de MongoDB.
#   3. Cambiar permisos al script.
#   4. Creación de un espacio de nombres.
#   5. Creación del servicio y cuenta de servicio.
#   6. Creación del 'StatefulSet' de MongoDB.
#   7. Creación del objeto 'Kustomization'.


# Laboratorio 60: "Despliegue de Prometeus con Operator"
# Archivo: lab-60.txt
#
# Ejercicios:
#   1. Instalación del Operator de Prometheus.


# Laboratorio 65-A: "Endpoints externos"
# Archivo: lab-65-A.txt
#
# Ejercicios:
#   1. Creación del endpoint.
#   2. Creación del servicio interno.
#   3. Prueba del endpoint.



# Laboratorio 65-B: "Monitorización del cluster."
# Archivo: lab-65-B.txt
#
# Ejercicios:
#   1. Sondas 'Readiness' y 'Liveness'.
#   2. Depuración de errores en el pull de imágenes.
#   3. Errores de la aplicación.
#   4. Sondas 'Readiness' y 'liveness'.
#   5. Experimentos con 'liveness' y 'readiness'.
#   6. Métricas simples.


# Laboratorio 65-C: "RBAC en AKS (Azure)"
# Archivo: lab-65-C.txt
#
# Ejercicios:
#   1. Introduccion a RBAC.
#   2. Habilitar la integración de Azure AD en AKS.
#   3. Añadir al usuario administrador del tenant al grupo 'aks admins'.
#   4. Crear un usuario y un grupo de seguridad para asignar roles.
#   5. Configurar RBAC en AKS.
#   6. Verificar RBAC para el usuario Luke.



