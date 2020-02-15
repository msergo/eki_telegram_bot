package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	redis := InitRedisWorker()
	_, err := redis.Ping()
	if err != nil {
		log.Fatalf("Redis connecting error %s", err)
	}
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(os.Getenv("WEBHOOK_ADDRESS")))
	if err != nil {
		log.Fatal(err)
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	updates := bot.ListenForWebhook("/" + os.Getenv("UUID_TOKEN"))
	go http.ListenAndServe("0.0.0.0:"+os.Getenv("PORT"), nil)

	var buttons []tgbotapi.InlineKeyboardButton

	for update := range updates {
		if update.Message == nil {
			conf := &tgbotapi.EditMessageTextConfig{}
			conf.ParseMode = "html"
			if update.CallbackQuery == nil || update.CallbackQuery.Message == nil { // hotfix for edited msgs TODO!
				continue
			}
			conf.MessageID = update.CallbackQuery.Message.MessageID
			conf.ChatID = update.CallbackQuery.Message.Chat.ID
			keysArr := strings.Split(update.CallbackQuery.Data, ",")
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

			if _, err := bot.Send(conf); err != nil {
				log.Print(err)
			}
			callbackConfig := tgbotapi.NewCallback(update.CallbackQuery.ID, "done")
			bot.AnswerCallbackQuery(callbackConfig)
			continue
		}
		var articles []string
		articles = redis.GetAllArticles(update.Message.Text)
		if len(articles) == 0 {
			articles = GetArticles(update.Message.Text)
			redis.StoreArticlesSet(update.Message.Text, articles)
		}
		buttons = buttons[:0]
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, articles[0])
		if len(articles) > 1 {
			msg.ReplyMarkup = MakeReplyMarkupSmart(update.Message.Text, len(articles), 0)
		}
		msg.ParseMode = "html"
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
