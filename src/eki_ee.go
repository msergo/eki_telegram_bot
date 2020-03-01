package main

import (
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// EkiEe interface for most bot functional
type EkiEe struct {
	redisWorker RedisWorker
	telegram    *tgbotapi.BotAPI
}

// InitEki a constructor
func InitEki() EkiEe {
	var err error
	redisWorker := InitRedisWorker()
	_, err = redisWorker.Ping() //TODO: think about Redis failure
	captureFatalErrorIfNotNull(err)
	bot, err := tgbotapi.NewBotAPI(environment.BotToken)
	captureFatalErrorIfNotNull(err)
	bot.Debug = false

	return EkiEe{redisWorker: redisWorker, telegram: bot}
}

// MakeNewSearchResponse prepare response for new search
func (e *EkiEe) MakeNewSearchResponse(update tgbotapi.Update) tgbotapi.Chattable {
	var articles []string
	searchWord := strings.ToLower(update.Message.Text)
	articles = e.redisWorker.GetAllArticles(searchWord)
	articlesLen := len(articles)
	if articlesLen == 0 {
		articles = FetchArticles(searchWord)
	}
	articlesLen = len(articles)
	if articlesLen == 0 {
		return nil
	}
	e.redisWorker.StoreArticlesSet(searchWord, articles)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, articles[0])
	msg.ParseMode = tgbotapi.ModeHTML
	if len(articles) > 1 {
		msg.ReplyMarkup = MakeReplyMarkup(searchWord, len(articles), "0")
	}
	return msg
}

// MakeArticleSwitchResponse prepare a response for switching articles
func (e *EkiEe) MakeArticleSwitchResponse(update tgbotapi.Update) tgbotapi.Chattable {
	keysArr := strings.Split(update.CallbackQuery.Data, ",")
	keyword := strings.ToLower(keysArr[0])
	buttonsLen := e.redisWorker.GetArticlesLenByKeyword(keyword)

	newText := e.redisWorker.GetArticleByIndex(keyword, keysArr[1])
	dataToSend := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, newText)
	dataToSend.ParseMode = tgbotapi.ModeHTML

	if buttonsLen > 1 {
		replyMarkup := MakeReplyMarkup(keyword, buttonsLen, keysArr[1])
		dataToSend.ReplyMarkup = &replyMarkup
	}

	return dataToSend
}
