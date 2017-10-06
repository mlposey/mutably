# First half of multi-stage build
FROM golang:1.9 as builder

RUN apt-get update && apt-get install -y git

WORKDIR /go/src/mutably/api
COPY . .

RUN go-wrapper download \
 && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api .

WORKDIR /external
RUN git clone https://github.com/vishnubob/wait-for-it.git

# Final Image
FROM alpine:latest

RUN apk --no-cache add bash

WORKDIR /app
COPY --from=builder /go/src/mutably/api/api .
COPY --from=builder /external/wait-for-it/wait-for-it.sh .

ENTRYPOINT /bin/bash wait-for-it.sh $DATABASE_HOST:5432 -t 0 -- ./api