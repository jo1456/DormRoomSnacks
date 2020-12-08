package main

import (
	"dormroomsnacks/structs"
	"encoding/json"
	"flag"
	"fmt"
	"sync"

	// "log"
	"net"
	"time"

	// "errors"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	// "github.com/joho/godotenv"
)

var (
	DB              *sql.DB
	foodMutex       *sync.RWMutex
	menuMutex       *sync.RWMutex
	orderMutex      *sync.RWMutex
	personMutex     *sync.RWMutex
	diningHallMutex *sync.RWMutex
	orderItemMutex  *sync.RWMutex
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
	foodMutex = &sync.RWMutex{}
	menuMutex = &sync.RWMutex{}
	orderMutex = &sync.RWMutex{}
	personMutex = &sync.RWMutex{}
	diningHallMutex = &sync.RWMutex{}
	orderItemMutex = &sync.RWMutex{}

	db, err := sql.Open("mysql", "b2766d1c91f7c7:0c0f617f@tcp(us-cdbr-east-02.cleardb.com)/heroku_5873df879639de6")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	DB = db

	// creates a variable with the passed flag. default value of 8080
	listenPort := flag.String("listen", "8090", "port to listen on")
	flag.Parse()

	for {
		// listen on passed port
		listener, err := net.Listen("tcp", ":"+*listenPort)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer listener.Close()

		for {
			// accept a connection
			connection, err := listener.Accept()
			if err != nil {
				fmt.Println(err)
				break
			}
			go ProcessConnection(connection)
		}
	}
}

// x
func ListLocations(encoder *json.Encoder, decoder *json.Decoder) {
	diningHallMutex.RLock()
	resp := structs.ListLocationsResponse{Locations: make([]structs.Location, 0)}
	respFail := structs.ListLocationsResponse{Locations: make([]structs.Location, 0)}

	rows, err := DB.Query("select * from DiningHalls;")
	if err != nil {
		fmt.Println(err)
		encoder.Encode(&respFail)
		diningHallMutex.RUnlock()
		return
	}
	defer rows.Close()

	for rows.Next() {
		var location structs.Location
		err := rows.Scan(&location.ID, &location.Name, &location.Address, &location.Phone, &location.MenuID, &location.Hours)
		if err != nil {
			fmt.Println(err)
			encoder.Encode(&respFail)
			diningHallMutex.RUnlock()
			return
		}
		resp.Locations = append(resp.Locations, structs.Location{ID: location.ID, Name: location.Name, Address: location.Address, Phone: location.Phone, MenuID: location.MenuID, Hours: location.Hours})
	}

	encoder.Encode(&resp)
	diningHallMutex.RUnlock()
}

// x
func GetMenu(encoder *json.Encoder, decoder *json.Decoder) {
	menuMutex.RLock()
	foodMutex.RLock()
	var req structs.GetMenuRequest
	decoder.Decode(&req)
	var menu structs.Menu
	var menuFail structs.Menu

	rows, err := DB.Query("select * from Menu where id = ?;", req.MenuID)
	if err != nil {
		fmt.Println(err)
		menuMutex.RUnlock()
		foodMutex.RUnlock()
		encoder.Encode(&menuFail)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&menu.ID, &menu.Name, &menu.LocationID)
		if err != nil {
			fmt.Println(err)
			menuMutex.RUnlock()
			foodMutex.RUnlock()
			encoder.Encode(&menuFail)
			return
		}
	}

	rows, err = DB.Query("select * from Foods where menuID = ?;", req.MenuID)
	if err != nil {
		fmt.Println(err)
		menuMutex.RUnlock()
		foodMutex.RUnlock()
		encoder.Encode(&menuFail)

		return
	}
	defer rows.Close()

	for rows.Next() {
		var item structs.FoodItem
		var num int
		err := rows.Scan(&item.ID, &num, &item.Name, &item.Description, &item.Cost, &item.IsAvailable, &item.NutritionFacts)
		if err != nil {
			fmt.Println(err)
			menuMutex.RUnlock()
			foodMutex.RUnlock()
			encoder.Encode(&menuFail)
			return
		}
		menu.Items = append(menu.Items, item)
	}
	encoder.Encode(&menu)
	menuMutex.RUnlock()
	foodMutex.RUnlock()
}

// x
func ViewItem(encoder *json.Encoder, decoder *json.Decoder) {
	foodMutex.RLock()
	var req structs.ViewItemRequest
	decoder.Decode(&req)
	var item structs.FoodItem
	var itemFail structs.FoodItem

	rows, err := DB.Query("select * from Foods where ID = ?;", req.ItemID)
	if err != nil {
		fmt.Println(err)
		foodMutex.RUnlock()
		encoder.Encode(&itemFail)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var num int
		err := rows.Scan(&item.ID, &num, &item.Name, &item.Description, &item.Cost, &item.IsAvailable, &item.NutritionFacts)
		if err != nil {
			fmt.Println(err)
			foodMutex.RUnlock()
			encoder.Encode(&itemFail)
			return
		}
	}

	encoder.Encode(&item)
	foodMutex.RUnlock()
}

// ORDER FLOW
// 1. Create an order when 1 item is added to the cart and set the order status to "cart"
// 2. Add items to order
// 3. Submit order which updates to status to "submitted"

// cart -> submitted ->

func CreateOrder(encoder *json.Encoder, decoder *json.Decoder) {
	orderMutex.Lock()
	var req structs.CreateOrderRequest
	decoder.Decode(&req)

	rows, err := DB.Query("SELECT personID from Orders where status = 'Cart' and personID = ?;", req.OrderRequest.UserID)
	if err != nil {
		fmt.Println(err)
		orderMutex.Unlock()
		return
	}
	defer rows.Close()

	var id int
	for rows.Next() {
		rows.Scan(&id)
		fmt.Println(err)
		orderMutex.Unlock()
		return
	}

	_, err = DB.Query("insert into Orders (personID, diningHallID, status, submitTime, lastStatusChange, swipeCost, centCost) values(?,?,?,?,?,?,?);",
		req.OrderRequest.UserID, req.OrderRequest.LocationID, "Cart", time.Now(), time.Now(), 0, 0)
	if err != nil {
		fmt.Println(err)
		orderMutex.Unlock()
		return
	}
	orderMutex.Unlock()
}

// x
func SubmitOrder(encoder *json.Encoder, decoder *json.Decoder) {
	orderMutex.Lock()
	personMutex.Lock()
	var req structs.UpdateOrderRequest
	decoder.Decode(&req)

	rows, err := DB.Query("select swipeCost, centCost, personID from Orders where ID = ?;", req.ID)
	if err != nil {
		fmt.Println(err)
		encoder.Encode("failure")
		orderMutex.Unlock()
		personMutex.Unlock()
		return
	}
	defer rows.Close()

	var swipeCost int
	var centCost int
	var userID int
	for rows.Next() {
		err := rows.Scan(&swipeCost, &centCost, &userID)
		if err != nil {
			fmt.Println(err)
			encoder.Encode("failure")
			orderMutex.Unlock()
			personMutex.Unlock()
			return
		}
	}

	rows, err = DB.Query("select dollarBalance, mealSwipeBalance from Persons where ID = ?;", userID)
	if err != nil {
		fmt.Println(err)
		encoder.Encode("failure")
		orderMutex.Unlock()
		personMutex.Unlock()
		return
	}
	defer rows.Close()

	var dollarBalance int
	var mealSwipeBalance int
	for rows.Next() {
		err := rows.Scan(&swipeCost, &centCost)
		if err != nil {
			fmt.Println(err)
			encoder.Encode("failure")
			orderMutex.Unlock()
			personMutex.Unlock()
			return
		}
	}

	if dollarBalance > centCost && mealSwipeBalance > swipeCost {

		_, err := DB.Query("update Personds set dollarBalance = dollarBalance - ?, mealSwipeBalance = mealSwipeBalance - ? where id = ?;",
			centCost, swipeCost, userID)
		if err != nil {
			fmt.Println(err)
			encoder.Encode("failure")
			orderMutex.Unlock()
			personMutex.Unlock()
			return
		}

		_, err = DB.Query("update Orders set status = \"submitted\", submitTime = ? where id = ?;",
			time.Now(), req.ID)
		if err != nil {
			fmt.Println(err)
			encoder.Encode("failure")
			orderMutex.Unlock()
			personMutex.Unlock()
			return
		}

		encoder.Encode("submitted")
		orderMutex.Unlock()
		personMutex.Unlock()
	}

	encoder.Encode("failure")
	orderMutex.Unlock()
	personMutex.Unlock()
}

// When Adding an item to an your must select if you will pay with the item with
// meal swipes or with dollars. This will build up 2 accumulators in the row in the
// orders table. These accumulators will then be used to see if the user has sufficent balance
// to pay for their order.
func AddItemToOrder(encoder *json.Encoder, decoder *json.Decoder) {
	orderMutex.Lock()
	orderItemMutex.Lock()
	foodMutex.RLock()
	var req structs.AddItemToOrderRequest
	decoder.Decode(&req)

	rows, err := DB.Query("select id from orders where personID = ? and status = 'Cart';", req.PersonID)
	if err != nil {
		fmt.Println(err)
		encoder.Encode("failure")
		orderMutex.Unlock()
		orderItemMutex.Unlock()
		foodMutex.RUnlock()
		return
	}
	defer rows.Close()

	var orderID int
	for rows.Next() {
		err := rows.Scan(&orderID)
		if err != nil {
			fmt.Println(err)
			encoder.Encode("failure")
			orderMutex.Unlock()
			orderItemMutex.Unlock()
			foodMutex.RUnlock()
			return
		}
	}

	_, err = DB.Query("insert into OrderItem (foodID, orderID, Customization, payWithSwipe) values(?,?,?,?);", req.Item.FoodID, orderID, req.Item.Customization, req.Item.PayWithSwipe)
	if err != nil {
		fmt.Println(err)
		encoder.Encode("failure")
		orderMutex.Unlock()
		orderItemMutex.Unlock()
		foodMutex.RUnlock()
		return
	}

	if req.Item.PayWithSwipe {
		_, err := DB.Query("update Orders set swipeCost = swipeCost + 1 where id = ?;", orderID)
		if err != nil {
			fmt.Println(err)
			encoder.Encode("failure")
			orderMutex.Unlock()
			orderItemMutex.Unlock()
			foodMutex.RUnlock()
			return
		}
	} else {
		rows, err := DB.Query("select price from Foods where ID = ?;", req.Item.FoodID)
		if err != nil {
			fmt.Println(err)
			encoder.Encode("failure")
			orderMutex.Unlock()
			orderItemMutex.Unlock()
			foodMutex.RUnlock()
			return
		}
		defer rows.Close()

		var price int
		for rows.Next() {
			err := rows.Scan(&price)
			if err != nil {
				fmt.Println(err)
				encoder.Encode("failure")
				orderMutex.Unlock()
				orderItemMutex.Unlock()
				foodMutex.RUnlock()
				return
			}
		}

		_, err = DB.Query("update Orders set centCost = centCost + ? where id = ?;", price, orderID)
		if err != nil {
			fmt.Println(err)
			encoder.Encode("failure")
			orderMutex.Unlock()
			orderItemMutex.Unlock()
			foodMutex.RUnlock()
			return
		}
	}
	encoder.Encode("success")
	orderMutex.Unlock()
	orderItemMutex.Unlock()
	foodMutex.RUnlock()
}

func GetCurrentUserCart(encoder *json.Encoder, decoder *json.Decoder) {
	orderMutex.RLock()
	orderItemMutex.RLock()
	foodMutex.RLock()

	var req structs.GetCartRequest
	decoder.Decode(&req)

	var orderAndItemsWithFood structs.OrderAndItemsWithFood
	var orderAndItemsWithFoodFail structs.OrderAndItemsWithFood

	rows, err := DB.Query("select * from Orders where personID = ? and status = 'Cart';", req.UserID)
	if err != nil {
		fmt.Println(err)
		orderMutex.RUnlock()
		orderItemMutex.RUnlock()
		foodMutex.RUnlock()
		encoder.Encode(&orderAndItemsWithFoodFail)
		return
	}
	defer rows.Close()

	var order structs.Order
	for rows.Next() {
		var lastStatusChange string
		err := rows.Scan(&orderAndItemsWithFood.Order.ID,
			&orderAndItemsWithFood.Order.UserID, &orderAndItemsWithFood.Order.LocationID,
			&orderAndItemsWithFood.Order.Status, &orderAndItemsWithFood.Order.SubmitTime,
			&lastStatusChange, &orderAndItemsWithFood.Order.SwipeCost, &orderAndItemsWithFood.Order.CentCost)
		if err != nil {
			fmt.Println(err)
			orderMutex.RUnlock()
			orderItemMutex.RUnlock()
			foodMutex.RUnlock()
			encoder.Encode(&orderAndItemsWithFoodFail)
			return
		}
	}

	rows, err = DB.Query("select * from OrderItem where orderID = ?;", order.ID)
	if err != nil {
		fmt.Println(err)
		orderMutex.RUnlock()
		orderItemMutex.RUnlock()
		foodMutex.RUnlock()
		encoder.Encode(&orderAndItemsWithFoodFail)
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
			orderMutex.RUnlock()
			orderItemMutex.RUnlock()
			foodMutex.RUnlock()
			encoder.Encode(&orderAndItemsWithFoodFail)
			return
		}

		foodRows, err := DB.Query("select * from Foods where id = ?;", item.FoodID)
		if err != nil {
			fmt.Println(err)
			orderMutex.RUnlock()
			orderItemMutex.RUnlock()
			foodMutex.RUnlock()
			encoder.Encode(&orderAndItemsWithFoodFail)
			return
		}
		defer foodRows.Close()
		for foodRows.Next() {
			var menuID int
			err = rows.Scan(&food.ID, &menuID, &food.Name, &food.Description, &food.Cost, &food.IsAvailable, &food.NutritionFacts)
			if err != nil {
				fmt.Println(err)
				orderMutex.RUnlock()
				orderItemMutex.RUnlock()
				foodMutex.RUnlock()
				encoder.Encode(&orderAndItemsWithFoodFail)
				return
			}
		}
		items = append(items, structs.OrderItemWithFood{Item: item, Food: food})
	}
	orderAndItemsWithFood.Items = items

	encoder.Encode(&orderAndItemsWithFood)
	orderMutex.RUnlock()
	orderItemMutex.RUnlock()
	foodMutex.RUnlock()
}

func DeleteItemFromOrder(encoder *json.Encoder, decoder *json.Decoder) {
	orderMutex.Lock()
	orderItemMutex.Lock()
	foodMutex.RLock()
	var req structs.DeleteItemFromOrderRequest
	decoder.Decode(&req)

	rows, err := DB.Query("select * from OrderItem where id = ?;", req.ItemID)
	if err != nil {
		fmt.Println(err)
		encoder.Encode("failure")
		orderMutex.Unlock()
		orderItemMutex.Unlock()
		foodMutex.RUnlock()
		return
	}
	defer rows.Close()

	var item structs.OrderItem
	for rows.Next() {
		var lastStatusChange string
		err := rows.Scan(&item.ID, &item.FoodID, &lastStatusChange, &item.Customization, &item.PayWithSwipe)
		if err != nil {
			fmt.Println(err)
			encoder.Encode("failure")
			orderMutex.Unlock()
			orderItemMutex.Unlock()
			foodMutex.RUnlock()
			return
		}
	}

	_, err = DB.Query("delete from OrderItem where id = ?;", req.ItemID)
	if err != nil {
		fmt.Println(err)
		encoder.Encode("failure")
		orderMutex.Unlock()
		orderItemMutex.Unlock()
		foodMutex.RUnlock()
		return
	}

	if item.PayWithSwipe {
		_, err := DB.Query("update Orders set swipeCost = swipeCost - 1 where id = ?;", req.OrderID)
		if err != nil {
			fmt.Println(err)
			encoder.Encode("failure")
			orderMutex.Unlock()
			orderItemMutex.Unlock()
			foodMutex.RUnlock()
			return
		}
	} else {
		rows, err := DB.Query("select price from Foods where ID = ?;", item.FoodID)
		if err != nil {
			fmt.Println(err)
			encoder.Encode("failure")
			orderMutex.Unlock()
			orderItemMutex.Unlock()
			foodMutex.RUnlock()
			return
		}
		defer rows.Close()

		var price int
		for rows.Next() {
			err := rows.Scan(&price)
			if err != nil {
				fmt.Println(err)
				encoder.Encode("failure")
				orderMutex.Unlock()
				orderItemMutex.Unlock()
				foodMutex.RUnlock()
				return
			}
		}

		_, err = DB.Query("update Orders set centCost = centCost - ? where id = ?;", price, req.OrderID)
		if err != nil {
			fmt.Println(err)
			encoder.Encode("failure")
			orderMutex.Unlock()
			orderItemMutex.Unlock()
			foodMutex.RUnlock()
			return
		}
	}

	encoder.Encode("success")
	orderMutex.Unlock()
	orderItemMutex.Unlock()
	foodMutex.RUnlock()
}

// x
func GetOrderHistory(encoder *json.Encoder, decoder *json.Decoder) {
	orderMutex.RLock()
	var req structs.GetOrderHistoryRequest
	decoder.Decode(&req)
	var orders []structs.Order
	var ordersFail []structs.Order

	rows, err := DB.Query("select * from Orders where personID = ?;", req.UserID)
	if err != nil {
		fmt.Println(err)
		orderMutex.RUnlock()
		encoder.Encode(&ordersFail)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var order structs.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.LocationID, &order.Status, &order.SubmitTime, &order.LastStatusChange, &order.SwipeCost, &order.CentCost)
		if err != nil {
			fmt.Println(err)
			orderMutex.RUnlock()
			encoder.Encode(&ordersFail)
			return
		}
		orders = append(orders, order)
	}
	encoder.Encode(&orders)
	orderMutex.RUnlock()
}

// x
func GetOrders(encoder *json.Encoder, decoder *json.Decoder) {
	orderMutex.RLock()
	var req structs.GetOrdersRequest
	decoder.Decode(&req)
	var orders []structs.Order
	var ordersFail []structs.Order

	rows, err := DB.Query("select * from Orders where diningHallID = ? and (status = \"submitted\" or status \"selected\");", req.LocationID)
	if err != nil {
		fmt.Println(err)
		orderMutex.RUnlock()
		encoder.Encode(&ordersFail)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var order structs.Order
		var lastStatusChange string
		err := rows.Scan(&order.ID, &order.UserID, &order.LocationID, &order.Status, &order.SubmitTime, &lastStatusChange)
		if err != nil {
			fmt.Println(err)
			orderMutex.RUnlock()
			encoder.Encode(&ordersFail)
			return
		}
		orders = append(orders, order)
	}
	encoder.Encode(&orders)
	orderMutex.RUnlock()
}

// x
func SelectOrder(encoder *json.Encoder, decoder *json.Decoder) {
	orderMutex.Lock()
	orderItemMutex.RLock()
	foodMutex.RLock()
	var req structs.SelectOrderRequest
	decoder.Decode(&req)

	var orderAndItems structs.OrderAndItemsWithFood
	var orderAndItemsFail structs.OrderAndItemsWithFood

	_, err := DB.Query("update Orders set status = \"selected\", lastStatusChange = ? where id = ?;",
		time.Now(), req.OrderID)
	if err != nil {
		fmt.Println(err)
		orderMutex.Unlock()
		orderItemMutex.RUnlock()
		foodMutex.RUnlock()
		encoder.Encode(&orderAndItemsFail)
		return
	}

	rows, err := DB.Query("select * from Orders where ID = ?;", req.OrderID)
	if err != nil {
		fmt.Println(err)
		orderMutex.Unlock()
		orderItemMutex.RUnlock()
		foodMutex.RUnlock()
		encoder.Encode(&orderAndItemsFail)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var lastStatusChange string
		err := rows.Scan(&orderAndItems.Order.ID, &orderAndItems.Order.UserID, &orderAndItems.Order.LocationID,
			&orderAndItems.Order.Status, &orderAndItems.Order.SubmitTime, &lastStatusChange)
		if err != nil {
			fmt.Println(err)
			orderMutex.Unlock()
			orderItemMutex.RUnlock()
			foodMutex.RUnlock()
			encoder.Encode(&orderAndItemsFail)
			return
		}
	}

	rows, err = DB.Query("select * from OrderItem, Foods where OrderItem.orderID = ? and Foods.id = OrderItem.foodID;", req.OrderID)
	if err != nil {
		fmt.Println(err)
		orderMutex.Unlock()
		orderItemMutex.RUnlock()
		foodMutex.RUnlock()
		encoder.Encode(&orderAndItemsFail)
		return
	}
	defer rows.Close()

	var items []structs.OrderItemWithFood
	for rows.Next() {
		var itemWithFood structs.OrderItemWithFood
		var orderID int
		var menuID int
		err := rows.Scan(&itemWithFood.Item.ID, &itemWithFood.Item.FoodID, &orderID,
			 &itemWithFood.Item.Customization, &itemWithFood.Item.PayWithSwipe, &itemWithFood.Food.ID,
			 &menuID, &itemWithFood.Food.Name, &itemWithFood.Food.Description,
			 &itemWithFood.Food.Cost, &itemWithFood.Food.IsAvailable, &itemWithFood.Food.NutritionFacts)
		if err != nil {
			fmt.Println(err)
			orderMutex.Unlock()
			orderItemMutex.RUnlock()
			foodMutex.RUnlock()
			encoder.Encode(&orderAndItemsFail)
			return
		}
		items = append(items, itemWithFood)
	}

	orderAndItems.Items = items

	encoder.Encode(&orderAndItems)
	orderMutex.Unlock()
	orderItemMutex.RUnlock()
	foodMutex.RUnlock()
}

// x
func CompleteOrder(encoder *json.Encoder, decoder *json.Decoder) {
	orderMutex.Lock()
	var req structs.CompelteOrderRequest
	decoder.Decode(&req)

	_, err := DB.Query("update Orders set status = \"complete\", submitTime = ? where id = ?;",
		time.Now(), req.OrderID)
	if err != nil {
		fmt.Println(err)
		orderMutex.Unlock()
		encoder.Encode("failure")
		return
	}

	encoder.Encode("complete")
	orderMutex.Unlock()
}

// x
func CreateItem(encoder *json.Encoder, decoder *json.Decoder) {
	foodMutex.Lock()
	var req structs.CreateItemRequest
	decoder.Decode(&req)

	_, err := DB.Query("insert into Foods (menuID, name, description, price, availability, nutritionFacts) values (?,?,?,?,?,?);",
		req.MenuID, req.NewItem.Name, req.NewItem.Description, req.NewItem.Cost, req.NewItem.IsAvailable, req.NewItem.NutritionFacts)
	if err != nil {
		fmt.Println(err)
		foodMutex.Unlock()
		encoder.Encode("failure")
		return
	}

	encoder.Encode("created")
	foodMutex.Unlock()
}

// x
func UpdateItem(encoder *json.Encoder, decoder *json.Decoder) {
	foodMutex.Lock()
	var req structs.UpdateItemRequest
	decoder.Decode(&req)

	_, err := DB.Query("update Foods set name = ?, description = ?, price = ?, availability = ?, nutritionFacts = ? where id = ?;",
		req.NewItem.Name, req.NewItem.Description, req.NewItem.Cost, req.NewItem.IsAvailable, req.NewItem.NutritionFacts, req.ItemID)
	if err != nil {
		fmt.Println(err)
		foodMutex.Unlock()
		encoder.Encode("failure")
		return
	}
	encoder.Encode("updated")
	foodMutex.Unlock()
}

// x
func DeleteItem(encoder *json.Encoder, decoder *json.Decoder) {
	foodMutex.Lock()
	var req structs.DeleteItemRequest
	decoder.Decode(&req)

	_, err := DB.Query("delete from Foods where id = ? AND menuID = ?;",
		req.ItemID, req.MenuID)
	if err != nil {
		fmt.Println(err)
		foodMutex.Unlock()
		encoder.Encode("failure")
		return
	}

	encoder.Encode("deleted")
	foodMutex.Unlock()
}

func SendMealSwipes(encoder *json.Encoder, decoder *json.Decoder) {
	personMutex.Lock()
	var req structs.SendMealSwipesRequest
	decoder.Decode(&req)
	var res structs.SendMealSwipesResponse
	res.Success = false

	rows, err := DB.Query("select mealSwipeBalance from persons where netID = ?;", req.FromID)
	if err != nil {
		fmt.Println(err)
		personMutex.Unlock()
		encoder.Encode(&res)
		return
	}
	defer rows.Close()

	var fromMealSwipeBalance int
	for rows.Next() {
		err := rows.Scan(&fromMealSwipeBalance)
		if err != nil {
			fmt.Println(err)
			personMutex.Unlock()
			encoder.Encode(&res)
			return
		}
	}

	if fromMealSwipeBalance >= req.NumSwipes && req.NumSwipes >= 0 {
		_, err := DB.Query("update persons set mealSwipeBalance = mealSwipeBalance + ? where ID = ?;", req.NumSwipes, req.ToID)
		if err != nil {
			fmt.Println(err)
			personMutex.Unlock()
			encoder.Encode(&res)
			return
		}
		defer rows.Close()

		_, err = DB.Query("update persons set mealSwipeBalance = mealSwipeBalance - ? where ID = ?;", req.NumSwipes, req.FromID)
		if err != nil {
			fmt.Println(err)
			personMutex.Unlock()
			encoder.Encode(&res)
			return
		}
		defer rows.Close()

		rows, err = DB.Query("select mealSwipeBalance from persons where ID = ?;", req.FromID)
		if err != nil {
			fmt.Println(err)
			personMutex.Unlock()
			encoder.Encode(&res)
			return
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&res.Balance)
			if err != nil {
				fmt.Println(err)
				personMutex.Unlock()
				encoder.Encode(&res)
				return
			}
		}

		res.Success = true
	} else {
		res.Success = false
	}

	encoder.Encode(&res)
	personMutex.Unlock()
}

// dollar amounts are in cents to avoid floating point
func GetPaymentBalances(encoder *json.Encoder, decoder *json.Decoder) {
	personMutex.RLock()
	var req structs.GetPaymentBalancesRequest
	decoder.Decode(&req)
	var res structs.GetPaymentBalancesResponse
	var failRes structs.GetPaymentBalancesResponse

	rows, err := DB.Query("select dollarBalance, mealSwipeBalance from persons where id = ?;", req.UserID)
	if err != nil {
		fmt.Println(err)
		personMutex.RUnlock()
		encoder.Encode(&failRes)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&res.CentsBalance, &res.MealSwipeBalance)
		if err != nil {
			fmt.Println(err)
			personMutex.RUnlock()
			encoder.Encode(&failRes)
			return
		}
	}

	encoder.Encode(&res)
	personMutex.RUnlock()
}

// x
func Login(encoder *json.Encoder, decoder *json.Decoder) {
	personMutex.RLock()
	var req structs.LoginRequest
	decoder.Decode(&req)

	res := structs.LoginResponse{Status: false, IsStudent: false, UserID: -1}
	failRes := structs.LoginResponse{Status: false, IsStudent: false, UserID: -1}

	rows, err := DB.Query("select id, student from persons where netID = ? and password = ?;", req.UserNetID, req.Password)
	if err != nil {
		fmt.Println(err)
		personMutex.RUnlock()
		encoder.Encode(&failRes)
		return
	}
	defer rows.Close()

	for rows.Next() {
		res.Status = true
		err := rows.Scan(&res.UserID, &res.IsStudent)
		if err != nil {
			fmt.Println(err)
			personMutex.RUnlock()
			encoder.Encode(&failRes)
			return
		}
	}

	encoder.Encode(&res)
	personMutex.RUnlock()
}

func ProcessConnection(connection net.Conn) {

	encoder := json.NewEncoder(connection)
	decoder := json.NewDecoder(connection)

	fmt.Println("got a connection")

	closed := false
	// infinite loop to accept and response to requests
	for !closed {

		// status of operation. In case of failure fe can redirect to homepage
		// status := Status{Success: true}

		var req structs.Request
		decoder.Decode(&req)

		switch req.FunctionName {
		case "ListLocations":
			ListLocations(encoder, decoder)
		case "GetMenu":
			GetMenu(encoder, decoder)
		case "ViewItem":
			ViewItem(encoder, decoder)
		case "CreateOrder":
			CreateOrder(encoder, decoder)
		case "SubmitOrder":
			SubmitOrder(encoder, decoder)
		case "AddItemToOrder":
			AddItemToOrder(encoder, decoder)
		case "GetOrderHistory":
			GetOrderHistory(encoder, decoder)
		case "GetOrders":
			GetOrders(encoder, decoder)
		case "SelectOrder":
			SelectOrder(encoder, decoder)
		case "CompleteOrder":
			CompleteOrder(encoder, decoder)
		case "UpdateItem":
			UpdateItem(encoder, decoder)
		case "CreateItem":
			CreateItem(encoder, decoder)
		case "DeleteItem":
			DeleteItem(encoder, decoder)
		case "SendMealSwipes":
			SendMealSwipes(encoder, decoder)
		case "GetPaymentBalances": // dollar amounts are in cents to avoid floating point
			GetPaymentBalances(encoder, decoder)
		case "Login":
			Login(encoder, decoder)
		case "DeleteItemFromOrder":
			DeleteItemFromOrder(encoder, decoder)
		case "GetCurrentUserCart":
			GetCurrentUserCart(encoder, decoder)
		case "Closed":
			closed = true
		}
	}
	connection.Close()
}
