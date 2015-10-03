// Telegram Integration (https://telegram.org/).
// This is service to interact with / to notify.

package services

import (
    "appengine"
    "appengine/datastore"
    "appengine/urlfetch"
    "encoding/json"
    "net/url"
    "strconv"
    "strings"
    "time"
)

const apiURL = "https://api.telegram.org/bot" + teleToken + "/sendMessage?"

// model

type TeleVerify struct {
    ChatId int
    Code   string
    Date   time.Time
    Name   string
}

// Telegram

type TeleUser struct {
    Id         int
    First_name string
    Last_name  string
    Username   string
    Title      string
}

type TeleMessage struct {
    Message_id     int
    Date           int
    Text           string
    From           TeleUser
    Chat           TeleUser
    New_chat_title string
}

type TelePayload struct {
    Update_id int
    Message   TeleMessage
}

// Return TeleVerify datastore key.
func teleVerifyKey(context appengine.Context, code string) *datastore.Key {
    return datastore.NewKey(context, "TeleVerify", code, 0, nil)
}

// Return Chat id from Code
func GetChatIdFromCode(context appengine.Context, code string) (int, string) {
    query := datastore.NewQuery("TeleVerify").Ancestor(
        teleVerifyKey(context, code)).Limit(1)
    teleVerify := make([]TeleVerify, 0, 1)
    keys, _ := query.GetAll(context, &teleVerify)
    if len(teleVerify) > 0 {
        chatId, name := teleVerify[0].ChatId, teleVerify[0].Name
        datastore.Delete(context, keys[0])
        return chatId, name
    }
    return 0, ""
}

// Send telegram message
func SendTeleMessage(context appengine.Context, text string, chat_id int) {
    uri := apiURL + "?parse_mode=Markdown&disable_web_page_preview=true"
    uri += "&chat_id=" + strconv.Itoa(chat_id)
    uri += "&text=" + url.QueryEscape(text)
    client := urlfetch.Client(context)
    resp, _ := client.Get(uri)
    defer resp.Body.Close()
}

// Telegram webhook
func Telegram(
    context appengine.Context, decoder *json.Decoder, token string) string {
    if token != teleToken {
        return "NOT OK"
    }
    var teleEvent TelePayload
    decoder.Decode(&teleEvent)
    message := teleEvent.Message
    if strings.Index(message.Text, "/getcode") > -1 {
        code := GetAlphaNumberic(6)
        teleVerify := TeleVerify{
            ChatId: message.Chat.Id,
            Code:   code,
            Date:   time.Now(),
            Name:   message.Chat.First_name,
        }
        if message.Chat.Id < 0 {
            teleVerify.Name = message.Chat.Title
        }
        key := datastore.NewIncompleteKey(
            context, "TeleVerify", teleVerifyKey(context, code))
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
