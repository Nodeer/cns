package main

import (
	"net/http"

	"github.com/cagnosolutions/web"
)

var DevComments = web.Route{"GET", "/dev/comment", func(w http.ResponseWriter, r *http.Request) {
	var comments []Comment
	db.All("comment", &comments)
	tc.Render(w, r, "all-comments.tmpl", web.Model{
		"comments": comments,
	})
}}
