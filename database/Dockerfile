FROM postgres:9.6-alpine

ENV POSTGRES_INITDB_ARGS "-E UTF8"

RUN apk add --no-cache bash

ADD schema.sql /docker-entrypoint-initdb.d/ 
