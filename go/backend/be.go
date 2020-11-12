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

	//local
	db, err := sql.Open("mysql", "root:Rusty123@tcp(127.0.0.1:3306)/DormRoomSnacks") 

	//shared
	// db, err := sql.Open("mysql", "b2766d1c91f7c7:0c0f617f@tcp(us-cdbr-east-02.cleardb.com:3306)/DormRoomSnacks")

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

// ORDER FLOW
// 1. Create an order when 1 item is added to the cart and set the order status to "cart"
// 2. Add items to order
// 3. Submit order which updates to status to "submitted"


func CreateOrder(req structs.Order) {
	_, err := DB.Query("insert into Orders (personID, diningHallID, status, submitTime, lastStatusChange) values(?,?,?,?,?);",
		req.UserID, req.LocationID, "In Queue", 0, time.Now())
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

	_, err := DB.Query("insert into OrderItem (foodID, orderID, Customization) values(?,?,?);", req.Item.ItemID, id, req.Item.Customization)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("order id:")
	fmt.Println(id)
	encoder.Encode(id)
}

func SubmitOrder(req structs.UpdateOrderRequest) {
	_, err := DB.Query("update Orders set status = \"submitted\", submitTime = ? where id = ?;",
		time.Now(), req.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	encoder.Encode("submitted")
}

func AddItemToOrder(req AddItemToOrderRequest) {
	_, err := DB.Query("insert into OrderItem (foodID, orderID, Customization) values(?,?,?);", req.Item.ItemID, req.OrderID, req.Item.Customization)
	if err != nil {
		fmt.Println(err)
		encoder.Encode("failure")
		return
	}
	encoder.Encode("success")
}

func CheckOrderStatus(req structs.CheckOrderStatusRequest) {
	rows, err := DB.Query("select * from Orders where ID = ?;", req.OrderID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	var order structs.Order
	for rows.Next() {
		var order structs.Order
		var lastStatusChange string
		err := rows.Scan(&order.ID, &order.UserID, &order.LocationID, &order.Status, &order.SubmitTime, &lastStatusChange)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	encoder.Encode(&order)
}

func GetOrders(structs.GetOrdersRequest)         {
	rows, err := DB.Query("select * from Orders where diningHallID = ?;", req.LocationID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	var orders []structs.Order
	for rows.Next() {
		var order structs.Order
		var lastStatusChange string
		err := rows.Scan(&order.ID, &order.UserID, &order.LocationID, &order.Status, &order.SubmitTime, &lastStatusChange)
		if err != nil {
			fmt.Println(err)
			return
		}
		orders = append(orders, order)
	}

	encoder.Encode(orders)
}

func SelectOrder(structs.SelectOrderRequest)     {
	_, err := DB.Query("update Orders set status = \"selected\", lastStatusChange = ? where id = ?;",
		time.Now(), req.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	var orderAndItems OrderAndItems

	rows, err := DB.Query("select * from Orders where ID = ?;", req.OrderID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	var order structs.Order
	for rows.Next() {
		var order structs.Order
		var lastStatusChange string
		err := rows.Scan(&order.ID, &order.UserID, &order.LocationID, &order.Status, &order.SubmitTime, &lastStatusChange)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	orderAndItems.Order = order


	rows, err := DB.Query("select * from OrderItem where orderID = ?;", req.OrderID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	var items []structs.ItemOrder
	for rows.Next() {
		var item structs.ItemOrder
		var lastStatusChange string
		err := rows.Scan(&item.ID, &item.ItemID, &num, &item.Customization)
		if err != nil {
			fmt.Println(err)
			return
		}
		items = append(items, item)
	}

	orderAndItems.Items = items

	encoder.Encode(&orderAndItems)
	
}

func CompleteOrder(req structs.CompelteOrderRequest) {
	_, err := DB.Query("update Orders set status = \"complete\", submitTime = ? where id = ?;",
		time.Now(), req.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	encoder.Encode("complete")
}


func CreateItem(req structs.CreateItemRequest)       {
	_, err := DB.Query("insert into Foods (menuID, name, description, price, availability, nutritionFacts) values (?,?,?,?,?,?);",
		req.MenuID, req.NewItem.Name, req.NewItem.Description, req.NewItem.Cost, req.NewItem.IsAvailable, req.NewItem.NutritionFacts)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	encoder.Encode("created")
}


func UpdateItem(req structs.UpdateItemRequest)       {
	_, err := DB.Query("update Foods set name = ?, description = ?, price = ?, availability = ?, nutritionFacts = ? where id = ?;",
		req.NewItem.Name, req.NewItem.Description, req.NewItem.Cost, req.NewItem.IsAvailable, req.NewItem.NutritionFacts, req.ItemID)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	encoder.Encode("updated")
}

func DeleteItem(req structs.DeleteItemRequest)       {
	_, err := DB.Query("delete from Foods where id = ? AND menuID = ?;",
		req.Items, req.MenuID)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	encoder.Encode("deleted")
}
