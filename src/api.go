package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
	log "github.com/sirupsen/logrus"
)

type EkiEeBot struct {
	redis RedisWorker
	telegram *tgbotapi.BotAPI
	logger log.Logger
}

func (e EkiEeBot) Init() {
	e.redis = InitRedisWorker()
	bot, _ := tgbotapi.NewBotAPI(environment.BotToken) // TODO: handle err
	e.telegram = bot
	e.logger.SetFormatter(&log.JSONFormatter{})
}

func (e EkiEeBot) ProcessNewMessage(update tgbotapi.Update) {
	//inrec, _ := json.Marshal(update)
	//json.Unmarshal(inrec, &updateInterface)
	//e.logger.WithFields(updateInterface).Info("article_switch")

	var articles []string
	searchWord := strings.ToLower(update.Message.Text)
	articles = e.redis.GetAllArticles(searchWord)
	if len(articles) == 0 {
		articles = FetchArticles(searchWord)
	}
	if len(articles) == 0 {
		return
	}
	e.redis.StoreArticlesSet(searchWord, articles)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, articles[0])
	if len(articles) > 1 {
		msg.ReplyMarkup = MakeReplyMarkupSmart(searchWord, len(articles), 0)
	}
	msg.ParseMode = "html"
	_, err := e.telegram.Send(msg)
	captureErrorIfNotNull(err)
}