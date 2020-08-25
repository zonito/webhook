/*
Package services provides Telegram Integration (https://telegram.org/).
This is service to interact with / to notify.
*/
package services

import (
	"context"
	"encoding/json"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var apiURL = "https://api.telegram.org/bot" + teleToken + "/sendMessage"

// TeleVerify model is a temporary database to store code from bot.
type teleVerify struct {
	ChatId int
	Code   string
	Date   time.Time
	Name   string
}

// teleUser is part of teleMessage to know, who messaged.
type teleUser struct {
	Id         int
	First_name string
	Last_name  string
	Username   string
	Title      string
}

// teleMessage is a message sent in telegram.
type teleMessage struct {
	Message_id     int
	Date           int
	Text           string
	From           teleUser
	Chat           teleUser
	New_chat_title string
}

// telePayload is a request body from telegram.
type telePayload struct {
	Update_id int
	Message   teleMessage
}

// Return TeleVerify datastore key.
func teleVerifyKey(context context.Context, code string) *datastore.Key {
	return datastore.NewKey(context, "teleVerify", code, 0, nil)
}

// GetChatIdFromCode Return Chat id from Code
func GetChatIdFromCode(context context.Context, code string) (int, string) {
	query := datastore.NewQuery("teleVerify").Ancestor(
		teleVerifyKey(context, code)).Limit(1)
	teleVerify := make([]teleVerify, 0, 1)
	keys, _ := query.GetAll(context, &teleVerify)
	if len(teleVerify) > 0 {
		chatId, name := teleVerify[0].ChatId, teleVerify[0].Name
		datastore.Delete(context, keys[0])
		return chatId, name
	}
	return 0, ""
}

// SendTeleMessage Send telegram message
func SendTeleMessage(context context.Context, text string, chat_id int) {
	uri := apiURL + "?parse_mode=Markdown&disable_web_page_preview=true"
	uri += "&chat_id=" + strconv.Itoa(chat_id)
	uri += "&text=" + url.QueryEscape(text)
	log.Infof(context, "%s", uri)
	client := urlfetch.Client(context)
	resp, _ := client.Get(uri)
	defer resp.Body.Close()
}

// Telegram webhook
func Telegram(
	context context.Context, decoder *json.Decoder, token string) string {
	if token != teleToken {
		return "!OK"
	}
	var teleEvent telePayload
	decoder.Decode(&teleEvent)
	message := teleEvent.Message
	if strings.Index(message.Text, "/getcode") > -1 {
		code := GetAlphaNumberic(6)
		teleVerify := teleVerify{
			ChatId: message.Chat.Id,
			Code:   code,
			Date:   time.Now(),
			Name:   message.Chat.First_name,
		}
		if message.Chat.Id < 0 {
			teleVerify.Name = message.Chat.Title
		}
		key := datastore.NewIncompleteKey(
			context, "teleVerify", teleVerifyKey(context, code))
		datastore.Put(context, key, &teleVerify)
		SendTeleMessage(context, code, message.Chat.Id)
	} else if strings.Index(message.Text, "/start") > -1 {
		SendTeleMessage(
			context, "Welcome! Next step is to get registered with webhook.co",
			message.Chat.Id)
	} else if strings.Index(message.Text, "/help") > -1 {
		SendTeleMessage(
			context, "Get registered with webhook.co", message.Chat.Id)
	}
	return "OK"
}
