FROM node:16.4.2
ENV APP_PATH /usr/src/app


WORKDIR $APP_PATH
COPY ["package.json", "package-lock.json", "./"]
RUN npm install
COPY . ./

EXPOSE 3333
CMD ["node", "ace", "serve", "--watch"]