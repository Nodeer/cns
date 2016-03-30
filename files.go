package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cagnosolutions/web"
)

func init() {
	mx.AddRoutes(files, filesApi)
}

type Node struct {
	Id       string `json:"id,omitempty"`
	Text     string `json:"text,omitempty"`
	Children bool   `json:"children,omitempty"`
	Type     string `json:"type,omitempty"`
}

var files = web.Route{"GET", "/files", func(w http.ResponseWriter, r *http.Request) {
	tc.Render(w, r, "files.tmpl", nil)
	return
}}

var filesApi = web.Route{"GET", "/api/files", func(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("id")
	files, err := ioutil.ReadDir("." + path)
	if err != nil {
		panic(err)
	}
	var nodes []Node
	for _, file := range files {
		if file.Name()[0] != '.' {
			n := Node{}
			n.Id = path + "/" + file.Name()
			n.Text = file.Name()
			if file.IsDir() {
				n.Type = "dir"
				n.Children = true
			}
			nodes = append(nodes, n)
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(nodes); err != nil {
		panic(err)
	}
	return
}}
