FROM golang:1.23-alpine
ENV CI=true

# gcc and musl-dev are required for CGO (go-sqlite3)
RUN apk add --no-cache gcc musl-dev

RUN mkdir -p $GOPATH/src/github.com/msergo/eki_telegram_bot
WORKDIR $GOPATH/src/github.com/msergo/eki_telegram_bot
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=1 go test -vet=off -v ./...

RUN CGO_ENABLED=1 GOOS=linux go build -o ./cmd/main


FROM alpine:3.19
RUN apk --no-cache add ca-certificates netcat-openbsd
WORKDIR /root/
COPY --from=0 /go/src/github.com/msergo/eki_telegram_bot/cmd/main .
COPY run.sh .
RUN chmod +x run.sh
CMD ["./run.sh", "./main"]
