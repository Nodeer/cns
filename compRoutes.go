package main

import (
	"net/http"

	"github.com/cagnosolutions/web"
)

func init() {

	mx.AddRoutes(compLogin)

}

var compLogin = web.Route{"GET", "/company/login", func(w http.ResponseWriter, r *http.Request) {
	tc.Render(w, r, "comp-login.tmpl", web.Model{})
}}

var compLoginPost = web.Route{"POST", "/company/login", func(w http.ResponseWriter, r *http.Request) {
	email, pass := r.FormValue("email"), r.FormValue("password")
	var company Company
	if !db.Auth("user", email, pass, &company) {
		web.SetErrorRedirect(w, r, "/company/login", "Incorrect username or password")
		return
	}
	sess := web.Login(w, r, company.Role)
	sess["id"] = company.Id
	sess["email"] = company.Email
	web.PutMultiSess(w, r, sess)
	web.SetSuccessRedirect(w, r, "/company/home", "Welcome "+company.FirstName)
	return
}}
