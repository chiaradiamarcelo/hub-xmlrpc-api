package client

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/kolo/xmlrpc"
)

type Client struct {
	connectTimeout, requestTimeout int
}

func NewClient(connectTimeout, requestTimeout int) *Client {
	return &Client{connectTimeout: connectTimeout, requestTimeout: requestTimeout}
}

func (c *Client) ExecuteCall(endpoint string, call string, args []interface{}) (response interface{}, err error) {
	client, err := getClientWithTimeout(endpoint, c.connectTimeout, c.requestTimeout)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	err = client.Call(call, args, &response)
	return response, err
}

func getClientWithTimeout(url string, connectTimeout, requestTimeout int) (*xmlrpc.Client, error) {
	transport := http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:  time.Duration(connectTimeout) * time.Second,
			Deadline: time.Now().Add(time.Duration(requestTimeout) * time.Second),
		}).DialContext,
	}
	return xmlrpc.NewClient(url, &transport)
}
