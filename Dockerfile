FROM golang:latest AS build

WORKDIR /usr/src/tech-db-forum

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build cmd/app/main.go

FROM ubuntu:20.04

MAINTAINER Oleg Gibadulin

RUN apt-get -y update && apt-get install -y tzdata

ENV TZ=Russia/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV PGVER 12
RUN apt-get -y update && apt-get install -y postgresql-$PGVER

USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER oleg WITH SUPERUSER PASSWORD 'postgres';" &&\
    createdb -O oleg forum &&\
    /etc/init.d/postgresql stop

EXPOSE 5432

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

COPY ./scripts/init.sql ./scripts/init.sql
COPY ./config.json ./config.json
COPY --from=build /usr/src/tech-db-forum/main .

EXPOSE 5000
ENV PGPASSWORD postgres
CMD service postgresql start && psql -h localhost -d forum -U oleg -p 5432 -a -q -f ./scripts/init.sql && ./main
