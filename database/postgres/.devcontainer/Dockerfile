FROM postgres:14.2-alpine


ENV POSTGRES_PASSWORD=root
ENV POSTGRES_DB=elven
ENV PGDATA=/var/lib/postgresql/data/pgdata

WORKDIR /app
COPY .devcontainer/postgres.conf .
RUN chmod 777 ./postgres.conf
COPY .devcontainer/run.sh .
RUN chmod 777 ./run.sh

RUN mkdir -p $PGDATA
RUN chown -R postgres:postgres $PGDATA
USER postgres