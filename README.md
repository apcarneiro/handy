Handy
==========================================

Handy is a fast and simple HTTP multiplexer for Golang. It implements a request
router that matches incoming requests against pre-registered Handlers. These
Handlers has methods for each HTTP verb. Registered paths can also have
variables.

## Creating a Handler
You just need to embed handy.DefaultHandler in your structure:

	type MyHandler struct {
		handy.DefaultHandler
	}

Now, override the HTTP verb:

	func (h *MyHandler) Get(ctx *handy.Context) {
		ctx.ResponseWriter.Write([]byte("Hello World"))
	}

And...

	package main

	import (
		"fmt"
		"handy"
		"net/http"
	)

	func main() {
		srv := handy.NewHandy()
		srv.HandleService("/hello/", new(MyHandler))
		fmt.Println(http.ListenAndServe(":8080", srv))
	}

## Path with variables
Path variables must be enclosed by braces.

	srv.HandleService("/hello/{foo}", new(MyHandler))

And you can read them using the Context:
	func (h *MyHandler) Get(ctx *handy.Context) {
		ctx.GetVar("foo")
		...
	}

## To create pre and post filters:
	func BeforeFilter(ctx *handy.Context) error {
		fmt.Printf("Hello %s\n", ctx.Request.RemoteAddr)
		return nil
	}

	func AfterFilter(ctx *handy.Context) error {
		fmt.Printf("Bye %s. x variable=%s\n", 
		ctx.Request.RemoteAddr, 
		ctx.GetVar("x"))
		return nil
	}

	func main() {
		srv := handy.NewHandy()
		srv.BeforeFilter(BeforeFilter)
		srv.AfterFilter(AfterFilter)
		...
		fmt.Println(http.ListenAndServe(":8080", srv))
}
