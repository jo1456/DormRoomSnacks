//Jeffrey Oberg jo1456

// Front end

package main

import (
	"dormroomsnacks/structs"
	"strconv"

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

func backendComm(req interface{}) interface{} {
	err := encoder.Encode(&req)
	if err != nil {
		panic(err.Error())
	}

	var response interface{}
	err = decoder.Decode(&response)
	if err != nil {
		panic(err.Error())
	}
	return response
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
	app.Get("/loginform", loginUser)
	app.Post("/signupform", registerUser)

	// paths that require auth
	app.Get("/logout", requiresLogin, logout)
	app.Get("/Menu", requiresLogin, getLocations) // students
	app.Get("/get-menu-student", requiresLogin, rediGetMenu)
	app.Get("/menu/{menuID:int}", requiresLogin, getMenu)
	app.Post("/add-item-cart", requiresLogin, addItemOrder)
	app.Get("/Cart", requiresLogin, getCart) // how do i do this?
	app.Post("/checkout", requiresLogin, submitOrder)
	app.Get("/Orders", requiresLogin, getLocationsTeller) // teller
	app.Get("/get-orders", requiresLogin, getOrders)

	// turn on the app
	app.Listen(":" + *listenPort)
}

func requiresLogin(ctx iris.Context) {
	session := sess.Start(ctx)
	auth, _ := session.GetBoolean("authenticated")
	if !auth {
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
	loginRes, err := backendComm(loginReq).(structs.LoginResponse)
	if !err {
		panic(1)
	}
	if !loginRes.Status {
		ctx.ViewData("error", true)
		ctx.ViewData("IsLogin", true)
		ctx.View("login.html")
	}
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
}

func getHomePage(ctx iris.Context) {
	session := sess.Start(ctx)

	userID, _ := session.GetInt("userID")
	isStudent, _ := session.GetBoolean("isStudent")
	if userID == -1 {
		ctx.ViewData("ClientName", "Guest")
		ctx.ViewData("LoggedIn", false)
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

// broken after, can't change selection
func getLocations(ctx iris.Context) {
	req := structs.Request{FunctionName: "ListLocations"}
	locations, ok := backendComm(req).(structs.ListLocationsResponse)
	if !ok {
		panic(1)
	}
	ctx.ViewData("isLocSelec", true)
	ctx.ViewData("Locations", locations.Locations)
	ctx.View("student.html")
}

func getLocationsTeller(ctx iris.Context) {
	req := structs.Request{FunctionName: "ListLocations"}
	locations, ok := backendComm(req).(structs.ListLocationsResponse)
	if !ok {
		panic(1)
	}
	ctx.ViewData("isLocSelec", true)
	ctx.ViewData("Locations", locations.Locations)
	ctx.View("teller.html")
}

func rediGetMenu(ctx iris.Context) {
	session := sess.Start(ctx)
	userID, _ := session.GetInt("userID")

	form := ctx.FormValues()
	menuID := form["menuID"][0]

	// check if there is already an active order else create one - in backend
	coReq := structs.CreateOrderRequest{UserID: userID}
	req := structs.Request{FunctionName: "GetMenu", Data: coReq}
	backendComm(req)

	redirectLink := fmt.Sprintf("%s%s", "/menu/", menuID)
	ctx.Redirect(redirectLink, iris.StatusFound)
}

func getMenu(ctx iris.Context) {
	params := ctx.Params()
	menuID, err := params.GetInt("menuID")
	if err != nil {
		panic(1)
	}

	menuReq := structs.GetMenuRequest{MenuID: menuID}
	req := structs.Request{FunctionName: "GetMenu", Data: menuReq}
	menu, ok := backendComm(req).(structs.Menu)
	if !ok {
		panic(1)
	}

	ctx.ViewData("isMenu", true)
	ctx.ViewData("Menu", menu)
	ctx.View("student.html")
}

func addItemOrder(ctx iris.Context) { // add pay with meal swipe
	form := ctx.FormValues()
	itemID, err := strconv.Atoi(form["itemID"][0])
	if err != nil {
		panic(1)
	}
	orderID, err := strconv.Atoi(form["orderID"][0])
	if err != nil {
		panic(1)
	}
	pws := form["mealSwipe"][0]
	pwsB := false
	if pws == "something" {
		pwsB = true
	} else {
		pwsB = false
	}
	someItem := structs.OrderItem{ID: orderID, FoodID: itemID, Customization: "none", PayWithSwipe: pwsB}
	addReq := structs.AddItemToOrderRequest{Item: someItem}
	req := structs.Request{FunctionName: "AddItemToOrder", Data: addReq}
	res, ok := backendComm(req).(string)
	if !ok {
		panic(1)
	}
	if res == "failure" {

	} else {

	}
}

func getCart(ctx iris.Context) {

}

func submitOrder(ctx iris.Context) {
	form := ctx.FormValues()
	orderID, err := strconv.Atoi(form["orderID"][0])
	if err != nil {
		panic(1)
	}
	req := structs.UpdateOrderRequest{ID: orderID}
	res, err2 := backendComm(req).(string)
	if err2 {
		panic(1)
	}
	if res == "success" {
		ctx.Redirect("/protected/", iris.StatusFound)
	} else {
		ctx.Redirect("/protected/Cart", iris.StatusFound)
	}
}

// new get active cart order (to be added)

func checkOrderStatus(ctx iris.Context) { // get order history

}

// for teller - only returns order IDs
func getOrders(ctx iris.Context) {
	req := structs.Request{FunctionName: "GetOrder"}
	backendComm(req)

	form := ctx.FormValues()
	locationID := form["locationID"][0]
	locationIDInt, _ := strconv.Atoi(locationID)
	ordersReq := structs.GetOrdersRequest{LocationID: locationIDInt}
	orders, ok := backendComm(ordersReq).([]structs.Order)
	if !ok {
		panic(1)
	}
	ctx.ViewData("isOrders", true)
	ctx.ViewData("Orders", orders)
	ctx.View("teller.html")
}

// get all details for a specific - returns food items - changes status - add new section for detailed food view
func selectOrder(ctx iris.Context) {

}

func completeOrder(ctx iris.Context) {

}

func createItem(req structs.CreateItemRequest) {

}

func updateItem(req structs.UpdateItemRequest) {

}

func deleteItem(req structs.DeleteItemRequest) {

}

// add form on home page
func sendMealSwipes(req structs.SendMealSwipesRequest) {

}

// add to homepage convert to string
// dollar amounts are in cents to avoid floating point
func getPaymentBalances(req structs.GetPaymentBalancesRequest) {

}
