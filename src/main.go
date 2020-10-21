package main

import (
	"encoding/json"
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
	log.WithFields(log.Fields{
		"event_type":          "app_event",
		"telegram_message":    "",
		"article_search_type": "",
	}).Error(err.Error())
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

	if environment.Env != "dev" {
		_, err = bot.SetWebhook(tgbotapi.NewWebhook(environment.WebhookAddress))
		captureErrorIfNotNull(err)
		info, err := bot.GetWebhookInfo()
		captureErrorIfNotNull(err)
		if info.LastErrorDate != 0 {
			log.WithFields(log.Fields{
				"event_type":          "app_event",
				"telegram_message":    "",
				"article_search_type": "",
			}).Error("Telegram callback failed: %s", info.LastErrorMessage)
		}
	}

	updates := bot.ListenForWebhook("/" + environment.UuidToken)
	go http.ListenAndServe("0.0.0.0:"+environment.AppPort, nil)

	for update := range updates {

		if update.Message == nil && update.CallbackQuery != nil {
			inrec, _ := json.Marshal(update)
			jsonStrMsg := strings.ReplaceAll(string(inrec), "null", "\"null\"")
			log.WithFields(log.Fields{
				"event_type":          "incoming_message",
				"telegram_message":    jsonStrMsg,
				"article_search_type": "article_switch",
			}).Info()
			conf := &tgbotapi.EditMessageTextConfig{}
			conf.ParseMode = "html"
			conf.MessageID = update.CallbackQuery.Message.MessageID
			conf.ChatID = update.CallbackQuery.Message.Chat.ID
			// fixing error
			//Bad Request: message is not modified: specified new message content and reply markup are exactly
			//the same as a current content and reply markup of the message
			if update.CallbackQuery.Data == "none" {
				_, err = bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "done"))
				continue
			}
			keysArr := strings.Split(update.CallbackQuery.Data, ",") // probleem,1
			keyword := keysArr[0]
			index, _ := strconv.ParseInt(keysArr[1], 10, 64)
			indexInt, _ := strconv.Atoi(keysArr[1])
			conf.Text = redis.GetArticleByIndex(keyword, index)
			buttonsLen := redis.GetArticlesLen(keyword)

			if buttonsLen > 1 {
				replyMarkup := MakeReplyMarkup(keyword, buttonsLen, indexInt)
				conf.ReplyMarkup = &replyMarkup
			}

			_, err := bot.Send(conf)
			captureErrorIfNotNull(err)
			callbackConfig := tgbotapi.NewCallback(update.CallbackQuery.ID, "done")
			_, err = bot.AnswerCallbackQuery(callbackConfig)
			captureErrorIfNotNull(err)
			continue
		}
		inrec, _ := json.Marshal(update)
		jsonStrMsg := strings.ReplaceAll(string(inrec), "null", "\"null\"")
		log.WithFields(log.Fields{
			"event_type":          "incoming_message",
			"telegram_message":    jsonStrMsg,
			"article_search_type": "new_search",
		}).Info()

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
		err = redis.pushToChannel(searchWord)
		captureErrorIfNotNull(err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, articles[0])
		if len(articles) > 1 {
			msg.ReplyMarkup = MakeReplyMarkup(searchWord, len(articles), 0)
		}
		msg.ParseMode = "html"
		_, err := bot.Send(msg)
		captureErrorIfNotNull(err)
	}
}
