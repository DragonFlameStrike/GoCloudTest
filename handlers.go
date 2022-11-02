package main

import (
	"fmt"
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
		s := r.URL.Query().Get("service")
		files := findFiles(s)
		if len(files) == 0 {
			iLog.Printf("NOT SUCCESS : config is not found")
			http.Error(w, "404 not found.", http.StatusNotFound)
			return
		}
		if s == "" {
			for i, file := range files {
				fmt.Fprintf(w, "%d. %s\n", i+1, file)
			}
			iLog.Printf("SUCCESS : configs is sent")
		} else {
			file := chooseNewestFile(files)
			iLog.Printf("SUCCESS : config is found - ", file)
			http.ServeFile(w, r, "./configs/"+file)
		}

	case "POST":
		iLog.Println("catch POST request...")
	default:
		iLog.Println(w, "Only GET and POST methods are supported.")
	}
}
