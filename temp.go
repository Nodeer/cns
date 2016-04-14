package main

import (
	"net/http"

	"github.com/cagnosolutions/web"
)

func init() {
	mx.AddRoutes(test, login, loginPost, logout)
}

var test = web.Route{"GET", "/test", func(w http.ResponseWriter, r *http.Request) {
	tc.Render(w, r, "test.tmpl", web.Model{})
}}

var login = web.Route{"GET", "/login", func(w http.ResponseWriter, r *http.Request) {
	tc.Render(w, r, "login.tmpl", web.Model{})
}}

var loginPost = web.Route{"POST", "/login", func(w http.ResponseWriter, r *http.Request) {
	email, pass := r.FormValue("email"), r.FormValue("password")
	var employee Employee
	if !db.Auth("user", email, pass, &employee) {
		web.SetErrorRedirect(w, r, "/login", "Incorrect username or password")
		return
	}
	sess := web.Login(w, r, employee.Role)
	sess["id"] = employee.Id
	sess["email"] = employee.Email
	web.PutMultiSess(w, r, sess)
	web.SetSuccessRedirect(w, r, "/", "Welcome "+employee.FirstName)
	return
}}

var logout = web.Route{"GET", "/logout", func(w http.ResponseWriter, r *http.Request) {
	web.Logout(w)
	web.SetSuccessRedirect(w, r, "/login", "Successfully logged out")
}}
