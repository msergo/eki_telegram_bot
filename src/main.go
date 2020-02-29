package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/Netflix/go-env"
	"github.com/getsentry/sentry-go"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

var environment Environment

type EkiEe struct {
	redisWorker RedisWorker
	bot         *tgbotapi.BotAPI
}

func (s *EkiEe) Init() {
	var err error
	redisWorker := InitRedisWorker()
	_, err = redisWorker.Ping() //TODO: think about Redis failure
	captureFatalErrorIfNotNull(err)
	s.redisWorker = redisWorker
}

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
	ekiEe := EkiEe{}
	ekiEe.Init()

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := ioutil.ReadAll(r.Body)

		var update tgbotapi.Update
		LogObject(update)
		json.Unmarshal(bytes, &update)
		var response tgbotapi.Chattable
		if !IsCallbackQuery(update) {
			response = ekiEe.MakeNewSearchResponse(update)
		}
		response = ekiEe.MakeArticleSwitchResponse(update)
		if response != nil {
			_, err = bot.Send(response)
			captureFatalErrorIfNotNull(err)
		}

	})

	go http.ListenAndServe("0.0.0.0:"+environment.AppPort, nil)

}
func (e *EkiEe) MakeNewSearchResponse(update tgbotapi.Update) tgbotapi.Chattable {
	var articles []string
	searchWord := strings.ToLower(update.Message.Text)
	articles = e.redisWorker.GetAllArticles(searchWord)
	if len(articles) == 0 {
		articles = FetchArticles(searchWord)
	}
	if len(articles) == 0 {
		return nil
	}
	e.redisWorker.StoreArticlesSet(searchWord, articles)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, articles[0])
	msg.ParseMode = "html"
	if len(articles) > 1 {
		msg.ReplyMarkup = MakeReplyMarkup(searchWord, len(articles), "0")
	}
	return msg
}

func (e *EkiEe) MakeArticleSwitchResponse(update tgbotapi.Update) tgbotapi.Chattable {
	keysArr := strings.Split(update.CallbackQuery.Data, ",") // TODO: refactor here
	keyword := strings.ToLower(keysArr[0])
	buttonsLen := e.redisWorker.GetArticlesLenByKeyword(keyword)

	newText := e.redisWorker.GetArticleByIndex(keyword, keysArr[1])
	dataToSend := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, newText)
	if buttonsLen > 1 {
		replyMarkup := MakeReplyMarkup(keyword, buttonsLen, keysArr[1])
		dataToSend.ReplyMarkup = &replyMarkup
	}

	return dataToSend
}
