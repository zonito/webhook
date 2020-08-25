package services

import (
	"encoding/json"
)

// Doorbell structs

type DBApplication struct {
	Name string
}

type DBData struct {
	Email       string
	Url         string
	member      User
	User_Agent  string
	Message     string
	Sentiment   string
	Application DBApplication
	Created     string
}

type DBPayload struct {
	Event string
	Data  DBData
}

// Return doorbell data.
func getDoorbellData(decoder *json.Decoder) (string, string) {
	var dEvent DBPayload
	decoder.Decode(&dEvent)
	data := dEvent.Data
	event := data.Application.Name + " --> " +
		data.Sentiment + " feedback - from " + data.Email
	desc := data.Message + "\n\n User Agent: " +
		data.User_Agent + "\n\n Reply: " + data.Url
	return event, desc
}
