package main

import (
	"net/http"

	"github.com/cagnosolutions/adb"
	"github.com/cagnosolutions/web"
)

func init() {
	mx.AddRoutes(driverLogin, driverLoginPost)
	//mx.AddSecureRoutes()
}

var driverLogin = web.Route{"GET", "/login/:slug", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	if !db.TestQueryOne("user", &company, adb.Eq("slug", r.FormValue(":slug")), adb.Eq("active", "true")) {
		web.SetErrorRedirect(w, r, "/login", "Error finding company")
		return
	}
	tc.Render(w, r, "driver-login.tmpl", web.Model{
		"company": company,
	})
}}

var driverLoginPost = web.Route{"POST", "/login/:slug", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	if !db.TestQueryOne("user", &company, adb.Eq("slug", r.FormValue(":slug")), adb.Eq("active", "true")) {
		web.SetErrorRedirect(w, r, "/login", "Error finding company")
		return
	}

	email, pass := r.FormValue("email"), r.FormValue("password")
	var driver Driver
	if !db.Auth("user", email, pass, &driver) || driver.CompanyId != company.Id {
		web.SetErrorRedirect(w, r, "/login/"+r.FormValue(":slug"), "Incorrect username or password")
		return
	}
	sess := web.Login(w, r, driver.Role)
	sess["id"] = driver.Id
	sess["email"] = driver.Email
	sess["companyId"] = driver.CompanyId
	web.PutMultiSess(w, r, sess)
	web.SetSuccessRedirect(w, r, "/"+company.Slug+"/home", "Welcome "+driver.FirstName)
	return
}}
