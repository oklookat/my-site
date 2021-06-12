FROM node:16.2-alpine3.12


ENV APP_PATH /usr/src/app
WORKDIR $APP_PATH
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 3333
CMD ["node", "ace", "serve", "--watch"]