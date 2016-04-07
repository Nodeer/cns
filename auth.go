package main

import "github.com/cagnosolutions/web"

var ADMIN = web.Auth{
	Roles:    []string{"ADMIN"},
	Redirect: "/login",
	Msg:      "you are not an admin",
}

var EMPLOYEE = web.Auth{
	Roles:    []string{"ADMIN", "EMPLOYEE"},
	Redirect: "/login/employee",
	Msg:      "Please register",
}
