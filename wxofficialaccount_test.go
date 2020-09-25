package wxofficialaccount

import (
	"fmt"
	"net/http"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"testing"
)

func Test_WxOfficialAccount(t *testing.T) {
	wx := NewWxOfficialAccount("appid", "appsecret", "token")

	// send custom message
	wx.SendCustomTextMessage("toUser", "custom text")

	// send template msg
	msg := &message.TemplateMessage{
		ToUser:     "toUser",
		TemplateID: "uRUXCN_s4Dn27rxSINVL_YX6sWoIomMz2HvSjO3p874",
		URL:        "https://www.baidu.com",
		Data:       map[string]*message.TemplateDataItem{
			"first": &message.TemplateDataItem {Value: "Hello"},
			"course": &message.TemplateDataItem {Value: "math"},
			"result": &message.TemplateDataItem {Value: "100"},
			// "xf": &message.TemplateDataItem {Value: "value"},
			// "ctype": &message.TemplateDataItem {Value: "value"},
			// "examtype": &message.TemplateDataItem {Value: "value"},
			"remark": &message.TemplateDataItem {Value: "world"},
		},
	}
	wx.SendTemplateMessage(msg)
	
	// serve wechat
	wx.SetMessageHandler(func(msg message.MixMessage) *message.Reply {
        text := message.NewText("Req:" + msg.Content)
        return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	})
	// http.HandleFunc("/", wx.ServeHTTP)
    fmt.Println("wechat server listener at", ":8001")
    err := http.ListenAndServe(":8001", wx)
    if err != nil {
        fmt.Printf("start server error , err=%v", err)
	}
}