package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
)

var (
	LOGFILE = "/tmp/GoCloudTest.log"
	iLog    *log.Logger
)

func main() {
	//Create logger
	f, err := os.OpenFile(LOGFILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	iLog = log.New(f, "GoCloudTestLog - ", log.LstdFlags)
	iLog.SetFlags(log.LstdFlags)
	iLog.Println("START LOGGER")

	sigCather()

	//Create http server
	http.HandleFunc("/config", configHandler)
	iLog.Print("Server is ready")
	http.ListenAndServe(":8080", nil)
	if err != nil {
		iLog.Print(err)
		return
	}
}

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

func chooseNewestFile(files []string) string {
	if len(files) == 0 {
		return ""
	}
	maxV := "0.0"
	maxVId := -1
	for i, file := range files {
		v := getV(file)
		if greater(v, maxV) {
			maxV = v
			maxVId = i
		}
	}
	return files[maxVId]
}

func greater(v string, maxV string) bool {
	vs := strings.Split(v, ".")
	mvs := strings.Split(maxV, ".")
	for i, e := range vs {
		if e > mvs[i] {
			return true
		}
	}
	return false
}

func getV(file string) string {
	s := strings.Split(file, "_v")
	return s[len(s)-1][:len(s[len(s)-1])-5] //last split element without last 5 symbols (.json)
}

func findFiles(serviceName string) []string {
	files, err := ioutil.ReadDir("./configs")
	if err != nil {
		iLog.Fatal(err)
	}
	var correctFiles []string
	for _, file := range files {
		f, err := os.OpenFile("./configs/"+file.Name(), os.O_RDONLY, 0644)
		if err != nil {
			iLog.Fatal(err)
		}
		byteValue, err := ioutil.ReadAll(f)
		if err != nil {
			iLog.Fatal(err)
		}
		f.Close()
		value := extractValue(string(byteValue[:]), "service")
		if value[1:] == serviceName { //value[1:] because first symbol everytime is SPACE
			correctFiles = append(correctFiles, file.Name())
		}
	}
	return correctFiles
}

// extracts the value for a key from a JSON-formatted string
// body - the JSON-response as a string. Usually retrieved via the request body
// key - the key for which the value should be extracted
// returns - the value for the given key
func extractValue(body string, key string) string {
	keystr := "\"" + key + "\":[^,;\\]}]*"
	r, _ := regexp.Compile(keystr)
	match := r.FindString(body)
	keyValMatch := strings.Split(match, ":")
	return strings.ReplaceAll(keyValMatch[1], "\"", "")
}

func sigCather() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		s := <-sigc
		iLog.Println("STOP LOGGER")
		panic(s)
	}()
}
