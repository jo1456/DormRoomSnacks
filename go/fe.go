//Jeffrey Oberg jo1456
//Project part 2 CRUD APP with front and back end servers
// Front end


package main

import (
  "github.com/kataras/iris/v12"

  "strings"
  "sort"
  "fmt"
  "flag"
  "net"
  "encoding/json"
)

var (
  connection net.Conn
  encoder *json.Encoder
  decoder *json.Decoder
  app *iris.Application
)

func main() {

    // creates a variable with the passed flag. default value of 8080
    listenPort := flag.String("listen", "8080", "port to listen on")
    backendHostandPort := flag.String("backend", ":8090", "host name and port of backend. (Format: hostName:port)")
    flag.Parse()

    // attempt to connect to the backend
    connection, err := net.Dial("tcp", *backendHostandPort)
    if err != nil {
      fmt.Println("Connection to passed backend host and port failed. Error:")
      fmt.Println(err)
      return
    }
    defer connection.Close()

    // assign the global encoder and decoder to the created connection
    encoder = json.NewEncoder(connection)
    decoder = json.NewDecoder(connection)

    app = iris.New()

    // turn on the app
    app.Listen(":"+*listenPort)
}
