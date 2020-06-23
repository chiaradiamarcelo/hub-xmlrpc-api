package integration_tests

import (
	"log"
	"reflect"
	"testing"

	"github.com/uyuni-project/hub-xmlrpc-api/client"
	"github.com/uyuni-project/hub-xmlrpc-api/controller"
)

func Test_Multicast(t *testing.T) {
	tt := []struct {
		name                      string
		username                  string
		password                  string
		expectedMulticastResponse *controller.MulticastResponse
		expectedError             string
	}{
		{
			name:                      "multicast.system.listSystems",
			username:                  "admin",
			password:                  "admin",
			expectedMulticastResponse: &controller.MulticastResponse{controller.MulticastStateResponse{}, controller.MulticastStateResponse{}},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			//setup env
			InitHubInfrastructure()
			gatewayServerURL := "http://localhost:2830/hub/rpc/api"
			client := client.NewClient(10, 10)
			//login
			loginResponse, err := client.ExecuteCall(gatewayServerURL, "hub.loginWithAutoconnectMode", []interface{}{"admin", "admin"})
			if err != nil && tc.expectedError != err.Error() {
				t.Fatalf("Error during executing request: %v", err)
			}
			hubSessionKey := loginResponse.(map[string]interface{})["SessionKey"].(string)
			serverIDsSlice := loginResponse.(map[string]interface{})["Successful"].(map[string]interface{})["ServerIds"].([]interface{})

			loggedInServerIDs := make([]int64, 0, len(serverIDsSlice))
			for _, serverID := range serverIDsSlice {
				loggedInServerIDs = append(loggedInServerIDs, serverID.(int64))
			}
			//execute multicast call
			systemsPerServer, err := client.ExecuteCall(gatewayServerURL, "multicast.system.listSystems", []interface{}{hubSessionKey, loggedInServerIDs})
			log.Printf("POPOPO: %v", systemsPerServer)
			if err != nil && tc.expectedError != err.Error() {
				t.Fatalf("Error during executing request: %v", err)
			}
			if err == nil && !reflect.DeepEqual(systemsPerServer, tc.expectedMulticastResponse) {
				t.Fatalf("Expected and actual values don't match, Expected value is: %v", tc.expectedMulticastResponse)
			}
		})
	}
}
