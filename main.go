package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
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
