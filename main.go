package main

import (
	"fmt"
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

	mx.AddRoutes(index, buttons)

	mx.AddRoutes(allCompany, viewCompany, saveCompany)
	mx.AddRoutes(allEmployee, viewEmployee, saveEmployee)
	mx.AddRoutes(allDriver, viewDriver, savedriver)

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

var allEmployee = web.Route{"GET", "/employee", func(w http.ResponseWriter, r *http.Request) {
	var employees []Employee
	ok := db.Match("user", `"role":"EMPLOYEE"`, &employees)
	if !ok {
		fmt.Println("error")
	}
	tc.Render(w, r, "all-employee.tmpl", web.Model{
		"employees": employees,
	})
}}

var viewEmployee = web.Route{"GET", "/employee/:id", func(w http.ResponseWriter, r *http.Request) {
	var employee Employee
	ok := db.Get("user", r.FormValue(":id"), &employee)
	if !ok || employee.Role != "EMPLOYEE" {
		web.SetErrorRedirect(w, r, "/employee", "Error finding employee")
		return
	}
	tc.Render(w, r, "employee.tmpl", web.Model{
		"employee": employee,
	})
}}

var saveEmployee = web.Route{"POST", "/employee/:id", func(w http.ResponseWriter, r *http.Request) {
	empId := r.FormValue(":id")
	var employee Employee
	if !db.Get("user", empId, &employee) {
		web.SetErrorRedirect(w, r, "/employee/"+empId, "Error saving employee")
		return
	}
	r.ParseForm()
	FormToStruct(&employee, r.Form, "")
	db.Set("user", empId, employee)
	web.SetSuccessRedirect(w, r, "/employee/"+empId, "Successfully saved employee")
	return
}}

var allCompany = web.Route{"GET", "/company", func(w http.ResponseWriter, r *http.Request) {
	var companies []Company
	ok := db.Match("user", `"role":"COMPANY"`, &companies)
	if !ok {
		fmt.Println("error")
	}
	tc.Render(w, r, "all-company.tmpl", web.Model{
		"companies": companies,
	})
}}

var viewCompany = web.Route{"GET", "/company/:id", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	ok := db.Get("user", r.FormValue(":id"), &company)
	if !ok || company.Role != "COMPANY" {
		web.SetErrorRedirect(w, r, "/company", "Error finding company")
		return
	}
	var drivers []Driver
	ok = db.Match("user", `"companyId":"`+company.Id+`"`, &drivers)
	tc.Render(w, r, "company.tmpl", web.Model{
		"company": company,
		"drivers": drivers,
	})
}}

var saveCompany = web.Route{"POST", "/company/:id", func(w http.ResponseWriter, r *http.Request) {
	compId := r.FormValue(":id")
	var company Company
	if !db.Get("user", compId, &company) {
		web.SetErrorRedirect(w, r, "/company/"+compId, "Error saving company")
		return
	}
	r.ParseForm()
	FormToStruct(&company, r.Form, "")
	db.Set("user", compId, company)
	web.SetSuccessRedirect(w, r, "/company/"+compId, "Successfully saved company")
	return
}}

var allDriver = web.Route{"GET", "/driver", func(w http.ResponseWriter, r *http.Request) {
	var drivers []Driver
	ok := db.Match("user", `"role":"DRIVER"`, &drivers)
	if !ok {
		fmt.Println("error")
	}
	tc.Render(w, r, "all-driver.tmpl", web.Model{
		"drivers": drivers,
	})
}}

var viewDriver = web.Route{"GET", "/driver/:id", func(w http.ResponseWriter, r *http.Request) {
	var driver Driver
	ok := db.Get("user", r.FormValue(":id"), &driver)
	if !ok || driver.Role != "DRIVER" {
		web.SetErrorRedirect(w, r, "/driver", "Error finding driver")
		return
	}
	tc.Render(w, r, "driver.tmpl", web.Model{
		"driver": driver,
	})
}}

var savedriver = web.Route{"POST", "/driver/:id", func(w http.ResponseWriter, r *http.Request) {
	driverId := r.FormValue(":id")
	var driver Driver
	if !db.Get("user", driverId, &driver) {
		web.SetErrorRedirect(w, r, "/driver/"+driverId, "Error saving driver")
		return
	}
	r.ParseForm()
	FormToStruct(&driver, r.Form, "")
	db.Set("user", driverId, driver)
	web.SetSuccessRedirect(w, r, "/driver/"+driverId, "Successfully saved driver")
	return
}}

func ToLowerFirst(s string) string {
	return strings.ToLower(string(s[0])) + s[1:len(s)]
}
