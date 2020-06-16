package main

import (
	"log"
	"net/http"

	"github.com/gorilla/rpc"
	"github.com/uyuni-project/hub-xmlrpc-api/controller"
	"github.com/uyuni-project/hub-xmlrpc-api/controller/parser"
	"github.com/uyuni-project/hub-xmlrpc-api/controller/xmlrpc"
)

type SystemInfo struct {
	Id   int64  `xmlrpc:"id"`
	Name string `xmlrpc:"name"`
}

type UyuniServer struct{}

var System_1 = SystemInfo{
	Id:   1000010000,
	Name: "server-1",
}
var System_2 = SystemInfo{
	Id:   1000010001,
	Name: "server-2",
}
var Systems = []SystemInfo{
	System_1,
	System_2,
}
var sessionkey = "300x2413800c14c02928568674dad9e71e0f061e2920be1d7c6542683d78de524bd4"

func (h *UyuniServer) Login(r *http.Request, args *struct{ Username, Password string }, reply *struct{ Data string }) error {
	log.Println("Hub -> auth.login", args.Username)
	if args.Username == "admin" && args.Password == "admin" {
		reply.Data = sessionkey
	} else {
		return controller.FaultInvalidCredentials
	}
	return nil
}

func (h *UyuniServer) ListUserSystems(r *http.Request, args *struct{ SessionKey, UserLogin string }, reply *struct{ Data []SystemInfo }) error {
	log.Println("Hub -> System.ListUserSystems", args.UserLogin)
	if args.SessionKey == sessionkey && args.UserLogin == "admin" {
		reply.Data = Systems
	}
	return nil
}

func (h *UyuniServer) ListSystems(r *http.Request, args *struct{ SessionKey string }, reply *struct{ Data []SystemInfo }) error {
	log.Println("Hub -> System.ListSystems", args.SessionKey)
	if args.SessionKey == sessionkey {
		reply.Data = Systems
	}
	return nil
}

func (h *UyuniServer) ListFqdns(r *http.Request, args *struct {
	Hubkey   string
	ServerId int64
}, reply *struct{ Data []string }) error {
	log.Println("Hub -> System.ListFqdns", args.ServerId)
	if args.Hubkey == sessionkey {
		if args.ServerId == 1000010000 {
			reply.Data = []string{"localhost:8002"}
		} else {
			reply.Data = []string{"localhost:8003"}
		}
	}
	return nil
}

func main() {
	rpcServer := rpc.NewServer()
	var codec = xmlrpc.NewCodec()

	codec.RegisterMapping("auth.login", "UyuniServer.Login", parser.LoginRequestParser)
	codec.RegisterMapping("system.listSystems", "UyuniServer.ListSystems", parser.LoginRequestParser)
	codec.RegisterMapping("system.listUserSystems", "UyuniServer.ListUserSystems", parser.LoginRequestParser)
	codec.RegisterMapping("system.listFqdns", "UyuniServer.ListFqdns", parser.LoginRequestParser)

	rpcServer.RegisterCodec(codec, "text/xml")
	rpcServer.RegisterService(new(UyuniServer), "")

	http.Handle("/rpc/api", rpcServer)
	log.Println("Starting XML-RPC server on localhost:8001/hub/rpc/api")
	log.Fatal(http.ListenAndServe(":8001", nil))
}
