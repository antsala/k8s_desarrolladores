FROM ubuntu:20.04
RUN apt-get update -yqq; apt-get -yqq install ruby ruby-dev build-essential redis-tools
RUN gem install sinatra redis:4.4.0
RUN mkdir -p /opt/webapp
EXPOSE 4567
CMD ["/opt/webapp/bin/webapp"]
