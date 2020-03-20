package client

import (
	"net"
	"net/http"
	"time"

	"github.com/kolo/xmlrpc"
	"github.com/uyuni-project/hub-xmlrpc-api/config"
)

const defaultConnectTimeout, defaultReadWriteTimeout = 1, 1

func timeoutDialer(connectTimeout, readWriteTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, connectTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(readWriteTimeout))
		return conn, nil
	}
}

func getClientWithTimeout(url string, connectTimeout, readWriteTimeout int) (*xmlrpc.Client, error) {
	transport := http.Transport{
		Dial: timeoutDialer(time.Duration(connectTimeout)*time.Second, time.Duration(readWriteTimeout)*time.Second),
	}
	return xmlrpc.NewClient(url, &transport)
}

func getDefaultTimeoutClient(url string) (*xmlrpc.Client, error) {
	return getClientWithTimeout(url, defaultConnectTimeout, defaultReadWriteTimeout)
}

type Client struct {
	Conf config.Config
}

func (c *Client) ExecuteXMLRPCCallWithURL(url string, method string, args []interface{}) (reply interface{}, err error) {
	client, err := getClientWithTimeout(url, c.Conf.ConnectTimeout, c.Conf.ReadWriteTimeout)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	err = client.Call(method, args, &reply)
	return reply, err
}

func (c *Client) ExecuteXMLRPCCallToHub(method string, args []interface{}) (reply interface{}, err error) {
	return c.ExecuteXMLRPCCallWithURL(c.Conf.Hub.SUMA_API_URL, method, args)
}
