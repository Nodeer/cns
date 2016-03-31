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
	mx.AddRoutes(index, buttons, allEmployee, allCompany, allDriver)
	mx.AddRoutes(makeEmployees, makeCompanies, makeDrivers, dt)
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

var viewEmployee = web.Route{"GET", "/employee/:id", func(w http.ResponseWriter, r *http.Request) {
	var employee Employee
	ok := db.Get("user", r.FormValue(":id"), &employee)
	if !ok || employee.Role != "EMPLOYEE" {
		web.SetErrorRedirect(w, r, "/contact/employee", "Error finding employee")
		return
	}
	tc.Render(w, r, "employee.tmpl", web.Model{
		"employee": employee,
	})
}}

var viewCompany = web.Route{"GET", "/company/:id", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	ok := db.Get("user", r.FormValue(":id"), &company)
	if !ok || company.Role != "COMPANY" {
		web.SetErrorRedirect(w, r, "/contact/company", "Error finding company")
		return
	}
	tc.Render(w, r, "company.tmpl", web.Model{
		"company": company,
	})
}}

var viewDriver = web.Route{"GET", "/driver/:id", func(w http.ResponseWriter, r *http.Request) {
	var driver Driver
	ok := db.Get("user", r.FormValue(":id"), &driver)
	if !ok || driver.Role != "DRIVER" {
		web.SetErrorRedirect(w, r, "/contact/driver", "Error finding driver")
		return
	}
	tc.Render(w, r, "driver.tmpl", web.Model{
		"driver": driver,
	})
}}

var dt = web.Route{"GET", "/dt", func(w http.ResponseWriter, r *http.Request) {
	tc.Render(w, r, "table-datatable.tmpl", nil)
}}
