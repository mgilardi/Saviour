package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/fsnotify/fsnotify"
)

func main() {
	fmt.Println("Checking For SSL/TLS Cert Update")
	startMonitoring()
}

func loadCoreOptions() map[string]interface{} {
	optionMap := make(map[string]interface{})
	raw, err := ioutil.ReadFile("/etc/saviour/config/core.json")
	if err != nil {
		log.Fatal("ReadFile::" + err.Error())
	}
	err = json.Unmarshal(raw, &optionMap)
	if err != nil {
		log.Fatal("JSON::" + err.Error())
	}
	return optionMap
}

func startMonitoring() {
	options := loadCoreOptions()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("Event:", event)
				if event.Op == fsnotify.Create {
					log.Println("CreateEventDetected")
					sendRestartCommand(options)
				}
			case err := <-watcher.Errors:
				log.Println("Error:", err)
			}
		}
	}()
	err = watcher.Add(options["CertLocation"].(string))
	log.Println("CheckingCertLocation::" + options["CertLocation"].(string))
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func sendRestartCommand(options map[string]interface{}) {
	servAddr := "localhost:" + options["PortIPC"].(string)
	conn, err := net.Dial("tcp", servAddr)
	defer conn.Close()
	if err != nil {
		log.Fatal("Saviour::Agent::ConnectionFailed::" + err.Error())
	}
	conn.Write([]byte("UpdateCert"))
	log.Println("Saviour::Agent::UpdateCertSignalSent")
}
