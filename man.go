package webhook

import (
    "fmt"
    "time"
    "net/http"
    "html/template"
    "appengine"
    "appengine/user"
    "appengine/datastore"
    // "appengine/urlfetch"
    // "io/ioutil"
    // "bytes"
)

func init() {
    http.HandleFunc("/", root)
    http.HandleFunc("/hooks/{handler}", hookHandler)
}

var listTmpl = template.Must(
  template.ParseFiles("templates/base.html", "templates/list.html"))

type Webhook struct {
	Handler string
	Email string
  AccessToken string
  BoardId string
  BoardName string
  ListId string
  ListName string
	Date time.Time
}

func webhookKey(context appengine.Context) *datastore.Key {
	return datastore.NewKey(context, "Webhook", "default_webhook", 0, nil)
}

func root(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
  appUser := user.Current(context)
	query := datastore.NewQuery("Webhook").Ancestor(
		webhookKey(context)).Filter(
      "Email =", appUser.Email).Order("-Date").Limit(10)
	webhooks := make([]Webhook, 0, 10)
	if _, err := query.GetAll(context, &webhooks); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := listTmpl.Execute(writer, webhooks); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func hookHandler(writer http.ResponseWriter, request *http.Request) {
  context := appengine.NewContext(request)
  context.Infof("Handler: %v", request.URL)
  fmt.Fprintf(writer, "Ur: %v", request.URL)
}

// func sign(writer http.ResponseWriter, request *http.Request) {
// 	context := appengine.NewContext(request)
// 	greeting := Greeting{
// 		Content: request.FormValue("content"),
// 		Date: time.Now(),
// 	}
// 	if appUser := user.Current(context); appUser != nil {
// 		greeting.Author = appUser.String()
// 	}
// 	key := datastore.NewIncompleteKey(
// 		context, "Greeting", guestbookKey(context))
// 	_, err := datastore.Put(context, key, &greeting)
// 	if err != nil {
// 		http.Error(writer, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	http.Redirect(writer, request, "/", http.StatusFound)
// }

// func tpl(writer http.ResponseWriter, request *http.Request) {
// 	fmt.Fprintf(writer, guestbookForm)
// }

// func show(writer http.ResponseWriter, request *http.Request) {
// 	err := signTemplate.Execute(writer, request.FormValue("content"))
// 	if err != nil {
// 		http.Error(writer, err.Error(), http.StatusInternalServerError)
// 	}
// }

// func handler(writer http.ResponseWriter, request *http.Request) {
//     context := appengine.NewContext(request)
//     app_user := user.Current(context)
//     if app_user == nil {
//         url, err := user.LoginURL(context, request.URL.String())
//         if err != nil {
//             http.Error(writer, err.Error(), http.StatusInternalServerError)
//             return
//         }
//         writer.Header().Set("Location", url)
//         writer.WriteHeader(http.StatusFound)
//         return
//     }
//     fmt.Fprintf(writer, "Hello, %v!", app_user)
// }

// func send_request(writer http.ResponseWriter, request *http.Request) {
//     context := appengine.NewContext(request)
//     client := urlfetch.Client(context)
//     url := ""
//     var jsonStr = []byte(`{"name":"Buy cheese and bread for breakfast."}`)
//     resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonStr))
//     if err != nil {
//         http.Error(writer, err.Error(), http.StatusInternalServerError)
//         return
//     }
//     defer resp.Body.Close()
//     fmt.Fprintf(writer, "URL:>", url)
//     fmt.Fprintf(writer, "response Status:", resp.Status)
//     fmt.Fprintf(writer, "response Headers:", resp.Header)
//     body, _ := ioutil.ReadAll(resp.Body)
//     fmt.Fprintf(writer, "response Body:", string(body))
// }
