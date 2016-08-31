package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/cagnosolutions/adb"
	"github.com/cagnosolutions/web"
)

/* --- Employee Management --- */

var allEmployee = web.Route{"GET", "/cns/employee", func(w http.ResponseWriter, r *http.Request) {
	var employees []Employee
	// get all "employees" except the default logins
	db.TestQuery("employee", &employees, adb.Gt("id", `"1"`))
	tc.Render(w, r, "employee-all.tmpl", web.Model{
		"employees": employees,
	})
}}

var viewEmployee = web.Route{"GET", "/cns/employee/:id", func(w http.ResponseWriter, r *http.Request) {
	var employee Employee
	employeeId := r.FormValue(":id")
	if employeeId != "new" && !db.Get("employee", employeeId, &employee) {
		web.SetErrorRedirect(w, r, "/employee", "Error finding employee")
		return
	}
	tc.Render(w, r, "employee.tmpl", web.Model{
		"employee": employee,
	})
}}

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

var delEmployee = web.Route{"POST", "/cns/employee/:id", func(w http.ResponseWriter, r *http.Request) {
	empId := r.FormValue(":id")
	db.Del("employee", empId)
	web.SetSuccessRedirect(w, r, "/cns/employee", "Successfully deleted employee")
	return
}}

var saveHomePage = web.Route{"POST", "/cns/employee/:id/homepage", func(w http.ResponseWriter, r *http.Request) {
	var employee Employee
	db.Get("employee", r.FormValue(":id"), &employee)
	employee.Home = r.FormValue("url")
	db.Set("employee", employee.Id, employee)
	ajaxResponse(w, `{"error":false}`)
	return
}}

/* ---  --- */

var emailTemplates = web.Route{"GET", "/admin/template", func(w http.ResponseWriter, r *http.Request) {
	var emailTemplate EmailTemplate
	var emailTemplates []EmailTemplate
	db.All("emailTemplate", &emailTemplates)
	tc.Render(w, r, "email-templates.tmpl", web.Model{
		"emailTemplate":  emailTemplate,
		"emailTemplates": emailTemplates,
	})
}}

var emailTemplatesView = web.Route{"GET", "/admin/template/:id", func(w http.ResponseWriter, r *http.Request) {
	var emailTemplate EmailTemplate
	if !db.Get("emailTemplate", r.FormValue(":id"), &emailTemplate) {
		web.SetErrorRedirect(w, r, "/admin/template", "Error finding template")
		return
	}
	var emailTemplates []EmailTemplate
	db.All("emailTemplate", &emailTemplates)
	tc.Render(w, r, "email-templates.tmpl", web.Model{
		"emailTemplate":  emailTemplate,
		"emailTemplates": emailTemplates,
	})
}}

var emailTemplateSave = web.Route{"POST", "/admin/template", func(w http.ResponseWriter, r *http.Request) {
	var emailTemplate EmailTemplate
	db.Get("emailTemplate", r.FormValue("id"), &emailTemplate)
	FormToStruct(&emailTemplate, r.Form, "")
	if emailTemplate.Id == "" {
		emailTemplate.Id = strconv.Itoa(int(time.Now().UnixNano()))
	}
	var emailTemplates []EmailTemplate
	db.TestQuery("emailTemplate", &emailTemplates, adb.Eq("name", emailTemplate.Name), adb.Ne("id", `"`+emailTemplate.Id+`"`))
	if len(emailTemplates) > 0 {
		web.SetErrorRedirect(w, r, "/admin/template/"+r.FormValue("id"), "Error saving email template. Name is already in use")
		return
	}
	db.Set("emailTemplate", emailTemplate.Id, emailTemplate)
	web.SetSuccessRedirect(w, r, "/admin/template/"+emailTemplate.Id, "Successfully saved email template")
	return

}}
