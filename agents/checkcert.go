package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"

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
		log.Fatal(err.Error())
	}
	err = json.Unmarshal(raw, &optionMap)
	if err != nil {
		log.Fatal(err.Error())
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
	err = watcher.Add(options["CertLocation"].(string) + "certs/")
	if err != nil {
		log.Fatal(err)
	}
	err = watcher.Add(options["CertLocation"].(string) + "keys/")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func sendRestartCommand(options map[string]interface{}) {
	packet := strings.NewReader(
		`{
	 "login": {},
	  "saviour": {
		  "message": "ReloadCert"
	 }
  }`)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	client.Post("https://localhost.saviour.diyccs.com:"+options["Port"].(string)+"/request/reloadcert", "application/json", packet)
	go exec.Command(options["SaviourLocation"].(string)).Run()
}
