# wxofficialaccount

A simple Wechat Official Account SDK based on [silenceper/wechat](https://github.com/silenceper/wechat).

Support Functions:

- server
- CustomerMessage
- TemplateMessage

## Usages

Install:

```sh
go get github.com/cdfmlr/wxofficialaccount
```

Start:

```go
import (
    "github.com/cdfmlr/wxofficialaccount"
    "github.com/silenceper/wechat/v2/officialaccount/message"
)

w := NewWxOfficialAccount("appid", "appsecret", "token")
```

Wechat Official Account HTTP service:

```go
w.SetMessageHandler(func(msg message.MixMessage) *message.Reply {
    text := message.NewText("Req:" + msg.Content)
    return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
})

err := http.ListenAndServe(":8001", w)
if err != nil {
    fmt.Printf("start server error , err=%v", err)
}
```

Or:

```go
http.HandleFunc("/path/", w.ServeHTTP)
http.ListenAndServe(":8001", nil)
```

Send a Customer Message:

```go
w.SendCustomTextMessage("toUserID", "custom text")
```

Send a Template Message:

```go
msg := &message.TemplateMessage{
    ToUser:     "toUser",
    TemplateID: "uRUXCN_s4Dn27rxSINVL_YX6sWoI0mMz2HvSjO3p874",
    URL:        "https://www.baidu.com",
    Data:       map[string]*message.TemplateDataItem{
        "first": &message.TemplateDataItem {Value: "Hello"},
        "remark": &message.TemplateDataItem {Value: "world"},
    },
}
w.SendTemplateMessage(msg)
```

## License

Copyright 2020 CDFMLR

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.