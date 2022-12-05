package main

import (
	"fmt"
	"net/http"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	type MailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var reqPayload MailMessage
	err := app.readJSON(w, r, &reqPayload)
	if err != nil {
		fmt.Println(err)
		app.errorJSON(w, err)
	}
	msg := Message{
		From:    reqPayload.From,
		To:      reqPayload.To,
		Subject: reqPayload.Subject,
		Data:    reqPayload.Message,
	}
	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		fmt.Println("error", err)
		app.errorJSON(w, err)
	}

	payload := jsonResponse{
		Error:   false,
		Message: "sent to" + reqPayload.To,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
