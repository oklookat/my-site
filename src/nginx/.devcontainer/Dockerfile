FROM nginx:1.21.6-alpine

# install mkcert
RUN apk update && apk add mkcert
WORKDIR /app/certs
RUN mkcert -key-file key.pem -cert-file cert.pem oklookat.ru "*.oklookat.ru" localhost 127.0.0.1 ::1

# install dev certs (mkcert)
#COPY .devcontainer/certs .
#RUN cp ./mkcert_root.pem /usr/local/share/ca-certificates/mkcert_root.pem
#RUN update-ca-certificates

# replace mime types
RUN rm /etc/nginx/mime.types
COPY ./mime.types /etc/nginx/mime.types

# replace config
RUN rm /etc/nginx/nginx.conf
COPY .devcontainer/nginx.conf /etc/nginx/nginx.conf

# run
CMD ["nginx", "-g", "daemon off;"]
EXPOSE 80 443
