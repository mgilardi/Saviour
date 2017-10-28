/*
This is the beginning entry point for the Saviour Server
*/
package main

import ("fmt"
    "config"
    "strconv")

func main() {
  fmt.Println("Saviour::Start...")
  fmt.Println("Saviour::Loading_Configuration...")
  settings := config.GetSettings()
  fmt.Println("Saviour::ModulesLoaded::" + strconv.Itoa(len(*settings)))
}