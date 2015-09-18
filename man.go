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


/***
  * Initialize Appengine.
  * Only routes, that's it!
  */
func init() {
    route := mux.NewRouter()
    route.HandleFunc("/", root)
    route.HandleFunc("/cb", callback)
    route.HandleFunc("/connect", connect)
    route.HandleFunc("/hooks/{handler}", hooks)
    route.HandleFunc("/redirect", redirect)
    route.HandleFunc("/save", save)
    route.HandleFunc("/trello/{type}/{boardid}", trelloList)
    http.Handle("/", route)
}


/***
  * Root handler (/), show for to create new and list of created hooks.
  */
func root(writer http.ResponseWriter, request *http.Request) {
    context := appengine.NewContext(request)
    appUser := user.Current(context)
    data := struct {
        AccessToken string
        WH []Webhook
    }{getAccessToken(context, appUser.Email),
      getWebhooks(context, appUser.Email)}
    if err := listTmpl.Execute(writer, data); err != nil {
        http.Error(writer, err.Error(), http.StatusInternalServerError)
    }
}


/***
  * Actual webhook handler, receive events and post it to connected services.
  */
func hooks(writer http.ResponseWriter, request *http.Request) {
    vars := mux.Vars(request)
    handler := vars["handler"]
    context := appengine.NewContext(request)
    webhook := getWebhookFromHandler(context, handler)
    if webhook != nil {
        decoder := json.NewDecoder(request.Body)
        hookType := getHookType(request)
        var event, desc string
        if hookType == "github" {
          event, desc = getGithubData(
            decoder, request.Header.Get("X-Github-Event"))
        } else if hookType == "doorbell" {
          event, desc = getDoorbellData(decoder)
        }
        url := "https://api.trello.com/1/lists/" + webhook.ListId +
          "/cards?key=" + trelloKey + "&token=" +
          getAccessTokenFromHandler(context, handler)
        payload := &TrelloPayLoad {
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
            http.Error(writer, err.Error(), http.StatusInternalServerError)
            return
        }
        defer resp.Body.Close()
        context.Infof("response Headers:", resp.Header)
        body, _ := ioutil.ReadAll(resp.Body)
        context.Infof("response Body:", string(body))
        fmt.Fprintf(writer, "OK")
    }
}


/***
  * Redirect use to get trello service approval.
  */
func connect(writer http.ResponseWriter, request *http.Request) {
    authorizeUrl :=
      "https://trello.com/1/OAuthAuthorizeToken" +
      "?key=" + trelloKey + "&callback_method=fragment&scope=read,write" +
      "&name=PGWebhook&scope=read,write&expiration=never" +
      "&return_url=http://pgwebhook.appspot.com/redirect"
    http.Redirect(writer, request, authorizeUrl, http.StatusFound)
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


/***
  * Callback with token in post payload.
  */
func callback(writer http.ResponseWriter, request *http.Request) {
    context := appengine.NewContext(request)
    appUser := user.Current(context)
    accessToken := AccessTokens {
        Email: appUser.Email,
        AccessToken: request.FormValue("token"),
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


/***
  * Get list of trello boards or lists.
  */
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


/***
  * Save new hook from web.
  */
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
