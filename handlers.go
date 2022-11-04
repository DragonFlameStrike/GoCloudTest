package main

import (
	"net/http"
)

func configReadCreate(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/config" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		iLog.Println("404 not found")
		return
	}

	switch r.Method {
	case "GET":
		iLog.Println("catch GET request...")
		RequestFile(w, r)
	case "POST":
		iLog.Println("catch POST request...")
		ReceiveFile(w, r)
	default:
		iLog.Println(w, "Only GET and POST methods are supported.")
	}
}
