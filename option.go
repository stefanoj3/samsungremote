package samsungremote

import (
	"crypto/tls"
	"net"
	"time"

	"golang.org/x/net/websocket"
)

type websocketConfigPass func(*websocket.Config) error

type Option func(*Client) error

func OptionAllowInsecureTLS(c *Client) error {
	pass := func(c *websocket.Config) error {
		if c.TlsConfig == nil {
			c.TlsConfig = &tls.Config{}
		}

		c.TlsConfig.InsecureSkipVerify = true

		return nil
	}

	c.websocketConfigsPass = append(c.websocketConfigsPass, pass)

	return nil
}

func OptionTimeout(t time.Duration) func(*Client) error {
	return func(c *Client) error {
		pass := func(c *websocket.Config) error {
			if c.Dialer == nil {
				c.Dialer = &net.Dialer{}
			}
			c.Dialer.Timeout = t

			return nil
		}

		c.websocketConfigsPass = append(c.websocketConfigsPass, pass)

		return nil
	}
}

func OptionTokenProvider(token string) func(*Client) error {
	return func(c *Client) error {
		c.token = token
		return nil
	}
}
