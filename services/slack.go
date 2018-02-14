/*
Package services provides Telegram Integration (https://telegram.org/).
This is service to interact with / to notify.
*/
package services

import (
    "appengine"
    "appengine/urlfetch"
    "bytes"
    "encoding/json"
)

type slackRequest struct {
    Text           string `json:"text"`
    Channel        string `json:"channel"`
}

// SendTeleMessage Send telegram message
func SendSlackMessage(
    context appengine.Context, text string, slack_url string,
    channel string) bool {
    context.Infof("%s", slack_url)
    client := urlfetch.Client(context)
    payload := &slackRequest{
        Text:          text,
        Channel:       channel,
    }
    str, _ := json.Marshal(payload)
    context.Infof(slack_url)
    context.Infof(string(str))
    resp, err := client.Post(
        slack_url, "application/json", bytes.NewBuffer(str))
    if err != nil {
        context.Infof("Slack client.Post: %v", err.Error())
        return false
    }
    defer resp.Body.Close()
    context.Infof("response Headers:", resp.Header)
    context.Infof("Slack Status: %s", resp.Status)
    if resp.Status == "200 OK" {
        return true
    }
    return false
}
