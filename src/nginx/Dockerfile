FROM nginx:latest

# uploads
ARG UPLOADS_DIR=/app/uploads
RUN mkdir -p ${UPLOADS_DIR}

# main
COPY nginx.conf /etc/nginx/nginx.conf
CMD ["nginx", "-g", "daemon off;"]
EXPOSE 443
