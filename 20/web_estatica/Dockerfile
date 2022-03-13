# Dockerfile para web estática.
FROM ubuntu:16.04
RUN apt-get -y update; apt-get -y install nginx
RUN echo 'Hola, estoy dentro de tu contenedor' > /var/www/html/index.html
EXPOSE 80