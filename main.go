package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/cagnosolutions/adb"
	"github.com/cagnosolutions/web"
)

// global vars
var mx *web.Mux = web.NewMux()
var tc *web.TmplCache = web.NewTmplCache()
var db *adb.DB = adb.NewDB()

// initialize routes
func init() {
	db.AddStore("user")
	mx.AddRoutes(index, buttons, contact, makeUsers)
}

// main http listener
func main() {
	log.Fatal(http.ListenAndServe(":8080", mx))
}

// current route controllers
var index = web.Route{"GET", "/", func(w http.ResponseWriter, r *http.Request) {
	tc.Render(w, r, "index.tmpl", web.Model{})
}}

var buttons = web.Route{"GET", "/buttons", func(w http.ResponseWriter, r *http.Request) {
	tc.Render(w, r, "buttons.tmpl", web.Model{})
}}

var contact = web.Route{"GET", "/contact", func(w http.ResponseWriter, r *http.Request) {
	var users []User
	ok := db.All("user", &users)
	if !ok {
		fmt.Println("error")
	}
	tc.Render(w, r, "contact.tmpl", web.Model{
		"users": users,
	})
}}

var makeUsers = web.Route{"GET", "/makeUsers", func(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 10; i++ {
		id := strconv.Itoa(int(time.Now().UnixNano()))
		user := User{
			Id:       id,
			Email:    fmt.Sprintf("%d@%d.com", i, i),
			Password: fmt.Sprintf("Password-%d", i),
			Role:     "EMPLOYEE",
			Active:   (i%2 == 0),
		}
		db.Add("user", id, user)
	}
	web.SetSuccessRedirect(w, r, "/", "success")
	return
}}
