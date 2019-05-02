package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
	"github.com/msergo/eki_telegram_bot/src/translation_fetcher"
	"os"
	"strconv"
	"github.com/msergo/eki_telegram_bot/src/redis_worker"
	"strings"
)

func main() {
	redis := redis_worker.InitRedisWorker()
	_, err := redis.Ping()
	if err != nil {
		log.Fatalf("Redis connecting error %s", err)
	}
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

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
			conf.MessageID = update.CallbackQuery.Message.MessageID
			conf.ChatID = update.CallbackQuery.Message.Chat.ID
			keysArr := strings.Split(update.CallbackQuery.Data, ",")
			keyword := keysArr[0]
			index, _ := strconv.ParseInt(keysArr[1], 10, 64)
			indexInt, _ := strconv.Atoi(keysArr[1])
			conf.Text = redis.GetArticleByIndex(keyword, index)
			buttons = buttons[:0]
			buttonsLen := redis.GetArticlesLen(keyword)
			if (buttonsLen > 1) {
				replyMarkup := MakeReplyMarkupNice(keyword, buttonsLen, indexInt)
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
		if (len(articles) == 0) {
			articles = translation_fetcher.GetArticles(update.Message.Text)
			redis.StoreArticlesSet(update.Message.Text, articles)
		}
		buttons = buttons[:0]
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, articles[0])
		if (len(articles) > 1) {
			msg.ReplyMarkup = MakeReplyMarkupNice(update.Message.Text, len(articles), 0)
		}
		msg.ParseMode = "html"
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

//func MakeReplyMarkup(keyword string, buttonsLen int) tgbotapi.InlineKeyboardMarkup {
//	var buttons []tgbotapi.InlineKeyboardButton
//	for i := 0; i < buttonsLen; i++ {
//		callbackData := keyword + "," + strconv.Itoa(i) //probleem,1
//		but := tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(i), callbackData)
//		buttons = append(buttons, but)
//	}
//	return tgbotapi.NewInlineKeyboardMarkup(buttons)
//}

func MakeReplyMarkupNice(keyword string, buttonsLen int, indexFrom int) tgbotapi.InlineKeyboardMarkup {
	var endPos int
	var startPos int
	if (buttonsLen <= 5) {
		startPos = 0
		endPos = buttonsLen
	} else if buttonsLen > 5 {
		startPos = indexFrom
		endPos = startPos + 5

		if endPos > buttonsLen {
			endPos = buttonsLen
			startPos = endPos - 5
		}
	}
	var buttons []tgbotapi.InlineKeyboardButton
	for i := startPos; i < endPos; i++ {
		callbackData := keyword + "," + strconv.Itoa(i) //probleem,1
		but := tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(i), callbackData)
		buttons = append(buttons, but)
	}

	if (startPos > 0) {
		buttons[0].Text = "<<" +  buttons[0].Text
	}

	if (endPos < buttonsLen) {
		buttons[len(buttons)-1].Text += ">>"
	}

	return tgbotapi.NewInlineKeyboardMarkup(buttons)
}
