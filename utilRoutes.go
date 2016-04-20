package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cagnosolutions/web"
)

const (
	CperE = 10
	COMP  = 500
	DperC = 20
)

func init() {
	mx.AddRoutes(makeUsers, upload, buttons, uploader, notify, alert, form)
}

var makeUsers = web.Route{"GET", "/makeUsers", func(w http.ResponseWriter, r *http.Request) {
	MakeEmployees()
	compIds := MakeCompanies()
	MakeDrivers(compIds)
	web.SetSuccessRedirect(w, r, "/", "Successfully made users")
	return
}}

func MakeEmployees() {
	for i := 0; i < (COMP / CperE); i++ {
		id := strconv.Itoa(int(time.Now().UnixNano()))

		user := Employee{
			FirstName: "John",
			LastName:  fmt.Sprintf("Smith the %dth", (i + 4)),
			Phone:     fmt.Sprintf("717-777-777%d", i),
		}

		user.Id = id
		user.Email = fmt.Sprintf("%d@cns.com", i)
		user.Password = fmt.Sprintf("Password-%d", i)
		user.Active = (i%2 == 0)
		user.Role = "EMPLOYEE"

		user.Street = fmt.Sprintf("12%d Main Street", 1)
		user.City = fmt.Sprintf("%dville", i)
		user.State = fmt.Sprintf("%d state", i)
		user.Zip = fmt.Sprintf("1234%d", i)

		db.Add("user", id, user)
	}
}

func MakeCompanies() [COMP]string {
	compIds := [COMP]string{}
	for i := 0; i < COMP; i++ {
		id := strconv.Itoa(int(time.Now().UnixNano()))
		compIds[i] = id
		user := Company{
			Name:    fmt.Sprintf("Company %d", i),
			Contact: fmt.Sprintf("Bobbi Sue the %dth", (i + 4)),
			Phone:   fmt.Sprintf("717-777-777%d", i),
		}

		user.Id = id
		user.Email = fmt.Sprintf("%d@company%d.com", i, i)
		user.Password = fmt.Sprintf("Password-%d", i)
		user.Active = (i%2 == 0)
		user.Role = "COMPANY"

		user.Street = fmt.Sprintf("12%d Main Street", 1)
		user.City = fmt.Sprintf("%dville", i)
		user.State = fmt.Sprintf("%d state", i)
		user.Zip = fmt.Sprintf("1234%d", i)
		user.CreateSlug()

		db.Add("user", id, user)
	}
	return compIds
}

func MakeDrivers(compIds [COMP]string) {
	for i := 0; i < (COMP * DperC); i++ {
		id := strconv.Itoa(int(time.Now().UnixNano()))
		compIdx := i / DperC
		user := Driver{
			FirstName:    "Daniel",
			LastName:     fmt.Sprintf("Jones the %dth", (i + 4)),
			Phone:        fmt.Sprintf("717-777-777%d", i),
			DOB:          fmt.Sprintf("198%d-03-1%d", i, i),
			LicenseNum:   fmt.Sprintf("1234567%d", i),
			LicenseState: fmt.Sprintf("%d state", i),
			CompanyId:    compIds[compIdx],
		}

		user.Id = id
		user.Email = fmt.Sprintf("%d@%d.com", i, i)
		user.Password = fmt.Sprintf("Password-%d", i)
		user.Active = (i%2 == 0)
		user.Role = "DRIVER"

		user.Street = fmt.Sprintf("12%d Main Street", 1)
		user.City = fmt.Sprintf("%dville", i)
		user.State = fmt.Sprintf("%d state", i)
		user.Zip = fmt.Sprintf("1234%d", i)

		db.Add("user", id, user)
	}
}

var makeEmployees = web.Route{"GET", "/makeEmployees", func(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 10; i++ {
		id := strconv.Itoa(int(time.Now().UnixNano()))

		user := Employee{
			FirstName: "John",
			LastName:  fmt.Sprintf("Smith the %dth", (i + 4)),
			Phone:     fmt.Sprintf("717-777-777%d", i),
		}

		user.Id = id
		user.Email = fmt.Sprintf("%d@cns.com", i)
		user.Password = fmt.Sprintf("Password-%d", i)
		user.Active = (i%2 == 0)
		user.Role = "EMPLOYEE"

		user.Street = fmt.Sprintf("12%d Main Street", 1)
		user.City = fmt.Sprintf("%dville", i)
		user.State = fmt.Sprintf("%d state", i)
		user.Zip = fmt.Sprintf("1234%d", i)

		db.Add("user", id, user)
	}
	web.SetSuccessRedirect(w, r, "/", "success")
	return
}}

var makeCompanies = web.Route{"GET", "/makeComps", func(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 10; i++ {
		id := strconv.Itoa(int(time.Now().UnixNano()))

		user := Company{
			Name:    fmt.Sprintf("Company %d", i),
			Contact: fmt.Sprintf("Bobbi Sue the %dth", (i + 4)),
			Phone:   fmt.Sprintf("717-777-777%d", i),
		}

		user.Id = id
		user.Email = fmt.Sprintf("%d@company%d.com", i, i)
		user.Password = fmt.Sprintf("Password-%d", i)
		user.Active = (i%2 == 0)
		user.Role = "COMPANY"

		user.Street = fmt.Sprintf("12%d Main Street", 1)
		user.City = fmt.Sprintf("%dville", i)
		user.State = fmt.Sprintf("%d state", i)
		user.Zip = fmt.Sprintf("1234%d", i)

		user.CreateSlug()
		db.Add("user", id, user)
	}
	web.SetSuccessRedirect(w, r, "/", "success")
	return
}}

var makeDrivers = web.Route{"GET", "/makeDrive", func(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 10; i++ {
		id := strconv.Itoa(int(time.Now().UnixNano()))

		user := Driver{
			FirstName:    "Daniel",
			LastName:     fmt.Sprintf("Jones the %dth", (i + 4)),
			Phone:        fmt.Sprintf("717-777-777%d", i),
			DOB:          fmt.Sprintf("198%d-03-1%d", i, i),
			LicenseNum:   fmt.Sprintf("1234567%d", i),
			LicenseState: fmt.Sprintf("%d state", i),
		}

		user.Id = id
		user.Email = fmt.Sprintf("%d@%d.com", i, i)
		user.Password = fmt.Sprintf("Password-%d", i)
		user.Active = (i%2 == 0)
		user.Role = "DRIVER"

		user.Street = fmt.Sprintf("12%d Main Street", 1)
		user.City = fmt.Sprintf("%dville", i)
		user.State = fmt.Sprintf("%d state", i)
		user.Zip = fmt.Sprintf("1234%d", i)

		db.Add("user", id, user)
	}
	web.SetSuccessRedirect(w, r, "/", "success")
	return
}}

var upload = web.Route{"GET", "/upload", func(w http.ResponseWriter, r *http.Request) {
	tc.Render(w, r, "form-uploads.tmpl", nil)
}}

var buttons = web.Route{"GET", "/buttons", func(w http.ResponseWriter, r *http.Request) {
	tc.Render(w, r, "buttons.tmpl", web.Model{})
}}

var uploader = web.Route{"POST", "/upd", func(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		println("ERROR")
	}
	fmt.Printf("%v\n", r.MultipartForm.File)
	fmt.Println(r.FormValue("id"))
}}

var notify = web.Route{"GET", "/notify", func(w http.ResponseWriter, r *http.Request) {
	tc.Render(w, r, "notification.tmpl", nil)
}}

var alert = web.Route{"GET", "/alert", func(w http.ResponseWriter, r *http.Request) {
	tc.Render(w, r, "sweet-alert.tmpl", nil)
}}

var form = web.Route{"GET", "/form", func(w http.ResponseWriter, r *http.Request) {
	tc.Render(w, r, "form-advanced.tmpl", nil)
}}
