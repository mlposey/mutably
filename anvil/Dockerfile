FROM golang:1.8

RUN apt-get update && apt-get install -y \
  git \
  bash \
  gcc \
  libonig2 libonig-dev

WORKDIR /page-archive
VOLUME ["/page-archive"]

WORKDIR /go/src/anvil
COPY . .

RUN go-wrapper download
RUN go-wrapper install

RUN git clone https://github.com/vishnubob/wait-for-it.git

ENTRYPOINT /bin/bash wait-for-it/wait-for-it.sh $DB_HOST:$DB_PORT -t 0 -- \
  go-wrapper run -import -host=$DB_HOST -port=$DB_PORT -d=$POSTGRES_DB \
  -u=$POSTGRES_USER -p=$POSTGRES_PASSWORD /page-archive/*.xml

