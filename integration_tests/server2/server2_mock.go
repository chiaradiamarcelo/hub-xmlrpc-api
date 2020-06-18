package main

import (
	"log"
	"net/http"

	"github.com/gorilla/rpc"
	"github.com/uyuni-project/hub-xmlrpc-api/controller/parser"
	"github.com/uyuni-project/hub-xmlrpc-api/controller/xmlrpc"
)

type SystemInfo struct {
	Id   int64  `xmlrpc:"id"`
	Name string `xmlrpc:"name"`
}
type System struct{}
type Auth struct{}

var System_1 = SystemInfo{
	Id:   1000010004,
	Name: "server2-minion-1",
}
var System_2 = SystemInfo{
	Id:   1000010005,
	Name: "server2-minion-2",
}
var Systems = []SystemInfo{
	System_1,
	System_2,
}

var sessionkey = "300x2413800c14c02928568674dad9e71e0f061e2920be1d7c6542683d78de524bd6"

func (h *Auth) Login(r *http.Request, args *struct{ Username, Password string }, reply *struct{ Data string }) error {
	log.Println("Server2 -> auth.login", args.Username)
	reply.Data = sessionkey
	return nil
}
func (h *System) ListSystems(r *http.Request, args *struct{ SessionK string }, reply *struct{ Data []SystemInfo }) error {
	log.Println("Server2 -> System.ListSystems", args.SessionK)
	if args.SessionK == sessionkey {
		reply.Data = Systems
	}
	return nil
}
func (h *System) ListUserSystems(r *http.Request, args *struct{ SessionK, UserLogin string }, reply *struct{ Data []SystemInfo }) error {
	log.Println("Server2 -> System.ListUserSystems", args.SessionK)
	if args.SessionK == sessionkey && args.UserLogin == "admin" {
		reply.Data = Systems
	}
	return nil
}

func main() {
	rpcServer := rpc.NewServer()
	var codec = xmlrpc.NewCodec()

	codec.RegisterMapping("auth.login", "Auth.Login", parser.LoginRequestParser)
	codec.RegisterMapping("system.listSystems", "System.ListSystems", parser.LoginRequestParser)
	codec.RegisterMapping("system.listUserSystems", "System.ListUserSystems", parser.LoginRequestParser)

	rpcServer.RegisterCodec(codec, "text/xml")
	rpcServer.RegisterService(new(Auth), "")
	rpcServer.RegisterService(new(System), "")

	http.Handle("/rpc/api", rpcServer)
	log.Println("Starting XML-RPC server on localhost:8003/rpc/api")
	log.Fatal(http.ListenAndServe(":8003", nil))
}
