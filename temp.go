package main

import (
	"net/http"

	"github.com/cagnosolutions/web"
)

func init() {
	mx.AddRoutes(test, login, loginPost, logout, driverDocuments)
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

var driverDocuments = web.Route{"GET", "/document/:id", func(w http.ResponseWriter, r *http.Request) {
	var document Document
	var driver Driver
	var company Company
	ok := db.Get("document", r.FormValue(":id"), &document)
	if !ok {
		web.SetErrorRedirect(w, r, "/", "Error, retrieving document.")
		return
	}
	ok = db.Get("user", document.DriverId, &driver)
	if !ok {
		web.SetErrorRedirect(w, r, "/", "Error, document is not associated with a driver.")
		return
	}
	ok = db.Get("user", document.CompanyId, &company)
	if !ok {
		web.SetErrorRedirect(w, r, "/", "Error, document is not associated with a company.")
		return
	}
	tc.Render(w, r, "dqf-"+document.DocumentId+".tmpl", web.Model{
		"document": document,
		"driver":   driver,
		"company":  company,
	})
	return
}}
