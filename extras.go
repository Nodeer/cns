package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cagnosolutions/web"
)

const (
	CperE = 1
	COMP  = 10
	DperC = 10
	VperC = 5
)

func init() {
	mx.AddRoutes(makeUsers, upload, buttons, uploader, notify, alert, form)
}

func defaultUsers() {

	developer := Employee{
		FirstName: "developer",
		LastName:  "developer",
	}

	developer.Id = "0"
	developer.Email = "developer@cns.com"
	developer.Password = "developer"
	developer.Active = true
	developer.Role = "DEVELOPER"

	admin := Employee{
		FirstName: "admin",
		LastName:  "admin",
	}

	admin.Id = "1"
	admin.Email = "admin@cns.com"
	admin.Password = "admin"
	admin.Active = true
	admin.Role = "ADMIN"

	db.Set("employee", "0", developer)
	db.Set("employee", "1", admin)
}

var makeUsers = web.Route{"GET", "/makeUsers", func(w http.ResponseWriter, r *http.Request) {
	MakeEmployees()
	compIds := MakeCompanies()
	MakeDrivers(compIds)
	MakeVehicles(compIds)
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

		company := Company{}
		company.Id = id
		company.Name = fmt.Sprintf("Company %d", i)
		company.ContactName = fmt.Sprintf("Bobbi Sue the %dth", (i + 4))
		company.ContactTitle = fmt.Sprintf("Worker #%d", i)
		company.ContactPhone = fmt.Sprintf("717-777-777%d", i)
		company.DOTNum = "DOT#" + id

		company.PhysicalAddress.Street = fmt.Sprintf("12%d Main Street", i)
		company.PhysicalAddress.City = fmt.Sprintf("%dville", i)
		company.PhysicalAddress.State = fmt.Sprintf("%d state", i)
		company.PhysicalAddress.Zip = fmt.Sprintf("1234%d", i)
		company.BusinessType = 1
		company.CarrierType = 2
		//company.CreateSlug()
		if i%2 == 0 {
			company.SameAddress = true

			company.MailingAddress.Street = fmt.Sprintf("12%d Main Street", i)
			company.MailingAddress.City = fmt.Sprintf("%dville", i)
			company.MailingAddress.State = fmt.Sprintf("%d state", i)
			company.MailingAddress.Zip = fmt.Sprintf("1234%d", i)

			company.SecondName = fmt.Sprintf("Terry Robinson the %dth", (i + 4))
			company.SecondTitle = fmt.Sprintf("Worker #%d", i+1)
			company.SecondPhone = fmt.Sprintf("717-965-435%d", i)

		} else {
			company.SameAddress = false

			company.MailingAddress.Street = fmt.Sprintf("12%d Main Street", i*10)
			company.MailingAddress.City = fmt.Sprintf("%dville", i*10)
			company.MailingAddress.State = fmt.Sprintf("%d state", i*10)
			company.MailingAddress.Zip = fmt.Sprintf("123%d", i*10)
		}

		company.MCNum = fmt.Sprintf("MC#%d", i*654)
		company.PUCNum = fmt.Sprintf("PUC#%d", i*789)
		company.Fax = fmt.Sprintf("515-555-555%d", i)
		company.Email = fmt.Sprintf("%d@company%d.com", i, i)
		company.EINNum = fmt.Sprintf("EIN#%d", i*425)

		//company.Id = id
		//company.Email = fmt.Sprintf("%d@company%d.com", i, i)
		//company.Password = fmt.Sprintf("Password-%d", i)
		//company.Active = (i%2 == 0)
		//company.Role = "COMPANY"
		db.Add("company", id, company)
	}
	return compIds
}

func MakeDrivers(compIds [COMP]string) {
	for i := 0; i < (COMP * DperC); i++ {
		id := strconv.Itoa(int(time.Now().UnixNano()))
		compIdx := i / DperC
		d := i % 10
		driver := Driver{
			FirstName:             "Daniel",
			LastName:              fmt.Sprintf("Jones the %dth", (i + 4)),
			Phone:                 fmt.Sprintf("717-777-777%d", d),
			EmergencyContactName:  "Samuel Johnson",
			EmergencyContactPhone: fmt.Sprintf("222-222-222%d", d),
			LicenseNum:            fmt.Sprintf("1234567%d", i),
			LicenseState:          fmt.Sprintf("%d state", i),
			LicenseExpire:         fmt.Sprintf("03/1%d/202%d", d, d),
			DOB:                   fmt.Sprintf("01/1%d/198%d", d, d),
			MedCardExpiry:         fmt.Sprintf("02/1%d/202%d", d, d),
			MVRExpiry:             fmt.Sprintf("03/1%d/202%d", d, d),
			ReviewExpiry:          fmt.Sprintf("04/1%d/202%d", d, d),
			OneEightyExpiry:       fmt.Sprintf("05/1%d/202%d", d, d),
			HireDate:              fmt.Sprintf("06/1%d/199%d", d, d),
			TermDate:              fmt.Sprintf("07/1%d/202%d", d, d),
			CompanyId:             compIds[compIdx],
		}

		driver.Id = id
		driver.Email = fmt.Sprintf("%d@%d.com", i, i)
		driver.Password = fmt.Sprintf("Password-%d", i)
		driver.Active = (i%2 == 0)
		driver.Role = "DRIVER"

		driver.Street = fmt.Sprintf("12%d Main Street", 1)
		driver.City = fmt.Sprintf("%dville", i)
		driver.State = fmt.Sprintf("%d state", i)
		driver.Zip = fmt.Sprintf("1234%d", d)

		db.Add("driver", id, driver)
	}
}

func MakeVehicles(compIds [COMP]string) {
	for i := 0; i < (COMP * VperC); i++ {
		id := strconv.Itoa(int(time.Now().UnixNano()))
		compIdx := i / VperC
		vType := "TRUCK"
		if i%3 == 0 {
			vType = "TRACTOR"
		} else if i%2 == 0 {
			vType = "TRAILER"
		}
		vehicle := Vehicle{
			Id:            id,
			CompanyId:     compIds[compIdx],
			VehicleType:   vType,
			UnitNumber:    fmt.Sprintf("%d", i),
			Make:          fmt.Sprintf("make-%d", i),
			VIN:           fmt.Sprintf("%d%d%d", i, i, i),
			Title:         fmt.Sprintf("title-%d", i),
			GVW:           fmt.Sprintf("GVW-%d", i),
			GCR:           fmt.Sprintf("GCR-%d", i),
			UnladenWeight: fmt.Sprintf("unladenWeight-%d", i),
			PurchasePrice: float32(i),
			PurchaseDate:  fmt.Sprintf("199%d-03-1%d", i, i),
			CurrentValue:  float32(i),
			//AxleAmmount:      fmt.Sprintf("axle-%d", i),
			FuelType: fmt.Sprintf("fuel-%d", i),
		}
		if i%3 == 0 {
			vehicle.AxleAmmount = 2
		} else if i%2 == 0 {
			vehicle.AxleAmmount = 3
		} else {
			vehicle.AxleAmmount = 4
		}

		db.Add("vehicle", id, vehicle)
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

/*
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
}}*/

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
	tc.Render(w, r, "form-elements.tmpl", nil)
}}
