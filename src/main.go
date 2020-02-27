package main

import (
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Netflix/go-env"
	"github.com/getsentry/sentry-go"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var environment Environment

func captureErrorIfNotNull(err error) {
	if err == nil {
		return
	}
	log.Fatal(err)
	sentry.CaptureException(err)
}
func main() {
	log.SetFormatter(&log.JSONFormatter{})
	es, err := env.UnmarshalFromEnviron(&environment)
	captureErrorIfNotNull(err)
	if err != nil {
		log.Panic(err)
	}
	// Remaining environment variables.
	environment.Extras = es

	sentry.Init(sentry.ClientOptions{
		Dsn:              environment.SentryDsn,
		AttachStacktrace: true,
		Environment:      environment.Env,
		ServerName:       environment.WebhookAddress,
	})
	redis := InitRedisWorker()
	_, err = redis.Ping()
	captureErrorIfNotNull(err)
	bot, err := tgbotapi.NewBotAPI(environment.BotToken)
	captureErrorIfNotNull(err)

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)
	if environment.Env != "dev" {
		_, err = bot.SetWebhook(tgbotapi.NewWebhook(environment.WebhookAddress))
		captureErrorIfNotNull(err)
		info, err := bot.GetWebhookInfo()
		captureErrorIfNotNull(err)
		if info.LastErrorDate != 0 {
			log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
		}
	}

	updates := bot.ListenForWebhook("/" + environment.UuidToken) // TODO: maybe remove
	go http.ListenAndServe("0.0.0.0:"+environment.AppPort, nil)

	for update := range updates {
		var dataToSend *tgbotapi.EditMessageTextConfig
		if IsCallbackQuery(update) {
			LogObject(update, "article_swtich")
			keysArr := strings.Split(update.CallbackQuery.Data, ",") // TODO: refactor here
			keyword := strings.ToLower(keysArr[0])
			buttonsLen := redis.GetArticlesLenByKeyword(keyword)
			dataToSend = &tgbotapi.EditMessageTextConfig{}

			if buttonsLen > 1 {
				replyMarkup := MakeReplyMarkupSmart(keyword, buttonsLen, keysArr[1])
				dataToSend.ReplyMarkup = &replyMarkup
			}
			dataToSend.ParseMode = "html"
			dataToSend.MessageID = update.CallbackQuery.Message.MessageID
			dataToSend.ChatID = update.CallbackQuery.Message.Chat.ID
			dataToSend.Text = redis.GetArticleByIndex(keyword, keysArr[1])

			_, err := bot.Send(dataToSend)
			captureErrorIfNotNull(err)
			continue
		}
		LogObject(update, "new_search")
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
			msg.ReplyMarkup = MakeReplyMarkupSmart(searchWord, len(articles), "0")
		}
		_, err := bot.Send(msg)
		captureErrorIfNotNull(err)
	}
}
