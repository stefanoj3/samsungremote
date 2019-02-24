package samsungremote_test

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stefanoj3/samsungremote"
	"golang.org/x/net/websocket"
)

func TestWillNotTryToAcquireANewTokenIfProvided(t *testing.T) {
	token := "my_awesome_token"

	c, err := samsungremote.NewClient(
		"somehost:8002",
		"MyApp",
		samsungremote.OptionAllowInsecureTLS,
		samsungremote.OptionTimeout(time.Millisecond*500),
		samsungremote.OptionTokenProvider(token),
	)
	if err != nil {
		t.Fatal(err.Error())
	}

	actualToken, err := c.AcquireToken()
	if err != nil {
		t.Fatal(err.Error())
	}

	if token != actualToken {
		t.Fatalf("unexpected value for token: expected %s got %s instead", token, actualToken)
	}
}

func TestAcquireToken(t *testing.T) {
	webSocketServer := websocket.Server{
		Handler: func(con *websocket.Conn) {
			con.Write([]byte(`{"data":{"clients":[{"attributes":{"name":"TXlBcHA="},"connectTime":1550433417482,"deviceName":"TXlBcHA=","id":"bfa15a3b-4163-4ca6-a4b3-f1f34bf3d415","isHost":false}],"id":"bfa15a3b-4163-4ca6-a4b3-f1f34bf3d415","token":"123123123"},"event":"ms.channel.connect"}`))
		},
	}

	srv := httptest.NewTLSServer(webSocketServer)
	defer srv.Close()

	c, err := samsungremote.NewClient(
		srv.Listener.Addr().String(),
		"MyApp",
		samsungremote.OptionAllowInsecureTLS,
		samsungremote.OptionTimeout(time.Millisecond*500),
	)
	if err != nil {
		t.Fatal(err.Error())
	}

	expectedToken := "123123123"
	actualToken, err := c.AcquireToken()
	if err != nil {
		t.Fatal(err.Error())
	}

	if expectedToken != actualToken {
		t.Fatalf("unexpected value for token: expected %s got %s instead", expectedToken, actualToken)
	}
}

func TestCanSendCommand(t *testing.T) {
	msgChannel := make(chan string, 1)
	defer close(msgChannel)

	webSocketServer := websocket.Server{
		Handler: func(con *websocket.Conn) {
			b := make([]byte, 800)
			readBytes, _ := con.Read(b)
			msgChannel <- string(b[:readBytes])
		},
	}

	srv := httptest.NewTLSServer(webSocketServer)
	defer srv.Close()

	c, err := samsungremote.NewClient(
		srv.Listener.Addr().String(),
		"MyApp",
		samsungremote.OptionAllowInsecureTLS,
		samsungremote.OptionTimeout(time.Millisecond*500),
		samsungremote.OptionTokenProvider("my_example_token"),
	)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = c.Send(samsungremote.NewKeyPayload(samsungremote.KEY_VOLUP))
	if err != nil {
		t.Fatal(err.Error())
	}

	expectedMessage := `{"method":"ms.remote.control","params":{"Cmd":"Click","DataOfCmd":"KEY_VOLUP","Option":false,"TypeOfRemote":"SendRemoteKey"}}`
	select {
	case message := <-msgChannel:
		if expectedMessage != message {
			t.Fatalf("expected x got %s instead", message)
		}
	case <-time.After(time.Millisecond * 200):
		t.Fatal("waiting for message, failed")
	}
}
