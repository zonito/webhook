package services

import (
	"bytes"
	"context"
	"encoding/json"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	"net/url"
)

// Send push over notification message
func SendPushoverMessage(
	context context.Context, message string, userKey string) int {
	apiUrl := "https://api.pushover.net/1/messages.json"
	parameters := url.Values{}
	parameters.Add("token", pushoverKey)
	parameters.Add("user", userKey)
	parameters.Add("message", message)
	parameters.Add("priority", "1")
	client := urlfetch.Client(context)
	resp, _ := client.Post(
		apiUrl, "application/x-www-form-urlencoded",
		bytes.NewBuffer([]byte(parameters.Encode())))
	defer resp.Body.Close()
	log.Infof(context, "response Headers:", resp.Header)
	decoder := json.NewDecoder(resp.Body)
	pEvent := struct {
		Status int
	}{0}
	decoder.Decode(&pEvent)
	return pEvent.Status
}
