package main

type User struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Role          string `json:"role"`
	Active        bool   `json:"active"`
	Name          string `json:"name"`
	Phone         string `json:"phone,omitempty"`
	Contact       string `json:"contact,omitempty"`
	DOB           string `json:"dob,omitempty"`
	LicenseNumber string `json:"licenseNumber, omitempty"`
}
