package services

import (
	"encoding/json"
	"strconv"
)

type FBPayloadApp struct {
	Name              string
	Bundle_identifier string
	Platform          string
}

type FBPayload struct {
	Display_id             int
	Title                  string
	Method                 string
	Crashes_count          int
	Impacted_devices_count int
	Impact_level           int
	Url                    string
	App                    FBPayloadApp
}

type FBMessage struct {
	Event        string
	Payload_type string
	Payload      FBPayload
}

func getFabricData(decoder *json.Decoder) (string, string) {
	var fbEvent FBMessage
	decoder.Decode(&fbEvent)
	event := "`Crashlytics: " + fbEvent.Payload_type + ", " +
		fbEvent.Payload.Title + " for " + fbEvent.Payload.App.Bundle_identifier + "`"
	payload := fbEvent.Payload
	desc := fbEvent.Event +
		"\nCrashes Count: " + strconv.Itoa(payload.Crashes_count) +
		"\nPlatform: " + payload.App.Platform +
		"\nName: " + payload.App.Name +
		"\nMethod: " + payload.Method +
		"\nURL: " + payload.Url +
		"\nMethod: " + payload.Method +
		"\nImpacted Devices Count: " + strconv.Itoa(payload.Impacted_devices_count) +
		"\nImpacted Level: " + strconv.Itoa(payload.Impact_level)
	return event, desc
}
