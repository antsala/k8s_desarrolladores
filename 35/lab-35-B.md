# Laboratorio 35-B: ***Asegurar el Ingress con TLS***
 
En este laboratorio aprenderemos a usar los objetos Ingress.

Requisitos:

1. Una máquina virtual con Ubuntu 20.04 LTS a la que poder hacer ssh o escritorio remoto.
2. ***TENER UNA SUBSCRIPCIÓN DE AZURE PARA DESPLEGAR AKS***
3. Desplegar el cluster AKS en Azure (Ver lab-00.md)

En el laboratorio anterior usamos ***nginx*** como controlador Ingress. Si bien es una solución interesante para los despliegues on-prem, no suele serla para los despluiegues de Kubernetes en el Cloud.

En la nube, el controlador de Ingress suele estar integrado en la solución de firewall de capa 7 propia del proveedor. En el caso de Azure es el ***Gateway de aplicación***. Al configurarlo como ***Controlador Ingress***, Kubernetes se encargará de traducir todas las reglas presentes en los diferentes objetos ***Ingress*** a las respectivas configuraciones del ***Application Gateway***, es decir, una vez creado servicio en el proveedor, nuestro trabajo consiste en crear los archivos YAML y aplicarlos. Digamos que para el usuario del cluster de Kubernetes es transparente.

Independientemente del ***Controlador Ingress*** configurado, desearemos ofrecer TLS a nuestros servidores. Podemos hacerlo instalando un certificado digital en el servidor apropiado (que corre dentro de un contenedor) o hacer que el ***Controlador Ingress*** sea el finalizador del túnel TLS. De esta forma, un solo certificado digital podrá proteger el tráfico que el controlador esté redirigiendo a los servicios.

Para automatizar la gestión de los certificados, instalaremos un add-on de Kubernetes llamado ***cert-manager***. que es un complemento (addon) que nos ayuda en la creación de los certificados TLS. Se encarga de la rotación de éstos cuando están próximos a expirar. Haremos que ***cert-manager*** se comunique con ***Let's Encrypt*** para solicitar los certificados.


## Ejercicio 1: ***Configuración de un Gateway de aplicación de Azure como Ingress de K8s***

Kubernetes no viene con un controlador ingress por defecto, por lo que hay que configurar alguno de los que sean compatibles. En Azure es ***AGIC (Application Gateway Ingress Controller)***, un balanceador de capa 7. Azure Application Gateway ofrece multitud de características avanzadas, como ***WAF (Web Application Firewall)***.

Hay dos formas de configurar ***AGIC***, o bien usando ***Helm*** (lo veremos más adelante) o, como un complemento al servicio ***AKS***. Esta última tiene la ventaja de que puede ser actualizado automáticamente por Microsoft, asegurando que el entorno está siempre actualizado.

Creamos un nuevo grupo de recursos para el AF en la misma ubicación del grupo de recursos de AKS:
```
az group create \
    --resource-group agic-rg \
    --location northeurope
```

Creamos los componentes de red necesarios para el Application Gateway. Primero un nombre dns para la IP pública. (Nota: debe ser único) (Nota: sustituir YYYYMMDD por la fecha actual)
```
PUBLIC_IP_DNS_NAME=pipnameasg<YYYYMMDD>
```

El registro A de esta IP tomará la forma: ***$PUBLIC_IP_DNS_NAME.northeurope.cloudapp.azure.com***:
```
az network public-ip create \
    --name agic-pip \
    --resource-group agic-rg \
    --allocation-method Static \
    --sku Standard \
    --dns-name $PUBLIC_IP_DNS_NAME
```

Y ahora una red virtual:
```
az network vnet create \
    --name agic-vnet \
    --resource-group agic-rg \
    --address-prefix 192.168.0.0/24 \
    --subnet-name agic-subnet \
    --subnet-prefix 192.168.0.0/24
```

Por último, creamos el Application Gateway. OJO. PUEDE TARDAR SOBRE 6 MINUTOS:
```
az network application-gateway create \
    --name agic \
    --location northeurope \
    --resource-group agic-rg \
    --sku Standard_v2 \
    --public-ip-address agic-pip \
    --vnet-name agic-vnet \
    --subnet agic-subnet \
    --priority "1"  
```

Una vez creado el ***AGIC (Application Gateway Ingress Controller)***, debemos configurarlo para que se integre en el cluster de Kubernetes, por medio del plug-in. También configuraremos el ***Virtual Network Peering*** para que el Application Gateway pueda enviar tráfico al cluster de K8s.

Para habilitar la integración entre el cluster y el Application Gateway, hacemos lo siguiente:
```
appgwId=$(az network application-gateway show \
            --name agic \
            --resource-group agic-rg  \
            --query id \
            --output tsv)
```
```
az aks enable-addons \
    --name myaks \
    --resource-group myaks-rg \
    --addons ingress-appgw \
    --appgw-id $appgwId
```

Ahora hay que conectar la red del Application Gateway con la red del cluster AKS:
```
nodeResourceGroup=$(az aks show \
                        --name myaks \
                        --resource-group myaks-rg \
                        --query "nodeResourceGroup" \
                        --output tsv)
```
```
aksVnetName=$(az network vnet list \
                --resource-group $nodeResourceGroup \
                --query [0].name \
                --output tsv) 
```
```
aksVnetId=$(az network vnet show \
                --name $aksVnetName \
                --resource-group $nodeResourceGroup \
                --query id \
                --output tsv)
```
```
az network vnet peering create \
    --name AppGWtoAKSVnetPeering \
    --resource-group agic-rg \
    --vnet-name agic-vnet \
    --remote-vnet $aksVnetId \
    --allow-vnet-access
```
```
appGWVnetId=$(az network vnet show \
                --name agic-vnet \
                --resource-group agic-rg \
                --query id \
                --output tsv)
```
```
az network vnet peering create \
    --name AKStoAppGWVnetPeering \
    --resource-group $nodeResourceGroup \
    --vnet-name $aksVnetName \
    --remote-vnet $appGWVnetId \
    --allow-vnet-access
```

Con esto se termina la integración entre el Application Gateway y el cluster de AKS.


## Ejercicio 2: ***Añadir una regla de entrada (Ingress) a la aplicación***

Cambiamos al directorio de trabajo.
```
cd ~/k8s_desarrolladores/35
```

Vamos a desplegar la aplicación ***Guestbook*** y exponerla por medio de una ingress. Editemos el archivo para conocer el despliegue. Es la app de ejemplo de Google.
```
nano lab-35-B-guestbook-all-in-one.yaml
```

Salimos sin modificar el archivo y creamos los objetos:
```
kubectl create -f lab-35-B-guestbook-all-in-one.yaml
```

En la carpeta del laboratorio, tenemos el archivo ***lab-35-B-simple-frontend-ingress.yaml***. Lo abrimos:
```
nano lab-35-B-simple-frontend-ingress.yaml
```

Las líneas más importantes a tener en cuenta son:

* *Línea 1*: Versión de la API de Kubernetes para el objeto que se crea.
* *Línea 2*: El objeto es un ***Ingress***.
* *Líneas 5-6*: El ingress es de la clase ***azure/applicacion-gateway***.
* *Líneas 8-12*: Se define el path en el que está escuchando el ingress.
* *Líneas 13-17*: Servicio al que se enviará el tráfico.

Creamos el objeto ingress:
```
kubectl apply -f lab-35-B-simple-frontend-ingress.yaml
```

Comprobamos que se ha creado el objeto:
```
kubectl get ingress
```

La salida monstrará algo así:
```
NAME                      CLASS    HOSTS   ADDRESS        PORTS   AGE
simple-frontend-ingress   <none>   *       20.126.12.28   80      45s
```

Probar la app. El tráfico ya entra por el AGIC. Conectar con un navegador a: ***http://<Poner aquí la IP del ingress>***. Aún no hay SSL, pero eso lo solucionaremos en breve. Es importante recordar que las reglas no se configuran desde Azure, sino desde archivos YAML. El driver AGIC se las pasará a Azure.

Azure crea un registro A para la IP del entrypoint del ingress:
```
URL_APP=http://$PUBLIC_IP_DNS_NAME.northeurope.cloudapp.azure.com
```

Comprobamos.
```
echo $URL_APP
```

Conectar con un navegador.

Comprobar que el servicio de frontend ***NO TIENE IP EXTERNA PÚBLICA***,  es un servicio ***INTERNO***. Esto será así porque es de tipo ClusterIP y no LoadBalancer:
```
kubectl get service
```

El tráfico queda así:

Navegador Usuario --> Ingress (IP Pública) --> Servicio frontednd (IP Privada)


Vamos a dar soporte HTTPS a la aplicación. Para ello necesitaremos un certificado, que lo proporcionará Let's Encrypt,  y el complemento ***cert-manager*** de Kubernetes (que pedirá certificados a Let's Encrypt). Los pasos a realizar son los siguiente:

1. Instalar ***cert-manager***.
2. Instalar el emisor (issuer) de certificados.
3. Crear el certificado SSL para un FQDN concreto.
4. Asegurar el servicio frontend creando un ingress con SSL.


## Ejercicio 3: ***Instalación de cert-manager***

Instalamos ***cert-manager***:
```
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.9.0/cert-manager.yaml
```

***cert-manager*** hace uso de una funcionalidad de Kubernetes llamada ***CustomResourceDefinition (CRD)***, que es usada para extender el API Server de Kubernetes y crear recursos personalizados. Algunos de ellos serán utilizados en breve.


## Ejercicio 4: ***Instalación del emisor de certificador (issuer)***

Hacemos una copia del archivo ***lab-35-B-certificate-issuer.yaml*** para no perder el original.
```
cp lab-35-B-certificate-issuer-initial.yaml lab-35-B-certificate-issuer.yaml
```

Instalación del emisor de certificador (issuer). El archivo ***lab-35-B-certificate-issuer.yaml*** contiene el código para este emisor. Lo abrimos para estudiarlo:
```
nano lab-35-B-certificate-issuer.yaml
```

Las líneas más importantes son:

* *Líneas 1-2*: Un objeto ***issuer*** es un enlace entre el cluster de k8s y la CA que crea el certificado, que en este caso es Let's Encrypt.
* *Líneas 6-10*: Aquí ponemos la configuración de Let's Encrypt y apuntamos al servidor staging. ¡¡¡¡¡MUY IMPORTANTE!!!!!. Poner el email.
* *Líneas 11-14*: Aquí se pone la configuración para que el cliente ACME certifique la propiedad del dominio. Hacemos apuntar a Let's Encrypt a la ingress del Application Gateway para que verifique que somos el propietario del dominio para el que solicitamos un certificado. (Más info: https://letsencrypt.org/es/docs/client-options/)

Guardarmos los cambios en el archivo y creamos el objeto:
```
kubectl create -f lab-35-B-certificate-issuer.yaml
```

Comprobamos que se haya creado el issuer. Observar que el campo READY ponga ***True***:
```
kubectl get issuer
```

Si hubiera problemas, lo describimos:
```
kubectl describe issuer letsencrypt-staging
```


## Ejercicio 5: ***Crear el certificado TLS y asegurar la Ingress***

Hay dos formas de configurar certificados: O bien creamos un certificado manualmente y lo enlazamos con la ingress, o podemos configurar al controlador ingress, de forma que ***cert-manager*** cree el certificado automáticamente. Utilizaremos el segundo método.
```
INGRESS_DNS=$PUBLIC_IP_DNS_NAME.northeurope.cloudapp.azure.com
```

Mostramos y copiamos el valor en el portapapeles:
```
echo $INGRESS_DNS
```

Hacemos una copia del archivo ***lab-35-B-ingress-with-tls-initial.yaml*** para no perder el original.
```
cp lab-35-B-ingress-with-tls-initial.yaml lab-35-B-ingress-with-tls.yaml
```

Abrimos el archivo 'lab-35-B-ingress-with-tls.yaml':
```
nano lab-35-B-ingress-with-tls.yaml 
```

Las líneas más interesantes son:

* *Líneas 7-8*: Se añaden dos anotaciones al ingress que apuntan al emisor de certificado y al ***acme-challenge*** para demostrar la propiedad del dominio.
* *Línea 20*: ¡¡¡¡¡¡¡¡IMPORTANTE!!!!!!Poner aquí el nombre de dominio del ingress $INGRESS_DNS. Esto es obligatorio porque Let's Encrypt solo entrega certificados a los dominios (no IPs)
* *Línea 23*: Esta es la configuración TLS del ingress. También hay que poner el nombre DNS anterior, $INGRESS_DNS.
* *Línea 24*: Este es el nombre del secreto que será creado para almacenar el certificado.

Guardar el archivo y actualizar el objeto ingress con el siguiente comando:
```
kubectl apply -f lab-35-B-ingress-with-tls.yaml
```

***cert-manager*** tardará aproximadamente 1 minuto en solicitar el certificado y configurar la ingress para que lo use. Los objetos intermedios que se crean son:

Un objeto ***certificate***, que podemos ver así (esperar a que ponga ***True***):
```
kubectl get certificate -w
```

Hay otro objeto que ***cert-manager*** creó para obtener el certificado. Es la ***petición***. Podemos verificar su estado así:
```
kubectl get certificaterequest
```

Podemos ver el flujo de la petición así:
```
kubectl describe certificaterequest
```

Comprobar que todo funciona conectando por TLS con un navegador. 
```
echo https://$INGRESS_DNS
```

IMPORTANTE: Dará un error de certificado porque estamos usando el servidor de staging de Let's Encrypt. https://letsencrypt.org/docs/staging-environment/ y así no le sobrecargamos sus sistemas. 


## Ejercicio 6: ***Cambiar al entorno de producción de Let´s Encrypt***

Cambiar del entorno de staging al de producción de Let's Encrypt.

Hacemos una copia del archivo ***lab-35-B-certificate-issuer-prod-initial.yaml*** para no perder el original.
```
cp lab-35-B-certificate-issuer-prod-initial.yaml lab-35-B-certificate-issuer-prod.yaml
```

Vamos a crear un nuevo ***issuer*** en el cluster, pero esta vez que pida certificados de producción. Editamos el archivo ***lab-35-B-certificate-issuer-prod.yaml***.
```
nano lab-35-B-certificate-issuer-prod.yaml
```

IMPORTANTE. En la línea 7, ponemos el email, Guardamos el archivo.

Copiar el valor.
```
echo $INGRESS_DNS
```

Copiamos el archivo ***lab-35-B-ingress-with-tls-prod-initial.yaml*** para no perder el original.
```
cp lab-35-B-ingress-with-tls-prod-initial.yaml lab-35-B-ingress-with-tls-prod.yaml
```

Editamos el archivo ***lab-35-B-ingress-with-tls-prod.yaml***:
```
nano lab-35-B-ingress-with-tls-prod.yaml
```

IMPORTANTE!!!. En la fila 20 y en la 23 ponemos la DNS de nuestra ingress. Guardamos.

Aplicamos los archivos yaml:
```
kubectl create -f lab-35-B-certificate-issuer-prod.yaml
kubectl apply -f lab-35-B-ingress-with-tls-prod.yaml
```

#Hay que esperar un minuto a que se descargue e instale el certificado. Esperar a que READY sea ***true***:
```
kubectl get certificates -w
```

Comprobar, de nuevo,  que todo funciona conectando por TLS con un navegador. 
```
echo https://$INGRESS_DNS
```


Para borrar los recursos que dimos de alta.
```
kubectl delete -f https://github.com/jetstack/cert-manager/releases/download/v1.9.0/cert-manager.yaml
```
```
az aks disable-addons \
    --name myaks \
    --resource-group myaks-rg \
    --addon ingress-appgw
```

Eliminamos las aplicación:
```
kubectl delete -f lab-35-B-guestbook-all-in-one.yaml
```

Comprobamos:
```
kubectl get all 
```

Solo debe quedar el servicio de Kubernetes.
```
NAME                 TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
service/kubernetes   ClusterIP   10.96.0.1    <none>        443/TCP   53m
```

IMPORTANTE: Recordar que el Application Gateway sigue estando en Azure, habría que eliminarlo. Mirar en el documento ***lab-00.md*** como se borra el cluster y sus recursos.


