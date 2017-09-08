FROM postgres:9.6-alpine

ENV POSTGRES_INITDB_ARGS "-E UTF8"

RUN apk add --no-cache \
  gcc \
  musl-dev \
  python3-dev \
  python3 \
  wget \
  bash

# Python scripts will need to parse data and insert it into
# tables defined in /docker-entrypoint-initdb.d/
RUN pip3 install psycopg2

WORKDIR /data

RUN wget --no-check-certificate \
   https://www.iana.org/assignments/language-subtag-registry/language-subtag-registry

ARG ENTRY="/docker-entrypoint-initdb.d/"
ADD 0schema.sql $ENTRY
ADD import-languages.py run-import-languages.sh $ENTRY
