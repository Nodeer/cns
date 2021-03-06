package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cagnosolutions/web"
)

const (
	CperE = 1
	COMP  = 10
	DperC = 10
	VperC = 5
)

func testDrivers() {
	var drivers []Driver
	fmt.Println("Getting all drivers...")
	db.All("driver", &drivers)
	fmt.Println("Compiling list of driver ids...")
	var ids []string
	for _, driver := range drivers {
		ids = append(ids, driver.Id)
	}
	fmt.Println("Waiting 2 seconds...")
	time.Sleep(2 * time.Second)
	fmt.Println("Getting all drivers individually by id...")
	i := 0
	for _, id := range ids {
		var driver Driver
		if !db.Get("driver", id, &driver) {
			fmt.Printf("Failed to get driver with id %s\n", id)
			i++
		}
	}
	fmt.Printf("\nFailed to get %d drivers\n\n", i)
}

func testEmployees() {
	var employees []Employee
	fmt.Println("Getting all employees...")
	db.All("employee", &employees)
	fmt.Println("Compiling list of employee ids...")
	var ids []string
	for _, employee := range employees {
		ids = append(ids, employee.Id)
	}
	fmt.Println("Waiting 2 seconds...")
	time.Sleep(2 * time.Second)
	fmt.Println("Getting all employees individually by id...")
	i := 0
	for _, id := range ids {
		var employee Employee
		if !db.Get("employee", id, &employee) {
			fmt.Printf("Failed to get employee with id %s\n", id)
			i++
		}
	}
	fmt.Printf("\nFailed to get %d employees\n\n", i)
}

func testCompanies() {
	var companies []Company
	fmt.Println("Getting all companies...")
	db.All("company", &companies)
	fmt.Println("Compiling list of company ids...")
	var ids []string
	for _, company := range companies {
		ids = append(ids, company.Id)
	}
	fmt.Println("Waiting 2 seconds...")
	time.Sleep(2 * time.Second)
	fmt.Println("Getting all companies individually by id...")
	i := 0
	for _, id := range ids {
		var company Company
		if !db.Get("company", id, &company) {
			fmt.Printf("Failed to get company with id %s\n", id)
			i++
		}
	}
	fmt.Printf("\nFailed to get %d companies\n\n", i)
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
			GVW:           1000 * i,
			GCR:           1155 * i,
			UnladenWeight: 1357 * i,
			PurchasePrice: float32(i),
			PurchaseDate:  fmt.Sprintf("03/1%d/199%d", i, i),
			CurrentValue:  float32(i),
			FuelType:      fmt.Sprintf("fuel-%d", i),
			Active:        i%2 == 0,
			Owner:         fmt.Sprintf("Vinny P number %d", i),
			Year:          fmt.Sprintf("%d", 1980+compIdx),
			PlateNum:      fmt.Sprintf("%d", 658231+compIdx),
			PlateExpire:   fmt.Sprintf("03/1%d/%d", 1992+compIdx, i),
		}
		if i%3 == 0 {
			vehicle.AxleAmount = "2"
		} else if i%2 == 0 {
			vehicle.AxleAmount = "3"
		} else {
			vehicle.AxleAmount = "4"
		}

		db.Add("vehicle", id, vehicle)
	}
}

func converVehicles() {
	var vehicle2s []Vehicle2
	db.All("vehicle", &vehicle2s)
	for _, vehicle2 := range vehicle2s {
		vehicle := Vehicle{
			Id:               vehicle2.Id,
			CompanyId:        vehicle2.CompanyId,
			VehicleType:      vehicle2.VehicleType,
			UnitNumber:       vehicle2.UnitNumber,
			Make:             vehicle2.Make,
			VIN:              vehicle2.VIN,
			Title:            vehicle2.Title,
			GVW:              vehicle2.GVW,
			GCR:              vehicle2.GCR,
			UnladenWeight:    vehicle2.UnladenWeight,
			PurchasePrice:    vehicle2.PurchasePrice,
			PurchaseDate:     vehicle2.PurchaseDate,
			CurrentValue:     vehicle2.CurrentValue,
			AxleAmount:       strconv.Itoa(vehicle2.AxleAmount),
			FuelType:         vehicle2.FuelType,
			Active:           vehicle2.Active,
			Owner:            vehicle2.Owner,
			Year:             vehicle2.Year,
			PlateNum:         vehicle2.PlateNum,
			PlateExpire:      vehicle2.PlateExpire,
			PlateExpireMonth: vehicle2.PlateExpireMonth,
			PlateExpireYear:  vehicle2.PlateExpireYear,
			BodyType:         vehicle2.BodyType,
			BodyTypeOther:    vehicle2.BodyTypeOther,
		}
		db.Set("vehicle", vehicle.Id, vehicle)
	}
	var vehicles []Vehicle
	db.All("vehicle", &vehicles)
	fmt.Printf("old vehicles: %d, converted vehicles %d\n", len(vehicle2s), len(vehicle2s))
}

type Vehicle2 struct {
	Id               string   `json:"id"`
	CompanyId        string   `json:"companyId,omitempty"`
	VehicleType      string   `json:"vehicleType,omitempty"`
	UnitNumber       string   `json:"unitNumber,omitempty"`
	Make             string   `json:"make,omitempty"`
	VIN              string   `json:"vin,omitempty"`
	Title            string   `json:"title,omitempty"`
	GVW              int      `json:"gvw,omitempty"`
	GCR              int      `json:"gcr,omitempty"`
	UnladenWeight    int      `json:"unladenWeight,omitempty"`
	PurchasePrice    float32  `json:"purchasePrice,omitempty"`
	PurchaseDate     string   `json:"purchaseDate,omitempty"`
	CurrentValue     float32  `json:"currentValue,omitempty"`
	AxleAmount       int      `json:"axleAmount,omitempty"`
	FuelType         string   `json:"fuelType,omitempty"`
	Active           bool     `json:"active"`
	Owner            string   `json:"owner,omitempty"`
	Year             string   `json:"year,omitempty"`
	PlateNum         string   `json:"plateNum,omitempty"`
	PlateExpire      string   `json:"plateExpire,omitempty"`
	PlateExpireMonth string   `json:"plateExpireMonth,omitempty"`
	PlateExpireYear  string   `json:"plateExpireYear,omitempty"`
	BodyType         BodyType `json:"bodyType,omitempty"`
	BodyTypeOther    string   `json:"bodyTypeOther,omitempty"`
}

func convertCCExpire() {
	var companies []Company
	db.All("company", &companies)
	for _, company := range companies {
		if company.CreditCard.ExpirationDate == "" {
			continue
		}
		ss := strings.Split(company.CreditCard.ExpirationDate, "/")
		if len(ss) != 3 {
			continue
		}
		if ss[1] == "" {
			continue
		}
		m, err := strconv.Atoi(ss[1])
		if err != nil {
			continue
		}
		if ss[2] == "" {
			continue
		}
		y, err := strconv.Atoi(ss[2])
		if err != nil {
			continue
		}
		if m < 1 || m > 12 || y < 0 {
			continue
		}
		company.CreditCard.ExpirationMonth = m
		company.CreditCard.ExpirationYear = y
		db.Set("company", company.Id, company)
	}
	var companies2 []Company
	db.All("company", &companies2)
	fmt.Printf("Old companies: %d, modified companies: %d\n", len(companies), len(companies2))
}

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

var httpError = web.Route{"GET", "/http/error", func(w http.ResponseWriter, r *http.Request) {
	tc.Render(w, r, "error.tmpl", nil)
}}
