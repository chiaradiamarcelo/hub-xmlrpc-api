package server_test

import (
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/uyuni-project/hub-xmlrpc-api/client"
	"github.com/uyuni-project/hub-xmlrpc-api/config"
	"github.com/uyuni-project/hub-xmlrpc-api/server"
	"github.com/uyuni-project/hub-xmlrpc-api/session"
)

func init() {
	/* load test data */
	os.Setenv("HUB_CONFIG_FILE", "../tests/config.json")
}

func TestLogin(t *testing.T) {

	tt := []struct {
		name     string
		username string
		password string
		err      string
	}{
		{name: "Invalid credentials", username: "unknown-user", password: "unknown-user", err: server.FaultInvalidCredentials.String},
		{name: "Valid credentials", username: "admin", password: "admin"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			conf := config.InitializeConfig()
			client := client.NewClient(conf.ConnectTimeout, conf.ReadWriteTimeout)
			session := session.NewSession(client, conf.Hub.SUMA_API_URL)
			hub := server.NewHubService(client, session, conf.Hub.SUMA_API_URL)

			req, err := http.NewRequest("GET", conf.Hub.SUMA_API_URL, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			reply := struct{ Data string }{""}
			err = hub.Login(req, &struct{ Username, Password string }{tc.username, tc.password}, &reply)
			if err != nil {
				if !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("Expected %v, Got %v", tc.err, err.Error())
				}
				return
			}
			// test the hubkey
			hubsessionkey := reply.Data
			if err != nil {
				if !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("Expected %v, Got %v", tc.err, err.Error())
				}
				return
			}
			// test the hubkey
			matched, _ := regexp.MatchString(`^[A-Za-z0-9]{68}$`, hubsessionkey)
			if !matched {
				t.Fatalf("Unexepected token pattern %v", hubsessionkey)
				return
			}
			username, password := session.GetUsernameAndPassword(hubsessionkey)
			if username != tc.username {
				t.Fatalf("User name doesn't match with the key, expected %v, got %v", tc.username, username)
			}
			if password != tc.password {
				t.Fatalf("User name doesn't match with the key, expected %v, got %v", tc.password, password)
			}
		})
	}

}
func TestLoginAutoconnect(t *testing.T) {
	tt := []struct {
		name     string
		username string
		password string
		err      string
	}{
		{name: "Valid credentials", username: "admin", password: "admin"},
		{name: "Invalid credentials", username: "unknown-user", password: "unknown-user", err: server.FaultInvalidCredentials.String},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			conf := config.InitializeConfig()
			client := client.NewClient(conf.ConnectTimeout, conf.ReadWriteTimeout)
			session := session.NewSession(client, conf.Hub.SUMA_API_URL)
			hub := server.NewHubService(client, session, conf.Hub.SUMA_API_URL)

			req, err := http.NewRequest("GET", conf.Hub.SUMA_API_URL, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			reply := struct{ Data string }{""}

			err = hub.LoginWithAutoconnectMode(req, &struct{ Username, Password string }{tc.username, tc.password}, &reply)
			if err != nil {
				if !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("Expected %v, Got %v", tc.err, err.Error())
				}
				return
			}
			// test the hubkey
			hubsessionkey := reply.Data
			matched, _ := regexp.MatchString(`^[A-Za-z0-9]{68}$`, hubsessionkey)
			if !matched {
				t.Fatalf("Unexepected token pattern %v", hubsessionkey)
				return
			}
			username, password := session.GetUsernameAndPassword(hubsessionkey)
			if username != tc.username {
				t.Fatalf("User name doesn't match with the key, expected %v, got %v", tc.username, username)
			}
			if password != tc.password {
				t.Fatalf("User name doesn't match with the key, expected %v, got %v", tc.password, password)
			}
			//test if servers attached to hub have also been authenticated automatically
			sessionKey := struct{ HubSessionKey string }{reply.Data}
			serverIdsreply := struct{ Data []int64 }{}
			err = hub.ListServerIds(req, &sessionKey, &serverIdsreply)
			serverIds := serverIdsreply.Data
			for _, s := range serverIds {
				url, severSessionkey := session.GetServerSessionInfoByServerID(sessionKey.HubSessionKey, s)
				if len(url) == 0 {
					t.Fatalf("Expected valid url for server with severId: %v, got empty instead %v", s, url)
				}
				if len(severSessionkey) <= 0 {
					t.Fatalf("Expected valid SessionKey for server with severId: %v, Got %v", s, severSessionkey)
				}
			}

		})
	}
}
func TestLoginWithAuthRelayMode(t *testing.T) {
	tt := []struct {
		name     string
		username string
		password string
		err      string
	}{
		{name: "Valid credentials", username: "admin", password: "admin"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			conf := config.InitializeConfig()
			client := client.NewClient(conf.ConnectTimeout, conf.ReadWriteTimeout)
			session := session.NewSession(client, conf.Hub.SUMA_API_URL)
			hub := server.NewHubService(client, session, conf.Hub.SUMA_API_URL)

			req, err := http.NewRequest("GET", conf.Hub.SUMA_API_URL, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			reply := struct{ Data string }{""}
			err = hub.LoginWithAuthRelayMode(req, &struct{ Username, Password string }{tc.username, tc.password}, &reply)
			if err != nil {
				if !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("Expected %v, Got %v", tc.err, err.Error())
				}
				return
			}
			// test the hubkey
			hubsessionkey := reply.Data
			matched, _ := regexp.MatchString(`^[A-Za-z0-9]{68}$`, hubsessionkey)
			if !matched {
				t.Fatalf("Unexepected token pattern %v", hubsessionkey)
				return
			}
			username, password := session.GetUsernameAndPassword(hubsessionkey)
			if username != tc.username {
				t.Fatalf("User name doesn't match with the key, expected %v, got %v", tc.username, username)
			}
			if password != tc.password {
				t.Fatalf("User name doesn't match with the key, expected %v, got %v", tc.password, password)
			}
		})
	}

}

func TestAttachToServers(t *testing.T) {
	tt := []struct {
		name     string
		username string
		password string
		err      string
	}{
		{name: "Valid credentials", username: "admin", password: "admin"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			conf := config.InitializeConfig()
			client := client.NewClient(conf.ConnectTimeout, conf.ReadWriteTimeout)
			session := session.NewSession(client, conf.Hub.SUMA_API_URL)
			hub := server.NewHubService(client, session, conf.Hub.SUMA_API_URL)

			req, err := http.NewRequest("GET", conf.Hub.SUMA_API_URL, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			reply := struct{ Data string }{""}
			//login
			err = hub.LoginWithAuthRelayMode(req, &struct{ Username, Password string }{tc.username, tc.password}, &reply)
			if err != nil {
				t.Fatalf("Login faied with error : %v", err)
			}
			sessionKey := struct{ HubSessionKey string }{reply.Data}

			// List server Ids
			serverIdsreply := struct{ Data []int64 }{}
			err = hub.ListServerIds(req, &sessionKey, &serverIdsreply)
			serverIds := serverIdsreply.Data

			srvArgs := server.MulticastArgs{sessionKey.HubSessionKey, serverIds, nil}
			err = hub.AttachToServers(req, &srvArgs, &struct{ Data []error }{})
			if err != nil && err.Error() != tc.err {
				t.Fatalf("Unexpected Result: Exepected %v, Got %v", tc.err, err.Error())
				return
			}
			for _, s := range serverIds {
				url, severSessionkey := session.GetServerSessionInfoByServerID(sessionKey.HubSessionKey, s)
				if len(url) == 0 {
					t.Fatalf("Expected valid url for server with severId: %v, got empty instead %v", s, url)
				}
				if len(severSessionkey) <= 0 {
					t.Fatalf("Expected valid SessionKey for server with severId: %v, Got %v", s, severSessionkey)
				}
			}
		})
	}
}

func TestIsHubSessionValid(t *testing.T) {
	const errorMessage = "is not valid"
	tt := []struct {
		name     string
		username string
		password string
		result   bool
	}{
		{name: "Valid credentials", username: "admin", password: "admin", result: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			conf := config.InitializeConfig()
			client := client.NewClient(conf.ConnectTimeout, conf.ReadWriteTimeout)
			session := session.NewSession(client, conf.Hub.SUMA_API_URL)
			hub := server.NewHubService(client, session, conf.Hub.SUMA_API_URL)

			req, err := http.NewRequest("GET", conf.Hub.SUMA_API_URL, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			reply := struct{ Data string }{""}

			err = hub.Login(req, &struct{ Username, Password string }{tc.username, tc.password}, &reply)
			if err != nil {
				t.Fatalf("Couldn't login with provided credentials")
				return
			}
			//Test if key is valid
			isvalid := session.IsHubSessionValid(reply.Data)
			if isvalid != tc.result {
				t.Fatalf("Unexpected Result: Exepected %v, Got %v", tc.result, isvalid)
			}
			//Append the key with some random string and test if it's invalid now
			isvalid = session.IsHubSessionValid(reply.Data + "invalid-part")
			if isvalid != false {
				t.Fatalf("Unexpected Result: Exepected %v, Got %v", tc.result, isvalid)
			}
		})
	}
}

func TestListServerIds(t *testing.T) {

	tt := []struct {
		name     string
		username string
		password string
		err      string
	}{
		{name: "Valid credentials", username: "admin", password: "admin"},
		{name: "With invalid  credentials", username: "unknownadmin", password: "unknownadmin", err: server.FaultInvalidCredentials.String},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			conf := config.InitializeConfig()
			client := client.NewClient(conf.ConnectTimeout, conf.ReadWriteTimeout)
			session := session.NewSession(client, conf.Hub.SUMA_API_URL)
			hub := server.NewHubService(client, session, conf.Hub.SUMA_API_URL)

			req, err := http.NewRequest("GET", conf.Hub.SUMA_API_URL, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			reply := struct{ Data string }{""}
			err = hub.Login(req, &struct{ Username, Password string }{tc.username, tc.password}, &reply)
			if err != nil {
				if !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("Expected %v, Got %v", tc.err, err.Error())
				}
				return
			}
			serverIdsreply := struct{ Data []int64 }{}
			err = hub.ListServerIds(req, &struct{ HubSessionKey string }{reply.Data}, &serverIdsreply)
			if err != nil && err.Error() != tc.err {
				t.Fatalf("Unexpected Result: Exepected %v, Got %v", tc.err, err.Error())
				return
			}
			serverIds := len(serverIdsreply.Data)
			if serverIds <= 0 {
				t.Fatalf("Unexpected Result: Expected some servers, got nothing")
			}
		})
	}

}