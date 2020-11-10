//Jeffrey Oberg jo1456
//Project part 2 CRUD APP with front and back end servers
// Front end

package main

import (
	"dormroomsnacks/structs"

	"github.com/kataras/iris/v12"

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
)

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
	app.Get("/Menu", menuMessageSetup, backendComm, getMenu)
	app.Get("/Cart", getCart)
	app.Post("/Checkout", checkoutCart)

	// turn on the app
	app.Listen(":" + *listenPort)
}

func backendComm(ctx iris.Context) {
	command, _ := ctx.Values().Get("BackendCommand").(string)
	err := encoder.Encode(&command)
	if err != nil {
		panic(err.Error())
	}
	var response interface{}
	err = decoder.Decode(&response)
	if err != nil {
		panic(err.Error())
	}
	ctx.Values().Set(command, response)
	ctx.Next()
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

func menuMessageSetup(ctx iris.Context) {
	ctx.Values().Set("BackendCommand", "GetMenu")
	ctx.Next()
}

func getMenu(ctx iris.Context) {
	menu, _ := ctx.Values().Get("GetMenu").(structs.Menu)
	ctx.ViewData("MenuItems", menu.Items)
}

func getCart(ctx iris.Context) {

}

func checkoutCart(ctx iris.Context) {

}
