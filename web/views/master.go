package views

import "html/template"

var MasterTemplate = template.Must(template.New("master").Parse(`<!DOCTYPE html>
<html>

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, shrink-to-fit=no">
    <title>{{template "title" .}} - XDumpGO</title>
    <link rel="stylesheet" href="assets/bootstrap/css/bootstrap.min.css">
    <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Nunito:200,200i,300,300i,400,400i,600,600i,700,700i,800,800i,900,900i">
    <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Montserrat">
    <link rel="stylesheet" href="assets/fonts/fontawesome-all.min.css">
    <link rel="stylesheet" href="assets/fonts/font-awesome.min.css">
    <link rel="stylesheet" href="assets/fonts/fontawesome5-overrides.min.css">
    <link rel="stylesheet" href="assets/css/styles.min.css">
</head>

<body id="page-top">
<div id="wrapper">
    <nav class="navbar navbar-dark align-items-start sidebar sidebar-dark accordion bg-gradient-primary p-0" style="background: #1e1e26!important;padding: 0px!important;">
        <div class="container-fluid d-flex flex-column p-0">
            <a class="navbar-brand d-flex justify-content-center align-items-center sidebar-brand m-0" href="#">
                <div class="sidebar-brand-text mx-3"><span id="xdg">XDG v{{.Version}}</span></div>
            </a>
            <hr class="sidebar-divider my-0">
            <ul class="nav navbar-nav text-light" id="accordionSidebar">
                <li class="nav-item" role="presentation"><a class="nav-link" href="/"><span><i class="fas fa-tachometer-alt" id="icon"></i>Dashboard</span></a></li>
                <li class="nav-item" role="presentation"><a class="nav-link" href="/lists"><span><i class="fas fa-th-list" id="icon"></i>Lists</span></a></li>
                <li class="nav-item" role="presentation"><a class="nav-link" href="#" style="height: 61px;"><i class="far fa-chart-bar" id="icon"></i><span>Data</span></a></li>
                <li class="nav-item" role="presentation"><a class="nav-link" href="/generator" style="height: 61px;"><i class="far fa-chart-bar" id="icon"></i><span>Generator</span></a></li>
                <li class="nav-item" role="presentation"><a class="nav-link" href="/stats"><i class="fas fa-chart-line" id="icon"></i><span>Stats</span></a></li>
                <li class="nav-item" role="presentation"><a class="nav-link" href="/settings"><i class="fa fa-gear" id="icon"></i><span>Settings</span></a></li>
                <li class="nav-item" role="presentation"><a class="nav-link" href="/profile"><i class="far fa-user" id="icon"></i><span>Profile</span></a></li>
                <li class="nav-item" role="presentation"><a class="nav-link" href="/bugs" style="height: 61px;"><i class="fas fa-bug" id="icon"></i><span>Report Bugs</span></a></li>
            </ul>
        </div>
    </nav>
    <div class="d-flex flex-column" id="content-wrapper">
        <div id="content">
            <div class="container-fluid" style="padding-top: 1.5rem;">
                <div class="d-sm-flex justify-content-between align-items-center mb-4">
                    <h3 class="text-dark mb-0" id="PageName" style="padding-left: 40px;">{{template "title" .}}</h3>
                </div>
                <div class="row row-cols-1 align-items-start">
                    <div class="col">
                        <div class="card" id="tip">
                            <div class="card-body">
                                <h4 class="card-title" id="updt">Update!</h4>
                                <p class="card-text" id="bottomupdt">SQLSniper is the shittiest tool on the market, if you use it you're a retarded pedophile.</p>
                            </div>
                        </div>
                    </div>
                </div>
               {{template "content" .}}
            </div>
        </div>
    </div>
</div>
<script src="assets/js/jquery.min.js"></script>
<script src="assets/bootstrap/js/bootstrap.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery-easing/1.4.1/jquery.easing.js"></script>
<script src="assets/js/script.min.js"></script>
</body>

</html>
`))