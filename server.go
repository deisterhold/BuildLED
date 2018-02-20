package main

import (
	"fmt"
	"net/http"
)

var servers [4]TFSServer
var credentials [4]TFSCredentials

var indexHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Testing")
}

var editHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	account := r.FormValue("account")
	project := r.FormValue("project")
	build := r.FormValue("build")
	username := r.FormValue("username")
	password := r.FormValue("password")
	domain := r.FormValue("domain")

	switch r.FormValue("index") {
	case "1":
		servers[0] = TFSHostedServer{account, TFSBuildDefinition{"", project, build}}
		credentials[0] = TFSCredentials{username, password, domain}
		fmt.Println("Updated server:", servers[0])
		fmt.Println("Updated credentials:", credentials[0])
		break
	case "2":
		servers[1] = TFSHostedServer{account, TFSBuildDefinition{"", project, build}}
		credentials[1] = TFSCredentials{username, password, domain}
		fmt.Println("Updated server:", servers[1])
		fmt.Println("Updated credentials:", credentials[1])
		break
	case "3":
		servers[2] = TFSHostedServer{account, TFSBuildDefinition{"", project, build}}
		credentials[2] = TFSCredentials{username, password, domain}
		fmt.Println("Updated server:", servers[2])
		fmt.Println("Updated credentials:", credentials[2])
		break
	case "4":
		servers[3] = TFSHostedServer{account, TFSBuildDefinition{"", project, build}}
		credentials[3] = TFSCredentials{username, password, domain}
		fmt.Println("Updated server:", servers[3])
		fmt.Println("Updated credentials:", credentials[3])
		break
	}
}

func startServer(exit <-chan bool) {
	http.Handle("/", indexHandler)
	http.Handle("/edit", editHandler)
	go http.ListenAndServe(":8080", nil)
	<-exit
}
