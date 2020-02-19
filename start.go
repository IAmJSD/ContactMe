package main

import (
	"encoding/json"
	"log"

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

}

// Defines the main function.
func main() {
	router := fasthttprouter.New()
	router.GET("/_forms", FormsRoute)
	router.POST("/", FormPost)

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}
