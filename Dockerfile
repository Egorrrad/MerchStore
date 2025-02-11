FROM golang:1.23

WORKDIR ${GOPATH}/avito-shop/
COPY . ${GOPATH}/avito-shop/

RUN go build -o /build ./src/cmd/service \
    && go clean -cache -modcache

EXPOSE 8080

CMD ["/build"]