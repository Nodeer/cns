package main

import "github.com/cagnosolutions/web"

var DEVELOPER = web.Auth{
	Roles:    []string{"DEVELOPER"},
	Redirect: "/login",
	Msg:      "Please Login",
}

var ADMIN = web.Auth{
	Roles:    []string{"DEVELOPER", "ADMIN"},
	Redirect: "/login",
	Msg:      "Please Login",
}

var EMPLOYEE = web.Auth{
	Roles:    []string{"DEVELOPER", "ADMIN", "EMPLOYEE"},
	Redirect: "/login",
	Msg:      "Please Login",
}

var COMPANY = web.Auth{
	Roles:    []string{"DEVELOPER", "COMPANY"},
	Redirect: "/company/login",
	Msg:      "Please Login",
}

var AJAX = web.Auth{
	Roles:    []string{"DEVELOPER", "ADMIN", "EMPLOYEE", "COMPANY"},
	Redirect: "/login",
	Msg:      "ERROR",
}
