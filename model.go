package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type Auth struct {
	Id       string `json:"id"`
	Email    string `json:"email" auth:"username"`
	Password string `json:"password" auth:"password"`
	Active   bool   `json:"active" auth:"active"`
	Role     string `json:"role"`
}

type Employee struct {
	Auth
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Home      string `json:"home,omitempty"`
	Address
}

type Company struct {
	Id              string  `json:"id"`
	DOTNum          string  `json:"dotNum,omitempty"`
	Name            string  `json:"name,omitempty"`
	ContactName     string  `json:"contactName,omitempty"`
	ContactTitle    string  `json:"contactTitle,omitempty"`
	ContactPhone    string  `json:"contactPhone,omitempty"`
	SecondName      string  `json:"secondName,omitempty"`
	SecondTitle     string  `json:"secondTitle,omitempty"`
	SecondPhone     string  `json:"secondPhone,omitempty"`
	SameAddress     bool    `json:"sameAddress"`
	PhysicalAddress Address `json:"pysicalAddress,omitempty"`
	MailingAddress  Address `json:"mailingAddress,omitempty"`
	MCNum           string  `json:"mcNum,omitempty"`
	PUCNum          string  `json:"pucNum,omitempty"`
	Fax             string  `json:"fax,omitempty"`
	Email           string  `json:"email,omitempty"`
	EINNum          string  `json:"einNum,omitempty"`
	//Slug            string  `json:"slug,omitempty"`
}

/*func (c *Company) CreateSlug() {
	// slug = title.replaceAll("[;/?:@&=+\\\$,\\{\\}\\|\\\\^\\[\\]`]", "").trim().replace(' ', '_').toLowerCase()
	r, err := regexp.Compile("[;/?:@&=+$,\\{\\}\\|^\\[\\]`]")
	if err != nil {
		log.Printf("model.go -> Company.creaeSlug() -> regexp.Compile() -> %v\n", err)
		return
	}
	c.Slug = strings.ToLower(strings.Replace(strings.Trim(r.ReplaceAllString(c.Name, ""), " "), " ", "-", -1))
}*/

type Driver struct {
	Auth
	FirstName    string `json:"firstName,omitempty"`
	LastName     string `json:"lastName,omitempty"`
	Phone        string `json:"phone,omitempty"`
	DOB          string `json:"dob,omitempty"`
	LicenseNum   string `json:"licenseNum,omitempty"`
	LicenseState string `json:"licenseState,omitempty"`
	CompanyId    string `json:"companyId,omitempty"`
	Address
}

func (d Driver) FormatDOB() string {
	ds := strings.Split(d.DOB, "-")
	if len(ds) != 3 {
		return ""
	}
	if ds[1][0] == '0' {
		ds[1] = ds[1][1:]
	}
	return fmt.Sprintf("%s/%s/%s", ds[1], ds[2], ds[0])
}

func (d Driver) GetAge() int32 {
	dobT, err := time.Parse("2006-1-02", d.DOB)
	if err != nil {
		return 0
	}
	dob := dobT.UnixNano()
	diff := time.Now().UnixNano() - dob
	return int32(math.Floor((float64(diff) / float64(1000) / float64(1000) / float64(1000) / float64(60) / float64(60) / float64(24) / float64(365.25))))
}

type Address struct {
	Street string `json:"street,omitempty"`
	City   string `json:"city,omitempty"`
	State  string `json:"state,omitempty"`
	Zip    string `json:"zip,omitempty"`
}

func (a Address) AddrHTML() string {
	return a.Street + "<br>" + a.City + ",  " + a.State + " " + a.Zip
}

type Document struct {
	Id         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	DocumentId string `json:"documentId,omitempty"`
	Complete   bool   `json:"complete"`
	Data       string `json:"data,omitempty"`
	CompanyId  string `json:"companyId,omitempty"`
	DriverId   string `json:"driverId,omitempty"`
}

var DQFS = [][]string{
	[]string{"100", "Driver's Application"},
	[]string{"180", "Certification of Violations"},
	//[]string{"200", "Annual Inquery & Review"},
	//[]string{"250", "Road Test Certication"},
	[]string{"300", "Previous Driver Inquires"},
	[]string{"400", "Drug & Alcohol Records Request"},
	//[]string{"450", "Drug & Alcohol Certified Receipt"},
	[]string{"500", "Certification Compliance"},
	[]string{"600", "Confictions for a Driver Violation"},
	[]string{"700", "New Hire Stmt On Duty Hours"},
	[]string{"750", "Other Ompensated Work"},
	[]string{"775", "Fair Credit Reporting Act"},
}

type Event struct {
	Id        string    `json:"id,omitempty"`
	Name      string    `json:"name"`
	Title     string    `json:"title,omitempty"`
	AllDay    bool      `json:"allDay,omitempty"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end,omitempty"`
	URL       string    `json:"url,omitempty"`
	ClassName string    `json:"className,omitempty"`
	Editable  bool      `json:"editable"`
	Overlap   bool      `json:"overlap,omitempty"`
}

type Note struct {
	Id        string `json:"id"`
	CompanyId string `json:"companyId,omitempty"`
	Body      string `json:"body,omitempty"`
}

type Comment struct {
	Id     string `json:"id"`
	Body   string `json:"body"`
	Url    string `json:"url"`
	Page   string `json:"page"`
	Closed bool   `json:"closed"`
}

type Vehicle struct {
	Id            string  `json:"id"`
	CompanyId     string  `json:"companyId,omitempty"`
	VehicleType   string  `json:"vehicleType,omitempty"`
	UnitNumber    string  `json:"unitNumber,omitempty"`
	Make          string  `json:"make,omitempty"`
	VIN           string  `json:"vin,omitempty"`
	Title         string  `json:"title,omitempty"`
	GVW           string  `json:"gvw,omitempty"`
	GCR           string  `json:"gcr,omitempty"`
	UnladenWeight string  `json:"unladenWeight,omitempty"`
	PurchasePrice float32 `json:"purchacePrice,omitempty"`
	PurchaseDate  string  `json:"purchaseDate,omitempty"`
	CurrentValue  float32 `json:"currentValue,omitempty"`
	AxleType      string  `json:"axleType,omitempty"`
	FuelType      string  `json:"fuelType,omitempty"`
}

func FormatDate(d string) string {
	ds := strings.Split(d, "-")
	if len(ds) != 3 {
		return ""
	}
	if ds[1][0] == '0' {
		ds[1] = ds[1][1:]
	}
	return fmt.Sprintf("%s/%s/%s", ds[1], ds[2], ds[0])
}
