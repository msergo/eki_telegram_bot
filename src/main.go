package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
	"github.com/msergo/eki_telegram_bot/src/translation_fetcher"
	"os"
	"strconv"
)

//var replykeyboard = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
//	tgbotapi.NewInlineKeyboardButtonData("<<<", "456"),
//	tgbotapi.NewInlineKeyboardButtonData(">>>", "456"),
//))
//
//
//func addButton(text string, dataIndex string) tgbotapi.InlineKeyboardButton {
//	return tgbotapi.NewInlineKeyboardButtonData(text, dataIndex)
//}

func main() {
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

	for update := range updates {
		if update.Message == nil {
			ed := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, update.CallbackQuery.Data)
			if _, err := bot.Send(ed); err != nil {
				log.Panic(err)
			}
			continue
		}

		articles := translation_fetcher.GetArticles(update.Message.Text)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, articles[0])
		var buttons []tgbotapi.InlineKeyboardButton
		for i := 0; i < len(articles); i++ {
			but := tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(i), strconv.Itoa(i)	)
			buttons = append(buttons, but)
		}
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons)
		msg.ParseMode = "html"
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}

		//if update.Message == nil {
		//	fmt.Println("aa")
		//
		//	ed := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "xxxxx")
		//	if _, err := bot.Send(ed); err != nil {
		//		log.Panic(err)
		//	}
		//	continue
		//} else {
		//
		//	articles := translation_fetcher.GetArticles(update.Message.Text)
		//	for i := 0; i < len(articles); i++ {
		//		msg := tgbotapi.NewMessage(update.Message.Chat.ID, articles[i])
		//		msg.ParseMode = "html"
		//		msg.ReplyMarkup = replykeyboard
		//		if _, err := bot.Send(msg); err != nil {
		//			log.Panic(err)
		//		}
		//	}
		//}

	}
}
