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
	case "POST", "PUT":
		// POST - create file on server
		// PUT - edit file on server
		if r.Method == "POST" {
			iLog.Println("catch CREATE request...")
		} else {
			iLog.Println("catch UPDATE request...")
		}
		err := ReceiveFile(w, r)
		if err != nil {
			iLog.Print(err)
			http.Error(w, "400 bad request.", http.StatusBadRequest)
		} else {
			iLog.Println("SUCCESS : file is received!")
		}
	case "DELETE":
		iLog.Println("catch DELETE request...")
		err := DeleteFile(w, r)
		if err != nil {
			iLog.Print(err)
			http.Error(w, "400 bad request.", http.StatusBadRequest)
		} else {
			iLog.Println("SUCCESS : file is deleted")
		}
	default:
		iLog.Println(w, "Only GET, POST, PUT and DELETE methods are supported.")
	}
}
