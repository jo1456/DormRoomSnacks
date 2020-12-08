package main

import (
	"dormroomsnacks/structs"
	"encoding/json"
	"flag"
	"fmt"

	// "log"
	"net"
	"time"

	// "errors"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	// "github.com/joho/godotenv"
)

var (
	connection net.Conn
	encoder    *json.Encoder
	decoder    *json.Decoder
	DB         *sql.DB
)

// func init() {
// 	// loads values from .env into the system
// 	if err := godotenv.Load(); err != nil {
// 		log.Print("No .env file found")
// 	}
// }

func main() {

	// DB_HOSTURL, dbexists1 := os.LookupEnv("DB_HOSTURL")
	// DB_USERNAME, dbexists2 := os.LookupEnv("DB_USERNAME")
	// DB_PASSWORD, dbexists3 := os.LookupEnv("DB_PASSWORD")
	// DB_NAME, dbexists4 := os.LookupEnv("DB_NAME")
	// if !dbexists1 || !dbexists2 || !dbexists3 || !dbexists4 {
	// 	panic(1)
	// }

	// databaseURI := fmt.Sprintf("%s:%s@%s/%s", DB_USERNAME, DB_PASSWORD, DB_HOSTURL, DB_NAME)

	db, err := sql.Open("mysql", "b2766d1c91f7c7:0c0f617f@tcp(us-cdbr-east-02.cleardb.com)/heroku_5873df879639de6")
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

		// status of operation. In case of failure fe can redirect to homepage
		// status := Status{Success: true}

		var req structs.Request
		decoder.Decode(&req)

		// fmt.Println(req.FunctionName)

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
		case "CreateOrder":
			var req structs.CreateOrderRequest
			decoder.Decode(&req)

			CreateOrder(req)
		case "SubmitOrder":
			var order structs.UpdateOrderRequest
			decoder.Decode(&order)

			SubmitOrder(order)
		case "AddItemToOrder":
			var req structs.AddItemToOrderRequest
			decoder.Decode(&req)

			AddItemToOrder(req)
		case "GetOrderHistory":
			var req structs.GetOrderHistoryRequest
			decoder.Decode(&req)

			GetOrderHistory(req)
		case "GetOrders":
			var req structs.GetOrdersRequest
			decoder.Decode(&req)

			GetOrders(req)
		case "SelectOrder":
			var req structs.SelectOrderRequest
			decoder.Decode(&req)

			SelectOrder(req)
		case "CompleteOrder":
			var req structs.CompelteOrderRequest
			decoder.Decode(&req)

			CompleteOrder(req)
		case "UpdateItem":
			var req structs.UpdateItemRequest
			decoder.Decode(&req)

			UpdateItem(req)
		case "CreateItem":
			var req structs.CreateItemRequest
			decoder.Decode(&req)

			CreateItem(req)
		case "DeleteItem":
			var req structs.DeleteItemRequest
			decoder.Decode(&req)

			DeleteItem(req)
		case "SendMealSwipes":
			var req structs.SendMealSwipesRequest
			decoder.Decode(&req)

			SendMealSwipes(req)
		case "GetPaymentBalances": // dollar amounts are in cents to avoid floating point
			var req structs.GetPaymentBalancesRequest
			decoder.Decode(&req)

			GetPaymentBalances(req)
		case "Login":
			var req structs.LoginRequest
			decoder.Decode(&req)

			Login(req)
		case "DeleteItemFromOrder":
			var req structs.DeleteItemFromOrderRequest
			decoder.Decode(&req)

			DeleteItemFromOrder(req)
		case "GetCurrentUserCart":
			var req structs.GetCartRequest
			decoder.Decode(&req)

			GetCurrentUserCart(req)
		}
	}
}

// x
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

	encoder.Encode(&resp)
	fmt.Println(resp)

}

// x
func GetMenu(req structs.GetMenuRequest) {

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

// x
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

// cart -> submitted ->

func CreateOrder(req structs.CreateOrderRequest) {

	rows, err := DB.Query("SELECT personID from Orders where status = 'Cart' and personID = ?;", req.OrderRequest.UserID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	var id int
	for rows.Next() {
		rows.Scan(&id)
		fmt.Println(err)
		return
	}

	_, err = DB.Query("insert into Orders (personID, diningHallID, status, submitTime, lastStatusChange, swipeCost, centCost) values(?,?,?,?,?,?,?);",
		req.OrderRequest.UserID, req.OrderRequest.LocationID, "Cart", time.Now(), time.Now(), 0, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// x
func SubmitOrder(req structs.UpdateOrderRequest) {

	rows, err := DB.Query("select swipeCost, centCost, personID from Orders where ID = ?;", req.ID)
	if err != nil {
		fmt.Println(err)
		encoder.Encode("failure")
		return
	}
	defer rows.Close()

	var swipeCost int
	var centCost int
	var userID int
	for rows.Next() {
		err := rows.Scan(&swipeCost, &centCost, &userID)
		if err != nil {
			fmt.Println("here1")
			fmt.Println(err)
			encoder.Encode("failure")
			return
		}
	}

	rows, err = DB.Query("select dollarBalance, mealSwipeBalance from Persons where ID = ?;", userID)
	if err != nil {
		fmt.Println(err)
		encoder.Encode("failure")
		return
	}
	defer rows.Close()

	var dollarBalance int
	var mealSwipeBalance int
	for rows.Next() {
		err := rows.Scan(&swipeCost, &centCost)
		if err != nil {
			fmt.Println(err)
			fmt.Println("here2")
			encoder.Encode("failure")
			return
		}
	}

	if dollarBalance > centCost && mealSwipeBalance > swipeCost {

		_, err := DB.Query("update Personds set dollarBalance = dollarBalance - ?, mealSwipeBalance = mealSwipeBalance - ? where id = ?;",
			centCost, swipeCost, userID)
		if err != nil {
			fmt.Println(err)
			fmt.Println("here3")
			encoder.Encode("failure")
			return
		}

		_, err = DB.Query("update Orders set status = \"submitted\", submitTime = ? where id = ?;",
			time.Now(), req.ID)
		if err != nil {
			fmt.Println(err)
			encoder.Encode("failure")
			return
		}

		encoder.Encode("submitted")
	}
	encoder.Encode("failure")
}

// When Adding an item to an your must select if you will pay with the item with
// meal swipes or with dollars. This will build up 2 accumulators in the row in the
// orders table. These accumulators will then be used to see if the user has sufficent balance
// to pay for their order.
func AddItemToOrder(req structs.AddItemToOrderRequest) {

	rows, err := DB.Query("select id from orders where personID = ? and status = 'Cart';", req.PersonID)
	if err != nil {
		fmt.Println(err)
		encoder.Encode("failure")
		return
	}
	defer rows.Close()

	var orderID int
	for rows.Next() {
		err := rows.Scan(&orderID)
		if err != nil {
			fmt.Println(err)
			encoder.Encode("failure")
			return
		}
	}

	_, err = DB.Query("insert into OrderItem (foodID, orderID, Customization, payWithSwipe) values(?,?,?,?);", req.Item.FoodID, orderID, req.Item.Customization, req.Item.PayWithSwipe)
	if err != nil {
		fmt.Println(err)
		encoder.Encode("failure")
		return
	}

	if req.Item.PayWithSwipe {
		_, err := DB.Query("update Orders set swipeCost = swipeCost + 1 where id = ?;", orderID)
		if err != nil {
			fmt.Println(err)
			encoder.Encode("failure")
			return
		}
	} else {
		rows, err := DB.Query("select price from Foods where ID = ?;", req.Item.FoodID)
		if err != nil {
			fmt.Println(err)
			encoder.Encode("failure")
			return
		}
		defer rows.Close()

		var price int
		for rows.Next() {
			err := rows.Scan(&price)
			if err != nil {
				fmt.Println(err)
				encoder.Encode("failure")
				return
			}
		}

		_, err = DB.Query("update Orders set centCost = centCost + ? where id = ?;", price, orderID)
		if err != nil {
			fmt.Println(err)
			encoder.Encode("failure")
			return
		}
	}

	encoder.Encode("success")
}

func GetCurrentUserCart(req structs.GetCartRequest) {
	var orderAndItemsWithFood structs.OrderAndItemsWithFood

	rows, err := DB.Query("select * from Orders where personID = ? and status = 'Cart';", req.UserID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	var order structs.Order
	for rows.Next() {
		var order structs.Order
		var lastStatusChange string
		err := rows.Scan(&order.ID, &order.UserID, &order.LocationID, &order.Status, &order.SubmitTime, &lastStatusChange, &order.SwipeCost, &order.CentCost)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	orderAndItemsWithFood.Order = order

	rows, err = DB.Query("select * from OrderItem where orderID = ?;", order.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	var items []structs.OrderItemWithFood
	for rows.Next() {
		var item structs.OrderItem
		var food structs.FoodItem
		err := rows.Scan(&item.ID, &item.FoodID, &item.OrderID, &item.Customization, &item.PayWithSwipe)
		if err != nil {
			fmt.Println(err)
			return
		}

		foodRows, err := DB.Query("select * from Foods where id = ?;", item.FoodID)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer foodRows.Close()
		err = rows.Scan(&food.ID, &food.Name, &food.Description, &food.Cost, &food.IsAvailable, &food.NutritionFacts)
		if err != nil {
			fmt.Println(err)
			return
		}

		items = append(items, structs.OrderItemWithFood{Item: item, Food: food})
	}

	orderAndItemsWithFood.Items = items

	encoder.Encode(&orderAndItemsWithFood)

}

func DeleteItemFromOrder(req structs.DeleteItemFromOrderRequest) {
	rows, err := DB.Query("select * from OrderItem where id = ?;", req.ItemID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	var item structs.OrderItem
	for rows.Next() {
		var lastStatusChange string
		err := rows.Scan(&item.ID, &item.FoodID, &lastStatusChange, &item.Customization, &item.PayWithSwipe)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	_, err = DB.Query("delete from OrderItem where id = ?;", req.ItemID)
	if err != nil {
		fmt.Println(err)
		return
	}

	if item.PayWithSwipe {
		_, err := DB.Query("update Orders set swipeCost = swipeCost - 1 where id = ?;", req.OrderID)
		if err != nil {
			fmt.Println(err)
			encoder.Encode("failure")
			return
		}
	} else {
		rows, err := DB.Query("select price from Foods where ID = ?;", item.FoodID)
		if err != nil {
			fmt.Println(err)
			encoder.Encode("failure")
			return
		}
		defer rows.Close()

		var price int
		for rows.Next() {
			err := rows.Scan(&price)
			if err != nil {
				fmt.Println(err)
				encoder.Encode("failure")
				return
			}
		}

		_, err = DB.Query("update Orders set centCost = centCost - ? where id = ?;", price, req.OrderID)
		if err != nil {
			fmt.Println(err)
			encoder.Encode("failure")
			return
		}
	}

	encoder.Encode("success")
	return
}

// x
func GetOrderHistory(req structs.GetOrderHistoryRequest) {
	rows, err := DB.Query("select * from Orders where personID = ?;", req.UserID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	var orders []structs.Order
	for rows.Next() {
		var order structs.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.LocationID, &order.Status, &order.SubmitTime, &order.LastStatusChange, &order.SwipeCost, &order.CentCost)
		if err != nil {
			fmt.Println("here4")
			fmt.Println(err)
			return
		}
		orders = append(orders, order)
	}

	encoder.Encode(&orders)
}

// x
func GetOrders(req structs.GetOrdersRequest) {
	rows, err := DB.Query("select * from Orders where diningHallID = ? and status = \"submitted\";", req.LocationID)
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

	encoder.Encode(&orders)
}

// x
func SelectOrder(req structs.SelectOrderRequest) {
	_, err := DB.Query("update Orders set status = \"selected\", lastStatusChange = ? where id = ?;",
		time.Now(), req.OrderID)
	if err != nil {
		fmt.Println(err)
		return
	}

	var orderAndItems structs.OrderAndItemsWithFood

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

	rows, err = DB.Query("select * from OrderItem, Foods where OrderItem.orderID = ? and Foods.id = OrderItem.foodID;", req.OrderID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	var items []structs.OrderItemWithFood
	for rows.Next() {
		var item structs.OrderItem
		var food structs.FoodItem
		var orderID int
		err := rows.Scan(&item.ID, &item.FoodID, &orderID, &item.Customization, &item.PayWithSwipe, &food.ID, &food.Name, &food.Description, &food.Cost, &food.IsAvailable, &food.NutritionFacts)
		if err != nil {
			fmt.Println(err)
			return
		}
		items = append(items, structs.OrderItemWithFood{Item: item, Food: food})
	}

	orderAndItems.Items = items

	encoder.Encode(&orderAndItems)

}

// x
func CompleteOrder(req structs.CompelteOrderRequest) {
	_, err := DB.Query("update Orders set status = \"complete\", submitTime = ? where id = ?;",
		time.Now(), req.OrderID)
	if err != nil {
		fmt.Println(err)
		return
	}

	encoder.Encode("complete")
}

// x
func CreateItem(req structs.CreateItemRequest) {
	_, err := DB.Query("insert into Foods (menuID, name, description, price, availability, nutritionFacts) values (?,?,?,?,?,?);",
		req.MenuID, req.NewItem.Name, req.NewItem.Description, req.NewItem.Cost, req.NewItem.IsAvailable, req.NewItem.NutritionFacts)
	if err != nil {
		fmt.Println(err)
		return
	}

	encoder.Encode("created")
}

// x
func UpdateItem(req structs.UpdateItemRequest) {
	_, err := DB.Query("update Foods set name = ?, description = ?, price = ?, availability = ?, nutritionFacts = ? where id = ?;",
		req.NewItem.Name, req.NewItem.Description, req.NewItem.Cost, req.NewItem.IsAvailable, req.NewItem.NutritionFacts, req.ItemID)
	if err != nil {
		fmt.Println(err)
		return
	}

	encoder.Encode("updated")
}

// x
func DeleteItem(req structs.DeleteItemRequest) {
	_, err := DB.Query("delete from Foods where id = ? AND menuID = ?;",
		req.ItemID, req.MenuID)
	if err != nil {
		fmt.Println(err)
		return
	}

	encoder.Encode("deleted")
}

func SendMealSwipes(req structs.SendMealSwipesRequest) {
	rows, err := DB.Query("select mealSwipeBalance from persons where ID = ?;", req.FromID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	var fromMealSwipeBalance int
	for rows.Next() {
		err := rows.Scan(&fromMealSwipeBalance)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	var res structs.SendMealSwipesResponse
	if fromMealSwipeBalance >= req.NumSwipes && req.NumSwipes >= 0 {
		_, err := DB.Query("update persons set mealSwipeBalance = mealSwipeBalance + ? where ID = ?;", req.NumSwipes, req.ToID)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer rows.Close()

		_, err = DB.Query("update persons set mealSwipeBalance = mealSwipeBalance - ? where ID = ?;", req.NumSwipes, req.FromID)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer rows.Close()

		rows, err = DB.Query("select mealSwipeBalance from persons where ID = ?;", req.FromID)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&res.Balance)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		res.Success = true
	} else {
		res.Success = false
	}

	encoder.Encode(&res)
}

// dollar amounts are in cents to avoid floating point
func GetPaymentBalances(req structs.GetPaymentBalancesRequest) {
	rows, err := DB.Query("select dollarBalance, mealSwipeBalance from persons where id = ?;", req.UserID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	var res structs.GetPaymentBalancesResponse

	for rows.Next() {
		err := rows.Scan(&res.CentsBalance, &res.MealSwipeBalance)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	encoder.Encode(&res)
}

// x
func Login(req structs.LoginRequest) {
	rows, err := DB.Query("select id, student from persons where netID = ? and password = ?;", req.UserNetID, req.Password)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	res := structs.LoginResponse{Status: false, IsStudent: false, UserID: -1}

	for rows.Next() {
		res.Status = true
		err := rows.Scan(&res.UserID, &res.IsStudent)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	encoder.Encode(&res)
}
