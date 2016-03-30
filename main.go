package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	mx.AddRoutes(index, buttons, contact, makeEmployees, profile, makeCompanies, dt)
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

var contact = web.Route{"GET", "/contact/:role", func(w http.ResponseWriter, r *http.Request) {
	role := r.FormValue(":role")
	if role == "" || (role != "employee" && role != "company" && role != "driver") {
		http.Redirect(w, r, "/contact/employee", 303)
		return
	}
	var users []User
	ok := db.Match("user", `"role":"`+strings.ToUpper(role)+`"`, &users)
	if !ok {
		fmt.Println("error")
	}
	tc.Render(w, r, "contact.tmpl", web.Model{
		"users": users,
		"role":  role,
	})
}}

var makeEmployees = web.Route{"GET", "/makeUsers", func(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 10; i++ {
		id := strconv.Itoa(int(time.Now().UnixNano()))
		user := User{
			Id:       id,
			Email:    fmt.Sprintf("%d@%d.com", i, i),
			Password: fmt.Sprintf("Password-%d", i),
			Role:     "EMPLOYEE",
			Active:   (i%2 == 0),
			Name:     fmt.Sprintf("John Smith the %dth", (i + 4)),
		}
		db.Add("user", id, user)
	}
	web.SetSuccessRedirect(w, r, "/", "success")
	return
}}

var makeCompanies = web.Route{"GET", "/makeComps", func(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 10; i++ {
		id := strconv.Itoa(int(time.Now().UnixNano()))
		user := User{
			Id:       id,
			Email:    fmt.Sprintf("company@company%d.com", i+1),
			Password: fmt.Sprintf("Password-%d", i),
			Role:     "COMPANY",
			Active:   (i%2 == 0),
			Name:     fmt.Sprintf("Company #%d", (i + 1)),
		}
		db.Add("user", id, user)
	}
	web.SetSuccessRedirect(w, r, "/", "success")
	return
}}

var profile = web.Route{"GET", "/profile/:id", func(w http.ResponseWriter, r *http.Request) {
	var user User
	db.Get("user", r.FormValue(":id"), &user)
	tc.Render(w, r, "profile.tmpl", map[string]interface{}{
		"user": user,
	})
}}

var dt = web.Route{"GET", "/dt", func(w http.ResponseWriter, r *http.Request) {
	tc.Render(w, r, "table-datatable.tmpl", nil)
}}
