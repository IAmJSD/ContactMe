package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/buaazp/fasthttprouter"
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
	// TODO: Handle this.
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
	router.GET("/styles.css", ServeStatic("./ui/dist/styles.css"))
	router.GET("/styles.css.map", ServeStatic("./ui/dist/styles.css.map"))

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}
