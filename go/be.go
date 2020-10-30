//Jeffrey Oberg jo1456
//Project part 2 CRUD APP with front and back end servers
// back end

package main

import (
        "fmt"
        "net"
        "encoding/json"
        "flag"
        // "errors"
        "database/sql"
        _ "github.com/go-sql-driver/mysql"
)

var (
  connection net.Conn
  encoder *json.Encoder
  decoder *json.Decoder
)

func main() {

        db, err := sql.Open("mysql", "root:<PUT YOUR PASSWORD HERE>@tcp(127.0.0.1:3306)/DormRoomSnacks")
        if err != nil {
            panic(err.Error())
        }

        defer db.Close()

        // creates a variable with the passed flag. default value of 8080
        listenPort := flag.String("listen", "8090", "port to listen on")
        flag.Parse()

        _, err = db.Query("INSERT INTO Foods VALUES ( 'New food', 'wicked good', 10.01, TRUE )")
        if err != nil {
          panic(err.Error())
        }

        // listen on passed port
        listener, err := net.Listen("tcp", ":" + *listenPort)
        if err != nil {
          fmt.Println(err)
          return
        }
        defer listener.Close()

        // accept a connection
        connection, err := listener.Accept()
        if err != nil {
          fmt.Println(err)
          return
        }

        // assign the global encoder and decoder to the created connection
        encoder = json.NewEncoder(connection)
        decoder = json.NewDecoder(connection)

        // infinite loop to accept and response to requests
        for {
          // status of operation. In case of failure fe can redirect to homepage
          status := Status{Success: true}

          var req Request
          decoder.Decode(&req)

          switch req.FunctionName {
          case "ListLocations":
          case "GetMenu":
          case "ViewItem":
          case "SubmitOrder":
          case "CheckOrderStatus":
          case "GetOrders":
          case "SelectOrder":
          case "CompleteOrder":
          case "UpdateItem":
          case "CreateItem":
          case "DeleteItem":

          }
          encoder.Encode(status)
        }
}

type Status struct {
  Success bool
}

type Request struct {
  FunctionName string
}

type Location struct {
  ID int
  Name string
  MenuID int
  Address string
  Hours string
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
  Ingredients string
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
  locations []Location
}

type GetMenuRequest struct {
  MenuID int
}

type GetMenuResponse struct {
  RequestedMenu Menu
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
  Success bool
  Status string
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

func ListLocations ()  () {}
func GetMenu (GetMenuRequest)  () {}
func ViewItem (ViewItemRequest)  () {}
func SubmitOrder (SubmitOrderRequest)  () {}
func CheckOrderStatus (CheckOrderStatusRequest)  () {}

func GetOrders (GetOrdersRequest)  () {}
func SelectOrder (SelectOrderRequest)  () {}
func CompleteOrder (CompelteOrderRequest)  () {}
func UpdateItem (UpdateItemRequest)  () {}
func CreateItem (CreateItemRequest)  () {}
func DeleteItem (DeleteItemRequest)  () {}
