package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/Netflix/go-env"
	"github.com/getsentry/sentry-go"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var environment Environment

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

	ekiEe := EkiEe{}
	ekiEe.Init()

	log.Printf("Authorized on account %s", ekiEe.telegram.Self.UserName)
	if environment.Env != "dev" {
		_, err = ekiEe.telegram.SetWebhook(tgbotapi.NewWebhook(environment.WebhookAddress))
		capturePanicErrorIfNotNull(err)
		info, err := ekiEe.telegram.GetWebhookInfo()
		capturePanicErrorIfNotNull(err)
		if info.LastErrorDate != 0 {
			capturePanicErrorIfNotNull(fmt.Errorf("Telegram callback failed: %s", info.LastErrorMessage))
		}
	}

	http.HandleFunc("/"+environment.UuidToken, func(w http.ResponseWriter, r *http.Request) {
		var update tgbotapi.Update
		var response tgbotapi.Chattable

		bytes, _ := ioutil.ReadAll(r.Body)
		LogObject(update)
		json.Unmarshal(bytes, &update)

		if !IsCallbackQuery(update) {
			response = ekiEe.MakeNewSearchResponse(update)
		} else {
			response = ekiEe.MakeArticleSwitchResponse(update)
		}

		if response != nil {
			_, err = ekiEe.telegram.Send(response)
			captureFatalErrorIfNotNull(err)
		}

	})

	http.ListenAndServe("0.0.0.0:"+environment.AppPort, nil)

}
