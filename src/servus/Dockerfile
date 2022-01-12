FROM golang:latest


# build
ARG BUILD_DIR=/app/build
RUN mkdir -p $BUILD_DIR
WORKDIR $BUILD_DIR
ADD . $BUILD_DIR
WORKDIR $BUILD_DIR/server
RUN go build
# start
ARG SERVE_DIR=/app
RUN mv server ${SERVE_DIR}/servus
WORKDIR $SERVE_DIR
RUN rm -rf ${BUILD_DIR}
CMD ["./servus"]
EXPOSE 3333
