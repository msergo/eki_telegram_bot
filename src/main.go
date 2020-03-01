package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"

	"time"

	"github.com/Netflix/go-env"
	"github.com/getsentry/sentry-go"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var environment Environment

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	sentryInitError := sentry.Init(sentry.ClientOptions{
		Dsn:              environment.SentryDsn,
		AttachStacktrace: true,
		Environment:      environment.Env,
		ServerName:       environment.WebhookAddress,
	})
	capturePanicErrorIfNotNull(sentryInitError)
	es, err := env.UnmarshalFromEnviron(&environment)
	capturePanicErrorIfNotNull(err)
	environment.Extras = es
	defer sentry.Flush(3 * time.Second)
	ekiEe := InitEki()

	log.Printf("Authorized on account %s", ekiEe.telegram.Self.UserName)
	if environment.Env != "dev" {
		_, err = ekiEe.telegram.SetWebhook(tgbotapi.NewWebhook(environment.WebhookAddress))
		capturePanicErrorIfNotNull(err)
	}

	http.HandleFunc("/"+environment.UuidToken, func(w http.ResponseWriter, r *http.Request) {
		var update tgbotapi.Update
		var plain map[string]interface{}
		var response tgbotapi.Chattable

		bytes, _ := ioutil.ReadAll(r.Body)

		json.Unmarshal(bytes, &plain)
		json.Unmarshal(bytes, &update)

		if !IsCallbackQuery(update) {
			log.WithFields(plain).Info("new_searh")
			response = ekiEe.MakeNewSearchResponse(update)
		} else {
			log.WithFields(plain).Info("callback_query")
			response = ekiEe.MakeArticleSwitchResponse(update)
		}

		if response != nil {
			_, err = ekiEe.telegram.Send(response)
			captureFatalErrorIfNotNull(err)
		}
	})

	http.ListenAndServe("0.0.0.0:"+environment.AppPort, nil)

}
