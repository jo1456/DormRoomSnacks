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
	app.Get("/Menu", requiresLogin, getLocations) // students
	app.Post("/create-order", requiresLogin, rediGetMenu)
	app.Get("/menu/{menuID:int}", requiresLogin, getMenu)
	app.Post("/add-item-cart", requiresLogin, addItemOrder)
	app.Get("/Cart", requiresLogin, getCart) // how do i do this?
	app.Post("/checkout", requiresLogin, submitOrder)
	app.Post("/sendMealSwipe", requiresLogin, sendMealSwipes)

	// staff only
	app.Get("/Orders", requiresStaffLogin, getLocationsTeller) // teller
	app.Get("/get-orders", requiresStaffLogin, getOrders)
	app.Get("/complete-order", completeOrder)

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
}

func getHomePage(ctx iris.Context) {
	session := sess.Start(ctx)

	userID, _ := session.GetInt("userID")
	isStudent, _ := session.GetBoolean("isStudent")
	if userID == -1 {
		ctx.ViewData("ClientName", "Guest")
		ctx.ViewData("LoggedIn", false)
	} else {
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
	ctx.ViewData("IsLocSelec", true)
	ctx.ViewData("Locations", res.Locations)
	ctx.View("student.html")
}

func getLocationsTeller(ctx iris.Context) {
	backendComm("ListLocations", nil)
	var res structs.ListLocationsResponse
	err := decoder.Decode(&res)
	if err != nil {
		panic(err.Error())
	}
	ctx.ViewData("IsLocSelec", true)
	ctx.ViewData("Locations", res.Locations)
	ctx.View("teller.html")
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
	fmt.Println("menuID:", menuID)
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

	ctx.ViewData("isMenu", true)
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
	pws := form["mealSwipe"][0]
	pwsB := false
	if pws == "something" { // to be changed
		pwsB = true
	} else {
		pwsB = false
	}
	// set the ID to 0 because i think its the orderID and i don't have access to it here
	someItem := structs.OrderItem{ID: 0, FoodID: itemID, Customization: "none", PayWithSwipe: pwsB}
	addReq := structs.AddItemToOrderRequest{Item: someItem}
	backendComm("AddItemToOrder", addReq)
	var res string
	err = decoder.Decode(&res)
	if err != nil {
		panic(err.Error())
	}
	if res == "failure" {
		menuID, err := session.GetInt("menuID")
		if err != nil {
			panic(1)
		}
		redirectLink := fmt.Sprintf("%s%d", "/menu/", menuID)
		ctx.Redirect(redirectLink, iris.StatusFound)
	} else {

	}
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

	ctx.ViewData("IsCheckout", true)
	ctx.ViewData("CartItems", res.Items) // are order item for rn
	ctx.View("student.html")
}

func updateCartItem(ctx iris.Context) {
	// formData := ctx.FormValues()
}

func deleteCartItem(ctx iris.Context) {
	// formData := ctx.FormValues()

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
	if res == "success" {
		ctx.Redirect("/", iris.StatusFound)
	} else {
		ctx.Redirect("/Cart", iris.StatusFound)
	}
}

// for teller - only returns order IDs
func getOrders(ctx iris.Context) {
	form := ctx.FormValues()
	locationID := form["locationID"][0]
	locationIDInt, _ := strconv.Atoi(locationID)
	ordersReq := structs.GetOrdersRequest{LocationID: locationIDInt}
	backendComm("GetOrders", ordersReq)
	var res []structs.Order
	err := decoder.Decode(&res)
	if err != nil {
		panic(err.Error())
	}
	ctx.ViewData("isOrders", true)
	ctx.ViewData("Orders", res)
	ctx.View("teller.html")
}

// get all details for a specific - returns food items - changes status - add new section for detailed food view
func selectOrder(ctx iris.Context) {

}

func completeOrder(ctx iris.Context) {

}

func createItem(ctx iris.Context) {

}

func updateItem(ctx iris.Context) {

}

func deleteItem(ctx iris.Context) {

}

// add form on home page
func sendMealSwipes(ctx iris.Context) {
	sessions := sess.Start(ctx)
	userID, _ := sessions.GetInt("userID")

	formData := ctx.FormValues()
	toID, _ := strconv.Atoi(formData["toID"][0])
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
}

// add to homepage convert to string
// dollar amounts are in cents to avoid floating point
func getPaymentBalances(req structs.GetPaymentBalancesRequest) {

}
