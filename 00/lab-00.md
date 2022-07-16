# Laboratorio 00: ***Herramientas de administración de Azure***
<br/>
En este laboratorio instalaremos las herramientas que necesitaremos para administrar Azure.

Los requisitos son:

1. Una máquina virtual con ***Ubuntu 20.04 LTS*** a la que poder hacer ssh o tener un escritorio remoto.
2. Una subscripción de Azure que ***permita*** crear clústeres de AKS
<br/>
<br/>
<br/>
<br/>
## Ejercicio 1: ***Instalación de Azure CLI***

En primer lugar desinstalamos versiones previas si estuvieran presentes:

```
sudo apt remove azure-cli -y
```
```
sudo apt autoremove -y
```




Actualizamos repositorios e instalamos dependencias:

```
sudo apt-get update
```
```
sudo apt-get install -y  ca-certificates curl apt-transport-https lsb-release gnupg
```




Descargamos la clave de firma de Microsoft:

```
curl -sL https://packages.microsoft.com/keys/microsoft.asc | \
    gpg --dearmor | \
    sudo tee /etc/apt/trusted.gpg.d/microsoft.gpg > /dev/null
```




Agregamos repositoritos de Azure-CLI:

```
AZ_REPO=$(lsb_release -cs) 
echo "deb [arch=amd64] https://packages.microsoft.com/repos/azure-cli/ $AZ_REPO main" | \
    sudo tee /etc/apt/sources.list.d/azure-cli.list
```




Actualizamos repos e instalamos azure-CLI:

```
sudo apt-get update
```
```
sudo apt-get install azure-cli
```




Comprobamos versión de ***azure-cli*** y actualizarla con ***az upgrade*** si se recomienda.

```
az version
```




# Ejercicio 2: ***Creación de AKS desde Azure CLI*** 

Iniciamos sesión con el usuario ***administrador*** de la subscripción de Azure.

```
az login
```




Creamos un grupo de recursos para el cluster.

```
az group create \
    --name myaks-rg \
    --location westeurope
```




Habilitamos la supervisión de clusteres.

```
az provider register \
    --namespace Microsoft.OperationsManagement

az provider register \
    --namespace Microsoft.OperationalInsights
```




Creamos el cluster. 

```
az aks create \
    --resource-group myaks-rg \
    --name myaks \
    --location westeurope \
    --node-count 2 \
    --node-vm-size Standard_DS2_v2 \
    --enable-addons monitoring \
    --generate-ssh-keys
```




Una vez conectado a la subscripción, el siguiente paso es conectar al servicio AKS. El siguiente comando descarga las credenciales y las almacena en ***./kube/config***.

```
az aks get-credentials \
    --resource-group myaks-rg \
    --name myaks \
    --admin \
    --overwrite-existing
```




Comprobamos el estado del cluster

```
az aks show \
    --resource-group myaks-rg \
    --name myaks
```




***Nota sobre los contextos***. Si se quiere cambiar de cluster, debemos cambiar el contexto. Primero listamos los contextos configurados.

```
kubectl config get-contexts
```




La salida del comando anterior mostrará algo como esto:
```
CURRENT   NAME                                               CLUSTER                                            AUTHINFO                                           NAMESPACE
          arn:aws:eks:eu-west-1:779450087377:cluster/myeks   arn:aws:eks:eu-west-1:779450087377:cluster/myeks   arn:aws:eks:eu-west-1:779450087377:cluster/myeks   
          minikube                                           minikube                                           minikube                                           default
*         myaks-admin                                        myaks                                              clusterAdmin_myaks-rg_myaks  
```

Podemos apreciar que el contexto actual es ***myaks-admin*** (El cluster AKS señalizado por el asterisco *).



Para conmutar al otro contexto (Minikube) usamos el siguiente comando.

```
kubectl config use-context minikube
```




La salida indicará lo siguiente:

```
Switched to context "minikube".
```




Para volver al contexto de Azure:

```
kubectl config use-context myaks-admin
```





## Ejercicio 3: Eliminación de AKS desde ***Azure CLI***

Eliminamos el cluster.

```
az aks delete \
    --name myaks \
    --resource-group myaks-rg \
    --yes
```




Eliminamos el grupo de recursos que contiene el cluster.

```
az group delete \
    --resource-group myaks-rg \
    --yes
```




Eliminamos el grupo de recursos que contiene los objetos de las ***Azure Functions***.

```
az group delete \
    --resource-group functions-rg \
    --yes
```




Si se ha creado el ***Application Gateway***, eliminamos el grupo de recursos. Tarda mucho, lo eliminamos de forma asíncrona con ***--no-wait***.

```
az group delete \
    --resource-group agic-rg \
    --yes \
    --no-wait
```




Si se han creado ***Grupos*** y ***usuarios*** en Azure AD para la integración de AKS con AAD, los eliminamos.

```
az ad group delete \
    --group "aks admins"

az ad group delete \
    --group "aks users"

az ad user delete \
    --id luke@antsalgrahotmail.onmicrosoft.com
```




Si se registró la característica ***EnablePodIdentityPreview***, la desregistramos de la subscripción.

```
az feature unregister \
    --name EnablePodIdentityPreview \
    --namespace Microsoft.ContainerService
```




Como el aviso indica, también hay que ejecutar el siguiente comando para que se ***propage el cambio***.

```
az provider register \
    --name Microsoft.ContainerService
```




Quitamos la extensión ***aks-preview*** de la CLI si estuviera instalada.

```
az extension remove \
    --name aks-preview
```




Eliminamos el ACR si fue creado.

```
az group delete \
    --resource-group myACR-rg \
    --yes
```