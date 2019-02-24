package samsungremote_test

import (
	"encoding/json"
	"testing"

	"github.com/stefanoj3/samsungremote"
)

func TestNewKeyPayload(t *testing.T) {
	expectedMessage := `{"method":"ms.remote.control","params":{"Cmd":"Click","DataOfCmd":"KEY_RETURN","Option":false,"TypeOfRemote":"SendRemoteKey"}}`
	b, err := json.Marshal(samsungremote.NewKeyPayload(samsungremote.KEY_RETURN))
	if err != nil {
		t.Fatalf("failed to marshal: %s", err.Error())
	}
	actualMessage := string(b)

	if expectedMessage != actualMessage {
		t.Fatalf("Expected %s, got %s instead", expectedMessage, actualMessage)
	}
}

func TestNewTizenOpenUrlPayload(t *testing.T) {
	expectedMessage := `{"method":"ms.channel.emit","params":{"event":"ed.apps.launch","to":"host","data":{"appId":"org.tizen.browser","action_type":"NATIVE_LAUNCH","metaTag":"http://github.com/stefanoj3"},"Option":false}}`
	b, err := json.Marshal(samsungremote.NewTizenOpenUrlPayload("http://github.com/stefanoj3"))
	if err != nil {
		t.Fatalf("failed to marshal: %s", err.Error())
	}
	actualMessage := string(b)

	if expectedMessage != actualMessage {
		t.Fatalf("Expected %s, got %s instead", expectedMessage, actualMessage)
	}
}
