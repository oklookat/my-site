FROM nginx
USER root
ENV CONTAINER_PATH /usr/src/nginx
ENV NGINX_CONF /etc/nginx
ENV NGINX_LOGS /etc/nginx/logs/

WORKDIR $NGINX_CONF
COPY ["nginx.conf", "./"]
RUN ["mkdir", "/etc/nginx/logs"]
WORKDIR $NGINX_LOGS
RUN ["touch", "error.log"]
RUN ["touch", "nginx.pid"]

WORKDIR $CONTAINER_PATH
COPY . .

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
