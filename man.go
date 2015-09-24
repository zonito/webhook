// Main file

package webhook

import (
    "appengine"
    "appengine/datastore"
    "appengine/urlfetch"
    "appengine/user"
    "encoding/json"
    "fmt"
    "github.com/gorilla/mux"
    "html/template"
    "io/ioutil"
    "net/http"
    "strings"
    "time"
)

var listTmpl = template.Must(
    template.ParseFiles("templates/base.html", "templates/list.html"))
var redirectTmpl = template.Must(
    template.ParseFiles("templates/callback.html"))

// Initialize Appengine.
// Only routes, that's it!
func init() {
    route := mux.NewRouter()
    route.HandleFunc("/", root)
    route.HandleFunc("/cb", callback)
    route.HandleFunc("/connect", connect)
    route.HandleFunc("/created.json", createList)
    route.HandleFunc("/w{handler}", hooks)
    route.HandleFunc("/redirect", redirect)
    route.HandleFunc("/save", save)
    route.HandleFunc("/telegram/{telegramToken}", telegramWebhook)
    route.HandleFunc("/trello/{type}/{boardid}", trelloList)
    http.Handle("/", route)
}

// Root handler (/), show for to create new and list of created hooks.
func root(writer http.ResponseWriter, request *http.Request) {
    context := appengine.NewContext(request)
    appUser := user.Current(context)
    url, _ := user.LogoutURL(context, "/")
    data := struct {
        AccessToken string
        Logout      string
    }{getAccessToken(context, appUser.Email), url}
    if err := listTmpl.Execute(writer, data); err != nil {
        http.Error(writer, err.Error(), http.StatusInternalServerError)
    }
}

// Return list of created webhooks
func createList(writer http.ResponseWriter, request *http.Request) {
    context := appengine.NewContext(request)
    appUser := user.Current(context)
    webhooks := getWebhooks(context, appUser.Email)
    list, _ := json.Marshal(webhooks)
    writer.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(writer, string(list))
}

// Redirect use to get trello service approval.
func connect(writer http.ResponseWriter, request *http.Request) {
    authorizeUrl :=
        "https://trello.com/1/OAuthAuthorizeToken" +
            "?key=" + trelloKey + "&callback_method=fragment&scope=read,write" +
            "&name=PGWebhook&scope=read,write&expiration=never" +
            "&return_url=http://webhook.co/redirect"
    http.Redirect(writer, request, authorizeUrl, http.StatusFound)
}

// Once approval from service is done, read the token, make post request
// to callback handler (/cb) to save token.
func redirect(writer http.ResponseWriter, request *http.Request) {
    if err := redirectTmpl.Execute(writer, nil); err != nil {
        http.Error(writer, err.Error(), http.StatusInternalServerError)
    }
}

// Callback with token in post payload.
func callback(writer http.ResponseWriter, request *http.Request) {
    context := appengine.NewContext(request)
    appUser := user.Current(context)
    accessToken := AccessTokens{
        Email:       appUser.Email,
        AccessToken: request.FormValue("token"),
    }
    key := datastore.NewIncompleteKey(
        context, "AccessTokens", accessTokenKey(context, appUser.Email))
    _, err := datastore.Put(context, key, &accessToken)
    if err != nil {
        http.Error(writer, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(writer, request, "/", http.StatusFound)
}

// Get list of trello boards or lists.
func trelloList(writer http.ResponseWriter, request *http.Request) {
    vars := mux.Vars(request)
    context := appengine.NewContext(request)
    appUser := user.Current(context)
    client := urlfetch.Client(context)
    url := "https://api.trello.com/1/members/me/boards"
    if vars["type"] == "lists" {
        url = "https://api.trello.com/1/boards/" + vars["boardid"] + "/lists"
    }
    url += "?fields=name&key=" + trelloKey + "&token=" +
        getAccessToken(context, appUser.Email)
    resp, err := client.Get(url)
    if err != nil {
        http.Error(writer, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    context.Infof("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    writer.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(writer, string(body))
}

// Save new hook from web.
func save(writer http.ResponseWriter, request *http.Request) {
    context := appengine.NewContext(request)
    appUser := user.Current(context)
    response := Response{
        Success: true,
        Reason:  "",
    }
    handler := "w" + getsrand(7)
    webhook := Webhook{
        User:    appUser.Email,
        Handler: handler,
        Date:    time.Now(),
        Count:   0,
    }
    if request.FormValue("service") == "trello" {
        webhook.Type = "Trello"
        webhook.BoardId = request.FormValue("board_id")
        webhook.BoardName = request.FormValue("board_name")
        webhook.ListId = request.FormValue("list_id")
        webhook.ListName = request.FormValue("list_name")
    } else if request.FormValue("service") == "telegram" {
        webhook.Type = "Telegram"
        webhook.TeleChatId, webhook.TeleChatName = getChatIdFromCode(
            context, request.FormValue("tele_code"))
        if webhook.TeleChatId == 0 {
            response.Success = false
            response.Reason = "Invalid code."
        } else {
            sendTeleMessage(
                context, "You are connected!", webhook.TeleChatId)
        }
    }
    if response.Success {
        key := datastore.NewIncompleteKey(
            context, "Webhook", webhookKey(context, handler))
        _, err := datastore.Put(context, key, &webhook)
        if err != nil {
            http.Error(writer, err.Error(), http.StatusInternalServerError)
            return
        }
        response.Handler = handler
    }
    writer.Header().Set("Content-Type", "application/json")
    resp, _ := json.Marshal(response)
    fmt.Fprintf(writer, string(resp))
}

// Telegram webhook
func telegramWebhook(writer http.ResponseWriter, request *http.Request) {
    vars := mux.Vars(request)
    if vars["telegramToken"] != teleToken {
        fmt.Fprintf(writer, "NOT OK")
        return
    }
    context := appengine.NewContext(request)
    decoder := json.NewDecoder(request.Body)
    var teleEvent TelePayload
    decoder.Decode(&teleEvent)
    message := teleEvent.Message
    if strings.Index(message.Text, "/getcode") > -1 {
        code := getsrand(6)
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
        sendTeleMessage(context, code, message.Chat.Id)
    } else if strings.Index(message.Text, "/start") > -1 {
        sendTeleMessage(
            context, "Welcome! Next step is to get registered with webhook.co",
            message.Chat.Id)
    } else if strings.Index(message.Text, "/help") > -1 {
        sendTeleMessage(
            context, "Get registered with webhook.co", message.Chat.Id)
    }
    fmt.Fprintf(writer, "OK")
}

// Actual webhook handler, receive events and post it to connected services.
func hooks(writer http.ResponseWriter, request *http.Request) {
    vars := mux.Vars(request)
    handler := "w" + vars["handler"]
    context := appengine.NewContext(request)
    webhook := getWebhookFromHandler(context, handler)
    if webhook != nil {
        event, desc := getEventData(request)
        if event != "" {
            if webhook.Type == "Trello" {
                pushToTrello(context, webhook, event, desc)
            } else if webhook.Type == "Telegram" {
                sendTeleMessage(
                    context, event+"\n\n"+desc, webhook.TeleChatId)
            }
        }
        fmt.Fprintf(writer, "OK")
    }
}
