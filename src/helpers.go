package main

import (
	"encoding/json"
	"github.com/getsentry/sentry-go"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"regexp"
)

func IsCallbackQuery(update tgbotapi.Update) bool {
	return update.Message == nil && update.CallbackQuery != nil
}
func IsRussian(searchWord string) bool {
	var rxCyrillic = regexp.MustCompile("^[\u0400-\u04FF\u0500-\u052F]+$")
	return rxCyrillic.MatchString(searchWord)
}
func LogObject(update interface{}) {
	var updateInterface map[string]interface{}
	inrec, _ := json.Marshal(update)
	json.Unmarshal(inrec, &updateInterface)
	log.WithFields(updateInterface).Info("")
}

func captureFatalErrorIfNotNull(err error) {
	if err == nil {
		return
	}
	sentry.CaptureException(err)
	log.Fatal(err)
}

func capturePanicErrorIfNotNull(err error) {
	if err == nil {
		return
	}
	sentry.CaptureException(err)
	log.Panic(err)
}
