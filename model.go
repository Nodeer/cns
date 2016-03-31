package main

type User struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Active   bool   `json:"active"`
	Role     string `json:"role"`
}

type Employee struct {
	User
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Address
}

type Company struct {
	User
	Name    string `json:"name,omitempty"`
	Contact string `json:"conatact,omityempty"`
	Phone   string `json:"phone,omitempty"`
	Address
}

type Driver struct {
	User
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

func (a Address) String() string {
	return a.Street + "<br>" + a.City + ",  " + a.State + " " + a.Zip
}
