package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func ReceiveFile(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseMultipartForm(32 << 20) // limit your max input length!
	if err != nil {
		return err
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		return err
	}
	defer file.Close()
	iLog.Printf("Get file - %s\n", header.Filename)
	name := strings.Split(header.Filename, ".")
	newFile, err := os.Create("./configs/" + name[0] + "_v1.0." + name[1])
	if err != nil {
		return err
	}
	defer newFile.Close()
	_, err = io.Copy(newFile, file)
	if err != nil {
		return err
	}
	return nil
}

func RequestFile(w http.ResponseWriter, r *http.Request) error {
	s := r.URL.Query().Get("service")
	files := findFiles(s)
	if len(files) == 0 {
		return errors.New("NOT SUCCESS : config is not found")
	}
	if s == "" { //Print all configs
		for i, file := range files {
			fmt.Fprintf(w, "%d. %s\n", i+1, file)
		}
		iLog.Printf("SUCCESS : configs is sent")
	} else { //Print config by service name
		file := chooseNewestFile(files)
		iLog.Printf("SUCCESS : config is found - ", file)
		http.ServeFile(w, r, "./configs/"+file)
	}
	return nil
}

func chooseNewestFile(files []string) string {
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
		if value[1:] == serviceName || serviceName == "" { //value[1:] because first symbol everytime is SPACE
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
