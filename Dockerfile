FROM golang:1.22-alpine
ENV CI=true

RUN mkdir -p $GOPATH/src/github.com/msergo/eki_telegram_bot
WORKDIR $GOPATH/src/github.com/msergo/eki_telegram_bot
COPY . .
RUN go mod tidy
RUN go test -vet=off -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o ./cmd/main 


FROM alpine:3.19
RUN apk --no-cache add ca-certificates netcat-openbsd
WORKDIR /root/
COPY --from=0 /go/src/github.com/msergo/eki_telegram_bot/cmd/main .
COPY run.sh .
RUN chmod +x run.sh
CMD ["./run.sh", "./main"]
