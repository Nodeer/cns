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
