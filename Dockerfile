FROM golang:1.8.5

ADD . /go/src/vault

RUN go get golang.org/x/crypto/acme/autocert && go install vault

WORKDIR /go/src/vault

VOLUME ~/dvolumes/vault

ENTRYPOINT /go/bin/vault

EXPOSE 8080