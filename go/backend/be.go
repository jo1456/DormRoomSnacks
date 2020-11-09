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

	// "errors"

	"github.com/joho/godotenv"
)

var (
	connection net.Conn
	encoder    *json.Encoder
	decoder    *json.Decoder
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

	// creates a variable with the passed flag. default value of 8080
	listenPort := flag.String("listen", "8090", "port to listen on")
	flag.Parse()

	// _, err = db.Query("INSERT INTO Foods VALUES ( 'New food', 'wicked good', 10.01, TRUE )")
	// if err != nil {
	// 	panic(err.Error())
	// }

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

	// infinite loop to accept and response to requests
	for {
		// status of operation. In case of failure fe can redirect to homepage
		status := structs.Status{Success: true}

		var req structs.Request
		decoder.Decode(&req)

		switch req.FunctionName {
		case "ListLocations":
		case "GetMenu":
		case "ViewItem":
		case "SubmitOrder":
		case "CheckOrderStatus":
		case "GetOrders":
		case "SelectOrder":
		case "CompleteOrder":
		case "UpdateItem":
		case "CreateItem":
		case "DeleteItem":

		}
		encoder.Encode(status)
	}
}

func ListLocations()                                   {}
func GetMenu(structs.GetMenuRequest)                   {}
func ViewItem(structs.ViewItemRequest)                 {}
func SubmitOrder(structs.SubmitOrderRequest)           {}
func CheckOrderStatus(structs.CheckOrderStatusRequest) {}

func GetOrders(structs.GetOrdersRequest)         {}
func SelectOrder(structs.SelectOrderRequest)     {}
func CompleteOrder(structs.CompelteOrderRequest) {}
func UpdateItem(structs.UpdateItemRequest)       {}
func CreateItem(structs.CreateItemRequest)       {}
func DeleteItem(structs.DeleteItemRequest)       {}
