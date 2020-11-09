//Jeffrey Oberg jo1456
//Project part 2 CRUD APP with front and back end servers
// Front end


package main

import (
  "github.com/kataras/iris/v12"

  // "strings"
  // "sort"
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

    encoder.Encode(Request{FunctionName: "ListLocations"})

    var list ListLocationsResponse
    err = decoder.Decode(&list)

    fmt.Println(list)

    fmt.Println("Requesting Menu for first location")

    encoder.Encode(Request{FunctionName: "GetMenu"})
    encoder.Encode(GetMenuRequest{MenuID: 1})

    var menu Menu
    err = decoder.Decode(&menu)

    fmt.Println(menu)

    encoder.Encode(Request{FunctionName: "ViewItem"})
    encoder.Encode(ViewItemRequest{ItemID: 1})

    var item FoodItem
    err = decoder.Decode(&item)

    fmt.Println(item)

    // type ItemOrder struct {
    //   ID int
    //   ItemID int
    //   Customization string
    //   Notes string
    // }
    //
    // type Order struct {
    //   UserID int
    //   LocationID int
    //   Items []ItemOrder
    // }

    encoder.Encode(Request{FunctionName: "SubmitOrder"})
    a := ItemOrder{ItemID: 1, Customization: "extra cheese", Notes: "extra cheese plz"}
    b := ItemOrder{ItemID: 2, Customization: "no yogurt", Notes: "extra toppings"}
    encoder.Encode(Order{UserID: 1, LocationID: 1, Items: []ItemOrder{a,b}})

    var num int
    err = decoder.Decode(&num)

    fmt.Print("ORDER ID: ")
    fmt.Println(num)

    app = iris.New()

    // turn on the app
    app.Listen(":"+*listenPort)

    // recieve the list of reports from the backend.
    // var list ReportsList
    // decoder.Decode(&list)
}


type Status struct {
  Success bool
}

type Request struct {
  FunctionName string
}

type Location struct {
  ID int `json:"id"`
  Name string `json:"name"`
  Address string `json:"address"`
  Phone string `json:"phone"`
  MenuID int `json:"menuid"`
  Hours string `json:"hours"`
}

type Menu struct {
  ID int
  Name string
  LocationID int
  Items []FoodItem
}

type FoodItem struct {
  ID int
  Name string
  Description string
  Cost int
  IsAvailable bool
  NutritionFacts string
}

type User struct {
  ID int
  DollarBalance int
  MealSwipeBalance int
}

type ItemOrder struct {
  ID int
  ItemID int
  Customization string
  Notes string
}

type Order struct {
  UserID int
  LocationID int
  Items []ItemOrder
}

type ListLocationsResponse struct {
  Locations []Location
}

type GetMenuRequest struct {
  MenuID int
}

type ViewItemRequest struct {
  ItemID int
}

type ViewItemResponse struct {
  Item FoodItem
}

type SubmitOrderRequest struct {
  OrderRequest Order
}

type SubmitOrderResponse struct {
  ID int
}

type CheckOrderStatusRequest struct {
  OrderID int
}

type CheckOrderStatusResponse struct {
  Status string
}

type GetOrdersRequest struct {
  LocationID int
}

type GetOrdersResponse struct {
  Orders []ItemOrder
}

type SelectOrderRequest struct {
  OrderID int
}

type SelectOrderResponse struct {
  Status string
}

type CompelteOrderRequest struct {
  OrderID int
}

type CompelteOrderResponse struct {
  Status string
}

type UpdateItemRequest struct {
  MenuID int
  ItemID int
  NewItem FoodItem
}

type UpdateItemResponse struct {
  Status string
}

type CreateItemRequest struct {
  MenuID int
  NewItem FoodItem
}

type CreateItemResponse struct {
  Status string
}

type DeleteItemRequest struct {
  MenuID int
  ItemID int
}

type DeleteItemResponse struct {
 Status string
}
