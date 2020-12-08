//Jeffrey Oberg jo1456

// Front end

package main

import (
	"dormroomsnacks/structs"
	"strconv"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"

	"encoding/json"
	"flag"
	"fmt"
	"net"
)

var (
	connection             net.Conn
	encoder                *json.Encoder
	decoder                *json.Decoder
	app                    *iris.Application
	cookieNameForSessionID = "mycookiesessionnameid"
	sess                   = sessions.New(sessions.Config{Cookie: cookieNameForSessionID})
)

func backendComm(functionName string, req interface{}) {
	subReq := structs.Request{FunctionName: functionName}
	err := encoder.Encode(&subReq)
	if err != nil {
		panic(err.Error())
	}

	if req != nil {
		err = encoder.Encode(&req)
		if err != nil {
			panic(err.Error())
		}
	}
}

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

	tmpl := iris.HTML("./views", ".html")
	tmpl.Delims("{{", "}}")
	app.RegisterView(tmpl)

	app.Get("/", getHomePage)
	app.Get("/Login", getLoginPage)
	app.Get("/Signup", getSignupPage)
	app.Post("/loginform", loginUser)
	app.Post("/signupform", registerUser)

	// paths that require auth
	app.Get("/logout", requiresLogin, logout)
	app.Get("/Menu", requiresLogin, getLocations)
	app.Post("/create-order", requiresLogin, rediGetMenu)
	app.Get("/menu/{menuID:int}", requiresLogin, getMenu)
	app.Post("/add-item-cart", requiresLogin, addItemOrder)
	app.Get("/Cart", requiresLogin, getCart)
	app.Post("/remove-cart-item", requiresLogin, deleteCartItem)
	app.Post("/checkout", requiresLogin, submitOrder)
	app.Post("/sendMealSwipe", requiresLogin, sendMealSwipes)

	// staff only
	app.Get("/Staff/Menu", requiresStaffLogin, getLocationsTeller)
	app.Post("/get-staff-menu", requiresStaffLogin, rediGetStaffMenu)
	app.Get("/Staff/menu/{menuID:int}", requiresStaffLogin, getStaffMenu)
	app.Post("/create-or-update-menu-item", requiresStaffLogin, createUpdateItem)
	app.Post("/delete-menu-item", requiresStaffLogin, deleteItem)
	app.Get("/Staff/Orders", requiresStaffLogin, getLocationsTellerOrder)
	app.Post("/get-staff-order", requiresStaffLogin, rediGetStaffOrders)
	app.Get("/Staff/orders/{locationID:int}", requiresStaffLogin, getOrders)
	app.Post("/select-order", requiresStaffLogin, selectOrder)
	app.Post("/complete-order", requiresStaffLogin, completeOrder)

	// turn on the app
	app.Listen(":"+*listenPort, iris.WithLogLevel("debug"))
}

func requiresLogin(ctx iris.Context) {
	session := sess.Start(ctx)
	auth, _ := session.GetBoolean("authenticated")
	if !auth {
		ctx.Redirect("/", iris.StatusFound)
	}
	ctx.Next()
}

func requiresStaffLogin(ctx iris.Context) {
	session := sess.Start(ctx)
	auth, _ := session.GetBoolean("authenticated")
	isStudent, _ := session.GetBoolean("isStudent")
	if !auth || isStudent {
		ctx.Redirect("/", iris.StatusFound)
	}
	ctx.Next()
}

func loginUser(ctx iris.Context) {
	session := sess.Start(ctx)

	formData := ctx.FormValues()
	userID := formData["userID"][0]
	password := formData["password"][0]

	loginReq := structs.LoginRequest{UserNetID: userID, Password: password}
	backendComm("Login", loginReq)

	var loginRes = structs.LoginResponse{}
	err := decoder.Decode(&loginRes)
	if err != nil {
		panic(err.Error())
	}

	if !loginRes.Status {
		fmt.Println("test2")
		ctx.ViewData("error", true)
		ctx.ViewData("IsLogin", true)
		ctx.View("login.html")
	}
	fmt.Println("test1")
	session.Set("authenticated", true)
	session.Set("userID", loginRes.UserID)
	session.Set("isStudent", loginRes.IsStudent)
	ctx.Redirect("/", iris.StatusFound)
}

func registerUser(ctx iris.Context) {
	return
}

func logout(ctx iris.Context) {
	session := sess.Start(ctx)
	session.Destroy()
	ctx.Redirect("/", iris.StatusFound)
}

func getHomePage(ctx iris.Context) {
	session := sess.Start(ctx)

	userID, _ := session.GetInt("userID")
	isStudent, _ := session.GetBoolean("isStudent")
	if userID == -1 {
		ctx.ViewData("ClientName", "Guest")
		ctx.ViewData("LoggedIn", false)
	} else if isStudent {
		subReq1 := structs.GetOrderHistoryRequest{UserID: userID}
		backendComm("GetOrderHistory", subReq1)

		var res1 []structs.Order
		err := decoder.Decode(&res1)
		if err != nil {
			panic(err.Error())
		}

		ctx.ViewData("OrderHistory", res1)

		subReq2 := structs.GetPaymentBalancesRequest{UserID: userID}
		backendComm("GetPaymentBalances", subReq2)
		var res2 structs.GetPaymentBalancesResponse
		err = decoder.Decode(&res2)
		if err != nil {
			panic(err.Error())
		}
		ctx.ViewData("MealSwipes", res2.MealSwipeBalance)
		cashString := fmt.Sprintf("%d.%d", res2.CentsBalance/100, res2.CentsBalance-(res2.CentsBalance/100))
		ctx.ViewData("Cash", cashString)

		req3 := structs.GetOrderHistoryRequest{UserID: userID}
		backendComm("GetOrderHistory", req3)
		var res3 []structs.Order
		err = decoder.Decode(&res3)
		if err != nil {
			panic(err.Error())
		}
		ctx.ViewData("OrderHistory", res3)
		ctx.ViewData("IsStudent", isStudent)
		ctx.ViewData("ClientName", userID)
		ctx.ViewData("LoggedIn", true)
	} else {
		ctx.ViewData("ClientName", userID)
		ctx.ViewData("IsStudent", isStudent)
		ctx.ViewData("LoggedIn", true)
	}

	ctx.View("index.html")
}

func getLoginPage(ctx iris.Context) {
	ctx.ViewData("IsLogin", true)
	ctx.View("login.html")
}

func getSignupPage(ctx iris.Context) {
	ctx.ViewData("IsLogin", false)
	ctx.View("login.html")
}

func getLocations(ctx iris.Context) {
	backendComm("ListLocations", nil)
	var res structs.ListLocationsResponse
	err := decoder.Decode(&res)
	if err != nil {
		panic(err.Error())
	}
	ctx.ViewData("IsMenu", true)
	ctx.ViewData("IsLocSelec", true)
	ctx.ViewData("Locations", res.Locations)
	ctx.View("student.html")
}

func rediGetMenu(ctx iris.Context) {
	session := sess.Start(ctx)
	userID, _ := session.GetInt("userID")

	form := ctx.FormValues()
	formRes := strings.Split(form["IDs"][0], "-")
	menuID := formRes[0]
	locationID, _ := strconv.Atoi(formRes[0])

	// check if there is already an active order else create one - in backend
	subreq := structs.Order{UserID: userID, LocationID: locationID}
	coReq := structs.CreateOrderRequest{OrderRequest: subreq}
	backendComm("CreateOrder", coReq)

	session.Set("menuID", menuID)

	redirectLink := fmt.Sprintf("%s%s", "/menu/", menuID)
	ctx.Redirect(redirectLink, iris.StatusFound)
}

func getMenu(ctx iris.Context) {
	session := sess.Start(ctx)

	params := ctx.Params()
	menuID, err := params.GetInt("menuID")
	session.Set("menuID", menuID)
	if err != nil {
		panic(1)
	}

	menuReq := structs.GetMenuRequest{MenuID: menuID}
	backendComm("GetMenu", menuReq)
	var res structs.Menu
	err = decoder.Decode(&res)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(res)

	ctx.ViewData("IsMenu", true)
	ctx.ViewData("MenuItems", res.Items)
	ctx.View("student.html")
}

func addItemOrder(ctx iris.Context) { // add pay with meal swipe
	session := sess.Start(ctx)

	form := ctx.FormValues()
	for key, value := range form {
		fmt.Println(key, value)
	}
	itemID, err := strconv.Atoi(form["itemID"][0])
	if err != nil {
		panic(1)
	}
	_, ok := form["mealSwipe"]
	pws := false
	if ok { // to be changed
		pws = true
	} else {
		pws = false
	}
	// set the ID to 0 because i think its the orderID and i don't have access to it here
	someItem := structs.OrderItem{ID: 0, FoodID: itemID, Customization: "none", PayWithSwipe: pws}

	userID, _ := session.GetInt("userID")

	addReq := structs.AddItemToOrderRequest{PersonID: userID, Item: someItem}
	backendComm("AddItemToOrder", addReq)
	var res string
	err = decoder.Decode(&res)
	if err != nil {
		panic(err.Error())
	}

	// need to check for error, and add an error on reload
	menuID, err := session.GetInt("menuID")
	if err != nil {
		panic(1)
	}
	redirectLink := fmt.Sprintf("%s%d", "/menu/", menuID)
	ctx.Redirect(redirectLink, iris.StatusFound)
}

func getCart(ctx iris.Context) {
	session := sess.Start(ctx)
	userID, _ := session.GetInt("userID")

	cartReq := structs.GetCartRequest{UserID: userID}
	backendComm("GetCurrentUserCart", cartReq)
	var res structs.OrderAndItemsWithFood
	err := decoder.Decode(&res)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(res.Order)
	fmt.Println(res.Items)
	ctx.ViewData("IsCheckout", true)
	ctx.ViewData("Order", res.Order)
	ctx.ViewData("CartItems", res.Items)
	ctx.View("student.html")
}

func updateCartItem(ctx iris.Context) {
	// formData := ctx.FormValues()
}

func deleteCartItem(ctx iris.Context) {
	formData := ctx.FormValues()
	orderItemID, _ := strconv.Atoi(formData["orderItemID"][0])
	orderID, _ := strconv.Atoi(formData["orderID"][0])

	req := structs.DeleteItemFromOrderRequest{ItemID: orderItemID, OrderID: orderID}
	backendComm("DeleteItemFromOrder", req)
	var res string
	err := decoder.Decode(&res)
	if err != nil {
		panic(err.Error())
	}

	// need to check for error, just reload for now
	// if res == "success" {
	ctx.Redirect("/Cart", iris.StatusFound)
	// }
}

func submitOrder(ctx iris.Context) {
	form := ctx.FormValues()
	orderID, err := strconv.Atoi(form["orderID"][0])
	if err != nil {
		panic(1)
	}
	req := structs.UpdateOrderRequest{ID: orderID}
	backendComm("SubmitOrder", req)
	var res string
	err = decoder.Decode(&res)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(res)
	if res == "submitted" {
		ctx.Redirect("/", iris.StatusFound)
	} else {
		ctx.Redirect("/Cart", iris.StatusFound)
	}
}

func getLocationsTeller(ctx iris.Context) {
	backendComm("ListLocations", nil)
	var res structs.ListLocationsResponse
	err := decoder.Decode(&res)
	if err != nil {
		panic(err.Error())
	}
	ctx.ViewData("IsMenu", true)
	ctx.ViewData("IsLocSelec", true)
	ctx.ViewData("Locations", res.Locations)
	ctx.View("teller.html")
}

func rediGetStaffMenu(ctx iris.Context) {
	session := sess.Start(ctx)

	form := ctx.FormValues()
	formRes := strings.Split(form["IDs"][0], "-")
	menuID := formRes[0]

	session.Set("menuID", menuID)

	redirectLink := fmt.Sprintf("%s%s", "/Staff/menu/", menuID)
	ctx.Redirect(redirectLink, iris.StatusFound)
}

func getStaffMenu(ctx iris.Context) {
	session := sess.Start(ctx)

	params := ctx.Params()
	menuID, err := params.GetInt("menuID")
	session.Set("menuID", menuID)
	if err != nil {
		panic(1)
	}

	menuReq := structs.GetMenuRequest{MenuID: menuID}
	backendComm("GetMenu", menuReq)
	var res structs.Menu
	err = decoder.Decode(&res)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(res)

	ctx.ViewData("IsMenu", true)
	ctx.ViewData("MenuItems", res.Items)
	ctx.View("teller.html")
}

func createUpdateItem(ctx iris.Context) {
	session := sess.Start(ctx)
	menuID, _ := session.GetInt("menuID")

	form := ctx.FormValues()
	foodName := form["foodName"][0] // get value of form input with name "taskname"
	price, _ := strconv.Atoi(form["price"][0])
	description := form["description"][0]
	isAvailable := false
	if _, ok := form["isAvailable"]; ok {
		isAvailable = true
	}
	nutritionFacts := form["nutritionFacts"][0]
	// if incomingTask == "" { // check if got a task name, empty task name in update taskname does nothing! theres a delete button right there!
	// 	ctx.Redirect("/", iris.StatusFound) // reload page for the client
	// 	return
	// }
	itemID, ok := form["itemID"]
	if !ok { // check if id field exists, if no value was recieved then its a create task
		newFood := structs.FoodItem{Name: foodName, Description: description, IsAvailable: isAvailable, NutritionFacts: nutritionFacts, Cost: price}
		req := structs.CreateItemRequest{MenuID: menuID, NewItem: newFood}
		backendComm("CreateItem", req)
		var res string
		err := decoder.Decode(&res)
		if err != nil {
			panic(err.Error())
		}

	} else {
		itemIDint, _ := strconv.Atoi(itemID[0])
		newFood := structs.FoodItem{Name: foodName, Description: description, IsAvailable: isAvailable, NutritionFacts: nutritionFacts, Cost: price}
		req := structs.UpdateItemRequest{ItemID: itemIDint, MenuID: menuID, NewItem: newFood}
		backendComm("UpdateItem", req)
		var res string
		err := decoder.Decode(&res)
		if err != nil {
			panic(err.Error())
		}
	}

	redirectLink := fmt.Sprintf("%s%d", "/Staff/menu/", menuID)
	ctx.Redirect(redirectLink, iris.StatusFound)
}

func deleteItem(ctx iris.Context) {
	session := sess.Start(ctx)
	menuID, _ := session.GetInt("menuID")

	form := ctx.FormValues()
	itemID := form["itemID"][0] // attempt to retreive the id of the task in the case its an update POST
	if itemID == "" {           // if no ID was specifed then undefined error, this does not happen unless client edits the HTML
		ctx.Redirect("/", iris.StatusFound) // reload page for the client
		return
	}
	itemIDint, _ := strconv.Atoi(itemID)
	req := structs.DeleteItemRequest{MenuID: menuID, ItemID: itemIDint}
	backendComm("DeleteItem", req)
	var res string
	err := decoder.Decode(&res)
	if err != nil {
		panic(err.Error())
	}

	redirectLink := fmt.Sprintf("%s%d", "/Staff/menu/", menuID)
	ctx.Redirect(redirectLink, iris.StatusFound)
}

func getLocationsTellerOrder(ctx iris.Context) {
	backendComm("ListLocations", nil)
	var res structs.ListLocationsResponse
	err := decoder.Decode(&res)
	if err != nil {
		panic(err.Error())
	}
	ctx.ViewData("IsOrders", true)
	ctx.ViewData("IsLocSelec", true)
	ctx.ViewData("Locations", res.Locations)
	ctx.View("teller.html")
}

func rediGetStaffOrders(ctx iris.Context) {
	session := sess.Start(ctx)

	form := ctx.FormValues()
	formRes := strings.Split(form["IDs"][0], "-")
	locationID := formRes[1]

	session.Set("locationID", locationID)

	redirectLink := fmt.Sprintf("%s%s", "/Staff/orders/", locationID)
	ctx.Redirect(redirectLink, iris.StatusFound)
}

// for teller - only returns order IDs
func getOrders(ctx iris.Context) {
	session := sess.Start(ctx)

	params := ctx.Params()
	locationID, err := params.GetInt("locationID")
	session.Set("locationID", locationID)
	if err != nil {
		panic(1)
	}

	ordersReq := structs.GetOrdersRequest{LocationID: locationID}
	backendComm("GetOrders", ordersReq)
	var res []structs.Order
	err = decoder.Decode(&res)
	if err != nil {
		panic(err.Error())
	}
	ctx.ViewData("IsOrders", true)
	ctx.ViewData("Orders", res)
	ctx.View("teller.html")
}

// get all details for a specific - returns food items - changes status - add new section for detailed food view
func selectOrder(ctx iris.Context) {
	formData := ctx.FormValues()
	orderID, _ := strconv.Atoi(formData["orderID"][0])
	req := structs.SelectOrderRequest{OrderID: orderID}
	backendComm("SelectOrder", req)
	var res structs.OrderAndItems
	err := decoder.Decode(&res)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(res)
	ctx.ViewData("IsSelected", true)
	ctx.ViewData("Order", res.Order)
	ctx.ViewData("Items", res.Items)
	ctx.View("teller.html")
}

func completeOrder(ctx iris.Context) {
	session := sess.Start(ctx)
	locationID, _ := session.GetInt("locationID")

	formData := ctx.FormValues()
	orderID, _ := strconv.Atoi(formData["orderID"][0])
	req := structs.CompelteOrderRequest{OrderID: orderID}
	backendComm("CompleteOrder", req)

	redirectLink := fmt.Sprintf("%s%d", "/Staff/menu/", locationID)
	ctx.Redirect(redirectLink, iris.StatusFound)
}

// add form on home page
func sendMealSwipes(ctx iris.Context) {
	sessions := sess.Start(ctx)
	userID, _ := sessions.GetInt("userID")

	formData := ctx.FormValues()
	toID := formData["toID"][0]
	numberSwipes, _ := strconv.Atoi(formData["numberSwipes"][0])

	subReq := structs.SendMealSwipesRequest{ToID: toID, FromID: userID, NumSwipes: numberSwipes}
	backendComm("SendMealSwipes", subReq)
	var res structs.SendMealSwipesResponse
	err := decoder.Decode(&res)
	if err != nil {
		panic(err.Error())
	}
	if res.Success {
		fmt.Println("swipe sent sucess")
	} else {
		fmt.Println("swipe sent fail")
	}

	ctx.Redirect("/", iris.StatusFound)
}
