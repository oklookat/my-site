FROM node:lts


# build
ARG BUILD_DIR=/app/build
RUN mkdir -p $BUILD_DIR
WORKDIR $BUILD_DIR
ADD . $BUILD_DIR
RUN npm cache verify
RUN npm install
RUN npm run build
# move builded
ARG ELVEN_DIR=/app/elven
WORKDIR $ELVEN_DIR
RUN mv ${BUILD_DIR}/dist $ELVEN_DIR
RUN rm -rf ${BUILD_DIR}