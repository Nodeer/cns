package main

type Auth struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Active   bool   `json:"active"`
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
	Contact string `json:"conatact,omityempty"`
	Phone   string `json:"phone,omitempty"`
	Address
}

type Driver struct {
	Auth
	FirstName    string `json:"firstName,omitempty"`
	LastName     string `json:"lastName,omitempty"`
	Phone        string `json:"phone,omitEmpty"`
	DOB          string `json:"dob,omitEmpty"`
	LicenseNum   string `json:"licenseNum,omitEmpty"`
	LicenseState string `json:"licenseState,omitEmpty"`
	Address
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
