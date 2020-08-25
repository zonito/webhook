package services

import (
	"bytes"
	"context"
	"encoding/json"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	"net/url"
)

type hCRequest struct {
	Color          string `json:"color"`
	Message        string `json:"message"`
	Notify         bool   `json:"notify"`
	Message_format string `json:"message_format"`
}

// Send telegram message
func SendHipchatMessage(
	context context.Context, text string, room_id string, token string,
	color string) bool {
	uri := "https://api.hipchat.com/v2/room/" + room_id
	uri += "/notification?auth_token=" + url.QueryEscape(token)
	client := urlfetch.Client(context)
	payload := &hCRequest{
		Color:          color,
		Message:        text,
		Notify:         true,
		Message_format: "text",
	}
	str, _ := json.Marshal(payload)
	log.Infof(context, uri)
	log.Infof(context, string(str))
	resp, err := client.Post(
		uri, "application/json", bytes.NewBuffer(str))
	if err != nil {
		log.Infof(context, "Hipchat client.Post: %v", err.Error())
		return false
	}
	defer resp.Body.Close()
	log.Infof(context, "response Headers:", resp.Header)
	log.Infof(context, "HC Status: %s", resp.Status)
	if resp.Status == "204 No Content" {
		return true
	}
	return false
}
