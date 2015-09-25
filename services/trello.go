package services

import (
    "appengine"
    "appengine/urlfetch"
    "bytes"
    "encoding/json"
    "io/ioutil"
    "strings"
)

const trelloApiURL = "https://api.trello.com/1/"

// To send as trello api payload.
type TrelloPayLoad struct {
    Name string
    Desc string
}

// Push to Trello
func PushToTrello(
    context appengine.Context, listId string, accessToken string,
    event string, desc string) {
    url := "https://api.trello.com/1/lists/" + listId +
        "/cards?key=" + trelloKey + "&token=" + accessToken
    payload := &TrelloPayLoad{
        Name: event,
        Desc: string(desc),
    }
    str, _ := json.Marshal(payload)
    jsonStr := strings.Replace(string(str), "Name", "name", 1)
    jsonStr = strings.Replace(jsonStr, "Desc", "desc", 1)
    client := urlfetch.Client(context)
    resp, err := client.Post(
        url, "application/json", bytes.NewBuffer([]byte(jsonStr)))
    if err != nil {
        context.Infof("PushToTrello client.Post: %v", err.Error())
        return
    }
    defer resp.Body.Close()
    context.Infof("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    context.Infof("response Body:", string(body))
}

// Return Trello Key
func GetAuthorizeUrl() string {
    return "https://trello.com/1/OAuthAuthorizeToken" +
        "?key=" + trelloKey + "&callback_method=fragment&scope=read,write" +
        "&name=PGWebhook&scope=read,write&expiration=never" +
        "&return_url=http://webhook.co/redirect"
}

// Get list of borads
func GetBoards(context appengine.Context, accessToken string) string {
    url := trelloApiURL + "members/me/boards?fields=name&key=" + trelloKey + "&token=" +
        accessToken
    return getResponse(context, url)
}

// Get list of borads
func GetBoardLists(
    context appengine.Context, boardId string, accessToken string) string {
    url := trelloApiURL + "boards/" + boardId + "/lists?fields=name&key=" +
        trelloKey + "&token=" + accessToken
    return getResponse(context, url)
}
