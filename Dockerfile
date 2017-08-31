FROM golang:1.9-alpine

RUN apk add --no-cache git

WORKDIR /go/src/anvil
COPY . .

RUN go-wrapper download
RUN go-wrapper install

ENTRYPOINT ["go-wrapper", "run"]
