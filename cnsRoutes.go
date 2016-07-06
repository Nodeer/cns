package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
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
	db.TestQuery("employee", &employees, adb.Gt("id", `"1"`))
	tc.Render(w, r, "employee-all.tmpl", web.Model{
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

// var settings = web.Route{"GET", "/cns/settings", func(w http.ResponseWriter, r *http.Request) {
// 	http.Redirect(w, r, "/cns/employee/"+web.GetSess(r, "id").(string), 303)
// }}

var saveEmployee = web.Route{"POST", "/cns/employee", func(w http.ResponseWriter, r *http.Request) {
	empId := r.FormValue("id")
	var employee Employee
	db.Get("employee", empId, &employee)
	FormToStruct(&employee, r.Form, "")
	if employee.Id == "" && empId == "" {
		employee.Id = strconv.Itoa(int(time.Now().UnixNano()))
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
	if !db.Get("company", compId, &company) && compId != "new" {
		web.SetErrorRedirect(w, r, "/cns/company", "Error finding company")
		return
	}
	var notes NoteRevSort
	var employees []Employee
	db.TestQuery("note", &notes, adb.Eq("companyId", `"`+company.Id+`"`))
	sort.Stable(notes)
	db.All("employee", &employees)
	tc.Render(w, r, "company.tmpl", web.Model{
		"company":       company,
		"notes":         notes,
		"employees":     employees,
		"quickNotes":    quickNotes,
		"userId":        web.GetSess(r, "id"),
		"companyConsts": GetCompanyConsts(),
	})
}}

var companyService = web.Route{"GET", "/cns/company/:id/service", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	compId := r.FormValue(":id")
	if !db.Get("company", compId, &company) {
		web.SetErrorRedirect(w, r, "/cns/company", "Error finding company")
		return
	}
	tc.Render(w, r, "company-service.tmpl", web.Model{
		"company": company,
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

var companyForm = web.Route{"GET", "/cns/company/:id/form", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	companyId := r.FormValue(":id")
	var docs []Document
	if !db.Get("company", companyId, &company) {
		web.SetErrorRedirect(w, r, "/cns/company", "Error finding company")
		return
	}
	db.TestQuery("document", &docs, adb.Eq("companyId", `"`+company.Id+`"`), adb.Eq("stateForm", "true"))
	// var d []Document
	// for _, doc := range docs {
	// 	if doc.DriverId == "" {
	// 		d = append(d, doc)
	// 	}
	// }
	var vehicles []Vehicle
	db.TestQuery("vehicle", &vehicles, adb.Eq("companyId", `"`+company.Id+`"`))
	tc.Render(w, r, "company-form.tmpl", web.Model{
		"company":  company,
		"docs":     docs,
		"forms":    CompanyForms,
		"vehicles": vehicles,
	})

}}

var companySaveNote = web.Route{"POST", "/cns/company/note", func(w http.ResponseWriter, r *http.Request) {
	var note Note
	r.ParseForm()
	FormToStruct(&note, r.Form, "")
	if note.Id == "" {
		note.Id = strconv.Itoa(int(time.Now().UnixNano()))
	}
	dt, err := time.Parse("01/02/2006 3:04 PM", r.FormValue("dateTime"))
	if err != nil {
		log.Printf("cnsRoutes.go >> companySaveNotes >> time.Parse() >> %v\n", err)
	}
	note.StartTime = dt.Unix()
	note.EndTime = dt.Unix()
	note.StartTimePretty = r.FormValue("dateTime")
	note.EndTimePretty = r.FormValue("dateTime")
	db.Set("note", note.Id, note)
	web.SetSuccessRedirect(w, r, "/cns/company/"+r.FormValue("companyId"), "Successfully saved note")
}}

var companySave = web.Route{"POST", "/cns/company", func(w http.ResponseWriter, r *http.Request) {
	compId := r.FormValue("id")
	var company Company
	db.Get("company", compId, &company)
	if compId == "" && company.Id == "" {
		//web.SetErrorRedirect(w, r, "/cns/company", "Error saving company. Please try again")
		//return
		company.Id = strconv.Itoa(int(time.Now().UnixNano()))
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
	if r.FormValue("from") == "vehicle" {
		web.SetSuccessRedirect(w, r, "/cns/company/"+company.Id+"/vehicle", "Successfully updated insurance information")
		return
	}
	if r.FormValue("from") == "service" {
		web.SetSuccessRedirect(w, r, "/cns/company/"+company.Id+"/service", "Successfully updated service information")
		return
	}
	web.SetSuccessRedirect(w, r, "/cns/company/"+company.Id, "Successfully saved company")
	return
}}

var companyVehicleView = web.Route{"GET", "/cns/company/:compId/vehicle/:vId", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	compId := r.FormValue(":compId")
	if !db.Get("company", compId, &company) {
		web.SetErrorRedirect(w, r, "/cns/company", "Error finding company")
		return
	}
	vehicleId := r.FormValue(":vId")
	var vehicle Vehicle
	if !db.Get("vehicle", vehicleId, &vehicle) && vehicleId != "new" {
		web.SetErrorRedirect(w, r, "/cns/company/"+compId+"/vehicle", "Error finding vehicle")
		return
	}

	tc.Render(w, r, "company-vehicle-view.tmpl", web.Model{
		"company":       company,
		"vehicle":       vehicle,
		"vehicleConsts": GetVehicleConsts(),
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
		vehicle.CompanyId = compId
	}
	FormToStruct(&vehicle, r.Form, "")
	db.Set("vehicle", vehicle.Id, vehicle)
	web.SetSuccessRedirect(w, r, "/cns/company/"+compId+"/vehicle/"+vehicle.Id, "Successfully saved vehicle")
	return
}}

var companyAddForm = web.Route{"POST", "/cns/company/:id/form", func(w http.ResponseWriter, r *http.Request) {
	compId := r.FormValue(":id")
	var company Company
	if !db.Get("company", compId, &company) {
		web.SetErrorRedirect(w, r, "/cns/company/", "Error finding company")
		return
	}

	id := strconv.Itoa(int(time.Now().UnixNano()))
	docId := r.FormValue("name")
	var vehicleIds []string
	if r.FormValue("vehicleIds") != "" {
		vehicleIds = strings.Split(r.FormValue("vehicleIds"), ",")
	}
	doc := Document{
		Id:         id,
		Name:       docId,
		DocumentId: strings.ToLower(strings.Replace(docId, " ", "_", -1)),
		Complete:   false,
		CompanyId:  compId,
		VehicleIds: vehicleIds,
		StateForm:  true,
	}
	db.Add("document", id, doc)

	// docIds := strings.Split(r.FormValue("docIds"), ",")
	// for _, docId := range docIds {
	// 	id := strconv.Itoa(int(time.Now().UnixNano()))
	// 	doc := Document{
	// 		Id:         id,
	// 		Name:       docId,
	// 		DocumentId: strings.Replace(docId, " ", "_", -1),
	// 		Complete:   false,
	// 		CompanyId:  compId,
	// 	}
	// 	db.Add("document", id, doc)
	// }

	web.SetSuccessRedirect(w, r, "/cns/company/"+company.Id+"/form", "Successfully added forms")
	return
}}

var companyFormDel = web.Route{"POST", "/cns/company/:companyId/form/:formId", func(w http.ResponseWriter, r *http.Request) {
	var form Document
	if !db.Get("document", r.FormValue(":formId"), &form) || form.CompanyId != r.FormValue(":companyId") {
		web.SetErrorRedirect(w, r, "/cns/company/"+r.FormValue(":companyId")+"/form", "Error deleting from")
		return
	}
	db.Del("document", form.Id)
	web.SetSuccessRedirect(w, r, "/cns/company/"+r.FormValue(":companyId")+"/form", "Successfully deleted form")
	return

}}

var testCompanyFormView = web.Route{"GET", "/test/form/:id", func(w http.ResponseWriter, r *http.Request) {
	var form Document
	db.Get("document", r.FormValue(":id"), &form)
	b, err := json.Marshal(form)
	if err != nil {
		fmt.Fprintf(w, "Error marshaling json")
		return
	}
	fmt.Fprintf(w, "%s", b)
	return
}}

/* --- Driver Management --- */

var allDriver = web.Route{"GET", "/cns/driver", func(w http.ResponseWriter, r *http.Request) {
	var drivers []Driver
	ok := db.TestQuery("driver", &drivers, adb.Eq("role", "DRIVER"))
	if !ok {
		fmt.Println("error")
	}
	tc.Render(w, r, "driver-all.tmpl", web.Model{
		"drivers": drivers,
	})
}}

var viewDriver = web.Route{"GET", "/cns/driver/:id", func(w http.ResponseWriter, r *http.Request) {
	var driver Driver
	driverId := r.FormValue(":id")
	if !db.Get("driver", driverId, &driver) && driverId != "new" {
		web.SetErrorRedirect(w, r, "/driver", "Error finding driver")
		return
	}
	companyId := r.FormValue("cid")
	var company Company
	if driverId == "new" && r.FormValue("cid") == "" {
		web.SetErrorRedirect(w, r, "/cns/company", "Error adding new driver. Please try again")
		return
	} else {
		db.Get("company", driver.CompanyId, &company)

	}

	tc.Render(w, r, "driver.tmpl", web.Model{
		"driver":    driver,
		"companyId": companyId,
		"company":   company,
	})
}}

var driverForms = web.Route{"GET", "/cns/driver/:id/form", func(w http.ResponseWriter, r *http.Request) {
	var driver Driver
	driverId := r.FormValue(":id")
	var docs []Document
	ok := db.Get("driver", driverId, &driver)
	if !ok || driver.Role != "DRIVER" {
		web.SetErrorRedirect(w, r, "/driver", "Error finding driver")
		return
	}
	db.TestQuery("document", &docs, adb.Eq("driverId", `"`+driver.Id+`"`))
	tc.Render(w, r, "driver-form.tmpl", web.Model{
		"driver": driver,
		"dqfs":   DQFS,
		"docs":   docs,
	})

}}

var driverFiles = web.Route{"GET", "/cns/driver/:id/file", func(w http.ResponseWriter, r *http.Request) {
	var driver Driver
	driverId := r.FormValue(":id")
	var files []map[string]interface{}
	if !db.Get("driver", driverId, &driver) {
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
	tc.Render(w, r, "driver-file.tmpl", web.Model{
		"driver": driver,
		"files":  files,
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
