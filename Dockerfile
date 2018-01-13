FROM golang:1.8.5

ADD . /go/src/vault

WORKDIR /go/src/vault

RUN go build ./

CMD ["./vault"]