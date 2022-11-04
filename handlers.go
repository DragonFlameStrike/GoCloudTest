package main

import (
	"net/http"
)

func configReadCreate(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/config" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		iLog.Println("Bad URL path")
		return
	}

	switch r.Method {
	case "GET":
		iLog.Println("catch READ request...")
		err := RequestFile(w, r)
		if err != nil {
			iLog.Print(err)
			http.Error(w, "404 not found.", http.StatusNotFound)
		} else {
			iLog.Println("SUCCESS : file(-s) is sent")
		}
	case "POST":
		iLog.Println("catch CREATE request...")
		err := ReceiveFile(w, r)
		if err != nil {
			iLog.Print(err)
			http.Error(w, "400 bad request.", http.StatusBadRequest)
		} else {
			iLog.Println("SUCCESS : file is created!")
		}
	default:
		iLog.Println(w, "Only GET and POST methods are supported.")
	}
}
