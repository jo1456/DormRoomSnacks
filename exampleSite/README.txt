Project phase 1

Name: Jeffrey Oberg jo1456

Running Instruction: Program can be run by navigating to the directory the be.go and fe.go files are in located and
running the command "go run be.go" to start the backend server and then "go run fe.go" to start the front end server.
This will make the application accessible on localhost:8080 with the back and front end servers communicating on
port 8090. If you would like to specify the port for the backend to run on, use the "--listen" flag followed by the port number.
If you would like to specify the port for the front end to listen on use the "--listen" flag followed by the port number.
If you would like to specify the ip and port of the backend that the front end should be connecting to, use the "--backend" flag
followed by the ip and port number in the format ip:port.

The github.com/kataras/iris/v12 v12.2.0-alpha library will also need to have been downloaded with the go get command
and then a go.mod file is also required to ensure the proper version of iris is being used (github.com/kataras/iris/v12 v12.2.0-alpha).
Lastly, this version of iris requires go 1.14 or above, so that is also required.

State of assignment: The assignment is complete. No known bugs.

Other resources: Iris documentation was used. "net" and "encoding/json" documentation
also used. All documentation found on https://godoc.org/

I have used gRPC for network communication in the past so I modeled my
approach off of my experience with that. I created structs for sending data as if
I were creating the proto files needed in grpc. The disadvantage was that I did not
know a great way to tell the backend which function I wanted to call so I would just
send a message with the function name to be used in a switch statement. I think this
ended up being pretty readable but I did not love the solution.

I was not totally sure whether or not the backend was supposed to be able to
accept connections from multiple instances of the front end until I asked in class.
Besides that, instructions were very clear. 
