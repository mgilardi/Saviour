/*
This is the beginning entry point for the Saviour Server
*/
package main

import ("fmt"
        "config"
  //"modules/logger"
)

func main() {
  fmt.Println("Saviour::Start...")
  fmt.Println("Saviour::Loading_Configuration...")
  err, value := config.FindValue("this", "shit")
  if (err != nil) {
    // handle error
    fmt.Println("Saviour::Error::" + err.Error())
  } else {
    fmt.Println("Saviour::Config::Key::" + value)
  }
}