package integration_tests

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/rpc"
	"github.com/uyuni-project/hub-xmlrpc-api/controller/parser"
	"github.com/uyuni-project/hub-xmlrpc-api/controller/xmlrpc"
)

type UyuniServer struct {
	username, password, sessionKey string
	mockLogin                      func(r *http.Request, args *struct{ Username, Password string }, reply *struct{ Data string }) error
	mockListUserSystems            func(r *http.Request, args *struct{ SessionKey, Username string }, reply *struct{ Data []SystemInfo }) error
	mockListSystems                func(r *http.Request, args *struct{ SessionKey string }, reply *struct{ Data []SystemInfo }) error
	mockListFqdns                  func(r *http.Request, args *struct {
		SessionKey string
		ServerId   int64
	}, reply *struct{ Data []string }) error
}

type SystemInfo struct {
	Id   int64  `xmlrpc:"id"`
	Name string `xmlrpc:"name"`
}

func (h *UyuniServer) Login(r *http.Request, args *struct{ Username, Password string }, reply *struct{ Data string }) error {
	return h.mockLogin(r, args, reply)
}

func (h *UyuniServer) ListUserSystems(r *http.Request, args *struct{ SessionKey, Username string }, reply *struct{ Data []SystemInfo }) error {
	return h.mockListUserSystems(r, args, reply)
}

func (h *UyuniServer) ListSystems(r *http.Request, args *struct{ SessionKey string }, reply *struct{ Data []SystemInfo }) error {
	return h.mockListSystems(r, args, reply)
}

func (h *UyuniServer) ListFqdns(r *http.Request, args *struct {
	SessionKey string
	ServerId   int64
}, reply *struct{ Data []string }) error {
	return h.mockListFqdns(r, args, reply)
}

func InitServer(port int64, uyuniServer *UyuniServer) {
	go func() {
		rpcServer := rpc.NewServer()
		var codec = xmlrpc.NewCodec()

		codec.RegisterMapping("auth.login", "UyuniServer.Login", parser.LoginRequestParser)
		codec.RegisterMapping("system.listSystems", "UyuniServer.ListSystems", parser.LoginRequestParser)
		codec.RegisterMapping("system.listUserSystems", "UyuniServer.ListUserSystems", parser.LoginRequestParser)
		codec.RegisterMapping("system.listFqdns", "UyuniServer.ListFqdns", parser.LoginRequestParser)

		rpcServer.RegisterCodec(codec, "text/xml")
		rpcServer.RegisterService(uyuniServer, "")

		mux := http.NewServeMux()
		mux.HandleFunc("/rpc/api", func(w http.ResponseWriter, r *http.Request) { rpcServer.ServeHTTP(w, r) })

		log.Printf("Starting XML-RPC server on localhost:%v/rpc/api", port)

		server := http.Server{
			Addr:    fmt.Sprintf(":%v", port),
			Handler: mux,
		}
		log.Fatal(server.ListenAndServe())
	}()
}
