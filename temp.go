package main

import (
	"net/http"

	"github.com/cagnosolutions/web"
)

func init() {
	mx.AddRoutes(test, dqf450)
}

var test = web.Route{"GET", "/test", func(w http.ResponseWriter, r *http.Request) {
	tc.Render(w, r, "test.tmpl", web.Model{})
}}

var dqf450 = web.Route{"GET", "/450", func(w http.ResponseWriter, r *http.Request) {
	tc.Render(w, r, "dqf-450.tmpl", web.Model{})
}}
