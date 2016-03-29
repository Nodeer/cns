package main

type User struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Active   bool   `json:"active"`
}
