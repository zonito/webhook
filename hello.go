package webhook

import (
    "fmt"
    "time"
    "net/http"
    "html/template"
    "appengine"
    "appengine/user"
    "appengine/datastore"
)

func init() {
    http.HandleFunc("/", root)
    http.HandleFunc("/tpl", tpl)
    http.HandleFunc("/hello", handler)
    http.HandleFunc("/sign", sign)
    http.HandleFunc("/show", show)
}

const guestbookForm = `
<html>l
  <body>
    <form action="/show" method="post">
      <div><textarea name="content" rows="3" cols="60"></textarea></div>
      <div><input type="submit" value="Sign Guestbook"></div>
    </form>
  </body>
</html>
`

const signTemplateHTML = `
<html>
  <body>
    <p>You wrote:</p>
    <pre>{{.}}</pre>
  </body>
</html>
`

type Greeting struct {
	Author string
	Content string
	Date time.Time
}

var signTemplate = template.Must(template.New("sign").Parse(signTemplateHTML))
var guestbookTemplate = template.Must(template.New("book").Parse(`
<html>
  <head>
    <title>Go Guestbook</title>
  </head>
  <body>
    {{range .}}
      {{with .Author}}
        <p><b>{{.}}</b> wrote:</p>
      {{else}}
        <p>An anonymous person wrote:</p>
      {{end}}
      <pre>{{.Content}} | {{.Date}}</pre>
    {{end}}
    <form action="/sign" method="post">
      <div><textarea name="content" rows="3" cols="60"></textarea></div>
      <div><input type="submit" value="Sign Guestbook"></div>
    </form>
  </body>
</html>
`))

func guestbookKey(context appengine.Context) *datastore.Key {
	return datastore.NewKey(context, "Guestbook", "default_guestbook", 0, nil)
}

func root(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	L(context, "test")
	query := datastore.NewQuery("Greeting").Ancestor(
		guestbookKey(context)).Order("-Date").Limit(10)
	greetings := make([]Greeting, 0, 10)
	if _, err := query.GetAll(context, &greetings); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := guestbookTemplate.Execute(writer, greetings); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func sign(writer http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	greeting := Greeting{
		Content: request.FormValue("content"),
		Date: time.Now(),
	}
	if appUser := user.Current(context); appUser != nil {
		greeting.Author = appUser.String()
	}
	key := datastore.NewIncompleteKey(
		context, "Greeting", guestbookKey(context))
	_, err := datastore.Put(context, key, &greeting)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(writer, request, "/", http.StatusFound)
}

func tpl(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, guestbookForm)
}

func show(writer http.ResponseWriter, request *http.Request) {
	err := signTemplate.Execute(writer, request.FormValue("content"))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func handler(writer http.ResponseWriter, request *http.Request) {
    context := appengine.NewContext(request)
    app_user := user.Current(context)
    if app_user == nil {
        url, err := user.LoginURL(context, request.URL.String())
        if err != nil {
            http.Error(writer, err.Error(), http.StatusInternalServerError)
            return
        }
        writer.Header().Set("Location", url)
        writer.WriteHeader(http.StatusFound)
        return
    }
    fmt.Fprintf(writer, "Hello, %v!", app_user)
}
