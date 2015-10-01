// Main file

package webhook

import (
    "appengine"
    "appengine/datastore"
    "appengine/user"
    "encoding/json"
    "fmt"
    "github.com/gorilla/mux"
    "html/template"
    "net/http"
    "services"
    "strings"
    "time"
)

var redirectTmpl = template.Must(
    template.ParseFiles("templates/callback.html"))

// Initialize Appengine.
// Only routes, that's it!
func init() {
    route := mux.NewRouter()
    route.HandleFunc("/", root)
    route.HandleFunc("/login", login)
    route.HandleFunc("/cb", callback)
    route.HandleFunc("/connect", connect)
    route.HandleFunc("/created.json", createdJson)
    route.HandleFunc("/w{handler}", hooks)
    route.HandleFunc("/redirect", redirect)
    route.HandleFunc("/save", save)
    route.HandleFunc("/telegram/{telegramToken}", telegramWebhook)
    route.HandleFunc("/trello/{type}/{boardid}", trelloList)
    http.Handle("/", route)
}

// Return list of created webhooks (/created.json)
func createdJson(writer http.ResponseWriter, request *http.Request) {
    context := appengine.NewContext(request)
    appUser := user.Current(context)
    webhooks := getWebhooks(context, appUser.Email)
    list, _ := json.Marshal(webhooks)
    writer.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(writer, string(list))
}

// Redirect use to get trello service approval. (/connect)
func login(writer http.ResponseWriter, request *http.Request) {
    http.Redirect(writer, request, "/", http.StatusFound)
}

// Root handler (/), show for to create new and list of created hooks.
func root(writer http.ResponseWriter, request *http.Request) {
    context := appengine.NewContext(request)
    appUser := user.Current(context)
    if appUser != nil {
        listTmpl := template.Must(
            template.ParseFiles("templates/base.html", "templates/list.html"))
        url, _ := user.LogoutURL(context, "/")
        data := struct {
            AccessToken string
            Logout      string
        }{getAccessToken(context, appUser.Email), url}
        if err := listTmpl.Execute(writer, data); err != nil {
            http.Error(writer, err.Error(), http.StatusInternalServerError)
        }
    } else {
        homeTmpl := template.Must(template.ParseFiles("templates/index.html"))
        homeTmpl.Execute(writer, nil)
    }
}

// Redirect use to get trello service approval. (/connect)
func connect(writer http.ResponseWriter, request *http.Request) {
    http.Redirect(writer, request, services.GetAuthorizeUrl(), http.StatusFound)
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
    writer.Header().Set("Content-Type", "application/json")
    accessToken := getAccessToken(context, appUser.Email)
    if vars["type"] == "lists" {
        fmt.Fprintf(
            writer, services.GetBoardLists(
                context, vars["boardid"], accessToken))
        return
    }
    fmt.Fprintf(writer, services.GetBoards(context, accessToken))
}

// Save new hook from web.
func save(writer http.ResponseWriter, request *http.Request) {
    context := appengine.NewContext(request)
    appUser := user.Current(context)
    response := Response{
        Success: true,
        Reason:  "",
    }
    handler := "w" + services.GetAlphaNumberic(7)
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
        services.PushToTrello(
            context, webhook.ListId,
            getAccessToken(context, webhook.User), "You are connected!", "")
    } else if request.FormValue("service") == "telegram" {
        webhook.Type = "Telegram"
        webhook.TeleChatId, webhook.TeleChatName = services.GetChatIdFromCode(
            context, request.FormValue("tele_code"))
        if webhook.TeleChatId == 0 {
            response.Success = false
            response.Reason = "Invalid code."
        } else {
            services.SendTeleMessage(
                context, "You are connected!", webhook.TeleChatId)
        }
    } else if request.FormValue("service") == "pushover" {
        webhook.Type = "Pushover"
        webhook.POUserKey = request.FormValue("po_userkey")
        status := services.SendPushoverMessage(
            context, "You are connected!", webhook.POUserKey)
        if status == 0 {
            response.Success = false
            response.Reason = "Invalid key."
        }
    } else if request.FormValue("service") == "hipchat" {
        webhook.Type = "Hipchat"
        webhook.HCToken = request.FormValue("hc_token")
        webhook.HCRoomId = request.FormValue("hc_roomid")
        status := services.SendHipchatMessage(
            context, "You are connected!", webhook.HCRoomId,
            webhook.HCToken, "green")
        if !status {
            response.Success = false
            response.Reason = "Invalid room id or token."
        }
    }
    if response.Success {
        key := datastore.NewIncompleteKey(
            context, "Webhook", webhookKey(context, handler))
        _, err := datastore.Put(context, key, &webhook)
        if err != nil {
            context.Infof("%v", err.Error())
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
    context := appengine.NewContext(request)
    decoder := json.NewDecoder(request.Body)
    fmt.Fprintf(
        writer, services.Telegram(context, decoder, vars["telegramToken"]))
}

// Actual webhook handler, receive events and post it to connected services.
func hooks(writer http.ResponseWriter, request *http.Request) {
    vars := mux.Vars(request)
    handler := "w" + vars["handler"]
    context := appengine.NewContext(request)
    webhook := getWebhookFromHandler(context, handler)
    if webhook != nil {
        event, desc := services.GetEventData(request)
        context.Infof("%s: %s \n %s", webhook.Type, event, desc)
        if event != "" {
            if webhook.Type == "Trello" {
                services.PushToTrello(
                    context, webhook.ListId,
                    getAccessToken(context, webhook.User), event, desc)
            } else if webhook.Type == "Telegram" {
                services.SendTeleMessage(
                    context, event+"\n"+desc, webhook.TeleChatId)
            } else if webhook.Type == "Pushover" {
                services.SendPushoverMessage(
                    context, event+"\n"+desc, webhook.POUserKey)
            } else if webhook.Type == "Hipchat" {
                color := "red"
                if strings.Index(event, " success ") > -1 ||
                    strings.Index(event, " merged ") > -1 ||
                    strings.Index(event, ": up ") > -1 ||
                    strings.Index(event, "Ping!") > -1 {
                    color = "green"
                } else if strings.Index(event, " pull ") > -1 {
                    color = "yellow"
                }
                services.SendHipchatMessage(
                    context, event+"\n"+desc, webhook.HCRoomId,
                    webhook.HCToken, color)
            }
        }
        fmt.Fprintf(writer, "OK")
    }
}
