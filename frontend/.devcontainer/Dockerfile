FROM node:lts-alpine3.15

# build
WORKDIR /app
COPY . .
RUN apk update

# copy entry
RUN chmod +x ./.devcontainer/entrypoint.sh

# run
ENTRYPOINT [".devcontainer/entrypoint.sh"]
CMD [""]
EXPOSE 3001