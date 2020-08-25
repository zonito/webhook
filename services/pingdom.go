package services

import (
	"encoding/json"
	"strconv"
)

type PDMessage struct {
	Action      string
	Check       string
	Checkname   string
	Description string
	Host        string
	IncidentId  int
}

func getPingdomData(decoder *json.Decoder) (string, string) {
	var pdEvent PDMessage
	decoder.Decode(&pdEvent)
	event := "Pingdom: " + pdEvent.Host + " " + pdEvent.Description
	desc := "More Info\nCheckname: " + pdEvent.Checkname +
		"\nCheck: " + pdEvent.Check +
		"\nIncident Id: " + strconv.Itoa(pdEvent.IncidentId) +
		"\nAction: " + pdEvent.Action
	return event, desc
}
