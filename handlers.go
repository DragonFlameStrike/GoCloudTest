package main

import (
	"net/http"
)

func configHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/config" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		iLog.Println("404 not found")
		return
	}

	switch r.Method {
	case "GET":
		iLog.Println("catch GET request...")
		s := r.URL.Query().Get("service")
		if s == "" {
			iLog.Println("ERROR : bad request without service field")
			http.Error(w, "404 not found.", http.StatusNotFound)
			return
		}
		files := findFiles(s)
		file := chooseNewestFile(files)
		if file == "" {
			iLog.Printf("NOT SUCCESS : config is not found")
			http.Error(w, "404 not found.", http.StatusNotFound)
			return
		} else {
			iLog.Printf("SUCCESS : config is found - ", file)
			http.ServeFile(w, r, "./configs/"+file)
		}
	case "POST":
		iLog.Println(w, "Post from website!\n")
	default:
		iLog.Println(w, "Only GET and POST methods are supported.")
	}
}
