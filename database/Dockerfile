FROM postgres:9.6-alpine

ENV POSTGRES_INITDB_ARGS "-E UTF8"

RUN apk add --no-cache bash

ADD 0_core_schema.sql 1_user_schema.sql /docker-entrypoint-initdb.d/ 
