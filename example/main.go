package main

import (
	"github.com/up1-io/zebra"
	"log"
)

func main() {
	r := zebra.NewRouter()

	r.On("/about", aboutHandler)
	r.On("/", homeHandler)
	r.On("/users/{id}/", userHandler)
	r.On("/users/{id}/details/{postId}", postHandler)
	r.On("/test/redirect-test/", redirectHandler)

	app, err := zebra.New()
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatalln(app.ListenAndServe(":8080", r))
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
			"Meta": struct {
				Title       string
				Description string
			}{
				Title:       "Zebra",
				Description: "Zebra is a minimalist web framework for Go that focuses on simplicity, performance, and ease of use.",
			},
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
				Id: r.PathVariables.Get("id"),
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
