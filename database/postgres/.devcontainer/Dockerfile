FROM postgres:14.2-alpine

# set args
ENV POSTGRES_CONFIG=/app/postgres.conf

ARG POSTGRES_USER
ENV POSTGRES_USER=${POSTGRES_USER}

ARG POSTGRES_PASSWORD
ENV POSTGRES_PASSWORD=${POSTGRES_PASSWORD}

ARG PGDATA
ENV PGDATA=${PGDATA}

ARG POSTGRES_DB
ENV POSTGRES_DB=${POSTGRES_DB}

#
WORKDIR /app

# copy entry
COPY .devcontainer/entrypoint.sh .
RUN chmod 777 ./entrypoint.sh

# copy config
COPY .devcontainer/postgres.conf .
RUN chmod 777 ./postgres.conf

# run
USER ${POSTGRES_USER}
ENTRYPOINT ["/app/entrypoint.sh"]

CMD ["postgres", "-c", "config_file=${POSTGRES_CONFIG}"]
EXPOSE 5432