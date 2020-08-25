/*
Package services provides Telegram Integration (https://telegram.org/).
This is service to interact with / to notify.
*/
package services

import (
	"bytes"
	"context"
	"encoding/json"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

type slackRequest struct {
	Text    string `json:"text"`
	Channel string `json:"channel"`
}

// SendTeleMessage Send telegram message
func SendSlackMessage(
	context context.Context, text string, slack_url string,
	channel string) bool {
	log.Infof(context, "%s", slack_url)
	client := urlfetch.Client(context)
	payload := &slackRequest{
		Text:    text,
		Channel: channel,
	}
	str, _ := json.Marshal(payload)
	log.Infof(context, slack_url)
	log.Infof(context, string(str))
	resp, err := client.Post(
		slack_url, "application/json", bytes.NewBuffer(str))
	if err != nil {
		log.Infof(context, "Slack client.Post: %v", err.Error())
		return false
	}
	defer resp.Body.Close()
	log.Infof(context, "response Headers:", resp.Header)
	log.Infof(context, "Slack Status: %s", resp.Status)
	if resp.Status == "200 OK" {
		return true
	}
	return false
}
