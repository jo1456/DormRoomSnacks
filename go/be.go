//Jeffrey Oberg jo1456
//Project part 2 CRUD APP with front and back end servers
// back end

package main

import (
        "time"
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
  encoder    *json.Encoder
  decoder    *json.Decoder
  DB         *sql.DB
)

func main() {

        db, err := sql.Open("mysql", "root:<>@tcp(127.0.0.1:3306)/DormRoomSnacks")
        // db, err := sql.Open("mysql", "root:<PUT YOUR PASSWORD HERE>@tcp(127.0.0.1:3306)/DormRoomSnacks")
        if err != nil {
            panic(err.Error())
        }

        defer db.Close()

        DB = db

        // creates a variable with the passed flag. default value of 8080
        listenPort := flag.String("listen", "8090", "port to listen on")
        flag.Parse()

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

        fmt.Println("got a connection")

        // infinite loop to accept and response to requests
        for {

          fmt.Println("in loop")
          // status of operation. In case of failure fe can redirect to homepage
          // status := Status{Success: true}

          var req Request
          decoder.Decode(&req)

          fmt.Println(req.FunctionName)

          switch req.FunctionName {
          case "ListLocations":
            ListLocations()
          case "GetMenu":
            var getMenuReq GetMenuRequest
            decoder.Decode(&getMenuReq)

            GetMenu(getMenuReq)
          case "ViewItem":
            var viewItemReq ViewItemRequest
            decoder.Decode(&viewItemReq)

            ViewItem(viewItemReq)
          case "SubmitOrder":
            var order Order
            decoder.Decode(&order)

            SubmitOrder(order)
          case "CheckOrderStatus":
          case "GetOrders":
          case "SelectOrder":
          case "CompleteOrder":
          case "UpdateItem":
          case "CreateItem":
          case "DeleteItem":

          }
        }
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

func ListLocations ()  () {
    rows, err := DB.Query("select * from DiningHalls;")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer rows.Close()

    resp := ListLocationsResponse{Locations: make([]Location, 0)}

    for rows.Next() {
        var location Location
        err := rows.Scan(&location.ID, &location.Name, &location.Address, &location.Phone, &location.MenuID, &location.Hours)
        if err != nil {
            fmt.Println(err)
            return
        }
         resp.Locations = append(resp.Locations, Location{ID:location.ID, Name:location.Name, Address:location.Address, Phone:location.Phone, MenuID:location.MenuID, Hours:location.Hours})
    }

    fmt.Println(resp)

    encoder.Encode(&resp)

}

func GetMenu (req GetMenuRequest)  () {
  //type Menu struct {
  //   ID int
  //   Name string
  //   LocationID int
  //   Items []FoodItem
  // }

  rows, err := DB.Query("select * from Menu where id = ?;", req.MenuID)
  if err != nil {
      fmt.Println(err)
      return
  }
  defer rows.Close()

  var menu Menu

  for rows.Next() {
      err := rows.Scan(&menu.ID, &menu.Name, &menu.LocationID)
      if err != nil {
          fmt.Println(err)
          return
      }
  }

  rows, err = DB.Query("select * from Foods where menuID = ?;", req.MenuID)
  if err != nil {
      fmt.Println(err)
      return
  }
  defer rows.Close()

  // type FoodItem struct {
  //   ID int
  //   Name string
  //   Description string
  //   Cost int
  //   IsAvailable bool
  //   Ingredients string
  //   NutritionFacts string
  // }

  for rows.Next() {
      var item FoodItem
      var num int
      err := rows.Scan(&item.ID, &num, &item.Name, &item.Description, &item.Cost, &item.IsAvailable, &item.NutritionFacts)
      if err != nil {
          fmt.Println(err)
          return
      }
       menu.Items = append(menu.Items, item)
  }

  fmt.Println(menu)

  encoder.Encode(&menu)
}


func ViewItem (req ViewItemRequest)  () {
  rows, err := DB.Query("select * from Foods where ID = ?;", req.ItemID)
  if err != nil {
      fmt.Println(err)
      return
  }
  defer rows.Close()

  var item FoodItem
  for rows.Next() {
      var num int
      err := rows.Scan(&item.ID, &num, &item.Name, &item.Description, &item.Cost, &item.IsAvailable, &item.NutritionFacts)
      if err != nil {
          fmt.Println(err)
          return
      }
  }

  encoder.Encode(&item)

}

func SubmitOrder (req Order)  () {
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

  _, err := DB.Query("insert into Orders (personID, diningHallID, status, submitTime, lastStatusChange) values(?,?,?,?,?);",
                        req.UserID, req.LocationID, "In Queue", time.Now(), time.Now())
  if err != nil {
      fmt.Println(err)
      return
  }
  rows, err := DB.Query("SELECT max(id) from Orders;")
  if err != nil {
      fmt.Println(err)
      return
  }
  defer rows.Close()

  var id int
  for rows.Next() {
      err := rows.Scan(&id)
      if err != nil {
          fmt.Println(err)
          return
      }
  }

  for _, item := range(req.Items) {
    _, err := DB.Query("insert into OrderItem (foodID, orderID, Customization) values(?,?,?);",item.ItemID, id, item.Customization)
    if err != nil {
        fmt.Println(err)
        return
    }
  }

  fmt.Println("order id:")
  fmt.Println(id)
  encoder.Encode(id)
}

func CheckOrderStatus (req CheckOrderStatusRequest)  () {
  rows, err := DB.Query("select * from Orders where ID = ?;", req.OrderID)
  if err != nil {
      fmt.Println(err)
      return
  }
  defer rows.Close()

  var item FoodItem
  for rows.Next() {
      var num int
      err := rows.Scan(&item.ID, &num, &item.Name, &item.Description, &item.Cost, &item.IsAvailable, &item.NutritionFacts)
      if err != nil {
          fmt.Println(err)
          return
      }
  }

  encoder.Encode(&item)
}

func GetOrders (GetOrdersRequest)  () {}
func SelectOrder (SelectOrderRequest)  () {}
func CompleteOrder (CompelteOrderRequest)  () {}
func UpdateItem (UpdateItemRequest)  () {}
func CreateItem (CreateItemRequest)  () {}
func DeleteItem (DeleteItemRequest)  () {}
