package main

import (
	"fmt"
	"net/http"

	"github.com/loger-service/data"
)

type JsonPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var reqPayload JsonPayload

	_ = app.readJSON(w, r, &reqPayload)
	fmt.Println("request payload is", reqPayload)

	event := &data.LogEntry{
		Name: reqPayload.Name,
		Data: reqPayload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	fmt.Println("error inserting event:", err)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}
