FROM golang:1.14-alpine
RUN apk add --no-cache wget curl git build-base
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN mkdir -p $GOPATH/src/github.com/msergo/eki_telegram_bot
COPY . $GOPATH/src/github.com/msergo/eki_telegram_bot
RUN cd $GOPATH/src/github.com/msergo/eki_telegram_bot && \
    dep ensure && \
    go test -vet=off -v ./...

WORKDIR $GOPATH/src/
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/src/github.com/msergo/eki_telegram_bot/cmd/main github.com/msergo/eki_telegram_bot/src
ENV CI=true
FROM alpine:3.7
RUN apk --no-cache add ca-certificates netcat-openbsd
WORKDIR /root/
COPY --from=0 /go/src/github.com/msergo/eki_telegram_bot/cmd/main .
COPY run.sh .
RUN chmod +x run.sh
CMD ["./run.sh", "./main"]
