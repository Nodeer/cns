package main

import (
	"net/http"

	"github.com/cagnosolutions/web"
)

func init() {
	mx.AddRoutes(test)
}

var test = web.Route{"GET", "/test", func(w http.ResponseWriter, r *http.Request) {
	tc.Render(w, r, "test.tmpl", web.Model{})
}}
