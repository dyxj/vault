FROM golang:1.17

ADD . /go/src/vault

WORKDIR /go/src/vault

RUN go build ./

CMD ["./vault"]