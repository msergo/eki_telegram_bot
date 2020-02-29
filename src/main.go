package main

import (
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Netflix/go-env"
	"github.com/getsentry/sentry-go"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"errors"
	"fmt"
)

var environment Environment

func main() {
	sentryInitError := sentry.Init(sentry.ClientOptions{
		Dsn:              environment.SentryDsn,
		AttachStacktrace: true,
		Environment:      environment.Env,
		ServerName:       environment.WebhookAddress,
	})
	capturePanicErrorIfNotNull(sentryInitError)
	log.SetFormatter(&log.JSONFormatter{})
	es, err := env.UnmarshalFromEnviron(&environment)
	capturePanicErrorIfNotNull(err)
	environment.Extras = es

	redis := InitRedisWorker()
	_, err = redis.Ping() //TODO: think about Redis failure
	captureFatalErrorIfNotNull(err)
	bot, err := tgbotapi.NewBotAPI(environment.BotToken)
	captureFatalErrorIfNotNull(err)

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)
	if environment.Env != "dev" {
		_, err = bot.SetWebhook(tgbotapi.NewWebhook(environment.WebhookAddress))
		capturePanicErrorIfNotNull(err)
		info, err := bot.GetWebhookInfo()
		capturePanicErrorIfNotNull(err)
		if info.LastErrorDate != 0 {
			capturePanicErrorIfNotNull(errors.New(fmt.Sprintf("Telegram callback failed: %s", info.LastErrorMessage)))
		}
	}

	updates := bot.ListenForWebhook("/" + environment.UuidToken) // TODO: maybe remove
	go http.ListenAndServe("0.0.0.0:"+environment.AppPort, nil)
	for update := range updates {
		LogObject(update)
		if update.Message.IsCommand() {
			continue
		}
		if IsCallbackQuery(update) {
			keysArr := strings.Split(update.CallbackQuery.Data, ",") // TODO: refactor here
			keyword := strings.ToLower(keysArr[0])
			buttonsLen := redis.GetArticlesLenByKeyword(keyword)

			newText := redis.GetArticleByIndex(keyword, keysArr[1])
			dataToSend := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, newText)
			if buttonsLen > 1 {
				replyMarkup := MakeReplyMarkup(keyword, buttonsLen, keysArr[1])
				dataToSend.ReplyMarkup = &replyMarkup
			}
			_, err := bot.Send(dataToSend)
			captureFatalErrorIfNotNull(err)
			continue
		}
		var articles []string
		searchWord := strings.ToLower(update.Message.Text)
		articles = redis.GetAllArticles(searchWord)
		if len(articles) == 0 {
			articles = FetchArticles(searchWord)
		}
		if len(articles) == 0 {
			continue
		}
		redis.StoreArticlesSet(searchWord, articles)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, articles[0])
		msg.ParseMode = "html"
		if len(articles) > 1 {
			msg.ReplyMarkup = MakeReplyMarkup(searchWord, len(articles), "0")
		}
		_, err := bot.Send(msg)
		captureFatalErrorIfNotNull(err)
	}
}
