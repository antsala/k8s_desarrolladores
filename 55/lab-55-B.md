# Laboratorio 55-B: ***MySQL replicado con StatefulSet***
 
En este laboratorio aprenderemos a usar el objeto ***StatefulSet***.

Vamos desplegar una topología replicada de MySQL (Basado en ejemplo de la web Kunernetes.io)

Requisitos:

1. Una máquina virtual con Ubuntu 20.04 LTS a la que poder hacer ssh o escritorio remoto.
2. Cluster ***Minikube*** iniciado.


En este laboratorio desplegaremos ***1 ConfigMap***, ***2 servicios*** ***1 StatefulSet***.

## Ejercicio 1: ***Crear el ConfigMap***


Cambiamos al directorio de trabajo:
```
cd ~/k8s_desarrolladores/55
```

Editamos el YAML del ConfigMap para estudiarlo.
```
code lab-55-B-MySQL-ConfigMap.yaml 
```

Las líneas más relevantes son:

* *Línea 2*: Es un ***ConfigMap***.
* *Líneas 7-15*: Define dos claves: ***primary.cnf*** y ***replica.cnf*** que serán tenidas en cuenta en función del rol que asignemos a cada pod de MySQL. Se aplican de la siguiente forma: ***primary.cnf***: Se aplicará al servidor MySQL que actúe como master o primario. En este ejemplo el servidor primario enviará logs de replicación a los servidores que actúan como réplicas. ***replica.cnf***: Las réplicas de MySQL leerán esta clave que significa que rechazarán todas las operaciones de escritura que no provengan desde logs de replicación del primario.

Creamos el ConfigMap.
```
kubectl apply -f lab-55-B-MySQL-ConfigMap.yaml
```

Comprobamos que se ha creado correctamente:
```
kubectl get configmap mysql 
kubectl describe configmap mysql 
```

## Ejercicio 2: ***Creación de los servicios***

Procedemos a crear los servicios de la app. El primer servicio servirá para enviar tráfico al pod maestro de MySQL. Si bien podemos usar un servicio de la forma que hemos aprendido hasta el momento, es preciso conocer formas más eficientes.

Para hacer operaciones de escritura en la base de datos, necesitamos enviar el tráfico al pod primario, es decir, a un ***único pod***. Por lo tanto, podríamos decir que el pod de frontend debería conectar directamente con el pod primario de MySQL. Esto contradice las buenas prácticas que hemos visto, y que 'obligan' a que el tráfico pase por un servicio.

Para solventar esta necesidad, Kubernetes permite crear un tipo de servicio especial llamado ***Headless***.

Un servicio ***Headless*** no asigna dirección IP ni reenvía tráfico, ni balancea. Lo que realmente ofrece un servicio Headless, y por esto es buena práctica su uso, es un ***nombre DNS*** que podrá usar el Frontend para alcanzar la instancia primaria (pod) de una base de datos.

En Kubernetes decimos que queremos un servicio ***Headless*** de estas características cambiando el tipo del  servicio de ***ClusterIP*** a ***None*** en el archivo de manifiesto.

El archivo ***lab-55-B-MySQL-services.yaml*** define los dos servicios que necesita esta app.
```
code lab-55-B-MySQL-services.yaml
```

Las líneas más interesantes son:

* *Líneas: 1-14*: Declaración del servicio ***Headless*** para alcanzar el servidor primario de MySQL.
* *Línea 5*: Este servicio se llamará ***mysql***.
* *Línea 11*: Se usará el puerto ***3306*** en el servicio.
* *Línea 12*: Al poner ***ClusterIP: None*** estamos definiendo un servicio ***Headless***.
* *Línea 16*: En YAML, los caracteres ***---*** indican que se procede a definir un nuevo objeto.
* *Líneas 13-14*: Este segundo servicio de asociará con pods que tenga definida la etiqueta ***app: mysql***.
* *Líneas: 18-30*: Definición del sergundo servicio que será usado para balancear las operaciones de lectura a los pods de la base de datos.
* *Línea 22*: El servicio se llamará ***mysql-read***.

Creamos los servicios:
```
kubectl apply -f lab-55-B-MySQL-services.yaml
```

Comprobamos:
```
kubectl get services
```

La salida será similar a esta: (Nota: Observar como el servicio Headless no tiene ***CLUSTER-IP***)
```
NAME         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
kubernetes   ClusterIP   10.96.0.1       <none>        443/TCP    5d15h
mysql        ClusterIP   None            <none>        3306/TCP   26s
mysql-read   ClusterIP   10.99.195.131   <none>        3306/TCP   26s
```

Como hemos comentado, un servicio ***Headless*** proporciona un sufijo de DNS común para los pods que crea el ***StatefulSet***. Puesto que al servicio ***Headless*** le hemos asignado el nombre ***mysql***, los pods pueden ser resueltos en la forma de ***<nombre_del_pod>.mysql*** DESDE CUALQUIER OTRO POD en el que esté en el mismo espacio de nombres.

En breve crearemos el ***StatefulSet***, que tendrá un nombre, por ejemplo ***mysql***. A diferencia del ***Deployment*** que genera identificadores aleatorios para el nombre del pod, el ***StatefulSet***
los ***NUMERA***. En consecuencia los nombres de los pods serán: ***mysql-0***, ***mysql-1***, ***mysql-2***, etc...

Pues ahora lo ponemos todo junto. Si el primer pod del StatefulSet es el primario de MySQL, su nombre DNS será ***mysql-0.mysql***, donde ***.mysql*** es el sufijo DNS que aporta el servicio HeadLess.

En consecuencia el Frontend deberá conectar con ***mysql-0.mysql*** (o ***mysql-0.mysql.default.svc.cluster.local*** si se prefiere) para hacer operaciones de escritura. Realmente estamos contactando con un pod, pero usando la DNS que proporciona el servicio Headless.

Para las consultas de lectura, los pods de Frontend contactan contra el servicio ***mysql-read***.


## Ejercicio 3: ***Creación del StatefulSet***

Procedemos a crear el ***StatefullSet***.

Editamos el archivo ***lab-55-B-statefulset.yaml***.
```
code lab-55-B-statefulset.yaml
```

Debido a la importancia de los contenidos que vamos a aprender, comentamos las líneas más importantes:

* *Línea 2*: Definimos un objeto ***StatefulSet***.
* *Líneas 6-8*: Este ***StatefulSet*** gestionará las réplicas del pod que tiene declarada la etiqueta ***app: mysql***.
* *Línea 9*: ***Nombre del servicio*** (mysql) con el que está asociado este StatefulSet.
* *Línea 10*: Deseamos 3 réplicas (1 primaria y 2 en RO)
* *Línea 11-Final*: Blueprint. Definición de los contenedores del pod.
* *Línea 16*: ***initContainers***. Kubernetes se asegura de iniciar los contenedores de esta sección antes de proceder con el inicio del resto de contenedores declarados en el pod.
* *Línea 17-63*: Se define el contenedor de inicio.
* *Línea 17*:  Se llama a ***init-mysql***. Su misión solo es la de copiar el archivo de configuración de MySQL en función de si el pod va a tener el rod de primario o réplica.
* *Líneas 19-35*: ***CMD*** Comando que se ejecuta al iniciar el contenedor. MUY IMPORTANTE!!! 
Observar cómo se copia al volumen el archivo de configuración del ConfigMap en función del rol de la instancia. La primera (id 0) será la maestra.
* *Línea 65*: Su nombre es ***mysql***.
* *Línea 72*: El contenedor se expone en el puerto ***3306***.        
* *Líneas 73-78*: Volúmenes que monta el contenedor.
* *Líneas 74-76*: Monta el archivo ***mysql*** en ***/var/lib/mysql***. Ver líneas 159-166.
* *Líneas 37 y 38*: Monta un volumen local de tipo ***EmptyDir*** en ***/mnt/conf.d***. Ver líneas 154 y 155.
* *Líneas 39 y 40*: Monta el volumen ***/mnt/config-map*** con el ConfigMap llamado ***config-map***.
* *Líneas 41-63*: Definición de otro contenedor de inicio.
* *Línea 41*: Su nombre es ***clone-mysql***. 
* *Líneas 43-57*: Determina si el pod se corresponde con la instancia primaria o con una réplica. Si es una réplica y el volumen ya contiene el archivo de la base de datos, entonces termina. Si el pod es el master, no hace nada y termina. En este último caso, el pod es una réplica y no tiene el archivo de la base de datos, así que se lo copia desde el pod anterior. Por último se asegura de que el archivo de la base de datos sea consistente para que sea restaurado cuando se levante el pod de la instancia.
* *Líneas 58-62*: Monta los volúmenes ***data*** en ***/var/lib/mysql*** y el ConfigMap en ***/etc/mysql/conf.d***. 

En este punto se ha terminado con la definición de los contenedores de incio. A continuación  tenemos los contenedores que quedarán funcionando en el pod.

* *Líneas 65-95*: Definición del contenedor llamado ***mysql***.
* *Líneas 73-78*: Monta el volumen ***data*** en ***/var/lib/mysql*** (el archivo de de la base de datos) y el volumen ***conf*** en ***/etc/mysql/conf.d*** que contiene el ConfigMap.
* *Líneas 79-82*: 500 milis de CPU y 1 Gi de RAM para el contenedor.
* *Líneas 83-9*5: Sondas ***Liveness*** y ***Readiness*** que serán explicadas más adelante.
* *Líneas 96-152*: Definición del contenedor sidecar o helper ***xtrabackup***. Su misión es iniciar/continuar la replicación desde el master (si el pod es el de una réplica). Luego monta un servidor de backup y se queda a la espera de que otro contenedor (de inicio) le pida el backup.

Aplicamos el ***StatefulSet***.
```
kubectl apply -f lab-55-B-statefulset.yaml 
```

Comprobamos:
```
kubectl get all
```

## Ejercicio 4: ***Enviar tráfico desde el Frontend***

Vamos a probar la conexión desde el Frontend contra el servidor primario de MySQL (***pod/mysql-0***). Lo haremos usando un contenedor temporal para realizar las pruebas:# (Nota: Se crea una tabla y se almacena un registro con el mensaje ***Hola***)
```
kubectl run mysql-client --image=mysql:5.7 -i --rm --restart=Never --\
  mysql -h mysql-0.mysql <<EOF
CREATE DATABASE test;
CREATE TABLE test.messages (message VARCHAR(250));
INSERT INTO test.messages VALUES ('Hola');
EOF
```

Ahora comprobamos con otro contenedor de Frontend operaciones de lectura. (Nota: Observar como el Frontend se contecta contra el servicio del StatefullSet (***mysql-read***))
```
kubectl run mysql-client --image=mysql:5.7 -i -t --rm --restart=Never --\
  mysql -h mysql-read -e "SELECT * FROM test.messages"
```

Probamos en balanceo del servicio ***mysql-read***: (Nota: Se crea un bucle infinito tomando info del servidor que atiende la consulta)
```
kubectl run mysql-client-loop --image=mysql:5.7 -i -t --rm --restart=Never --\
  bash -ic "while sleep 1; do mysql -h mysql-read -e 'SELECT @@server_id,NOW()'; done"
```

Vamos a probar la HA frente a la caída de un pod de réplica de MySQL.

La sonda ***readiness*** del contenedor ***mysql*** ejecuta el comando ***mysql -h 127.0.0.1 -e SELECT 1*** para comprobar que el servidor está funcionando. Una forma de hacer que la sonda falle es la siguiente: (Nota: Vamos a renombrar el ejecutable del comando ***mysql*** en la instancia ***mysql-2*** del StatefulSet)

En otra terminal:
```
kubectl exec mysql-2 -c mysql -- mv /usr/bin/mysql /usr/bin/mysql.off
```

El la terminal principal solo responderán las instancias ***100*** y ***101***, que se corresponden con ***mysql-0*** (master) y ***mysql-1*** (primera réplica). El servicio ***mysql-read*** no envía tráfico al pod ***mysql-2*** (segunda instancia de réplica) porque la sonda no responde.

Esto se puede comprobar con
```
kubectl get pod mysql-2
```

La salida será como esta:
```
NAME      READY   STATUS    RESTARTS   AGE
mysql-2   1/2     Running   0          20m
```

Se puede apreciar como ***1/2*** indica que e1 contenedor no está respondiendo a la sonda ***Readiness***.

Volvemos a renombrar el comando.
```
kubectl exec mysql-2 -c mysql -- mv /usr/bin/mysql.off /usr/bin/mysql
```

El la terminal principal, observar como se empieza a enviar tráfico al pod ***mysql-2*** (segunda réplica)

Ahora probamos la HA eliminando un pod: (Nota: El StatefulSet detectará que falta un pod y volverá a crear uno nuevo)

En la segunda terminal:
```
kubectl delete pod mysql-2
kubectl get statefulset mysql
```

Probamos las aplicación aumentando el número de réplicas del StatefulSet ***mysql***.
```
kubectl scale statefulset mysql --replicas=5
kubectl get pods
```

La salida es como la siguiente. (Nota: Observar como tenemos 5 pods de MySQL.)
```
NAME                READY   STATUS    RESTARTS   AGE
mysql-0             2/2     Running   0          29m
mysql-1             2/2     Running   0          28m
mysql-2             2/2     Running   0          2m32s
mysql-3             2/2     Running   0          36s
mysql-4             2/2     Running   0          20s
mysql-client-loop   1/1     Running   0          11m
```

En la terminal principal podemos ver como se reparte el tráfico entre las 5 instancias.
```
kubectl run mysql-client --image=mysql:5.7 -i -t --rm --restart=Never --\
  mysql -h mysql-3.mysql -e "SELECT * FROM test.messages"
``` 

Podemos comprobar como las nuevas instancias han replicado la base de datos:
```
kubectl run mysql-client --image=mysql:5.7 -i -t --rm --restart=Never --\
  mysql -h mysql-3.mysql -e "SELECT * FROM test.messages"
```

Cuando se crean réplicas, la Persistent Volume Claim crea nuevos volúmenes:
```
kubectl get pvc
```

La salida es como la siguiente: (Nota: puede verse la PVC y el volumen)
```
NAME           STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
data-mysql-0   Bound    pvc-91e328d7-db30-4c99-97d5-ce47b518187e   10Gi       RWO            standard       33m
data-mysql-1   Bound    pvc-4186206f-ce33-432f-9cdc-60974acf1be9   10Gi       RWO            standard       32m
data-mysql-2   Bound    pvc-a06b0957-9397-48e7-a81b-5e8374f21b1e   10Gi       RWO            standard       32m
data-mysql-3   Bound    pvc-a75a2c39-b97d-439f-a7e2-48486fda3279   10Gi       RWO            standard       4m26s
data-mysql-4   Bound    pvc-7d283858-efda-46f1-87ad-a010cc439dac   10Gi       RWO            standard       4m10s
```

Mostramos los volúmenes persistentes (PVs)
```
kubectl get pv
```

La salida es:
```
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                  STORAGECLASS   REASON   AGE
pvc-3f34e81d-8bd8-431c-8d1e-0b67955bd00c   10Gi       RWO            Delete           Bound    default/data-mysql-4   standard                66s
pvc-4186206f-ce33-432f-9cdc-60974acf1be9   10Gi       RWO            Delete           Bound    default/data-mysql-1   standard                39m
pvc-8ea5d05d-48cc-4e13-9c2c-8199776bdcf6   10Gi       RWO            Delete           Bound    default/data-mysql-3   standard                82s
pvc-91e328d7-db30-4c99-97d5-ce47b518187e   10Gi       RWO            Delete           Bound    default/data-mysql-0   standard                40m
pvc-a06b0957-9397-48e7-a81b-5e8374f21b1e   10Gi       RWO            Delete           Bound    default/data-mysql-2   standard                39m
```

Si un pod vuelve a reclamar el volumen, la ***RECLAIM POLICY*** a ***Delete*** garantiza que se borrarán los datos previos.

Realizamos un scale-in. Pasamos a 3 réplicas.
```
kubectl scale statefulset mysql  --replicas=3
```

Volvemosa mostrar las PVCs
```
kubectl get pvc
```

Comprobar que siguen existiendo las 5 PVCs, a pesar de que se han destruidos dos pods. Es muy importante entender que la eliminación de las PVCs (y sus respectivos PV) deberá realizarse de forma explícita. (Nota: Si se borra la PVC, se borrará el PV asociado)
```
kubectl delete pvc data-mysql-3
kubectl delete pvc data-mysql-4
```
```
kubectl get pvc
```
```
kubectl get pv
```

Cerramos la segunda terminal.

En la terminal principal CTRL+C.

Limpiamos los recursos:
```
kubectl delete statefulset mysql
kubectl delete service mysql
kubectl delete service mysql-read
```

No olvidemos borrar las PVCs.
```
kubectl get pvc
```
```
kubectl delete pvc data-mysql-0
kubectl delete pvc data-mysql-1
kubectl delete pvc data-mysql-2
```

Comprobamos:
```
kubectl get all
kubectl get pv
```
