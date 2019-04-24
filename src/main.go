package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
	"github.com/msergo/eki_telegram_bot/src/translation_fetcher"
	"os"
)

func main() {
	//result := translation_fetcher.GetTranslations("saama")
	//fmt.Print(result)
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
	go http.ListenAndServe("0.0.0.0:" + os.Getenv("PORT"), nil)

	for update := range updates {
		articles := translation_fetcher.GetTranslations(update.Message.Text)
		for i := 0; i < len(articles); i++ {
			article := articles[i]
			articleHeader := article.ArticleHeader
			articleText := ""
			for j := 0; j < len(article.Meanings); j ++ {
				articleText += articles[i].Meanings[j].Translation + "\r\n"
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, articleHeader+"\r\n"+articleText)
			//msg.ReplyToMessageID = update.Message.MessageID
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}

		}
	}
}
