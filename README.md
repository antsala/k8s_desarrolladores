# k8s_desarrolladores

Para poder realizar los laboratorios es necesario clonar el repo desde GitHub.Para ello realizamos abrimos una consola de comandos e instalamos ***git*** en la máquina virtual.

```
sudo apt-get -y update
sudo apt-get -y install git
```

Ahora clonamos el repositorio con los laboratorios del curso.

```
cd ~
git clone https://github.com/antsala/k8s_desarrolladores.git
cd ~/k8s_desarrolladores
```

Instalamos ***Visual Studio Code*** (Si lo prefieres instala tu editor preferido)

```
sudo snap install code --classic
```

El repositorio GIT está organizado en una serie de carpetas. Las presentaciones de PowerPoint irán indicando en cada momento el laboratorio a realizar. No obstante, te presento la lista de actividades que se realizarán en el curso, y la carpeta en la que se ubica.

## Carpeta 00 (AZURE)

1. Laboratorio 00: ***Herramientas de administración en Azure***. Los ejercicios a realizar son:
   - Instalación de Azure CLI.
   - Creación de AKS desde Azure CLI.
   - Eliminación de AKS desde Azure CLI.


##  Carpeta 01 (AWS)

1. Laboratorio 01: ***Herramientas de administración en AWS***. Ejercicios:
   - Instalación de ***AWS CLI***.
   - Configuración de la credencial AWS para la ***CLI***.
   - Instalación y configuración de ***eksctl***.
   - Creación de ***EKS*** desde ***AWS CLI***.
   - Eliminación de ***EKS*** desde ***AWS CLI***.


## Carpeta 03

1. Laboratorio 03-A: ***Creación de contenedores con Docker***. Ejercicios:
   - Instalación de ***Docker***.
   - Primeros ***contenedores*** con Docker.
   - ***Imágenes*** con Docker.

2. Laboratorio 03-B: ***Construir imágenes desde Dockerfile***. Ejercicios:
   - Creación de imagen desde ***Dockerfile***.
   - ***Publicación*** de puertos en el host.
   - ***ENTRYPOINT*** en ***Dockerfile***.
   - ***ENTRYPOINT*** y ***CMD*** en ***Dockerfile***.

3. Laboratorio 03-C: ***Volúmenes***. Ejercicios:
   - Publicar aplicación en el contenedor.

4. Laboratorio 03-D: ***Frontend-Backend***. Ejercicios:
   - Creación del ***Frontend***.
   - Creación del ***Backend***.
   - Creación de una ***red***.
   - Recreación del Backend conectado a la nueva red.
   - Despliegue de la versión de Frontend que conecta con Backend.

5. Laboratorio 03-E: ***Micro servicios***. Ejercicios:
   - Instalación de ***Go***.
   - ***Compilar*** una app en Go.
   - ***Contenerizar*** la app Go.
   - Desplegar servicio en ***Swarm***.


## Carpeta 06

1. Laboratorio 06-A: ***Instalar Podman***. Ejercicios:
   - Desinstalación de Docker.
   - Instalación de ***Podman***.

2. Laboratorio 06-B: ***Frontend-Backend con POD***. Ejercicios:
   - Descargar imágenes de contenedor.
   - Archivo con ***variables de entorno***.
   - ***Creación*** del pod.
   - ***Eliminación*** del pod.


## Carpeta 20

1. Laboratorio 20-A: ***Instalación de Minikube***. Ejercicios:
   - Instalación de ***Minikube***.
   - Instalación de ***kubectl***.
   - Modificar ***sudoers*** e instalar ***uidmap***.
   - Iniciar ***Minikube***.

2. Laboratorio 20-B: ***Comandos básicos de kubectl***. Ejercicios:
   - Primera toma de contacto con ***kubectl***.
   - Crear un ***deployment*** con kubectl.
   - El ***ReplicaSet***.
   - ***Editar*** un deployment con kubectl.
   - ***Rollout undo*** del deployment con kubectl.
   - ***Describir*** un objeto con kubectl.
   - Ver la ***salida estándar*** del contenedor con kubectl.
   - ***Ejecutar comandos*** en el contenedor con kubectl.
   - ***Eliminar objetos*** del cluster con kubectl.


## Carpeta 25

1. Laboratorio 25-A: ***Despliegue de archivos YAML***. Ejercicios:
   - Aplicar un ***deployment*** desde archivo YAML.
   - Aplicar un ***servicio*** desde archivo YAML.
   - Obtener el ***estado de Kubernetes***.
   - Crear un servicio de tipo ***LoadBalancer*** desde archivo YAML.

2. Laboratorio 25-B: ***Despliegue de MongoDB***. Ejercicios:
   - Descripción del sistema.
   - Crear la base de datos ***MongoDB***.
   - Crear un ***secreto*** en Kubernetes.
   - Aplicar el ***deployment*** de MongoDB.
   - Crear el deployment ***Mongo Express***.
   - Aplicar el deployment de ***Mongo Express***.

3. Laboratorio 25-C: ***Backend de Redis con un master y dos réplicas***
   - Despliegue del ***maestro de Redis***.
   - Creación de ***ConfigMap*** desde un archivo.
   - Despliegue de las ***réplicas de Redis***.
   - Despliegue del ***Frontend***.
   - Despliegue del ***balanceador*** para el Frontend.


## Carpeta 30

1. Laboratorio 30: ***Espacios de Nombres***. Ejercicios:
   - Creación de un ***espacio de nombres***.
   - ***Aplicar*** un archivo YAML en un espacio de nombres.
   - ***Predeterminar*** el espacio de nombres.
   - Predeterminar el espacio de nombres en el ***archivo YAML***.


## Carpeta 35

1. Laboratorio 35-A: ***Ingress***. Ejercicios:
   - Despliegue de ***helloContainer*** y ***mongodb***.
   - Creación del ***objeto Ingress***.
   - Instalar el ***Controlador Ingress***.
   - Configurar el ***Registro de Recurso*** de DNS.
   - Configurar el ***Default Backend***.

2. Laboratorio 35-B: ***Asegurar el Ingress con TLS***. Ejercicios:
   - Configuración de un ***Gateway de aplicación de Azure*** como ***Ingress*** de K8s.
   - Añadir una ***regla de entrada*** (Ingress) a la aplicación.
   - Instalación de ***cert-manager***.
   - Instalación del ***emisor de certificados*** (issuer).
   - ***Crear*** el certificado TLS y ***asegurar*** la Ingress.
   - Cambiar al entorno de producción de Let´s Encrypt.


## Carpeta 45

1. Laboratorio 45-A: ***Horizontal POD Autoscaler (HPA) en Azure y autoescalado de nodos***. Ejercicios:
   - Desplegar la aplicación de ejemplo.
   - ***Escalar*** el frontend de GuestBook.
   - ***Autoescalado*** de nodos.

2. Laboratorio 45-B: ***Horizontal POD Autoscaler (HPA) en AWS***. Ejercicios:
   - Desplegar la aplicación de ejemplo.
   - ***Escalar*** el frontend de GuestBook.


## Carpeta 50

1. Laboratorio 50: ***Instalar aplicaciones usando HELM***. Ejercicios:
   - Instalar ***Helm***.
   - ***Instalar WordPress*** con Helm.
   - ***Desinstalar WordPress*** con Helm.


## Carpeta 55

1. Laboratorio 55-A: ***Usar volúmenes en las aplicaciones***. Ejercicios:
   - Despliegue del servidor Redis.
   - Creación de un ***ConfigMap*** desde un archivo.
   - Creación de un ConfigMap desde archivo YAML.
   - Instalar aplicaciones ***con estado (statefull)*** en el cluster.
   - ***Persistent Volume Claims*** (PVCs).

2. Laboratorio 55-B: ***MySQL replicado con StatefulSet***. Ejercicios:
   - Crear el ***ConfigMap***.
   - Creación de los ***servicios***.  
   - Creación del ***StatefulSet***.
   - ***Enviar tráfico*** desde el Frontend.

3. Laboratorio 55-C: ***MongoDB con StatefulSet y Sidecar***. Ejercicios:
   - Crear ***secretos*** y ***script de inicio***.
   - ***Proteger*** las comunicaciones de MongoDB.
   - Cambiar permisos al script.
   - Creación de un ***espacio de nombres***.
   - Creación del ***servicio*** y ***cuenta de servicio***.
   - Creación del ***StatefulSet*** de MongoDB.
   - Creación del objeto ***Kustomization***.


## Carpeta 60

1. Laboratorio 60: ***Despliegue de Prometeus con Operator***. Ejercicios:
   - Instalación del ***Operator*** de Prometheus.

## Carpeta 65

1. Laboratorio 65-A: ***Endpoints externos***. Ejercicios:
   - Creación del ***endpoint***.
   - Creación del ***servicio interno***.
   - Prueba del endpoint.

2. Laboratorio 65-B: ***Monitorización del cluster***. Ejercicios: 
   - Sondas ***Readiness*** y ***Liveness***.
   - ***Depuración de errores*** en el pull de imágenes.
   - Errores de la aplicación.
   - Sondas ***Readiness*** y ***liveness*** (Revisión).
   - Experimentos con ***liveness*** y ***readiness***.
   - ***Métricas simples***.

3. Laboratorio 65-C: ***RBAC en AKS (Azure)***. Ejercicios: 
   - Introduccion a ***RBAC***.
   - Habilitar la ***integración de Azure AD*** en ***AKS***.
   - Añadir al usuario ***administrador del tenant*** al grupo ***aks admins***.
   - Crear un usuario y un grupo de seguridad para asignar roles.
   - ***Configurar*** RBAC en AKS.
   - ***Verificar*** RBAC para el usuario Luke.

