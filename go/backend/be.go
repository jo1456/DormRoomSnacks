//Jeffrey Oberg jo1456
//Project part 2 CRUD APP with front and back end servers
// back end

package main

import (
	"dormroomsnacks/structs"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	// "errors"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	connection net.Conn
	encoder    *json.Encoder
	decoder    *json.Decoder
	DB         *sql.DB
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	// DB_HOSTURL, dbexists1 := os.LookupEnv("DB_HOSTURL")
	// DB_USERNAME, dbexists2 := os.LookupEnv("DB_USERNAME")
	// DB_PASSWORD, dbexists3 := os.LookupEnv("DB_PASSWORD")
	// DB_NAME, dbexists4 := os.LookupEnv("DB_NAME")
	// if !dbexists1 || !dbexists2 || !dbexists3 || !dbexists4 {
	// 	panic(1)
	// }

	// databaseURI := fmt.Sprintf("%s:%s@%s/%s", DB_USERNAME, DB_PASSWORD, DB_HOSTURL, DB_NAME)

	// db, err := sql.Open("postgres", databaseURI)
	// if err != nil {
	// 	panic(err.Error())
	// }

	// defer db.Close()

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
	listener, err := net.Listen("tcp", ":"+*listenPort)
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

		var req structs.Request
		decoder.Decode(&req)

		fmt.Println(req.FunctionName)

		switch req.FunctionName {
		case "ListLocations":
			ListLocations()
		case "GetMenu":
			var getMenuReq structs.GetMenuRequest
			decoder.Decode(&getMenuReq)

			GetMenu(getMenuReq)
		case "ViewItem":
			var viewItemReq structs.ViewItemRequest
			decoder.Decode(&viewItemReq)

			ViewItem(viewItemReq)
		case "SubmitOrder":
			var order structs.Order
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

func ListLocations() {
	rows, err := DB.Query("select * from DiningHalls;")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	resp := structs.ListLocationsResponse{Locations: make([]structs.Location, 0)}

	for rows.Next() {
		var location structs.Location
		err := rows.Scan(&location.ID, &location.Name, &location.Address, &location.Phone, &location.MenuID, &location.Hours)
		if err != nil {
			fmt.Println(err)
			return
		}
		resp.Locations = append(resp.Locations, structs.Location{ID: location.ID, Name: location.Name, Address: location.Address, Phone: location.Phone, MenuID: location.MenuID, Hours: location.Hours})
	}

	fmt.Println(resp)

	encoder.Encode(&resp)

}

func GetMenu(req structs.GetMenuRequest) {
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

	var menu structs.Menu

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
		var item structs.FoodItem
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

func ViewItem(req structs.ViewItemRequest) {
	rows, err := DB.Query("select * from Foods where ID = ?;", req.ItemID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	var item structs.FoodItem
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

func SubmitOrder(req structs.Order) {
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

	for _, item := range req.Items {
		_, err := DB.Query("insert into OrderItem (foodID, orderID, Customization) values(?,?,?);", item.ItemID, id, item.Customization)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("order id:")
	fmt.Println(id)
	encoder.Encode(id)
}

func CheckOrderStatus(req structs.CheckOrderStatusRequest) {
	rows, err := DB.Query("select * from Orders where ID = ?;", req.OrderID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	var item structs.FoodItem
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

func GetOrders(structs.GetOrdersRequest)         {}
func SelectOrder(structs.SelectOrderRequest)     {}
func CompleteOrder(structs.CompelteOrderRequest) {}
func UpdateItem(structs.UpdateItemRequest)       {}
func CreateItem(structs.CreateItemRequest)       {}
func DeleteItem(structs.DeleteItemRequest)       {}
