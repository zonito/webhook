package services

import (
	"encoding/json"
)

type SDIncident struct {
	Incident_id    string
	Resource_id    string
	Resource_name  string
	State          string
	Started_at     int
	Ended_at       int
	Policy_name    string
	Condition_name string
	Url            string
	Summary        string
}

type SDMessage struct {
	Incident SDIncident
	Version  int
}

func getStackDriverData(decoder *json.Decoder) (string, string) {
	var sdnEvent SDMessage
	decoder.Decode(&sdnEvent)
	event := "StackDriver: " + sdnEvent.Incident.Policy_name +
		", Condition: " + sdnEvent.Incident.Condition_name
	desc := "URL: " + sdnEvent.Incident.Url +
		"\nSummary: " + sdnEvent.Incident.Summary +
		"\nState: " + sdnEvent.Incident.State +
		"\nResource: " + sdnEvent.Incident.Resource_name
	return event, desc
}
