FROM golang:1.19.2 AS builder
COPY . /go/src/build
WORKDIR /go/src/build
RUN update-ca-certificates &&\
    adduser --disabled-password --disabled-login --no-create-home --quiet --system -u 2003 mylittleboxy-backup &&\
    go get &&\
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/src/build/mylittleboxy-backup &&\
    chown 2003:2003 /go/src/build/mylittleboxy-backup &&\
    mkdir /mylittleboxy &&\
    chown mylittleboxy-backup:2003 /mylittleboxy

FROM debian
COPY --from='builder' /go/src/build/mylittleboxy-backup /go/bin/mylittleboxy-backup
COPY --from='builder' /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
#COPY --from='builder' /etc/passwd /etc/passwd
#COPY --from='builder' --chown=2003 /mylittleboxy /mylittleboxy
#USER mylittleboxy-backup
ENTRYPOINT ["/go/bin/mylittleboxy-backup"]