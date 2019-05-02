FROM iron/go:dev
RUN echo $GOPATH/src/github.com/msergo/eki_telegram_bot
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN mkdir -p $GOPATH/src/github.com/msergo/eki_telegram_bot
COPY . $GOPATH/src/github.com/msergo/eki_telegram_bot
RUN cd $GOPATH/src/github.com/msergo/eki_telegram_bot && \
    dep ensure && \
    go test -v ./...

WORKDIR $GOPATH/src/
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/src/github.com/msergo/eki_telegram_bot/cmd/main github.com/msergo/eki_telegram_bot/src

FROM msergo/redis_go:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/msergo/eki_telegram_bot/cmd/main .
CMD ["./main"]
