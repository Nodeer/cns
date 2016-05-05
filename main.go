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

	db.AddStore("employee")
	db.AddStore("company")
	db.AddStore("driver")
	db.AddStore("vehicle")
	db.AddStore("document")
	db.AddStore("event")
	db.AddStore("note")
	db.AddStore("comment")

	mx.AddRoutes(login, loginPost, logout, viewDocument, saveDocument, GetComment, PostComent)

	mx.AddSecureRoutes(EMPLOYEE, index)

	mx.AddSecureRoutes(EMPLOYEE, allCompany, viewCompany, saveCompany, saveNote)
	mx.AddSecureRoutes(EMPLOYEE, allEmployee, viewEmployee, saveEmployee, settings)
	mx.AddSecureRoutes(EMPLOYEE, allDriver, viewDriver, saveDriver)

	mx.AddSecureRoutes(AJAX, uploadDriverFile, addDriverDocument, viewDriverFile, delDriverFile, documentDel)

	mx.AddRoutes(calendar, calendarEvents, calendarEvent)

	web.Funcs["lower"] = strings.ToLower
	web.Funcs["size"] = PrettySize
	tc = web.NewTmplCache()
	defaultUsers()
}

// main http listener
func main() {
	log.Fatal(http.ListenAndServe(":8080", mx))
}

var logout = web.Route{"GET", "/logout", func(w http.ResponseWriter, r *http.Request) {
	redirect := "/login"
	switch web.GetRole(r) {
	case "EMPLOYEE":
		redirect = "/login"
	case "COMPANY":
		redirect = "/company/login"
	case "DRIVER":
		var company Company
		db.Get("company", web.GetSess(r, "companyId").(string), &company)
		redirect = "/login/" + company.Slug
	}
	web.Logout(w)
	web.SetSuccessRedirect(w, r, redirect, "Successfully logged out")
}}

var addDriverDocument = web.Route{"POST", "/driver/document", func(w http.ResponseWriter, r *http.Request) {
	driverId := r.FormValue("id")
	var driver Driver
	if !db.Get("driver", driverId, &driver) {
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
	db.TestQuery("document", &docs, adb.Eq("driverId", `"`+driver.Id+`"`))
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

	//db.Match("document", `"driverId":"`+r.FormValue(":driverId")+`"`, &docs)
	db.TestQuery("document", &docs, adb.Eq("driverId", `"`+r.FormValue(":driverId")+`"`))

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

var calendar = web.Route{"GET", "/cns/calendar", func(w http.ResponseWriter, r *http.Request) {
	tc.Render(w, r, "calendar.tmpl", nil)
	return
}}

var calendarEvents = web.Route{"GET", "/cns/calendar/events", func(w http.ResponseWriter, r *http.Request) {
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

var calendarEvent = web.Route{"POST", "/cns/calendar/event", func(w http.ResponseWriter, r *http.Request) {
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

var viewDocument = web.Route{"GET", "/document/:id", func(w http.ResponseWriter, r *http.Request) {
	var document Document
	var driver Driver
	var company Company
	ok := db.Get("document", r.FormValue(":id"), &document)
	if !ok {
		web.SetErrorRedirect(w, r, "/", "Error, retrieving document.")
		return
	}
	ok = db.Get("driver", document.DriverId, &driver)
	if !ok {
		web.SetErrorRedirect(w, r, "/", "Error, document is not associated with a driver.")
		return
	}
	db.Get("company", driver.CompanyId, &company)
	tc.Render(w, r, "dqf-"+document.DocumentId+".tmpl", web.Model{
		"document": document,
		"driver":   driver,
		"company":  company,
	})
	return
}}

var saveDocument = web.Route{"POST", "/document/save", func(w http.ResponseWriter, r *http.Request) {
	var document Document
	db.Get("document", r.FormValue("id"), &document)
	document.Data = r.FormValue("data")
	db.Set("document", document.Id, document)
	ajaxErrorResponse(w, `{"status":"success","msg":"Successfully saved document"}`)
	return
}}

var GetComment = web.Route{"GET", "/comment", func(w http.ResponseWriter, r *http.Request) {
	tc.Render(w, r, "comment.tmpl", web.Model{
		"return":  r.FormValue("return"),
		"comment": true,
		"page":    r.FormValue("page"),
	})
}}

var PostComent = web.Route{"POST", "/comment", func(w http.ResponseWriter, r *http.Request) {
	id := strconv.Itoa(int(time.Now().UnixNano()))
	var comment Comment
	r.ParseForm()
	FormToStruct(&comment, r.Form, "")
	comment.Id = id
	db.Set("comment", id, comment)
	web.SetSuccessRedirect(w, r, comment.Url, "Successfully added comment")
}}
