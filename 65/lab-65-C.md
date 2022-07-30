# Laboratorio 65-C: ***RBAC en AKS (AZURE)***
 
En este laboratorio aprenderemos a integrar RBAC de Kubernetes en AKS.

Requisitos:

1. Una máquina virtual con Ubuntu 20.04 LTS a la que poder hacer ssh o escritorio remoto.
2. Cluster AKS. (Ver lab-00.md)

## Ejercicio 1. ***Introduccion a RBAC***

Hasta el momento, hemos tenido permisos para crear, leer, actualizar y eliminar objetos en el cluster. Esto funciona bien en un entorno de prueba, pero no es recomendable para uno de producción.

En los clusters de producción, la recomendación es aprovechar RBAC y conceder un conjunto limitado de permisos a los usuarios. Será necesario configurar RBAC en Kubernetes e intregrarlo con ***Azure AD***.

RBAC tiene 3 conceptos importantes:

* *Role*: Contiene un conjunto de permisos. Por defecto, el rol no tiene ningún permiso y en consecuencia hay que especificarlos. Los permisos son del tipo ***get***, ***watch***, ***list***... Se les llama ***verbos***. Aquí los verbos admitidos por el API Server: https://kubernetes.io/docs/reference/access-authn-authz/authorization/#determine-the-request-verb. El rol también contiene los recursos a los que se aplican esos permisos. Los recursos pueden ser todos los ***pods***, los ***deployments***, etc, o pueden ser un objeto concreto, como ***pod/mypod***.

* *Subject*: Se refiere a la persona o a la cuenta de servicio a la que se asigna el rol. En los clusteres de AKS integrados con AAD, el subject puede ser un ***usuario*** o un ***grupo*** de AAD.

* *RoleBinding*: Sirve para enlazar un subject a un rol para un contexto. Si es ***CusterRoleBinding***, se refiere a la totalidad del cluster.

Un concepto importante a comprender es que hay dos capas de RBAC: Los RBAC de Azure y los RBAC de Kubernetes.


Los RBAC de Azure tienen que ver con los roles asignados a las personas para hacer ***cambios en Azure***, como crear, modificar o borrar clústeres. Los RBAC de Kubernetes tienen que ver con el ***derecho de acceso*** a los recursos del cluster.

Los RBACs de Kubernetes son una característica OPCIONAL. Por defecto los clusters que se crean tienen RBAC habilitado, sin embargo no están integrados con Azure AD. Esto significa que no se pueden dar permisos de Kubernetes a usuarios de AAD, y habría que integrarlo.


## Ejercicio 2. ***Habilitar la integración de Azure AD en AKS***

Una vez que el cluster ha sido integrado con Azure AD, esta funcionalidad no puede ser deshabilitada Empezamos creando un grupo en Azure AD al que le asignaremos permisos en AKS:
```
AKS_ADMIN_GROUP_ID=$(az ad group create \
    --display-name "aks admins" \
    --mail-nickname aksadmins \
    --description "Administradores de clusteres AKS" \
    --query id \
    --output tsv)
```

Mostramos el ID del grupo.
```
echo $AKS_ADMIN_GROUP_ID
```

Actualizamos la integración de Azure AD para el cluster:
```
az aks update \
    --resource-group myaks-rg \
    --name myaks \
    --enable-aad \
    --aad-admin-group-object-ids $AKS_ADMIN_GROUP_ID 
```

Esta acción podemos verla en la GUI en: ***Home/Kubernetes Services/myaks/Cluster configuration/Kubernetes authentication and authotization***.

## Ejercicio 3. ***Añadir al usuario administrador del tenant al grupo 'aks admins'***


Al habilitar la integración con Azure AD, es necesario poner al administrador del tenant en el grupo de administradores del cluster, de lo contrario no podrá administrarlo por la ***GUI*** ni por la ***CLI***.

Tomamos el ***ID del administrador del tenant***: (Nota: En el UPN poner el dominio verificado apropiado)
```
TENANT_ADMIN_USER=antsalgra_hotmail.com#EXT#@antsalgrahotmail.onmicrosoft.com 
ADMIN_USER_ID=$(az ad user show \
                    --id $TENANT_ADMIN_USER\
                    --query id \
                    --output tsv)
```

Mostramos el ID.
```
echo $ADMIN_USER_ID
```

Agregamos al admin al grupo de administradores del cluster.
```
az ad group member add \
    --group $AKS_ADMIN_GROUP_ID \
    --member-id $ADMIN_USER_ID
```

## Ejercicio 4. ***Crear un usuario y un grupo de seguridad para asignar roles***


Creamos un usuario del cluster: (Nota:) En el UPN poner el dominio verificado apropiado.
```
LUKE_USER_ID=$(az ad user create \
                --display-name "Luke Skywalker" \
                --password useTheForce# \
                --user-principal-name luke@antsalgrahotmail.onmicrosoft.com \
                --mail-nickname luke \
                --query id \
                --output tsv)
```

Mostramos el ID del usuario ***Luke***.
```
echo $LUKE_USER_ID
```

Creamos otro grupo de seguridad, llamado ***aks users***, donde pondremos a los usuarios del cluster sin rol administrativo:
```
AKS_USERS_GROUP_ID=$(az ad group create \
                        --display-name "aks users" \
                        --mail-nickname aksusers \
                        --description "usuarios de clusteres AKS" \
                        --query id \
                        --output tsv)
```

Mostramos el ID del grupo ***aks users***
```
echo $AKS_USERS_GROUP_ID
```

Agregamos a ***Luke*** al grupo de usuarios del cluster:
```
az ad group member add \
    --group $AKS_USERS_GROUP_ID \
    --member-id $LUKE_USER_ID
```

Ahora necesitamos hacer que Luke sea un ***usuario de cluster*** en el ***RBAC de AKS***. Esto lo habilitará para usar la Azure CLI y conseguir acceso al cluster. Lo primero es tomar el identificador de recurso del cluster AKS en Azure:
```
AKS_ID=$(az aks show \
            --resource-group myaks-rg \
            --name myaks \
            --query id \
            --output tsv)
```

Es algo así: ***/subscriptions/5d72e184-55f6-4093-838e-3d0f7506881a/resourcegroups/myaks-rg/providers/Microsoft.ContainerService/managedClusters/myaks***

Lo consultamos:
```
echo $AKS_ID
```

Creamos una asignación de rol de Azure para el grupo ***aks users***, del que es miembro ***Luke***. Si diera error indicando que el principal de seguridad no existe, esperar unos segundos porque el grupo aún no se ha creado:
```
sleep 60
az role assignment create \
    --assignee $AKS_USERS_GROUP_ID \
    --role "Azure Kubernetes Service Cluster User Role" \
    --scope $AKS_ID
```

El resultado de esta acción se puede ver en la GUI en: ***Home/Kubernetes services/myaks/Access 
control (IAM)/Role Assignments***

El rol ***Azure Kubernetes Service Cluster Role***, tiene como descripción ***List cluster user credential action***, que permite tomar las credenciales de ese usuario en el cluster y almacenarlas en ***.kube/config***, para que posteriormente ***kubectl*** pueda usarlas.

NOTA: Si se hubiera usado la ***CloudShell***, también habría que dar permisos a ***AKS_USERS_GROUP_ID*** para la cuenta de almacenamiento donde reside la CloudShell. En este ejemplo no la usamos.


## Ejercicio 5. ***Configurar RBAC en AKS***

Para hacer la demo, crearemos dos ***namespaces*** y desplegaremos la aplicación de voto de Azure en cada espacio de nombres. Asignaremos al grupo que creamos acceso de solo lectura de ámbito de cluster a los pods. Al usuario de asignaremos la capacidad de eliminar pods solo en uno de los espacios de nombres.

Crearemos los siguientes objetos en k8s.

* *Un ClusterRole*: Para dar el acceso de solo lectura a todos los pods del cluster 
* *ClusterRoleBinding*: Para asignar al grupo ***aks users*** el rol anterior de solo lectura.
* *Otro ClusterRole*: Para dar permisos de eliminación en el espacio de nombres ***delete-access***.
* *Otro ClusterRoleBinding*: Para asignar al usuario ***Luke*** el rol de eliminación anterior.

Creamos los dos espacios de nombres. ***no-access*** y ***delete-access***. La idea es que el usuario que vamos a crear pueda borrar pods en ***delete-access*** y no pueda hacerlo en ***no-access***:
```
kubectl create ns no-access
kubectl create ns delete-access
```

Cambiamos al directorio de trabajo:
```
cd ~/k8s_desarrolladores/65
```

Desplegamos la aplicación de voto en los espacios de nombres:
```
kubectl create -f lab-65-C-azure-vote.yaml --namespace no-access
kubectl create -f lab-65-C-azure-vote.yaml --namespace delete-access
```

Comprobamos:
```
kubectl get all --namespace no-access 
kubectl get all --namespace delete-access 
```

Ahora vamos a crear el objeto ***ClusterRole***, que asignará permisos de solo lectura en todo el cluster. Editamos el archivo ***lab-65-C-clusterRole.yaml***:
```
code lab-65-C-clusterRole.yaml
```

* *Línea 2*:  Define la creación de una instancia ***ClusterRole***.
* *Línea 4*:  A la que le asigna el nombre ***readOnly***.
* *Línea 6*:  Concede acceso a todos los grupos de la API. Ver https://kubernetes.io/docs/reference/using-api/#api-groups
* *Línea 7*:  Concede acceso a todos los pods.
* *Línea 8*:  Concede acceso a las acciontes ***get***, ***watch*** y ***list***. Ver https://kubernetes.io/docs/reference/access-authn-authz/authorization/#determine-the-request-verb

Creamos el objeto ***ClusterRole***:
```
kubectl create -f lab-65-C-clusterRole.yaml
```

Comprobamos:
```
kubectl get clusterRole
```

Podemos inspeccionalo:
```
kubectl describe clusterRole readOnly
```

Mostramos para copiar.
```
echo $AKS_USERS_GROUP_ID
```

Ahora vamos a crear un objeto ***ClusterRoleBinding*** que enlaza el rol a un usuario o grupo. Editamos el archivo ***lab-65-C-clusterRoleBinding.yaml***:
```
code lab-65-C-clusterRoleBinding.yaml
```

* *Línea 2*: Define que estamos creando una instancia ***ClusterRoleBinding***.
* *Línea 4*: Le asigna el nombre ***readOnlyBinding***.
* *Líneas 5-8*: Hace referencia al objeto ***ClusterRole*** que creamos anteriormente.
* *Líneas 9-12*: Se refiere al grupo de AAD que creamos antes (***aks users***). IMPORTANTE!!!! Sustituir en la línea 12 el id que hemos visualizado para ***AKS_USERS***.

Guardar los cambios y salir.

Creamos el objeto ***ClusterRoleBinding***:
```
kubectl create -f lab-65-C-clusterRoleBinding.yaml
```

Comprobamos:
```
kubectl get clusterRoleBinding
```

Lo inspeccionamos:
```
kubectl describe clusterRoleBinding readOnlyBinding
```

A continuación crearemos un rol que permite la eliminación en el espacio de nombres ***delete-access***. Editamos el archivo ***lab-65-C-role.yaml***.
```
code lab-65-C-role.yaml
```

* *Línea 2*: Se indica que se está creando una instancia de ***Role*** y no de ***ClusterRole***. La instancia de ***Role*** no se aplica a todo el cluster.
* *Línea 5*: Aquí ponemos el espacio de nombres al que se aplica.
* *Líneas 7-9*: Los tipos de recursos afectados y los verbos permitidos.

Creamos el rol:
```
kubectl create -f lab-65-C-role.yaml
```

Comprobamos (ojo con el espacio de nombres):
```
kubectl get role --namespace delete-access
```

Inspeccionamos:
```
kubectl describe role deleteRole  --namespace delete-access
```

Por último creamos una instancia de ***RoleBinding*** para asignar el role al usuario ***Luke***. Editamos el archivo ***lab-65-C-roleBinding.yaml***:
```
code lab-65-C-roleBinding.yaml
```

* *Línea 2*: Crea una instancia de un ***RoleBinding*** y no de un ***ClusterRoleBinding*** porque lo que se está asociando es un ***Role*** y no un ***ClusterRole***.
* *Línea 5*: Indica el espacio de nombres en el que se crea este rol.
* *Línea 7*: Hacer referencia a una instancia ***Role*** y no a una ***ClusterRole***.
* *Líneas 11-13*: Define un usuario en lugar de un grupo. IMPORTANTE!!!! Poner en la línea 13 el usuario que hemos creado ***luke@antsalgrahotmail.onmicrosoft.com***.

Creamos la instancia:
```
kubectl create -f lab-65-C-roleBinding.yaml
```

Comprobamos:
```
kubectl get roleBinding --namespace delete-access
```

Inspeccionamos:
```
kubectl describe roleBinding deleteBinding --namespace delete-access
```

## Ejercicio 5. ***Verificar RBAC para el usuario Luke***

Cerramos la sesión de Azure del usuario actual.
```
az logout
```

Borramos la caché de cuentas de Azure:
```
az account clear
```

IMPORTANTE: Para que Luke pueda interactuar con la subscripción hay que asignarle un rol en la misma. Lo hacemos con la interfaz web de azure: ***Home/Subscriptions/MSDN Platforms/Access Control (IAM)***

Asignar a ***Luke*** el rol de ***contibuidor*** a la subscripción.

Iniciamos sesión con el usuario ***Luke@antsalgrahotmail.onmicrosoft.com*** y password ***useTheForce#***
```
az login 
```


¡¡¡¡MUY IMPORTANTE!!!!

***az aks get-credentials*** toma las credenciales de un cluster administrado de AKS y las almacena en el archivo ***~/.kube/config*** de forma que ***kubectl*** pueda usarlas. Hasta el momento, hemos usado la CLI como ***administrador***, usando el parámetro ***--admin*** del comando anterior. Este parámetro se utiliza para:

"***--admin -a***: Get cluster administrator credentials.  Default: cluster user credentials. On clusters with Azure Active Directory integration, this bypasses normal Azure AD authentication and can be used if you're permanently blocked by not having access to a valid Azure AD group with access to your cluster. Requires ***Azure Kubernetes Service Cluster Admin*** role."

Es decir, que en lugar de poner a un usuario en un grupo asignado al rol ***Azure Kubernetes Service Cluster Admin*** hemos utilizado ***--admin*** para saltarnos el RBAC de Azure.

Como ahora mismo es esto lo que probamos, descargamos las credenciales para el usuario ***Luke*** SIN USAR ***--admin***.

Otro parámetro interesante es ***--overwrite-existing*** cuya finalidad es:

"***--overwrite-existing***: Overwrite any existing cluster entry with the same name", que debemos utilizarlo si le cambiamos los permisos al usuario en Azure.

En definitiva, para descargar las credenciales de ***Luke*** debemos poner:
```
az aks get-credentials \
    --resource-group myaks-rg \
    --name myaks \
    --overwrite-existing
```

Podemos ver que estamos logados a la subscripción de Azure con el usuario ***Luke*** con el comando:
```
az account show
```

Vamos a verificar si el usuario tiene permiso para ver los pods en todos los espacios de nombres. Debe ver los pods en los dos espacios:
```
kubectl get pods --namespace no-access
kubectl get pods --namespace delete-access
```

Esto es debido al objeto ***ClusterRole*** asignado al grupo (aks users). De hecho, al aplicarse al cluster tendría acceso a todos los espacios de nombres, como puede comprobarse con el siguiente comando:
```
kubectl get pods --all-namespaces
```

Ahora comprobamos los permisos de eliminación. Solo debe poder eliminar del espacio de nombres ***delete-access***:
```
kubectl delete pod --all --namespace delete-access
kubectl delete pod --all --namespace no-access
```

Para limpiar, INICIAMOS SESIÓN CON EL USUARIO ADMINISTRADOR:
```
az logout
```

Borramos la caché de cuentas de Azure:
```
az account clear
```

Iniciamos sesión con el usuario ***Luke@antsalgrahotmail.onmicrosoft.com*** y password ***useTheForce#***
```
az login
```

Actualmente el cluster tiene configurado RBAC y la autenticación de AAD. Actualizamos credenciales:
```
az aks get-credentials \
    --resource-group myaks-rg \
    --name myaks \
    --overwrite-existing
```

Alternativamente podemos saltarnos la autenticación de AAD y tener credenciales de administrador en el cluster así:
```
az aks get-credentials \
    --resource-group myaks-rg \
    --name myaks \
    --overwrite-existing \
    --admin
```

Limpiamos
```
kubectl delete -f lab-65-C-azure-vote.yaml -n no-access
kubectl delete -f lab-65-C-azure-vote.yaml -n delete-access
kubectl delete -f lab-65-C-clusterRoleBinding.yaml 
kubectl delete -f lab-65-C-clusterRole.yaml 
kubectl delete -f lab-65-C-roleBinding.yaml 
kubectl delete -f lab-65-C-role.yaml 
kubectl delete ns no-access
kubectl delete ns delete-access
```
