package main

import (
	"bytes"
	"encoding/json"
	"html"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/buaazp/fasthttprouter"
	"github.com/mailgun/mailgun-go"
	"github.com/valyala/fasthttp"
)

// FormsRoute defines the route to show the forms.
func FormsRoute(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	j, err := json.Marshal(&Config)
	if err != nil {
		panic(err)
	}
	ctx.Response.SetBody(j)
}

// FormPost is used to handle a form being posted.
func FormPost(ctx *fasthttp.RequestCtx) {
	// Get the JSON.
	var m map[string]string
	err := json.Unmarshal(ctx.Request.Body(), &m)
	if err != nil {
		ctx.Response.SetStatusCode(400)
		return
	}

	// Get/delete __formName
	FormName, ok := m["__formName"]
	if !ok {
		ctx.Response.SetStatusCode(400)
		return
	}
	delete(m, "__formName")

	// Get the form.
	var f *Form
	for _, v := range Config.Forms {
		if v.Name == FormName {
			f = &v
			break
		}
	}
	if f == nil {
		ctx.Response.SetStatusCode(400)
		return
	}

	// Do a sanity check of the form.
	for i, v := range *f.Children {
		Required, _ := v["required"].(bool)
		_, ok := m[i]
		if Required && !ok {
			ctx.Response.SetStatusCode(400)
			return
		}
	}

	// Get the IP address.
	m["IP Address"] = string(ctx.Request.Header.Peek("CF-Connecting-IP"))

	// Create the e-mail.
	Email := "The form \"" + FormName + "\" has had an submission:\n"
	for i, v := range m {
		Email += "<hr /><p><b>" + html.EscapeString(i) + ":</b> " + html.EscapeString(v) + "</p>"
	}

	// Send the e-mail.
	mg := mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_API_KEY"))
	mg.SetAPIBase("https://api.eu.mailgun.net/v3")
	msg := mg.NewMessage(os.Getenv("FROM_ADDRESS"), "Form submitted: "+FormName, "", os.Getenv("TO_ADDRESS"))
	msg.SetHtml(Email)
	_, _, err = mg.Send(msg)
	if err != nil {
		println(err.Error())
		ctx.Response.SetStatusCode(500)
	}
}

// ServeStatic is used to serve the static content.
func ServeStatic(FilePath string) func(ctx *fasthttp.RequestCtx) {
	File, err := ioutil.ReadFile(FilePath)
	if err != nil {
		panic(err)
	}
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(200)
		ctx.Response.Header.Set("Content-Type", http.DetectContentType(File))
		_, _ = ctx.Write(File)
	}
}

// ServeHome is used to serve the homepage.
func ServeHome(ctx *fasthttp.RequestCtx) {
	File, err := ioutil.ReadFile("./index.template.html")
	if err != nil {
		panic(err)
	}
	t, err := template.New("home").Parse(string(File))
	if err != nil {
		panic(err)
	}
	ctx.SetStatusCode(200)
	ctx.Response.Header.Set("Content-Type", "text/html")
	b := bytes.Buffer{}
	err = t.Execute(&b, map[string]string{
		"Title":       strings.ReplaceAll(Config.PageDescription.HTMLTitle, "{title}", Config.PageDescription.Title),
		"Description": Config.PageDescription.Description,
	})
	if err != nil {
		panic(err)
	}
	ctx.Response.SetBody(b.Bytes())
}

// Defines the main function.
func main() {
	router := fasthttprouter.New()
	router.GET("/_forms", FormsRoute)
	router.POST("/", FormPost)
	router.GET("/", ServeHome)
	router.GET("/mount.js", ServeStatic("./ui/dist/mount.js"))
	router.GET("/mount.js.map", ServeStatic("./ui/dist/mount.js.map"))

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}
