package services

import (
	"encoding/json"
)

type C1Message struct {
	Message string
}

func getCustom1Data(decoder *json.Decoder) (string, string) {
	var cEvent C1Message
	decoder.Decode(&cEvent)
	event := cEvent.Message
	return event, ""
}
