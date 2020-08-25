package services

import (
	"encoding/json"
)

type CMessage struct {
	Message string `json:"message"`
}

func getCustomData(decoder *json.Decoder) (string, string) {
	var cEvent CMessage
	decoder.Decode(&cEvent)
	event := cEvent.Message
	return event, ""
}
