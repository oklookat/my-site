FROM node:16.2-alpine3.12


# NODE_MODULES CACHE
WORKDIR /usr/src/cache
COPY ["package.json", "package-lock.json*", "./"]
RUN npm install


ENV APP_PATH /usr/src/app
WORKDIR $APP_PATH
COPY . .
EXPOSE 8080
