# Laboratorio 65-D: ***Despliegue Canary***
 
En este laboratorio aprenderemos a usar los despliegues Canary.

Requisitos:

1. Una máquina virtual con Ubuntu 20.04 LTS a la que poder hacer ssh o escritorio remoto.
2. Cluster ***Minikube***.

## Ejercicio 1. ***¿Qué es un despliegue Canary?***

Supongamos que tenemos desplegada una app en el cluster de producción. En el workflow de CI/CD es habitual tener preparada la próxima versión de nuestra aplicación. Ésta se despliega y se prueba en el entorno de desarrollo o de pre-producción. La cuestión que se nos plantea es obvia: ¿Un despliegue exitoso en el entorno de desarrollo/preprod garantiza lo mismo en el de producción?

La respuesta dependerá de la experiencia de cada uno de nosotros, pero podemos decir de forma categórica que no. Aunque el entorno de pre-producción se intenta hacer lo más parecido al de producción, siempre existirán ligeras diferencias que podrían jugarnos una mala pasada.

En consecuencia, la única forma de garantizar que la nueva versión funcionará correctamente, es desplegarla en el entorno de producción. Para que no afecte realmente, debemos hacer un depsliegue de tipo ***canary***.

Supongamos que tenemos ***3 pods*** que usan la ***versión 1.0*** de la imagen de la aplicación. Debe existir un ***servicio*** que enrute el tráfico hacia ellos. En los despliegues ***canary***, desplegaremos un nuevo pod (con la versión nueva, p.e. 1.1) y nos apoyaremos en el uso de etiquetas para que el servicio actual ***no envíe tráfico*** al pod de versión ***1.1***. Otro servicio, con los selectores de etiquetas apropiados, enviarán tráfico al pod de versión 1.1.

Con este tipo de despliegues, los usuarios de la aplicación seguirán usando la versión ***1.0***, mientras que otros usuarios seleccionados, podrán usar la versión ***1.1*** en el entorno de producción. Cuando nos aseguremos que la nueva versión funciona correctamente, actualizaremos la versión de la imagen de los pod para el despliegue de producción.

Como curiosidad aclaramos el significado de ***canary***. El canario es un pájaro que se solía usar antiguamente en la minas de carbón. Este ave es muy delicada frente a la presencia de gases nocivos,  peligrosos, o explosivos (como el grisú) que se acumulan en las minas. Los mineros solian usar estos pájaros que, frente a la presencia de estos gases, morían, indicando de esta forma un peligro inminente para los mineros.

## Ejercicio 2. ***Ejemplo de un despliegue Canary***.

Cambiamos al directorio de trabajo.
```
cd ~/k8s_desarrolladores/65
```

Vamos a usar archivos html sencillos que representarán las dos versiones de la aplicación. Haremos uso de configmaps para inyectarlos.

El archivo ***lab-65-D-app-v1-configmap.yaml*** contiene el código fuente de la ***versión 1.0***. Lo inspeccionamos:
```
code lab-65-D-app-v1-configmap.yaml
```

De forma análoga abrimos el configmap para la ***versión 2.0*** de la aplicación.
```
code lab-65-D-app-v2-configmap.yaml
```

Creamos los configmaps
```
kubectl apply -f lab-65-D-app-v1-configmap.yaml
kubectl apply -f lab-65-D-app-v2-configmap.yaml
```

Comprobamos
```
kubectl describe configmap app-v1-configmap
```
```
kubectl describe configmap app-v2-configmap
```

Vamos a hacer un ejemplo de este tipo de despliegue. Para ello usaremos una imagen de ***nginx***. Abrimos el archivo ***lab-65-D-nginx-deployment.yaml***.
```
code lab-65-D-nginx-deployment.yaml
```

Creamos el deployment.
```
kubectl apply -f lab-65-D-nginx-deployment.yaml
```

Comprobamos que se despliega correctamente.
```
kubectl get all
```

Procedemos a crear un ***servicio*** para balancear el tráfico contra los tres pods que acabamos de iniciar. Abrimos el archivo ***lab-65-D-nginx-service.yaml***
```
code lab-65-D-nginx-service.yaml
```

Desplegamos el servicio.
```
kubectl apply -f lab-65-D-nginx-service.yaml
```

