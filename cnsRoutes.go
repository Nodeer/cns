package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/cagnosolutions/adb"
	"github.com/cagnosolutions/web"
)

var index = web.Route{"GET", "/", func(w http.ResponseWriter, r *http.Request) {
	tc.Render(w, r, "index.tmpl", web.Model{})
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
	redirect := "/cns/company"
	if employee.Home != "" {
		redirect = employee.Home
	}
	web.SetSuccessRedirect(w, r, redirect, "Welcome "+employee.FirstName)
	return
}}

var allEmployee = web.Route{"GET", "/cns/employee", func(w http.ResponseWriter, r *http.Request) {
	var employees []Employee
	//ok := db.Match("user", `"role":"EMPLOYEE"`, &employees)
	ok := db.TestQuery("user", &employees, adb.Eq("role", "EMPLOYEE"))
	if !ok {
		fmt.Println("error")
	}
	tc.Render(w, r, "all-employee.tmpl", web.Model{
		"employees": employees,
	})
}}

var viewEmployee = web.Route{"GET", "/cns/employee/:id", func(w http.ResponseWriter, r *http.Request) {
	var employee Employee
	employeeId := r.FormValue(":id")
	if employeeId != "add" {
		ok := db.Get("user", r.FormValue(":id"), &employee)
		if !ok || employee.Role != "EMPLOYEE" {
			web.SetErrorRedirect(w, r, "/employee", "Error finding employee")
			return
		}
	}
	tc.Render(w, r, "employee.tmpl", web.Model{
		"employee": employee,
	})
}}

var settings = web.Route{"GET", "/cns/settings", func(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/cns/employee/"+web.GetSess(r, "id").(string), 303)
}}

var saveEmployee = web.Route{"POST", "/cns/employee", func(w http.ResponseWriter, r *http.Request) {
	empId := r.FormValue("id")
	var employee Employee
	db.Get("user", empId, &employee)
	FormToStruct(&employee, r.Form, "")
	if employee.Id == "" && empId == "" {
		employee.Id = strconv.Itoa(int(time.Now().UnixNano()))
		employee.Password = employee.Email
		employee.Role = "EMPLOYEE"
		employee.Active = true
	}

	var users []interface{}
	db.TestQuery("user", &users, adb.Eq("email", employee.Email), adb.Ne("id", `"`+employee.Id+`"`))
	if len(users) > 0 {
		web.SetErrorRedirect(w, r, "/cns/employee/"+employee.Id, "Error saving employee. Email is already registered")
		return
	}
	db.Set("user", employee.Id, employee)
	web.SetSuccessRedirect(w, r, "/cns/employee/"+employee.Id, "Successfully saved employee")
	return
}}

var allCompany = web.Route{"GET", "/cns/company", func(w http.ResponseWriter, r *http.Request) {
	var companies []Company
	//ok := db.Match("user", `"role":"COMPANY"`, &companies)
	ok := db.TestQuery("user", &companies, adb.Eq("role", "COMPANY"))
	if !ok {
		fmt.Println("error")
	}
	tc.Render(w, r, "all-company.tmpl", web.Model{
		"companies": companies,
	})
}}

var viewCompany = web.Route{"GET", "/cns/company/:id", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	var drivers []Driver
	compId := r.FormValue(":id")
	if compId != "add" {
		ok := db.Get("user", compId, &company)
		if !ok || company.Role != "COMPANY" {
			web.SetErrorRedirect(w, r, "/company", "Error finding company")
			return
		}
		//db.Match("user", `"companyId":"`+company.Id+`"`, &drivers)
		//db.Query("user", &drivers, "companyId="+company.Id)
		db.TestQuery("user", &drivers, adb.Eq("companyId", `"`+company.Id+`"`))
	}
	tc.Render(w, r, "company.tmpl", web.Model{
		"company": company,
		"drivers": drivers,
	})
}}

var saveCompany = web.Route{"POST", "/cns/company", func(w http.ResponseWriter, r *http.Request) {
	compId := r.FormValue("id")
	var company Company
	db.Get("user", compId, &company)
	if compId == "" && company.Id == "" {
		company.Id = strconv.Itoa(int(time.Now().UnixNano()))
		company.Password = company.Email
		company.Role = "COMPANY"
		company.CreateSlug()
	}
	FormToStruct(&company, r.Form, "")
	var users []interface{}
	db.TestQuery("user", &users, adb.Eq("email", company.Email), adb.Ne("id", `"`+company.Id+`"`))
	if len(users) > 0 {
		web.SetErrorRedirect(w, r, "/cns/company/"+company.Id, "Error saving company. Email is already registered")
		return
	}
	company.Active = r.FormValue("auth.Active") == "true"
	db.Set("user", company.Id, company)
	web.SetSuccessRedirect(w, r, "/cns/company/"+company.Id, "Successfully saved company")
	return
}}

var allDriver = web.Route{"GET", "/cns/driver", func(w http.ResponseWriter, r *http.Request) {
	var drivers []Driver
	//ok := db.Match("user", `"role":"DRIVER"`, &drivers)
	ok := db.TestQuery("user", &drivers, adb.Eq("role", "DRIVER"))
	if !ok {
		fmt.Println("error")
	}
	tc.Render(w, r, "all-driver.tmpl", web.Model{
		"drivers": drivers,
	})
}}

var viewDriver = web.Route{"GET", "/cns/driver/:id", func(w http.ResponseWriter, r *http.Request) {
	var driver Driver
	driverId := r.FormValue(":id")
	var files []map[string]interface{}
	var docs []Document
	if driverId != "add" {
		ok := db.Get("user", driverId, &driver)
		if !ok || driver.Role != "DRIVER" {
			web.SetErrorRedirect(w, r, "/driver", "Error finding driver")
			return
		}
		if fileInfos, err := ioutil.ReadDir("upload/driver/" + driverId); err == nil {
			for _, fileInfo := range fileInfos {
				var info = make(map[string]interface{})
				info["name"] = fileInfo.Name()
				info["size"] = fileInfo.Size()
				files = append(files, info)
			}
		}
		//db.Match("document", `"driverId":"`+driver.Id+`"`, &docs)
		db.TestQuery("document", &docs, adb.Eq("driverId", `"`+driver.Id+`"`))
	}
	tc.Render(w, r, "driver.tmpl", web.Model{
		"driver": driver,
		"files":  files,
		"dqfs":   DQFS,
		"docs":   docs,
	})
}}

var saveDriver = web.Route{"POST", "/cns/driver", func(w http.ResponseWriter, r *http.Request) {
	driverId := r.FormValue("id")
	var driver Driver
	db.Get("user", driverId, &driver)
	FormToStruct(&driver, r.Form, "")
	if driver.Id == "" && driverId == "" {
		driver.Id = strconv.Itoa(int(time.Now().UnixNano()))
		driver.Password = driver.Email
		driver.Role = "DRIVER"
	}
	db.Set("user", driver.Id, driver)
	web.SetSuccessRedirect(w, r, "/cns/driver/"+driver.Id, "Successfully saved driver")
	return
}}
