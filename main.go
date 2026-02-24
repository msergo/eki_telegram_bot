package main

import (
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/Netflix/go-env"
	"github.com/getsentry/sentry-go"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var environment Environment

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	extras, err := env.UnmarshalFromEnviron(&environment)
	captureErrorIfNotNull(err)

	if err != nil {
		log.Panic(err)
	}

	// Remaining environment variables.
	environment.Extras = extras

	sentry.Init(sentry.ClientOptions{
		Dsn:              environment.SentryDsn,
		AttachStacktrace: true,
		Environment:      environment.Env,
		ServerName:       environment.WebhookAddress,
	})

	cacheEngine, err := NewCacheEngine(environment.CacheEngine)
	if err != nil {
		captureErrorIfNotNull(err)
		log.Panic(err)
	}
	defer cacheEngine.Close()

	// Graceful shutdown on SIGINT / SIGTERM.
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Info("Shutting down – closing cache engine")
		cacheEngine.Close()
		os.Exit(0)
	}()

	bot, err := tgbotapi.NewBotAPI(environment.BotToken)
	captureErrorIfNotNull(err)

	bot.Debug = false

	if isProduction() {
		_, err = bot.SetWebhook(tgbotapi.NewWebhook(environment.WebhookAddress))
		captureErrorIfNotNull(err)

		info, err := bot.GetWebhookInfo()
		captureErrorIfNotNull(err)

		if info.LastErrorDate != 0 {
			log.WithFields(log.Fields{
				"event_type":          "app_event",
				"telegram_message":    "",
				"article_search_type": "",
			}).Error("Telegram callback failed: " + info.LastErrorMessage)
		}
	}

	updates := bot.ListenForWebhook("/")

	go http.ListenAndServe("0.0.0.0:"+environment.AppPort, nil)

	for update := range updates {
		logIncomingMessage(update)

		if isNewSearchRequest(update) {
			var articles []string

			keyword := strings.ToLower(update.Message.Text)
			articles, err = cacheEngine.GetAllArticles(keyword)
			captureErrorIfNotNull(err)

			if len(articles) == 0 {
				articles = GetArticles(keyword, nil)
			}

			if len(articles) == 0 {
				continue
			}

			err = cacheEngine.StoreArticles(keyword, articles)
			captureErrorIfNotNull(err)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, articles[0])

			if len(articles) > 1 {
				msg.ReplyMarkup = MakeReplyMarkup(keyword, len(articles), 0)
			}

			msg.ParseMode = "html"

			_, err := bot.Send(msg)
			captureErrorIfNotNull(err)

			continue
		}

		if isCallbackQuery(update) {
			if update.CallbackQuery.Data == "none" {
				_, err = bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "done"))
				captureErrorIfNotNull(err)

				continue
			}

			// Instead of sending a new message, we edit the existing one so it'll have smooth transition in the screen
			conf := &tgbotapi.EditMessageTextConfig{}
			conf.ParseMode = "html"
			conf.MessageID = update.CallbackQuery.Message.MessageID
			conf.ChatID = update.CallbackQuery.Message.Chat.ID

			keyword, index := getKeywordAndIndex(update.CallbackQuery.Data)
			conf.Text, err = cacheEngine.GetArticleByIndex(keyword, index)
			captureErrorIfNotNull(err)

			replyButtonsCnt, err := cacheEngine.GetArticlesLen(keyword)
			captureErrorIfNotNull(err)

			if replyButtonsCnt > 1 {
				replyMarkup := MakeReplyMarkup(keyword, replyButtonsCnt, int(index))

				conf.ReplyMarkup = &replyMarkup
			}

			_, err = bot.Send(conf)
			captureErrorIfNotNull(err)

			_, err = bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "done"))
			captureErrorIfNotNull(err)

			continue
		}

	}
}
