package main

import (
	"regexp"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

func IsCallbackQuery(update tgbotapi.Update) bool {
	return update.Message == nil && update.CallbackQuery != nil
}
func IsRussian(searchWord string) bool {
	var rxCyrillic = regexp.MustCompile("^[\u0400-\u04FF\u0500-\u052F]+$")
	return rxCyrillic.MatchString(searchWord)
}
func LogObject(update interface{}, msgType string) {
	var updateInterface map[string]interface{}
	inrec, _ := json.Marshal(update)
	json.Unmarshal(inrec, &updateInterface)
	log.WithFields(updateInterface).Info(msgType)
}
func CreateCallbackQueryResponse(update tgbotapi.Update, text string, replyMarkup tgbotapi.InlineKeyboardMarkup) *tgbotapi.EditMessageTextConfig {
	conf := &tgbotapi.EditMessageTextConfig{}
	conf.ParseMode = "html"
	conf.MessageID = update.CallbackQuery.Message.MessageID
	conf.ChatID = update.CallbackQuery.Message.Chat.ID
	conf.Text = text
	conf.ReplyMarkup = &replyMarkup
	return conf
}
