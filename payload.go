package samsungremote

type AppData struct {
	AppId      string `json:"appId,omitempty"`
	ActionType string `json:"action_type,omitempty"`
	MetaTag    string `json:"metaTag,omitempty"`
}

type Params struct {
	Event        string      `json:"event,omitempty"`
	To           string      `json:"to,omitempty"`
	Data         interface{} `json:"data,omitempty"`
	Cmd          string      `json:"Cmd,omitempty"`
	DataOfCmd    string      `json:"DataOfCmd,omitempty"`
	Option       bool        `json:"Option"`
	TypeOfRemote string      `json:"TypeOfRemote,omitempty"`
}

// Payload represent the data structure that you
// need to emit when sending a command to the TV
type Payload struct {
	Method string `json:"method"`
	Params Params `json:"params"`
}

// NewKeyPayload builds a payload to execute a KEY_* command
func NewKeyPayload(key string) Payload {
	return Payload{
		Method: "ms.remote.control",
		Params: Params{Cmd: "Click", DataOfCmd: key, Option: false, TypeOfRemote: "SendRemoteKey"},
	}
}

// NewTizenOpenUrlPayload builds a payload that will command the TV to open the specified url on the tizen browser
func NewTizenOpenUrlPayload(url string) Payload {
	return Payload{
		Method: "ms.channel.emit",
		Params: Params{
			Event: "ed.apps.launch",
			To:    "host",
			Data: AppData{
				AppId:      "org.tizen.browser",
				ActionType: "NATIVE_LAUNCH",
				MetaTag:    url,
			},
		},
	}
}
