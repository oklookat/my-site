FROM postgres:14.2-alpine


ENV POSTGRES_PASSWORD=root
ENV POSTGRES_DB=elven
ENV PGDATA=/var/lib/postgresql/data/pgdata

WORKDIR /app
COPY .devcontainer/run.sh .
RUN chmod +x ./run.sh

USER postgres