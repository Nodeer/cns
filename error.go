package main

var HTTP_ERROR_ADMIN = `<html> <head> <title>%d</title> <meta charset="utf-8"> <meta name="viewport" content="width=device-width, initial-scale=1.0"> <meta name="description" content=""> <meta name="author" content=""> <link rel="shortcut icon" href="/static/img/favicon_1.ico"> <link href='http://fonts.googleapis.com/css?family=Source+Sans+Pro:100,300,400,600,700,900,400italic' rel='stylesheet'> <link href="/static/css/bootstrap.min.css" rel="stylesheet"> <link href="/static/css/bootstrap-reset.css" rel="stylesheet"> <link href="/static/css/animate.css" rel="stylesheet"> <link href="/static/font-awesome/css/font-awesome.css" rel="stylesheet"/> <link href="/static/ionicon/css/ionicons.min.css" rel="stylesheet"/> <link href="/static/css/style.css" rel="stylesheet"> <link href="/static/css/helper.css" rel="stylesheet"> <link href="/static/css/style-responsive.css" rel="stylesheet"/> <link href="/static/notifications/notification.css" rel="stylesheet"/> <link rel="stylesheet" href="/static/css/custom.css" charset="utf-8"><!--[if lt IE 9]> <script src="/static/js/html5shiv.js"></script> <script src="/static/js/respond.min.js"></script><![endif]--> </head> <body> <aside class="left-panel"> <div class="logo"> <a class="logo-expanded"> <span class="nav-label">CNS Internal</span> </a> </div><nav class="navigation"> <ul class="list-unstyled"> <li> <a href="/cns/company"><i class="fa fa-building-o text-center"></i><span class="nav-label">Customer</span></a> </li><li> <a href="/cns/driver"><i class="ion-man text-center"></i><span class="nav-label">Driver</span></a> </li><li> <a href="/cns/employee"><i class="ion-person text-center"></i><span class="nav-label">Employee</span></a> </li><li> <a href="/logout"><i class="fa fa-sign-out text-center"></i><span class="nav-label">Logout</span></a> </li></ul> </nav> </aside> <section class="content"> <div class="wraper container-fluid"> <div class="row"> <br><br><div class="col-sm-offset-3 col-sm-6 text-center"> <img src="/static/img/cns-logo.png" alt="CNS Truck Licensing Logo"/> <h1>Sorry for the inconvenience</h1> <p>%s</p><p>HTTP Status %d</p></div></div></div></section> </body></html>`

var HTTP_ERROR_EMPLOYEE = `<html> <head> <title>%d</title> <meta charset="utf-8"> <meta name="viewport" content="width=device-width, initial-scale=1.0"> <meta name="description" content=""> <meta name="author" content=""> <link rel="shortcut icon" href="/static/img/favicon_1.ico"> <link href='http://fonts.googleapis.com/css?family=Source+Sans+Pro:100,300,400,600,700,900,400italic' rel='stylesheet'> <link href="/static/css/bootstrap.min.css" rel="stylesheet"> <link href="/static/css/bootstrap-reset.css" rel="stylesheet"> <link href="/static/css/animate.css" rel="stylesheet"> <link href="/static/font-awesome/css/font-awesome.css" rel="stylesheet"/> <link href="/static/ionicon/css/ionicons.min.css" rel="stylesheet"/> <link href="/static/css/style.css" rel="stylesheet"> <link href="/static/css/helper.css" rel="stylesheet"> <link href="/static/css/style-responsive.css" rel="stylesheet"/> <link href="/static/notifications/notification.css" rel="stylesheet"/> <link rel="stylesheet" href="/static/css/custom.css" charset="utf-8"><!--[if lt IE 9]> <script src="/static/js/html5shiv.js"></script> <script src="/static/js/respond.min.js"></script><![endif]--> </head> <body> <aside class="left-panel"> <div class="logo"> <a class="logo-expanded"> <span class="nav-label">CNS Internal</span> </a> </div><nav class="navigation"> <ul class="list-unstyled"> <li> <a href="/cns/company"><i class="fa fa-building-o text-center"></i><span class="nav-label">Customer</span></a> </li><li> <a href="/cns/driver"><i class="ion-man text-center"></i><span class="nav-label">Driver</span></a> </li><li> <a href="/logout"><i class="fa fa-sign-out text-center"></i><span class="nav-label">Logout</span></a> </li></ul> </nav> </aside> <section class="content"> <div class="wraper container-fluid"> <div class="row"> <br><br><div class="col-sm-offset-3 col-sm-6 text-center"> <img src="/static/img/cns-logo.png" alt="CNS Truck Licensing Logo"/> <h1>Sorry for the inconvenience</h1> <p>%s</p><p>HTTP Status %d</p></div></div></div></section> </body></html>`

var HTTP_ERROR_DEFAULT = `<html> <head> <title>%d</title> <meta charset="utf-8"> <meta name="viewport" content="width=device-width, initial-scale=1.0"> <meta name="description" content=""> <meta name="author" content=""> <link rel="shortcut icon" href="/static/img/favicon_1.ico"> <link href='http://fonts.googleapis.com/css?family=Source+Sans+Pro:100,300,400,600,700,900,400italic' rel='stylesheet'> <link href="/static/css/bootstrap.min.css" rel="stylesheet"> <link href="/static/css/bootstrap-reset.css" rel="stylesheet"> <link href="/static/css/animate.css" rel="stylesheet"> <link href="/static/font-awesome/css/font-awesome.css" rel="stylesheet"/> <link href="/static/ionicon/css/ionicons.min.css" rel="stylesheet"/> <link href="/static/css/style.css" rel="stylesheet"> <link href="/static/css/helper.css" rel="stylesheet"> <link href="/static/css/style-responsive.css" rel="stylesheet"/> <link href="/static/notifications/notification.css" rel="stylesheet"/> <link rel="stylesheet" href="/static/css/custom.css" charset="utf-8"><!--[if lt IE 9]> <script src="/static/js/html5shiv.js"></script> <script src="/static/js/respond.min.js"></script><![endif]--> </head> <body> <aside class="left-panel"> <div class="logo"> <a class="logo-expanded"> <span class="nav-label">CNS Internal</span> </a> </div><nav class="navigation"> <ul class="list-unstyled"> <li> <a href="/login"><i class="fa fa-sign-out text-center"></i><span class="nav-label">Login</span></a> </li></ul> </nav> </aside> <section class="content"> <div class="wraper container-fluid"> <div class="row"> <br><br><div class="col-sm-offset-3 col-sm-6 text-center"> <img src="/static/img/cns-logo.png" alt="CNS Truck Licensing Logo"/> <h1>Sorry for the inconvenience</h1> <p>%s</p><p>HTTP Status %d</p></div></div></div></section> </body></html>`
