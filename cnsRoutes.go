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
	if !db.Auth("employee", email, pass, &employee) {
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

/* --- Employee Management --- */

var allEmployee = web.Route{"GET", "/cns/employee", func(w http.ResponseWriter, r *http.Request) {
	var employees []Employee
	//ok := db.Match("user", `"role":"EMPLOYEE"`, &employees)
	ok := db.TestQuery("employee", &employees, adb.Eq("role", "EMPLOYEE"))
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
		ok := db.Get("employee", r.FormValue(":id"), &employee)
		if !ok {
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
	db.Get("employee", empId, &employee)
	FormToStruct(&employee, r.Form, "")
	if employee.Id == "" && empId == "" {
		employee.Id = strconv.Itoa(int(time.Now().UnixNano()))
		employee.Password = employee.Email
		employee.Role = "EMPLOYEE"
		employee.Active = true
	}

	var employees []Employee
	db.TestQuery("employee", &employees, adb.Eq("email", employee.Email), adb.Ne("id", `"`+employee.Id+`"`))
	if len(employees) > 0 {
		web.SetErrorRedirect(w, r, "/cns/employee/"+employee.Id, "Error saving employee. Email is already registered")
		return
	}
	db.Set("employee", employee.Id, employee)
	web.SetSuccessRedirect(w, r, "/cns/employee/"+employee.Id, "Successfully saved employee")
	return
}}

/* --- Company Management --- */

var companyAll = web.Route{"GET", "/cns/company", func(w http.ResponseWriter, r *http.Request) {
	var companies []Company
	db.All("company", &companies)
	tc.Render(w, r, "company-all.tmpl", web.Model{
		"companies": companies,
	})
}}

var companyView = web.Route{"GET", "/cns/company/:id", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	compId := r.FormValue(":id")
	ok := db.Get("company", compId, &company)
	if !ok {
		web.SetErrorRedirect(w, r, "/cns/company", "Error finding company")
		return
	}

	tc.Render(w, r, "company.tmpl", web.Model{
		"company": company,
	})
}}

var companyDriver = web.Route{"GET", "/cns/company/:id/driver", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	var drivers []Driver
	compId := r.FormValue(":id")
	ok := db.Get("company", compId, &company)
	if !ok {
		web.SetErrorRedirect(w, r, "/cns/company", "Error finding company")
		return
	}
	db.TestQuery("driver", &drivers, adb.Eq("companyId", `"`+company.Id+`"`))
	tc.Render(w, r, "company-driver.tmpl", web.Model{
		"company": company,
		"drivers": drivers,
	})
}}

var companyVehicle = web.Route{"GET", "/cns/company/:id/vehicle", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	var vehicles []Vehicle
	compId := r.FormValue(":id")
	ok := db.Get("company", compId, &company)
	if !ok {
		web.SetErrorRedirect(w, r, "/cns/company", "Error finding company")
		return
	}
	db.TestQuery("vehicle", &vehicles, adb.Eq("companyId", `"`+company.Id+`"`))

	tc.Render(w, r, "company-vehicle.tmpl", web.Model{
		"company":  company,
		"vehicles": vehicles,
	})
}}

var companyNote = web.Route{"GET", "/cns/company/:id/note", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	var notes []Note
	compId := r.FormValue(":id")
	ok := db.Get("company", compId, &company)
	if !ok {
		web.SetErrorRedirect(w, r, "/cns/company", "Error finding company")
		return
	}
	db.TestQuery("note", &notes, adb.Eq("companyId", `"`+company.Id+`"`))

	tc.Render(w, r, "company-note.tmpl", web.Model{
		"company": company,
		"notes":   notes,
	})
}}

var companySetting = web.Route{"GET", "/cns/company/:id/setting", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	compId := r.FormValue(":id")
	ok := db.Get("company", compId, &company)
	if !ok {
		web.SetErrorRedirect(w, r, "/cns/company", "Error finding company")
		return
	}
	tc.Render(w, r, "company-setting.tmpl", web.Model{
		"company": company,
	})
}}

var companySaveNote = web.Route{"POST", "/cns/company/note", func(w http.ResponseWriter, r *http.Request) {
	id := strconv.Itoa(int(time.Now().UnixNano()))
	note := Note{
		Id:        id,
		CompanyId: r.FormValue("id"),
		Body:      r.FormValue("body"),
	}
	db.Add("note", id, note)
	web.SetSuccessRedirect(w, r, "/cns/company/"+r.FormValue("id"), "Successfully saved note")
}}

var companySave = web.Route{"POST", "/cns/company", func(w http.ResponseWriter, r *http.Request) {
	compId := r.FormValue("id")
	var company Company
	db.Get("company", compId, &company)
	if compId == "" && company.Id == "" {
		web.SetErrorRedirect(w, r, "/cns/company", "Error saving company. Please try again")
		return
		//company.Id = strconv.Itoa(int(time.Now().UnixNano()))
		//company.Password = company.Email
		//company.Role = "COMPANY"
		//company.CreateSlug()
	}
	FormToStruct(&company, r.Form, "")
	var companies []Company
	db.TestQuery("company", &companies, adb.Eq("email", company.Email), adb.Ne("id", `"`+company.Id+`"`))
	if len(companies) > 0 {
		web.SetErrorRedirect(w, r, "/cns/company/"+company.Id, "Error saving company. Email is already registered")
		return
	}
	if company.SameAddress {
		company.MailingAddress = company.PhysicalAddress
	}
	db.Set("company", company.Id, company)
	web.SetSuccessRedirect(w, r, "/cns/company/"+company.Id, "Successfully saved company")
	return
}}

var companyVehicleView = web.Route{"GET", "/cns/company/:compId/vehicle/:vId", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	var vehicle Vehicle
	compId := r.FormValue(":compId")
	if !db.Get("company", compId, &company) {
		web.SetErrorRedirect(w, r, "/cns/company", "Error finding company")
		return
	}
	if !db.Get("vehicle", r.FormValue(":vId"), &vehicle) {
		web.SetErrorRedirect(w, r, "/cns/company/"+compId, "Error finding vehicle")
		return
	}

	tc.Render(w, r, "company-vehicle-view.tmpl", web.Model{
		"company": company,
		"vehicle": vehicle,
	})
}}

var companyVehicleSave = web.Route{"POST", "/cns/company/:compId/vehicle", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	compId := r.FormValue(":compId")
	if !db.Get("company", compId, &company) {
		web.SetErrorRedirect(w, r, "/cns/company", "Error finding company")
		return
	}
	var vehicle Vehicle
	vehicleId := r.FormValue("id")
	db.Get("vehicle", vehicleId, &vehicle)
	if vehicleId == "" || vehicle.Id == "" {
		vehicle.Id = strconv.Itoa(int(time.Now().UnixNano()))
	}
	FormToStruct(&vehicle, r.Form, "")
	db.Set("vehicle", vehicle.Id, vehicle)
	web.SetSuccessRedirect(w, r, "/cns/company/"+compId+"/vehicle/"+vehicle.Id, "Successfully saved vehicle")
	return
}}

/* --- Driver Management --- */

var allDriver = web.Route{"GET", "/cns/driver", func(w http.ResponseWriter, r *http.Request) {
	var drivers []Driver
	//ok := db.Match("user", `"role":"DRIVER"`, &drivers)
	ok := db.TestQuery("driver", &drivers, adb.Eq("role", "DRIVER"))
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
		ok := db.Get("driver", driverId, &driver)
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
	db.Get("driver", driverId, &driver)
	FormToStruct(&driver, r.Form, "")
	if driver.Id == "" && driverId == "" {
		driver.Id = strconv.Itoa(int(time.Now().UnixNano()))
		driver.Password = driver.Email
		driver.Role = "DRIVER"
	}
	db.Set("driver", driver.Id, driver)
	web.SetSuccessRedirect(w, r, "/cns/driver/"+driver.Id, "Successfully saved driver")
	return
}}
