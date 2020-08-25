package services

import (
	"encoding/json"
)

// Grafana structs

type EvalMatch struct {
	Value  int
	Metric string
	Tags   []string
}

type GrafanaPayload struct {
	EvalMatches []EvalMatch
	ImageUrl    string
	Message     string
	RuleId      int
	RuleName    string
	RuleUrl     string
	State       string
	Title       string
}

// Return grafana data.
func getGrafanaData(decoder *json.Decoder) (string, string) {
	var gfEvent GrafanaPayload
	decoder.Decode(&gfEvent)
	return gfEvent.Title, gfEvent.Message
}
