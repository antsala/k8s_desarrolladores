# Laboratorio 01: ***Herramientas de administración de EKS***
 
En este laboratorio instalaremos las herramientas que necesitaremos para administrar EKS.

Requisitos:

1. Una máquina virtual con Ubuntu 20.04 LTS a la que poder hacer ssh o escritorio remoto.
2. Tener instalada la herramienta ***kubectl***.

## Ejercicio 1: ***Instalación de AWS CLI***


Creamos un directorio para descargar la herramienta ***AWS CLI***.
```
mkdir -p ~/awscli
cd ~/awscli
```

Descargamos la aplicación desde la web de Amazon:
```
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
ls -l
```

La salida debe mostrar el archivo descargado.
```
-rw-rw-r-- 1 antonio antonio 46271162 mar 19 10:40 awscliv2.zip
```

Descomprimimos e instalamos la herramienta:
```
unzip awscliv2.zip
sudo ./aws/install
```

Una vez instalada, probamos si funciona:
```
aws --version 
```


## Ejercicio 2: ***Configuración de la credencial AWS para la CLI***

Ahora es necesario configurar las credencial de acceso a AWS en la CLI. Para ello desplegamos el menú usuario y elegimos ***Mis credenciales de seguridad***, en la web de AWS.

Seleccionamos la opción ***Crear una clave de acceso para CLI, SDK y API*** y, luego el botón ***Crear una clave de acceso***.

Descargamos el archivo CSV.

Procedemos a configurar la herramienta: (Nota: Proporcionar ***Access Key ID*** y ***Secret Access Key*** desde el archivo CSV. En región poner ***eu-west-1*** y formato de salida en blanco)
``
aws configure
```

## Ejercicio 3: ***Instalación y configuración de eksctl***

Creamos un directorio para la descarga:
```
mkdir -p ~/eksctl
cd ~/eksctl
```

Descargamos la herramienta y la movemos a un directorio apropiado:
```
curl --silent --location "https://github.com/weaveworks/eksctl/releases/latest/download/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
sudo mv /tmp/eksctl /usr/bin
```

Probamos que funcione:
```
eksctl version
```

## Ejercicio 4: ***Creación de EKS desde AWS CLI***


El siguiente comando creará un cluster y realizará lo siguiente:

1. Creará un cluster llamado ***myeks*** con nodos de tamaño ***t3.small***. Este tamaño permite 11 pods en el nodo. En número de pods depende del número de ENIs (Elastic Network Interface y este depende del tipo de instancia EC2. En el siguiente link se pueden comprobar estos limites: https://github.com/awslabs/amazon-eks-ami/blob/master/files/eni-max-pods.txt
2. Para el escalado del grupo de nodos indicamos que queremos un mínimo de ***1*** y un máximo de ***4***. Al grupo de nodos que se creará le llamaremos ***workers***.
```
eksctl create cluster \
    --name myeks \
    --version 1.21 \
    --region eu-west-1 \
    --nodegroup-name workers \
    --node-type t3.small \
    --nodes 3 \
    --nodes-min 1 \
    --nodes-max 4 \
    --managed
```

Se creará un stack de CloudFormation llamado ***eksctl-myeks-cluster***. El proceso puede tardar 15 minutos. Una vez terminado, solo queda actualizar ***kubeconfig*** para almacenar el nuevo contexto:
```
aws eks update-kubeconfig --name myeks --region eu-west-1
```

La salida del comando anterior será similar a esta:
```
Added new context arn:aws:eks:eu-west-1:779450087377:cluster/myeks to /home/antonio/.kube/config
```

## Ejercicio 5: ***Eliminación de EKS desde AWS CLI.***

Para eliminar el cluster hay que quitar todos los servicios que tengan asociados una ***EXTERNAL-IP***, ya que están balanceados por un balanceador ***ELB*** de AWS. Si no se borran, entonces la eliminación del cluster no eliminará el ELB asociado.
```
kubectl get services --all-namespaces
```

Tomar nota de los servicios y eliminarlos con el comando:
```
kubectl delete service <Poner aquí el nombre del servicio>
```

Ahora ya podemos borrar el cluster: (Nota: Tarda unos 5 minutos.)
```
eksctl delete cluster --name myeks
```

Comprobar a través de la GUI que no existe el cluster y borrar las credenciales.
