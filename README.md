## samsungremote

[![CircleCI](https://circleci.com/gh/stefanoj3/samsungremote/tree/master.svg?style=svg)](https://circleci.com/gh/stefanoj3/samsungremote/tree/master)

samsungremote is a library that allows to send commands to your samsung TV via websocket.

I have tested this library only against a samsung model `UEMU6199UXZG`, 
however based on the research I have done, any tv exposing this `/api/v2/channels/samsung.remote.control` 
endpoint on the 8002 port should work just fine.

Example usage:
```go
package main

import (
    "fmt"
    "time"
    
    "github.com/stefanoj3/samsungremote"
)

func main() {
    c, err := samsungremote.NewClient(
        "192.168.3.6:8002",
        "MyCoolApplication",
        samsungremote.OptionAllowInsecureTLS,
        samsungremote.OptionTimeout(time.Millisecond*500),
    )
    if err != nil {
        panic(err)
    }
    
    token, err := c.AcquireToken()
    if err != nil {
        panic(err)
    }
    
    fmt.Println("my application have been autorized with the token:", token)
    // you can save the token for later,
    // so next time you create a client you can provide it to the client
    // (passing the samsungremote.OptionTokenProvider(token) when creating the client)
    // and you wont need to authorize the application every time you trigger a command 
    
    err = c.Send(samsungremote.NewKeyPayload(samsungremote.KEY_VOLUP))
    if err != nil {
        panic(err)
    }
}
``` 

#### Notes
When scanning the TV using nmap (`nmap -p "*" <my-tv-ip-address>`) those are the ports I found open:
```
Starting Nmap 7.70 ( https://nmap.org ) at 2019-02-24 18:52 CET
Nmap scan report for XXX
Host is up (0.018s latency).
Not shown: 8296 closed ports
PORT      STATE SERVICE
7676/tcp  open  imqbrokerd
8001/tcp  open  vcom-tunnel
8002/tcp  open  teradataordbms
8080/tcp  open  http-proxy
9119/tcp  open  mxit
9197/tcp  open  unknown
32768/tcp open  filenet-tms
32769/tcp open  filenet-rpc
32770/tcp open  sometimes-rpc3
32771/tcp open  sometimes-rpc5
```


#### Legal note

SAMSUNG is a trademark of Samsung Electronics Co., Ltd.. and/or other respective owners. This software is 
not created by Samsung, and is for purely educational and research purposes. 
It is your sole responsibility to follow copyright law.
The creators hold no responsibility for the consequences of use of this software.