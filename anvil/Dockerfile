# First half of multi-stage build
FROM golang:1.9 as builder

RUN apt-get update && apt-get install -y git

WORKDIR /go/src/mutably/anvil
COPY . .

RUN go-wrapper download \
 && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o anvil .

WORKDIR /external
RUN git clone https://github.com/vishnubob/wait-for-it.git

# Final Image
FROM alpine:latest

RUN apk --no-cache add bash

# The .xml dump should be here.
WORKDIR /archive
VOLUME ["/archive"]

WORKDIR /app
COPY --from=builder /go/src/mutably/anvil/anvil .
COPY --from=builder /external/wait-for-it/wait-for-it.sh .

ENTRYPOINT /bin/bash wait-for-it.sh -h $DATABASE_HOST -p 5432 -t 0 -- \
  ./anvil import -host=$DATABASE_HOST -port=5432 -d=$DATABASE_NAME \
  -u=$DATABASE_USER -p=$DATABASE_PASSWORD /archive/*.xml

