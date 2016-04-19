package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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
	db.AddStore("document")
	db.AddStore("event")

	mx.AddRoutes(login, loginPost, logout, driverDocuments)

	mx.AddSecureRoutes(EMPLOYEE, index)

	mx.AddSecureRoutes(EMPLOYEE, allCompany, viewCompany, saveCompany)
	mx.AddSecureRoutes(EMPLOYEE, allEmployee, viewEmployee, saveEmployee, settings)
	mx.AddSecureRoutes(EMPLOYEE, allDriver, uploadDriverFile, addDriverDocument, viewDriver, savedriver, viewDriverFile, delDriverFile, documentDel)

	mx.AddRoutes(calendar, calendarEvents, calendarEvent)

	web.Funcs["lower"] = strings.ToLower
	web.Funcs["size"] = PrettySize
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
	redirect := "/company"
	if employee.Home != "" {
		redirect = employee.Home
	}
	web.SetSuccessRedirect(w, r, redirect, "Welcome "+employee.FirstName)
	return
}}

var logout = web.Route{"GET", "/logout", func(w http.ResponseWriter, r *http.Request) {
	web.Logout(w)
	web.SetSuccessRedirect(w, r, "/login", "Successfully logged out")
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

var settings = web.Route{"GET", "/settings", func(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/employee/"+web.GetSess(r, "id").(string), 303)
}}

var saveEmployee = web.Route{"POST", "/employee", func(w http.ResponseWriter, r *http.Request) {
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
	//db.Query("user", &users, "email="+employee.Email, "id^"+employee.Id)
	db.TestQuery("users", &users, adb.Eq("email", employee.Email), adb.Ne("id", employee.Id))
	fmt.Println(users)
	if len(users) > 0 {
		web.SetErrorRedirect(w, r, "/employee/"+employee.Id, "Error saving employee. Email is already registered")
		return
	}
	db.Set("user", employee.Id, employee)
	web.SetSuccessRedirect(w, r, "/employee/"+employee.Id, "Successfully saved employee")
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
	var drivers []Driver
	compId := r.FormValue(":id")
	if compId != "add" {
		ok := db.Get("user", compId, &company)
		if !ok || company.Role != "COMPANY" {
			web.SetErrorRedirect(w, r, "/company", "Error finding company")
			return
		}
		//db.Match("user", `"companyId":"`+company.Id+`"`, &drivers)
		db.Query("user", &drivers, "companyId="+company.Id)
	}
	tc.Render(w, r, "company.tmpl", web.Model{
		"company": company,
		"drivers": drivers,
	})
}}

var saveCompany = web.Route{"POST", "/company", func(w http.ResponseWriter, r *http.Request) {
	compId := r.FormValue("id")
	var company Company
	db.Get("user", compId, &company)
	if compId == "" && company.Id == "" {
		company.Id = strconv.Itoa(int(time.Now().UnixNano()))
		company.Password = company.Email
		company.Role = "COMPANY"
	}
	FormToStruct(&company, r.Form, "")
	db.Set("user", company.Id, company)
	web.SetSuccessRedirect(w, r, "/company/"+company.Id, "Successfully saved company")
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
		db.Match("document", `"driverId":"`+driver.Id+`"`, &docs)
	}
	tc.Render(w, r, "driver.tmpl", web.Model{
		"driver": driver,
		"files":  files,
		"dqfs":   DQFS,
		"docs":   docs,
	})
}}

var savedriver = web.Route{"POST", "/driver", func(w http.ResponseWriter, r *http.Request) {
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
	web.SetSuccessRedirect(w, r, "/driver/"+driver.Id, "Successfully saved driver")
	return
}}

var addDriverDocument = web.Route{"POST", "/driver/document", func(w http.ResponseWriter, r *http.Request) {
	driverId := r.FormValue("id")
	var driver Driver
	if !db.Get("user", driverId, &driver) {
		ajaxErrorResponse(w, `{"status":"error", "msg":"Error adding documents"}`)
		return
	}
	docIds := strings.Split(r.FormValue("docIds"), ",")
	for _, docId := range docIds {
		id := strconv.Itoa(int(time.Now().UnixNano()))
		doc := Document{
			Id:         id,
			Name:       "dqf-" + docId,
			DocumentId: docId,
			Complete:   false,
			CompanyId:  driver.CompanyId,
			DriverId:   driver.Id,
		}
		db.Add("document", id, doc)
	}
	var docs []Document
	db.Match("document", `"driverId":"`+driver.Id+`"`, &docs)
	resp := make(map[string]interface{}, 0)
	resp["status"] = "success"
	resp["msg"] = "Successfully added documents"
	resp["data"] = docs
	b, err := json.Marshal(resp)
	if err != nil {
		ajaxErrorResponse(w, `{"status":"error","type":"marshal","msg":"Successfully added documents. Please refresh the page to view"}`)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", b)
	return
}}

var documentDel = web.Route{"POST", "/document/del/:driverId/:docId", func(w http.ResponseWriter, r *http.Request) {
	db.Del("document", r.FormValue(":docId"))

	var docs []Document

	db.Match("document", `"driverId":"`+r.FormValue(":driverId")+`"`, &docs)

	resp := make(map[string]interface{}, 0)
	resp["status"] = "success"
	resp["msg"] = "Successfully deleted document"
	resp["data"] = docs
	b, err := json.Marshal(resp)

	if err != nil {
		ajaxErrorResponse(w, `{"status":"error","type":"marshal","msg":"Successfully deleted document. Please refresh the page to view"}`)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", b)
	return
}}

var uploadDriverFile = web.Route{"POST", "/driver/upload", func(w http.ResponseWriter, r *http.Request) {
	driverId := r.FormValue("id")
	if driverId == "" {
		log.Printf("main.go -> uploadDriverFile() -> os.MkdirAll() -> no dirver id specified")
		ajaxErrorResponse(w, `{"status":"error","msg":"Error uploading file"}`)
		return
	}
	path := "upload/driver/" + driverId + "/"
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Printf("main.go -> uploadDriverFile() -> os.MkdirAll() -> %v\n", err)
		ajaxErrorResponse(w, `{"status":"error","msg":"Error uploading file"}`)
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Printf("main.go -> uploadDriverFile() -> r.FormFile() -> %v\n", err)
		ajaxErrorResponse(w, `{"status":"error","msg":"Error uploading file `+handler.Filename+`"}`)
		return
	}
	defer file.Close()
	f, err := os.OpenFile(path+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("main.go -> uploadDriverFile() -> os.OpenFile() -> %v\n", err)
		ajaxErrorResponse(w, `{"status":"error","msg":"Error uploading file `+handler.Filename+`}`)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		ajaxErrorResponse(w, `{"status":"error","msg":"Successfully uploaded file `+handler.Filename+`. Please refresh the page to view","type":"readDir"}`)
		return
	}

	var fileInfo []map[string]interface{}

	for _, file := range files {
		f := make(map[string]interface{})
		f["name"] = file.Name()
		f["size"] = file.Size()
		fileInfo = append(fileInfo, f)
	}

	resp := make(map[string]interface{}, 0)
	resp["status"] = "success"
	resp["msg"] = "Successfully uploaded file " + handler.Filename
	resp["data"] = fileInfo
	b, err := json.Marshal(resp)
	if err != nil {
		ajaxErrorResponse(w, `{"status":"error","msg":"Successfully uploaded file `+handler.Filename+`. Please refresh the page to view","type":"marshal"}`)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", b)
	return
}}

var viewDriverFile = web.Route{"GET", "/driver/file/:id/:file", func(w http.ResponseWriter, r *http.Request) {
	server := http.StripPrefix("/driver/file/", http.FileServer(http.Dir("upload/driver/")))
	server.ServeHTTP(w, r)
}}

var delDriverFile = web.Route{"POST", "/driver/file/:id/:file", func(w http.ResponseWriter, r *http.Request) {
	driverId := r.FormValue(":id")
	filename := r.FormValue(":file")
	os.Remove("upload/driver/" + driverId + "/" + filename)

	files, err := ioutil.ReadDir("upload/driver/" + driverId)
	if err != nil {
		ajaxErrorResponse(w, `{"status":"error","msg":"Successfully uploaded file `+filename+`. Please refresh the page to view","type":"readDir"}`)
		return
	}

	var fileInfo []map[string]interface{}

	for _, file := range files {
		f := make(map[string]interface{})
		f["name"] = file.Name()
		f["size"] = file.Size()
		fileInfo = append(fileInfo, f)
	}

	resp := make(map[string]interface{}, 0)
	resp["status"] = "success"
	resp["msg"] = "Successfully uploaded file " + filename
	resp["data"] = fileInfo
	b, err := json.Marshal(resp)
	if err != nil {
		ajaxErrorResponse(w, `{"status":"error","msg":"Successfully uploaded file `+filename+`. Please refresh the page to view","type":"marshal"}`)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", b)
	return
}}

var calendar = web.Route{"GET", "/calendar", func(w http.ResponseWriter, r *http.Request) {
	tc.Render(w, r, "calendar.tmpl", nil)
	return
}}

var calendarEvents = web.Route{"GET", "/calendar/events", func(w http.ResponseWriter, r *http.Request) {
	/*
		var events []string
		for i := 0; i < 5; i++ {
			events = append(events, fmt.Sprintf(`{"title":"Event #%d","start":%q,"allDay":true}`, i, time.Now().AddDate(0, 0, i).Format(time.RFC3339)))
		}
	*/
	var events []Event
	db.All("event", &events)
	b, err := json.Marshal(events)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Fprintf(w, "%s", b)
	return
}}

var calendarEvent = web.Route{"POST", "/calendar/event", func(w http.ResponseWriter, r *http.Request) {
	event := Event{
		Id:     r.FormValue("id"),
		Title:  r.FormValue("title"),
		AllDay: true,
	}
	t, err := time.Parse("2006-01-02", r.FormValue("start"))
	if err != nil {
		log.Println("error parsing the time...")
		ajaxErrorResponse(w, `{"err":true,"code":500,"msg":"There was an issue saving the event to the database"}`)
		return
	}
	event.Start = t
	if !db.Add("event", event.Id, event) {
		ajaxErrorResponse(w, `{"err":true,"code":500,"msg":"There was an issue saving the event to the database"}`)
		return
	}
	ajaxErrorResponse(w, `{"err":false,"code":200,"msg":"Successfully added the event to the database"}`)
	return
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
