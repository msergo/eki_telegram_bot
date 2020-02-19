package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Netflix/go-env"
	"github.com/getsentry/sentry-go"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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

	var buttons []tgbotapi.InlineKeyboardButton

	for update := range updates {
		if update.Message == nil && update.CallbackQuery != nil {
			conf := &tgbotapi.EditMessageTextConfig{}
			conf.ParseMode = "html"
			conf.MessageID = update.CallbackQuery.Message.MessageID
			conf.ChatID = update.CallbackQuery.Message.Chat.ID
			keysArr := strings.Split(update.CallbackQuery.Data, ",") // TODO: refactor here
			keyword := keysArr[0]
			index, _ := strconv.ParseInt(keysArr[1], 10, 64)
			indexInt, _ := strconv.Atoi(keysArr[1])
			conf.Text = redis.GetArticleByIndex(keyword, index)
			buttons = buttons[:0]
			buttonsLen := redis.GetArticlesLen(keyword)

			if buttonsLen > 1 {
				replyMarkup := MakeReplyMarkupSmart(keyword, buttonsLen, indexInt)
				conf.ReplyMarkup = &replyMarkup
			}

			_, err := bot.Send(conf)
			captureErrorIfNotNull(err)
			callbackConfig := tgbotapi.NewCallback(update.CallbackQuery.ID, "done")
			_, err = bot.AnswerCallbackQuery(callbackConfig)
			captureErrorIfNotNull(err)
			continue
		}
		var articles []string
		searchWord := strings.ToLower(update.Message.Text)
		articles = redis.GetAllArticles(searchWord)
		if len(articles) == 0 {
			articles = GetArticles(searchWord)
		}
		if len(articles) == 0 {
			continue
		}
		redis.StoreArticlesSet(searchWord, articles)
		buttons = buttons[:0]
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, articles[0])
		if len(articles) > 1 {
			msg.ReplyMarkup = MakeReplyMarkupSmart(searchWord, len(articles), 0)
		}
		msg.ParseMode = "html"
		_, err := bot.Send(msg)
		captureErrorIfNotNull(err)
	}
}
