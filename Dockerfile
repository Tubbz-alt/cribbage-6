FROM golang:1.14.3-alpine3.11

ENV REST_PORT=8081
# These happen to coincide with a local myql server on macOS with a root user and no password
ENV DSN_HOST=host.docker.internal
ENV DSN_USER=root
ENV DSN_PASSWORD=

ENV GOPATH=/go
RUN mkdir -p $GOPATH/src/github.com/joshprzybyszewski/cribbage
WORKDIR $GOPATH/src/github.com/joshprzybyszewski/cribbage

COPY vendor vendor
COPY model model
COPY logic logic
COPY utils utils
COPY jsonutils jsonutils
COPY network network
COPY server server
COPY main.go main.go

EXPOSE 80

CMD go run main.go \
    -restPort=$REST_PORT \
    -dsn_host=$DSN_HOST \
    -dsn_user=$DSN_USER \
    -dsn_password=$DSN_PASSWORD