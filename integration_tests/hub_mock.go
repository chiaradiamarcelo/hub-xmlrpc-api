package integration_tests

import (
	"log"
	"net/http"
	"strconv"

	"github.com/uyuni-project/hub-xmlrpc-api/controller"
)

func InitHubInfrastructure() {
	var system_1 = SystemInfo{
		Id:   1000010000,
		Name: "server-1",
	}
	var system_2 = SystemInfo{
		Id:   1000010001,
		Name: "server-2",
	}
	var systems = []SystemInfo{
		system_1,
		system_2,
	}

	hub := new(UyuniServer)
	sessionKey := "300x2413800c14c02928568674dad9e71e0f061e2920be1d7c6542683d78de524bd4"
	hub.mockLogin = func(r *http.Request, args *struct{ Username, Password string }, reply *struct{ Data string }) error {
		log.Println("Hub -> auth.login", args.Username)
		if args.Username == "admin" && args.Password == "admin" {
			reply.Data = sessionKey
		} else {
			return controller.FaultInvalidCredentials
		}
		return nil
	}
	hub.mockListSystems = func(r *http.Request, args *struct{ SessionKey string }, reply *struct{ Data []SystemInfo }) error {
		log.Println("Hub -> System.ListSystems", args.SessionKey)
		if args.SessionKey == sessionKey {
			reply.Data = systems
		}
		return nil
	}
	hub.mockListUserSystems = func(r *http.Request, args *struct{ SessionKey, Username string }, reply *struct{ Data []SystemInfo }) error {
		log.Println("Hub -> System.ListUserSystems", args.Username)
		if args.SessionKey == sessionKey && args.Username == "admin" {
			reply.Data = systems
		}
		return nil
	}
	hub.mockListFqdns = func(r *http.Request, args *struct {
		SessionKey string
		ServerId   int64
	}, reply *struct{ Data []string }) error {
		log.Println("Hub -> System.ListFqdns", args.ServerId)
		if args.SessionKey == sessionKey {
			if args.ServerId == 1000010000 {
				reply.Data = []string{"localhost:8002"}
			} else {
				reply.Data = []string{"localhost:8003"}
			}
		}
		return nil
	}
	InitServer(8001, hub)
	InitPeripheralServers([]int64{8002, 8003}, systems)
}

func InitPeripheralServers(ports []int64, systems []SystemInfo) {
	for i, systemInfo := range systems {
		var System_1 = SystemInfo{
			Id:   systemInfo.Id + 1,
			Name: systemInfo.Name + "-minion-1",
		}
		var System_2 = SystemInfo{
			Id:   systemInfo.Id + 2,
			Name: systemInfo.Name + "-minion-2",
		}
		var Systems = []SystemInfo{
			System_1,
			System_2,
		}

		serverNumber := strconv.Itoa(i)
		server := new(UyuniServer)
		sessionKey := "300x2413800c14c02928568674dad9e71e0f061e2920be1d7c6542683d78de524bd" + serverNumber
		server.mockLogin = func(r *http.Request, args *struct{ Username, Password string }, reply *struct{ Data string }) error {
			log.Println("Server"+serverNumber+" -> auth.login", args.Username)
			reply.Data = sessionKey
			return nil
		}
		server.mockListSystems = func(r *http.Request, args *struct{ SessionKey string }, reply *struct{ Data []SystemInfo }) error {
			log.Println("Server"+serverNumber+" -> System.ListSystems", args.SessionKey)
			if args.SessionKey == sessionKey {
				reply.Data = Systems
			}
			return nil
		}
		server.mockListUserSystems = func(r *http.Request, args *struct{ SessionKey, Username string }, reply *struct{ Data []SystemInfo }) error {
			log.Println("Server"+serverNumber+" -> System.ListUserSystems", args.Username)
			if args.SessionKey == sessionKey && args.Username == "admin" {
				reply.Data = Systems
			}
			return nil
		}
		InitServer(ports[i], server)
	}
}
