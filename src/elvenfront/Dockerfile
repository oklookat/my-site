FROM node:16.4.2
ENV APP_PATH /usr/src/app


WORKDIR $APP_PATH
COPY package*.json ./
RUN npm install

COPY . .

EXPOSE 8080
CMD ["npm", "run", "serve"]