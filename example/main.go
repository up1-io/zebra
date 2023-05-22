package main

import (
	"github.com/up1-io/zebra"
	"log"
)

func main() {
	app, err := zebra.New()
	if err != nil {
		log.Fatalln(err)
	}

	app.Router.On("/about", aboutHandler)
	app.Router.On("/", homeHandler)
	app.Router.On("/users/{id}/", userHandler)
	app.Router.On("/users/{id}/details/{postId}", postHandler)
	app.Router.On("/test/redirect-test/", redirectHandler)

	log.Fatalln(app.ListenAndServe(":8080"))
}

func aboutHandler(ctx zebra.Request, callback zebra.Callback) {
	callback(nil, zebra.Result{
		Data: map[string]interface{}{
			"Content": struct {
				Title string
			}{
				Title: "Dynamic title set on server side",
			},
		},
	})
}

func homeHandler(ctx zebra.Request, callback zebra.Callback) {
	callback(nil, zebra.Result{
		Data: map[string]interface{}{
			"Content": struct {
				Status string
			}{
				Status: "Hello World",
			},
		},
	})
}

func userHandler(r zebra.Request, callback zebra.Callback) {
	callback(nil, zebra.Result{
		Data: map[string]interface{}{
			"Content": struct {
				Id string
			}{
				Id: r.PathVariables["id"],
			},
		},
	})
}

func postHandler(r zebra.Request, callback zebra.Callback) {
	callback(nil, zebra.Result{
		Data: map[string]interface{}{
			"Content": struct {
				Id     string
				PostId string
			}{
				Id:     r.PathVariables.Get("id"),
				PostId: r.PathVariables.Get("postId"),
			},
		},
	})
}

func redirectHandler(r zebra.Request, callback zebra.Callback) {
	callback(nil, zebra.Result{
		Redirect: "/",
	})
}
