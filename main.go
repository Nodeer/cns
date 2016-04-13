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

	mx.AddSecureRoutes(EMPLOYEE, index)

	mx.AddSecureRoutes(EMPLOYEE, allCompany, viewCompany, saveCompany)
	mx.AddSecureRoutes(EMPLOYEE, allEmployee, viewEmployee, saveEmployee)
	mx.AddSecureRoutes(EMPLOYEE, allDriver, uploadDriverFile, addDriverDocument, viewDriver, savedriver, viewDriverFile)

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
	driverId := r.FormValue(":id")
	ok := db.Get("user", driverId, &driver)
	if !ok || driver.Role != "DRIVER" {
		web.SetErrorRedirect(w, r, "/driver", "Error finding driver")
		return
	}
	var files []map[string]interface{}
	if fileInfos, err := ioutil.ReadDir("upload/driver/" + driverId); err == nil {
		for _, fileInfo := range fileInfos {
			var info = make(map[string]interface{})
			info["name"] = fileInfo.Name()
			info["size"] = fileInfo.Size()
			files = append(files, info)
		}
	}
	var docs []Document
	db.Match("document", `"driverId":"`+driver.Id+`"`, &docs)

	tc.Render(w, r, "driver.tmpl", web.Model{
		"driver": driver,
		"files":  files,
		"dqfs":   DQFS,
		"docs":   docs,
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
	resp["files"] = fileInfo
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
	resp["docs"] = docs
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
