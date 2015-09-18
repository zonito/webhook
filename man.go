package webhook

import (
    "appengine"
    "appengine/datastore"
    "appengine/urlfetch"
    "appengine/user"
    "bytes"
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

func init() {
    route := mux.NewRouter()
    route.HandleFunc("/", root)
    route.HandleFunc("/cb", callback)
    route.HandleFunc("/connect", connect)
    route.HandleFunc("/hooks/{handler}", hookHandler)
    route.HandleFunc("/redirect", redirect)
    route.HandleFunc("/save", save)
    route.HandleFunc("/trello/{type}/{boardid}", trelloList)
    http.Handle("/", route)
}

func root(writer http.ResponseWriter, request *http.Request) {
    context := appengine.NewContext(request)
    appUser := user.Current(context)
    data := struct {
        AccessToken string
        WH []Webhook
    }{getAccessToken(context, appUser.Email),
        getWebhooks(context, appUser.Email)}
    context.Infof("Token: %v", data.AccessToken)
    if err := listTmpl.Execute(writer, data); err != nil {
        http.Error(writer, err.Error(), http.StatusInternalServerError)
    }
}

func hookHandler(writer http.ResponseWriter, request *http.Request) {
    vars := mux.Vars(request)
    handler := vars["handler"]
    context := appengine.NewContext(request)
    context.Infof("Handler: %v", handler)
    fmt.Fprintf(writer, "Url: %v", request.URL)
    event := request.Header.Get("X-Github-Event")
    event += " Event"
    desc, _ := ioutil.ReadAll(request.Body)
    context.Infof("Body: %s", desc)
    webhook := getWebhookFromHandler(context, handler)
    if webhook != nil {
        context.Infof("id: %v", webhook.ListId)
        var url bytes.Buffer
        url.WriteString("https://api.trello.com/1/lists/")
        url.WriteString(webhook.ListId)
        url.WriteString("/cards?key=")
        url.WriteString(trelloKey)
        url.WriteString("&token=")
        url.WriteString(getAccessTokenFromHandler(context, handler))
        context.Infof("url: %v", url.String())
        type PayLoad struct {
            Name string
            Desc string
        }
        context.Infof("des: ", string(desc))
        payload := &PayLoad {
            Name: event,
            Desc: string(desc),
        }
        str, err := json.Marshal(payload)
        if err != nil {
            context.Infof("Error", err)
            return
        }
        jsonStr := string(str)
        jsonStr = strings.Replace(jsonStr, "Name", "name", 1)
        jsonStr = strings.Replace(jsonStr, "Desc", "desc", 1)
        context.Infof("json:", jsonStr)
        client := urlfetch.Client(context)
        resp, err := client.Post(
            url.String(), "application/json", bytes.NewBuffer([]byte(jsonStr)))
        if err != nil {
            http.Error(writer, err.Error(), http.StatusInternalServerError)
            return
        }
        defer resp.Body.Close()
        fmt.Fprintf(writer, "response Status:", resp.Status)
        fmt.Fprintf(writer, "response Headers:", resp.Header)
        body, _ := ioutil.ReadAll(resp.Body)
        fmt.Fprintf(writer, "response Body:", string(body))
    }
}

func connect(writer http.ResponseWriter, request *http.Request) {
    var buffer bytes.Buffer
    buffer.WriteString("https://trello.com/1/OAuthAuthorizeToken?key=")
    buffer.WriteString(trelloKey)
    buffer.WriteString("&callback_method=fragment&scope=read,write")
    buffer.WriteString("&expiration=never&return_url=")
    buffer.WriteString("http://pgwebhook.appspot.com/redirect")
    http.Redirect(writer, request, buffer.String(), http.StatusFound)
}


/***
  * Once approval from service is done, read the token, make post request
  * to callback handler (/cb) to save token.
  */
func redirect(writer http.ResponseWriter, request *http.Request) {
    if err := redirectTmpl.Execute(writer, nil); err != nil {
        http.Error(writer, err.Error(), http.StatusInternalServerError)
    }
}

func callback(writer http.ResponseWriter, request *http.Request) {
    context := appengine.NewContext(request)
    url := request.URL.String()
    index := strings.Index(url, "#")
    token := request.FormValue("token")
    fmt.Fprintf(writer, "URL: %v, %d", url, index)
    context.Infof("Token: %v", token)
    appUser := user.Current(context)
    accessToken := AccessTokens {
        Email: appUser.Email,
        AccessToken: token,
    }
    key := datastore.NewIncompleteKey(
        context, "AccessTokens", accessTokenKey(context))
    _, err := datastore.Put(context, key, &accessToken)
    if err != nil {
        http.Error(writer, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(writer, request, "/", http.StatusFound)
}

func trelloList(writer http.ResponseWriter, request *http.Request) {
    vars := mux.Vars(request)
    context := appengine.NewContext(request)
    appUser := user.Current(context)
    client := urlfetch.Client(context)
    context.Infof("URL: %v", request.URL)
    var buffer bytes.Buffer
    url := "https://trello.com/1/members/me/boards"
    if vars["type"] == "lists" {
        boardId := vars["boardid"]
        var urlBuffer bytes.Buffer
        urlBuffer.WriteString("https://api.trello.com/1/boards/")
        urlBuffer.WriteString(boardId)
        urlBuffer.WriteString("/lists")
        url = urlBuffer.String()
    }
    buffer.WriteString(url)
    buffer.WriteString("?fields=name&key=")
    buffer.WriteString(trelloKey)
    buffer.WriteString("&token=")
    buffer.WriteString(getAccessToken(context, appUser.Email))
    context.Infof("Test: %v", buffer.String())
    resp, err := client.Get(buffer.String())
    if err != nil {
        http.Error(writer, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    context.Infof("URL:>", buffer.String())
    context.Infof("response Status:", resp.Status)
    context.Infof("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    writer.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(writer, string(body))
}

func save(writer http.ResponseWriter, request *http.Request) {
    context := appengine.NewContext(request)
    appUser := user.Current(context)
    webhook := Webhook {
        Email: appUser.Email,
        Handler: request.FormValue("handler"),
        BoardId: request.FormValue("boards"),
        ListId: request.FormValue("lists"),
        Date: time.Now(),
    }
    key := datastore.NewIncompleteKey(
        context, "Webhook", webhookKey(context))
    _, err := datastore.Put(context, key, &webhook)
    if err != nil {
        http.Error(writer, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(writer, request, "/", http.StatusFound)
}
