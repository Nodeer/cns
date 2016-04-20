package main

import "github.com/cagnosolutions/web"

var ADMIN = web.Auth{
	Roles:    []string{"ADMIN"},
	Redirect: "/login",
	Msg:      "Please Login",
}

var EMPLOYEE = web.Auth{
	Roles:    []string{"ADMIN", "EMPLOYEE"},
	Redirect: "/login",
	Msg:      "Please Login",
}

var COMPANY = web.Auth{
	Roles:    []string{"COMPANY"},
	Redirect: "/company/login",
	Msg:      "Please Login",
}

var AJAX = web.Auth{
	Roles:    []string{"ADMIN", "EMPLOYEE", "COMPANY"},
	Redirect: "/login",
	Msg:      "ERROR",
}
