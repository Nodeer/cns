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

		employee := Employee{
			FirstName: "John",
			LastName:  fmt.Sprintf("Smith the %dth", (i + 4)),
			Phone:     fmt.Sprintf("717-777-777%d", i),
		}

		employee.Id = id
		employee.Email = fmt.Sprintf("%d@cns.com", i)
		employee.Password = fmt.Sprintf("Password-%d", i)
		employee.Active = (i%2 == 0)
		employee.Role = "EMPLOYEE"

		employee.Street = fmt.Sprintf("12%d Main Street", 1)
		employee.City = fmt.Sprintf("%dville", i)
		employee.State = fmt.Sprintf("%d state", i)
		employee.Zip = fmt.Sprintf("1234%d", i)

		db.Add("employee", id, employee)
	}
}

func MakeCompanies() [COMP]string {
	compIds := [COMP]string{}
	for i := 0; i < COMP; i++ {
		id := strconv.Itoa(int(time.Now().UnixNano()))
		compIds[i] = id
		company := Company{
			Name:    fmt.Sprintf("Company %d", i),
			Contact: fmt.Sprintf("Bobbi Sue the %dth", (i + 4)),
			Phone:   fmt.Sprintf("717-777-777%d", i),
		}

		company.Id = id
		company.Email = fmt.Sprintf("%d@company%d.com", i, i)
		company.Password = fmt.Sprintf("Password-%d", i)
		company.Active = (i%2 == 0)
		company.Role = "COMPANY"

		company.Street = fmt.Sprintf("12%d Main Street", 1)
		company.City = fmt.Sprintf("%dville", i)
		company.State = fmt.Sprintf("%d state", i)
		company.Zip = fmt.Sprintf("1234%d", i)
		company.CreateSlug()

		db.Add("company", id, company)
	}
	return compIds
}

func MakeDrivers(compIds [COMP]string) {
	for i := 0; i < (COMP * DperC); i++ {
		id := strconv.Itoa(int(time.Now().UnixNano()))
		compIdx := i / DperC
		driver := Driver{
			FirstName:    "Daniel",
			LastName:     fmt.Sprintf("Jones the %dth", (i + 4)),
			Phone:        fmt.Sprintf("717-777-777%d", i),
			DOB:          fmt.Sprintf("198%d-03-1%d", i, i),
			LicenseNum:   fmt.Sprintf("1234567%d", i),
			LicenseState: fmt.Sprintf("%d state", i),
			CompanyId:    compIds[compIdx],
		}

		driver.Id = id
		driver.Email = fmt.Sprintf("%d@%d.com", i, i)
		driver.Password = fmt.Sprintf("Password-%d", i)
		driver.Active = (i%2 == 0)
		driver.Role = "DRIVER"

		driver.Street = fmt.Sprintf("12%d Main Street", 1)
		driver.City = fmt.Sprintf("%dville", i)
		driver.State = fmt.Sprintf("%d state", i)
		driver.Zip = fmt.Sprintf("1234%d", i)

		db.Add("driver", id, driver)
	}
}

var makeEmployees = web.Route{"GET", "/makeEmployees", func(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 10; i++ {
		id := strconv.Itoa(int(time.Now().UnixNano()))

		employee := Employee{
			FirstName: "John",
			LastName:  fmt.Sprintf("Smith the %dth", (i + 4)),
			Phone:     fmt.Sprintf("717-777-777%d", i),
		}

		employee.Id = id
		employee.Email = fmt.Sprintf("%d@cns.com", i)
		employee.Password = fmt.Sprintf("Password-%d", i)
		employee.Active = (i%2 == 0)
		employee.Role = "EMPLOYEE"

		employee.Street = fmt.Sprintf("12%d Main Street", 1)
		employee.City = fmt.Sprintf("%dville", i)
		employee.State = fmt.Sprintf("%d state", i)
		employee.Zip = fmt.Sprintf("1234%d", i)

		db.Add("employee", id, employee)
	}
	web.SetSuccessRedirect(w, r, "/", "success")
	return
}}

var makeCompanies = web.Route{"GET", "/makeComps", func(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 10; i++ {
		id := strconv.Itoa(int(time.Now().UnixNano()))

		company := Company{
			Name:    fmt.Sprintf("Company %d", i),
			Contact: fmt.Sprintf("Bobbi Sue the %dth", (i + 4)),
			Phone:   fmt.Sprintf("717-777-777%d", i),
		}

		company.Id = id
		company.Email = fmt.Sprintf("%d@company%d.com", i, i)
		company.Password = fmt.Sprintf("Password-%d", i)
		company.Active = (i%2 == 0)
		company.Role = "COMPANY"

		company.Street = fmt.Sprintf("12%d Main Street", 1)
		company.City = fmt.Sprintf("%dville", i)
		company.State = fmt.Sprintf("%d state", i)
		company.Zip = fmt.Sprintf("1234%d", i)

		company.CreateSlug()
		db.Add("company", id, company)
	}
	web.SetSuccessRedirect(w, r, "/", "success")
	return
}}

var makeDrivers = web.Route{"GET", "/makeDrive", func(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 10; i++ {
		id := strconv.Itoa(int(time.Now().UnixNano()))

		driver := Driver{
			FirstName:    "Daniel",
			LastName:     fmt.Sprintf("Jones the %dth", (i + 4)),
			Phone:        fmt.Sprintf("717-777-777%d", i),
			DOB:          fmt.Sprintf("198%d-03-1%d", i, i),
			LicenseNum:   fmt.Sprintf("1234567%d", i),
			LicenseState: fmt.Sprintf("%d state", i),
		}

		driver.Id = id
		driver.Email = fmt.Sprintf("%d@%d.com", i, i)
		driver.Password = fmt.Sprintf("Password-%d", i)
		driver.Active = (i%2 == 0)
		driver.Role = "DRIVER"

		driver.Street = fmt.Sprintf("12%d Main Street", 1)
		driver.City = fmt.Sprintf("%dville", i)
		driver.State = fmt.Sprintf("%d state", i)
		driver.Zip = fmt.Sprintf("1234%d", i)

		db.Add("driver", id, driver)
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
