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
	Auth
	Name    string `json:"name,omitempty"`
	Contact string `json:"contact,omityempty"`
	Phone   string `json:"phone,omitempty"`
	Address
}

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
