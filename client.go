package samsungremote

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"

	"golang.org/x/net/websocket"
)

// NewClient creates a new Client for samsung tv
//
// host should be the address used to connect to the tv - EG `192.168.1.6:8002`
// applicationName should the the human readable string that represents your
// application (will be displayed on the TV)
func NewClient(
	host string,
	applicationName string,
	options ...Option,
) (*Client, error) {
	c := &Client{
		host:            host,
		applicationName: applicationName,
	}

	for _, o := range options {
		err := o(c)
		if err != nil {
			return nil, fmt.Errorf("client.NewClient failed to apply option: %s", err)
		}
	}

	return c, nil
}

type Client struct {
	host                 string
	port                 int
	applicationName      string
	token                string
	websocketConfigsPass []websocketConfigPass

	websocketConnection *websocket.Conn
	connectionMx        sync.Mutex
}

// Send forwards the provided Payload to the device the client is connected to
func (c *Client) Send(command Payload) error {
	b, err := json.Marshal(command)
	if err != nil {
		return err
	}

	ws, err := c.getConnection()
	if err != nil {
		return err
	}

	return websocket.Message.Send(ws, string(b))
}

func (c *Client) AcquireToken() (string, error) {
	if c.token != "" {
		return c.token, nil
	}

	ws, err := c.getConnection()
	if err != nil {
		return "", err
	}

	b := make([]byte, 800)
	read, err := ws.Read(b)
	if err != nil {
		return "", fmt.Errorf("client.AcquireToken failed to read from token: %s", err)
	}

	received := struct {
		Data struct {
			Clients []struct {
				Attributes struct {
					Name string `json:"name"`
				} `json:"attributes"`
				ConnectTime int64  `json:"connectTime"`
				DeviceName  string `json:"deviceName"`
				ID          string `json:"id"`
				IsHost      bool   `json:"isHost"`
			} `json:"clients"`
			ID    string `json:"id"`
			Token string `json:"token"`
		} `json:"data"`
		Event string `json:"event"`
	}{}

	err = json.Unmarshal(b[:read], &received)
	if err != nil {
		return "", fmt.Errorf("client.AcquireToken failed to unmarshal message: %s - message: %s", err, b[:read])
	}

	c.token = received.Data.Token
	c.reinitializeConnection()

	return c.token, nil
}

func (c *Client) reinitializeConnection() {
	c.connectionMx.Lock()
	defer c.connectionMx.Unlock()
	c.websocketConnection = nil
}

func (c *Client) getConnection() (*websocket.Conn, error) {
	c.connectionMx.Lock()
	defer c.connectionMx.Unlock()

	if c.websocketConnection != nil {
		return c.websocketConnection, nil
	}

	wc, err := c.buildWebsocketConfig()
	if err != nil {
		return nil, err
	}

	ws, err := websocket.DialConfig(wc)
	if err != nil {
		return nil, fmt.Errorf("client.getConnection failed to dial: %s", err)
	}

	c.websocketConnection = ws

	return ws, err
}

func (c *Client) buildWebsocketConfig() (*websocket.Config, error) {
	tvURL, err := configToUrl(c.host, c.applicationName, c.token)
	if err != nil {
		return nil, fmt.Errorf("client.buildWebsocketConfig: invalid config for client: %s", err.Error())
	}

	config, err := websocket.NewConfig(tvURL.String(), "http://localhost/")
	if err != nil {
		return nil, fmt.Errorf("client.buildWebsocketConfig: failed to create websocket config: %s", err)
	}

	for _, o := range c.websocketConfigsPass {
		err := o(config)
		if err != nil {
			return nil, fmt.Errorf("client.buildWebsocketConfig failed to apply option: %s", err)
		}
	}

	return config, nil
}

func configToUrl(host string, applicationName, token string) (*url.URL, error) {
	u, err := url.Parse(
		fmt.Sprintf(
			"wss://%s/api/v2/channels/samsung.remote.control",
			host,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("client.configToUrl: unable to convert config to url: %s", err.Error())
	}

	q := u.Query()
	q.Set(
		"name",
		base64.StdEncoding.EncodeToString([]byte(applicationName)),
	)

	if len(token) > 0 {
		q.Set("token", token)
	}

	u.RawQuery = q.Encode()

	return u, nil
}
