//Jeffrey Oberg jo1456
//Project part 2 CRUD APP with front and back end servers
// Front end

package main

import (
	"dormroomsnacks/structs"
	"strconv"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"

	"encoding/json"
	"flag"
	"fmt"
	"net"
)

var (
	connection net.Conn
	encoder    *json.Encoder
	decoder    *json.Decoder
	app        *iris.Application
	secret     = []byte("signature_hmac_secret_shared_key")
)

type authClaims struct {
	UserID   string `json:"userID"`
	Username string `json:"username"`
}

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

	signer := jwt.NewSigner(jwt.HS256, secret, 10*time.Minute)
	app.Get("/", generateToken(signer))

	verifier := jwt.NewVerifier(jwt.HS256, secret)
	verifier.WithDefaultBlocklist()
	verifyMiddleware := verifier.Verify(func() interface{} {
		return new(authClaims)
	})

	tmpl := iris.HTML("./views", ".html")
	tmpl.Delims("{{", "}}")
	app.RegisterView(tmpl)

	app.Get("/", getHomePage)
	app.Get("/Login", getLoginPage)
	app.Get("/Signup", getSignupPage)
	app.Get("/loginform", loginUser)
	app.Post("/signupform", registerUser)

	protectedAPI := app.Party("/protected")
	protectedAPI.Use(verifyMiddleware)

	protectedAPI.Get("/", protected)
	protectedAPI.Get("/logout", logout)
	protectedAPI.Get("/Menu", getLocations)
	protectedAPI.Get("/get-menu-student", getMenu)
	protectedAPI.Post("/add-item-cart", addItemCart)
	// protectedAPI.Get("/Cart", getCart) // how do i do this?
	protectedAPI.Get("/Orders", getLocationsTeller)
	protectedAPI.Get("/get-orders", getOrders)
	protectedAPI.Post("/checkout")

	// turn on the app
	app.Listen(":" + *listenPort)
}

func generateToken(signer *jwt.Signer) iris.Handler {
	return func(ctx iris.Context) {
		claims := authClaims{UserID: "bar"}

		token, err := signer.Sign(claims)
		if err != nil {
			ctx.StopWithStatus(iris.StatusInternalServerError)
			return
		}

		ctx.Write(token)
	}
}

func protected(ctx iris.Context) {
	claims := jwt.Get(ctx).(*authClaims)

	standardClaims := jwt.GetVerifiedToken(ctx).StandardClaims
	expiresAtString := standardClaims.ExpiresAt().
		Format(ctx.Application().ConfigurationReadOnly().GetTimeFormat())
	timeLeft := standardClaims.Timeleft()

	ctx.Writef("foo=%s\nexpires at: %s\ntime left: %s\n", claims.UserID, expiresAtString, timeLeft)
}

func logout(ctx iris.Context) {
	err := ctx.Logout()
	if err != nil {
		ctx.WriteString(err.Error())
	} else {
		ctx.Writef("token invalidated, a new token is required to access the protected API")
	}
}

func getHomePage(ctx iris.Context) {
	ctx.ViewData("ClientName", "Guest") // to be changed later, for dynamic response
	ctx.ViewData("LoggedIn", false)
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

func loginUser(ctx iris.Context) {
	return
}

func registerUser(ctx iris.Context) {
	return
}

func getLocations(ctx iris.Context) {
	req := structs.Request{FunctionName: "ListLocations"}
	locations, ok := backendComm(req).(structs.ListLocationsResponse)
	if !ok {
		panic(1)
	}
	ctx.ViewData("Locations", locations.Locations)
	ctx.View("student.html")
}

func getLocationsTeller(ctx iris.Context) {
	req := structs.Request{FunctionName: "ListLocations"}
	locations, ok := backendComm(req).(structs.ListLocationsResponse)
	if !ok {
		panic(1)
	}
	ctx.ViewData("Locations", locations.Locations)
	ctx.View("teller.html")
}

func getMenu(ctx iris.Context) {
	req := structs.Request{FunctionName: "GetMenu"}
	backendComm(req)

	form := ctx.FormValues()
	menuID := form["menuID"][0]
	menuIDInt, _ := strconv.Atoi(menuID)
	menuReq := structs.GetMenuRequest{MenuID: menuIDInt}
	menu, ok := backendComm(menuReq).(structs.Menu)
	if !ok {
		panic(1)
	}
	ctx.ViewData("isMenu", true)
	ctx.ViewData("Menu", menu)
	ctx.View("student.html")
}

func addItemCart(ctx iris.Context) {
	req := structs.Request{FunctionName: "AddItemToOrder"}
	backendComm(req)

	form := ctx.FormValues()
	itemID, err := strconv.Atoi(form["itemID"][0])
	if err != nil {
		panic(1)
	}
	someItem := structs.ItemOrder{ItemID: itemID} // why do we need this? lets just change directly to itemID
	addReq := structs.AddItemToOrderRequest{OrderID: 0, Item: someItem, PayWithMealSwipe: false}
	backendComm(addReq)
}

func submitOrder(ctx iris.Context) {

}

// can we add items to cart and when checkout, we can package them into an order

// why do we need this
func checkOrderStatus(ctx iris.Context) {

}

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

// idk what this is
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

func sendMealSwipes(req structs.SendMealSwipesRequest) {

}

// dollar amounts are in cents to avoid floating point
func getPaymentBalances(req structs.GetPaymentBalancesRequest) {

}
