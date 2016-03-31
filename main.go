package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/cagnosolutions/adb"
	"github.com/cagnosolutions/web"
)

// global vars
var mx *web.Mux = web.NewMux()
var tc *web.TmplCache
var db *adb.DB = adb.NewDB()

// initialize routes
func init() {
	db.AddStore("user")
	mx.AddRoutes(index, buttons, makeEmployees, makeCompanies, makeDrivers, dt, contactRedirect)
	mx.AddRoutes(viewCompany, viewDriver, viewEmployee)

	web.Funcs["lower"] = strings.ToLower
	tc = web.NewTmplCache()
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

var contactRedirect = web.Route{"GET", "/contact", func(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/contact/employee", 303)
}}

// var contact = web.Route{"GET", "/contact/:role", func(w http.ResponseWriter, r *http.Request) {
// 	role := r.FormValue(":role")
// 	if role == "" || (role != "employee" && role != "company" && role != "driver") {
// 		http.Redirect(w, r, "/contact/employee", 303)
// 		return
// 	}
// 	var users []User
// 	ok := db.Match("user", `"role":"`+strings.ToUpper(role)+`"`, &users)
// 	if !ok {
// 		fmt.Println("error")
// 	}
// 	tc.Render(w, r, "contact.tmpl", web.Model{
// 		"users": users,
// 		"role":  role,
// 	})
// }}

var viewEmployee = web.Route{"GET", "/employee/:id", func(w http.ResponseWriter, r *http.Request) {
	var user Employee
	ok := db.Get("user", r.FormValue(":id"), &user)
	if !ok || user.Role != "EMPLOYEE" {
		web.SetErrorRedirect(w, r, "/contact/employee", "Error finding employee")
		return
	}
	tc.Render(w, r, "employee.tmpl", web.Model{
		"user": user,
	})
}}

var viewCompany = web.Route{"GET", "/company/:id", func(w http.ResponseWriter, r *http.Request) {
	var user Company
	ok := db.Get("user", r.FormValue(":id"), &user)
	if !ok || user.Role != "COMPANY" {
		web.SetErrorRedirect(w, r, "/contact/company", "Error finding company")
		return
	}
	tc.Render(w, r, "company.tmpl", web.Model{
		"user": user,
	})
}}

var viewDriver = web.Route{"GET", "/driver/:id", func(w http.ResponseWriter, r *http.Request) {
	var user Driver
	ok := db.Get("user", r.FormValue(":id"), &user)
	if !ok || user.Role != "DRIVER" {
		web.SetErrorRedirect(w, r, "/contact/driver", "Error finding driver")
		return
	}
	tc.Render(w, r, "driver.tmpl", web.Model{
		"user": user,
	})
}}

var dt = web.Route{"GET", "/dt", func(w http.ResponseWriter, r *http.Request) {
	tc.Render(w, r, "table-datatable.tmpl", nil)
}}

/*var profile = web.Route{"GET", "/profile/:id", func(w http.ResponseWriter, r *http.Request) {
	var user User
	db.Get("user", r.FormValue(":id"), &user)
	tc.Render(w, r, "profile.tmpl", web.Model{
		"user": user,
	})
}}*/
