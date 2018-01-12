# FROM golang:1.8.5

# ADD . /go/src/vault

# # RUN go get golang.org/x/crypto/acme/autocert && go install vault
# RUN go install vault

# WORKDIR /go/src/vault

# ENTRYPOINT /go/bin/vault

# EXPOSE 8080

FROM golang:1.8.5

ADD . /go/src/vault

WORKDIR /go/src/vault

# RUN go build ./

# CMD ["./vault"]