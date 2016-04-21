package main

import (
	"io/ioutil"
	"net/http"

	"github.com/cagnosolutions/adb"
	"github.com/cagnosolutions/web"
)

func init() {

	mx.AddRoutes(compLogin, compLoginPost)

	mx.AddSecureRoutes(COMPANY, compHome, compDriver, compSave)

}

var compLogin = web.Route{"GET", "/company/login", func(w http.ResponseWriter, r *http.Request) {
	tc.Render(w, r, "comp-login.tmpl", web.Model{})
}}

var compLoginPost = web.Route{"POST", "/company/login", func(w http.ResponseWriter, r *http.Request) {
	email, pass := r.FormValue("email"), r.FormValue("password")
	var company Company
	if !db.Auth("company", email, pass, &company) {
		web.SetErrorRedirect(w, r, "/company/login", "Incorrect username or password")
		return
	}
	sess := web.Login(w, r, company.Role)
	sess["id"] = company.Id
	sess["email"] = company.Email
	web.PutMultiSess(w, r, sess)
	web.SetSuccessRedirect(w, r, "/company/home", "Welcome "+company.Contact)
	return
}}

var compHome = web.Route{"GET", "/company/home", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	compId := web.GetSess(r, "id").(string)
	db.Get("company", compId, &company)
	var drivers []Driver
	db.TestQuery("driver", &drivers, adb.Eq("companyId", `"`+company.Id+`"`))
	tc.Render(w, r, "comp-home.tmpl", web.Model{
		"company": company,
		"drivers": drivers,
	})
}}

var compDriver = web.Route{"GET", "/company/driver/:id", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	compId := web.GetSess(r, "id").(string)
	db.Get("company", compId, &company)
	var driver Driver
	if !db.Get("driver", r.FormValue(":id"), &driver) || driver.CompanyId != company.Id {
		web.SetErrorRedirect(w, r, "/company/home", "Error retrieving driver")
		return
	}
	var files []map[string]interface{}
	var docs []Document
	if fileInfos, err := ioutil.ReadDir("upload/driver/" + driver.Id); err == nil {
		for _, fileInfo := range fileInfos {
			var info = make(map[string]interface{})
			info["name"] = fileInfo.Name()
			info["size"] = fileInfo.Size()
			files = append(files, info)
		}
	}
	db.TestQuery("document", &docs, adb.Eq("driverId", `"`+driver.Id+`"`))
	tc.Render(w, r, "comp-driver.tmpl", web.Model{
		"company": company,
		"driver":  driver,
		"files":   files,
		"dqfs":    DQFS,
		"docs":    docs,
	})
}}

var compSave = web.Route{"POST", "/company", func(w http.ResponseWriter, r *http.Request) {
	compId := r.FormValue("id")
	var company Company
	db.Get("company", compId, &company)
	FormToStruct(&company, r.Form, "")
	var companies []Company
	db.TestQuery("company", &companies, adb.Eq("email", company.Email), adb.Ne("id", `"`+company.Id+`"`))
	if len(companies) > 0 {
		web.SetErrorRedirect(w, r, "/company/home", "Error saving company. Email is already registered")
		return
	}
	company.Active = r.FormValue("auth.Active") == "true"
	db.Set("company", company.Id, company)
	web.SetSuccessRedirect(w, r, "/company/home", "Successfully saved company")
	return
}}
