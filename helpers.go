package main

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/getsentry/sentry-go"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

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

func isCallbackQuery(update tgbotapi.Update) bool {
	return update.Message == nil && update.CallbackQuery != nil
}

func isNewSearchRequest(update tgbotapi.Update) bool {
	return update.Message != nil && update.CallbackQuery == nil
}

func logIncomingMessage(update tgbotapi.Update) {
	jsonUpdate, _ := json.Marshal(update)
	jsonStringEscaped := strings.Replace(string(jsonUpdate), "null", "\"null\"", -1)
	var articleSearchType string

	if isCallbackQuery(update) {
		articleSearchType = "article_switch"
	} else if isNewSearchRequest(update) {
		articleSearchType = "new_search"
	}

	log.WithFields(log.Fields{
		"event_type":          "incoming_message",
		"telegram_message":    jsonStringEscaped,
		"article_search_type": articleSearchType,
	}).Info()
}

func getKeywordAndIndex(data string) (string, int64) {
	keysArr := strings.Split(data, ",") // probleem,1
	keyword := keysArr[0]
	index, _ := strconv.ParseInt(keysArr[1], 10, 64)

	return keyword, index
}

func isProduction() bool {
	return environment.Env == "prod"
}
