FROM nginx:1.21.6-alpine

# copy root cert from host (mkcert)
WORKDIR /app/certs
COPY .devcontainer/certs .
RUN chmod 777 ./*
RUN cp ./mkcert_root.pem /usr/local/share/ca-certificates/mkcert_root.pem
RUN update-ca-certificates

# replace mime types
RUN rm /etc/nginx/mime.types
COPY ./mime.types /etc/nginx/mime.types

# replace config
RUN rm /etc/nginx/nginx.conf
COPY .devcontainer/nginx.conf /etc/nginx/nginx.conf

# run
WORKDIR /app

# due on vscode remote specific we run nginx in devcontainer.json file
EXPOSE 80 443