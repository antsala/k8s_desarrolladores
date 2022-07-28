# Laboratorio 65-B: "Monitorización del cluster"
 
# Este laboratorio aprenderemos a conectar pods de Kubernetes con servicios que corren fuera del
# cluster de Kubernetes.

# Requisitos:
#
#   1) Una máquina virtual con Ubuntu 20.04 LTS a la que poder hacer ssh.
#
#   2) Minikube iniciado.


#################################################
# Ejercicio 1: Sondas 'Readiness' y 'Liveness'. #
#################################################

# Exploraremos cómo Kubernetes se asegura que nuestra aplicación está funcionando correctamente
# usando las sondas 'readiness' y 'liveness'.
#

# Cambiamos al directorio de trabajo.
cd ~/k8s_desarrolladores/65


# Desplegamos la aplicación:

kubectl apply -f lab-65-B-guestbook-all-in-one.yaml


# A continuación comprobamos los pods:

kubectl get pods -w

# De la salida del comando anterior, las columnas y su interpretación es la siguiente:
#
# Columna 'NAME':       Nombre del pod.
#
# Columna 'READY':      Indica cuántos contenedores del pod están listos, comparando contra 
#                       el número total de contenedores del pod. Tiene que ver con las sondas 
#                       'readiness' y 'liveness' que serán tratadas en breve.
#
# Columna 'STATUS':     Indica el estado. Por ejemplo 'ContainerCreating', 'Running', ...
#
# Columna 'RESTARTS':   Número de reinicios.
#
# Columna 'AGE':        Edad del pod desde su creación.


# Podemos mostrar columnas adicionales con:

kubectl get pods -o wide


# El significado de las nuevas columnas es:
#
# Columna 'IP':                 Dirección IP del pod.
#
# Columna 'NODE':               Nodo en el que está corriendo.
#
# Columna 'NOMINATED NODE':     Un nodo nominado solo se pone cuando un pod de mayor prioridad
#                               se anticipa a otro de baja prioridad. El campo 'nominated node'
#                               debería configurarse en el pod de alta prioridad. Señala el nodo
#                               donde se planificará el pod de mayor prioridad una vez que el pod
#                               de baja prioridad haya finalizado limpiamente.
#
# Columna 'READINESS GATES':    Es una forma de presentar componentes del sistema externo al
#                               'readiness' del pod.

# Podemos añadir más información por medio del uso de 'labels'. Se utilizan para enlazar objetos,
# como un servicio a un pod, o un deployment a un ReplicaSet y a un pod. Si las etiquetas no
# coinciden los recursos no se conectarán.



##############################################################
# Ejercicio 2: Depuración de errores en el pull de imágenes. #
##############################################################


# Vamos a elegir un tag de imagen que no existe.

# Editamos el deployment 'frontend':

kubectl edit deployment/frontend

# Cambiamos la etiqueta de la imagen de 'v4' a 'v_no_existente'
# Escribimos '/gb-frontend' y 'Enter' para localizar el texto.
# Ir hasta el final de la línea con la tecla derecha. Pulsar 'i' para entrar en modo de inserción.
# Borramos 'v4' y escribimos 'v_no_existente'
# 'x' borra el caracter en la posición del cursor. 'X' borra el caracter anterior.
# Guardamos con ':wq!'

# Miramos los pods:

kubectl get pods

# Aparecerá el error 'ErrImagePull' o el error 'ImagePullBackOff'. Ambos errores se refieren al 
# hecho de que Kubernetes no puede descargar la imagen desde el registro.
#
# 'ErrImagePull' indica eso exactamente. 'ImagePullBackOff' indica que k8s esperará antes de 
# reintentar la descarga. El intervalo de espera es exponencial, empezando en 10 segundos, 
# luego 20, 40, ..., hasta 5 minutos.

# Más info al describir:

FAILED_POD=<Poner aquí el nombre de un pod que esté dando error>

kubectl describe pod $FAILED_POD


# Editar el deployment y cambiamos de nuevo a 'v4':

kubectl edit deployment/frontend


# Escribimos '/gb-frontend' y 'Enter' para posicionarnos.
# Ir hasta el final de la línea con la tecla derecha. Pulsar 'i' para entrar en modo de inserción.
# Borramos 'v_no_existente' y volvemos a poner 'v4'
# Guardamos con ':wq!'


# Esto deberia resolver el problema:

kubectl get pods


# NOTA: Debido a que Kubernetes hizo una 'rolling update', el frontend estuvo disponible sin downtime.
# K8s reconoció el problema al cambiar la imagen y detuvo el 'rolling update' automáticamente.



##########################################
# Ejercicio 3: Errores de la aplicación. #
##########################################

# Veamos como depurar un error en la aplicación.

# Si se usa Minukube, en otra terminal levantar el tunel.

minikube tunnel

# Obtenermos la IP pública del servicio 'frontend':

kubectl get service

# Almacenamos la IP Externa.
IP_EXTERNA=<Poner aquí la IP Externa del servicio 'frontend'

# La mostramos para copiarla.
echo http://$IP_EXTERNA

# Conectamos con un navegador a 'http://$IP_EXTERNA' y comprobamos que funciona.

# Vamos a escalar el frontend para dejarlo con una sola réplica:

kubectl scale --replicas=1 deployment/frontend


# Comprobamos:

kubectl get pods -w


# Esperar a que solo quede un pod del frontend y capturar su nombre.
FRONTEND_POD_NAME=<Poner aquí el nombre del único pod del frontend.

# Vamos a usar 'kubectl exec' para simular un error de la app:
# (Nota: los siguiente comandos se ejecutan dentro del contenedor)

kubectl exec -it $FRONTEND_POD_NAME -- bash

# Actualizamos el repo e instalamos editor:

apt update -y

apt install -y nano


# Editamos el 'guestkook.php':

nano guestbook.php

# Localizar la línea que tiene 'if ($_GET['cmd'] == 'set') {' en la línea 17.
# debajo de ella insertar el siguiente código:

$host = 'localhost';
if(!defined('STDOUT')) define('STDOUT', fopen('php://stdout', 'w'));
fwrite(STDOUT, "hostname al principio del comando 'set' :"); 
fwrite(STDOUT, $host);
fwrite(STDOUT, "\n");


# Guardamos CTRL+X + Y + Enter.


# Se ha introducido un error donde la lectura de los mensajes seguirá funcionando,
# pero no la escritura. Se ha hecho pidiendo al frontend que se conecte al 
# 'redis master' en un servidor 'localhost' que no es el correcto, así que la
# escritura fallará.
#
# Salimos del contenedor con 'exit'


echo http://$IP_EXTERNA

# Con el navegador conectar a 'http://$IP_EXTERNA'

# Escribir algún mensaje.
# Al refrescar el navegador, el mensaje ya no está, debido a que la escritura falló.

# Obtenermos el log del pod:

kubectl logs $FRONTEND_POD_NAME

# Se puede ver el mensaje de depuración que pusimos:
# 'hostname al principio del comando 'set': localhost',
# que provocará el error de escritura, porque el cliente de redis no puede conectar
# con un servidor de redis llamado 'localhost'.
#
# Por lo tanto, habilitar la depuración el contenedor, ayuda a solucionar problemas.

# La solución del error en este caso es sencilla, puesto que la provocamos al entrar
# en el contenedor y modificar el código. Vamos a eliminar el pod.
#
# Al eliminar el pod, se borra el contenedor, y la capa reescribible de éste, que es
# la que contenía los cambios problemáticos. Como hay un ReplicaSet, se instanciará un
# nuevo pod desde la imagen de contenedor original:

kubectl delete pod $FRONTEND_POD_NAME

# Probar desde el navegador que los mensajes se guardan.

# Limpiamos recursos.

kubectl delete -f lab-65-B-guestbook-all-in-one.yaml



#################################################
# Ejercicio 4: Sondas 'Readiness' y 'liveness'. #
#################################################


# Kubernetes utiliza las sondas para monitorizar la disponibilidad de la 
# aplicación. Estas son:
#
# 'liveness probe':     Monitoriza la disponibilidad de la app mientras está en ejecución.
#                       Si el 'liveness' falla, k8s reiniciará el pod. Suele ser útil para
#                       recuperarse de bucles infinitos o de una app que se quede colgada.
#
# 'readiness probe':    Monitoriza cuando la aplicación no responde. En este caso k8s no
#                       enviará tráfico adicional al pod. Esta sonda es interesante cuando 
#                       la aplicación tiene que inicializarse, o cuando está sufriendo una
#                       sobrecarga puntual de la que se está recuperando.
#
# Las sondas 'liveness' y 'readiness' no necesitan ser servidas desde el mismo endpoint en
# la aplicación.

# Vamos a crear dos despliegues de nginx. Cada uno con una página index y una página health.
# La página de índice servirá como sonda 'liveness'.

# Editamos el erchivo 'index1.html':

code lab-65-B-index1.html


# Editamos el erchivo 'index2.html':

code lab-65-B-index2.html


# Editamos el archivo 'healthy.html':

code lab-65-B-healthy.html


# Vamos a montar estos archivos en el despliegue. Usaremos 'ConfigMap' para cada uno, lo que
# nos permitirá conectarlos a los pods:

kubectl create configmap server1 --from-file=lab-65-B-index1.html

kubectl create configmap server2 --from-file=lab-65-B-index2.html

kubectl create configmap healthy --from-file=lab-65-B-healthy.html


# A modo de recordatorio, para el primer configmap tenemos:

kubectl describe configmap/server1


# Ahora vamos a crear dos despliegues web, muy parecidos uno al otro, cambiando solo el ConfigMap.

# Editamos el archivo 'lab-65-B-webdeploy1.yaml'

code lab-65-B-webdeploy1.yaml


# Las líneas más destacables son:
#
# Líneas 23-28:     Es la sonda 'liveness'. Apunta a la página 'healthy.html'. Recordemos 
#                   que si la página 'healthy.html falla, el pod (contenedor) se reiniciará.
#
# Líneas 29-32:     Es la sonda 'readiness'. Apunta a 'index.html'. Si esta página falla, el 
#                   pod no recibirá temporalmente ningún tráfico, pero permanecerá en 
#                   ejecución.
#
# Líneas 44-45:     Al lanzar nginx, se copian las páginas a su ubicación correcta. Luego 
#                   inicia nginx.

# Vamos a lanzar estos dos despliegues:

kubectl apply -f lab-65-B-webdeploy1.yaml

kubectl apply -f lab-65-B-webdeploy2.yaml


# Comprobamos:

kubectl get deployments -w

# Por último, un creamos un servicio que enrute el tráfico a los dos deployments.
# Comprobamos el manifiesto:

code lab-65-B-webservice.yaml


# Lo creamos:

kubectl apply -f lab-65-B-webservice.yaml


# Comprobamos:

kubectl get svc -w



###########################################################
# Ejercicio 5: Experimentos con 'liveness' y 'readiness'. #
###########################################################

# Obtenemos la IP Pública del servicio:

kubectl get service web

# Capturamos la IP Externa del servicio 'web'
IP_EXTERNA=<Poner aquí la IP Externa del servicio 'web'>

# Mostramos para copiar.
echo http://$IP_EXTERNA

# Conectamos con el navegador a 'http://$IP_EXTERNA'
# Con CTRL+F5 se puede ver el balanceo.

# Vamos a usar un pequeño script ('testWeb.sh') que lo que hace es conectar 
# 50 veces con el servicio, y así poder ver el balanceo más comodamente.
# Lo hacemos ejecutable:

chmod +x lab-65-B-testWeb.sh

# Lo lanzamos contra la IP Externa:

./lab-65-B-testWeb.sh $IP_EXTERNA


# Empecemos con un fallo en la sonda 'readiness'. Esto provocará que se pare
# de forma temporal el tráfico que se envía al contenedor. Usaremos 'exec' para
# cambiar de ubicación el archivo 'index.html' y así fallará la sonda:


kubectl get pods 

# Capturamos el nombre del pod 'server1':

SERVER1_POD_NAME=<Poner aquí el nombre del POD 'server1-XXXXXXXXXX-XXXXX'>


# Cambiamos el nombre del archivo 'index.html' del por 'server1.html':

kubectl exec $SERVER1_POD_NAME -- mv /usr/share/nginx/html/index.html /usr/share/nginx/html/server1.html

# Vemos el cambio de estado del pod. El estado de readiness (READY) del pod del server 1
# cambia a 0/1:

kubectl get pods -w


# Esto provoca que no se envíe más trafico al pod del server1. Lo comprobamos:

./lab-65-B-testWeb.sh $IP_EXTERNA


# Restauramos el estado del server1 volviendo a renombrar el archivo.
kubectl exec $SERVER1_POD_NAME -- mv /usr/share/nginx/html/server1.html /usr/share/nginx/html/index.html


# Vemos el cambio de estado del pod. El estado de readiness (READY) del pod del server1
# debe ser de nuevo 1/1:

kubectl get pods -w


# Y el balanceo debe volver a producirse:

./lab-65-B-testWeb.sh $IP_EXTERNA


# Ahora volvemos a hacer el experimento pero con la sonda 'liveness'. Si falla, k8s reiniciará el pod.

# Capturamos el nombre del pod 'server2':

SERVER2_POD_NAME=<Poner aquí el nombre del POD 'server2-XXXXXXXXXX-XXXXX'


# Borramos el archivo 'healthy.html' del 'server2':

kubectl exec $SERVER2_POD_NAME -- rm /usr/share/nginx/html/healthy.html


# Miramos. Habrá que esperar un tiempo para ver que se reinicia:

kubectl get pods -w


# Podemos verlo con más detalles así:

kubectl describe pod $SERVER2_POD_NAME


# Limpiamos recursos:

kubectl delete deployment server1 server2

kubectl delete service web



##################################
# Ejercicio 6: Métricas simples. #
##################################


# Para listar los nodos del cluster:

kubectl get nodes

# Capturar el nombre de un nodo, si se trabaja con Minikube solo hay uno:

NODE_NAME=<Copiar aquí el nombre de uno de los nodos>


# Más información:

kubectl get -o wide nodes


# Consumo en los nodos:

kubectl top nodes


# Ver eventos en un nodo:

kubectl describe node $NODE_NAME


# Ver consumos del pod.
# Pods corriendo en el espacio de nombres 'kube-system':

kubectl get pods -n kube-system


# Capturar el nombre de un pod de 'coredns':

COREDNS_POD_NAME=<Poner aquí el nombre de un pod de 'coredns'>


# Veamos los 'requests' y 'limits' de un pod concreto:

kubectl describe pod $COREDNS_POD_NAME -n kube-system


# Si el pod pide más memoria que su límite, k8s reinicia el pod.
# Consumo actual de los pods:

kubectl top pods --all-namespaces


# Consumo actual de los pods del espacio de nombres 'default':

kubectl top pods 


#######################
# FIN DEL LABORATORIO #
#######################
