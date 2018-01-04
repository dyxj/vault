FROM golang:1.8.5

ADD . /go/src/vault

RUN go install vault

WORKDIR /go/src/vault

ENTRYPOINT /go/bin/vault

EXPOSE 8080