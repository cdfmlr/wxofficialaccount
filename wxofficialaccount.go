package wxofficialaccount

import (
	"log"
	"net/http" 

	wechat "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	"github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

// WxOfficialAccount 是对 github.com/silenceper/wechat/v2/officialaccount.OfficialAccount 的简单封装。
// 
// WxOfficialAccount “继承”了 silenceper/wechat/v2/officialaccount.OfficialAccount，
// 同时补充实现了字段 messageHandler, 以及方法 SetMessageHandler, ServeHTTP, SendCustomTextMessage, SendTemplateMessage。
// WxOfficialAccount 实现了 http.Handler 接口，可以直接用来开 http 服务。
//
// 该结构体应该由 NewWxOfficialAccount 函数完成构造
//
// - messageHandler 字段是一个函数，用来处理公众号被动消息回复接收到的消息，使用 SetMessageHandler 方法进行设置。
// - ServeHTTP 方法是一个 http.HandlerFunc，用来处理接收到的消息（具体逻辑由 messageHandler 提供）。
// - SendCustomTextMessage 方法发送纯文本的客服消息。
// - SendTemplateMessage 方法发送模版消息。
type WxOfficialAccount struct {
	officialaccount.OfficialAccount
	messageHandler func(msg message.MixMessage) *message.Reply
}

// NewWxOfficialAccount 构造一个 WxOfficialAccount，初始化其中的 OfficialAccount、messageHandler。
// 这里初始化的 messageHandler 是一个「空方法」，只是打印一条日志，并响应给用户 "Server Error: Not Yet Implemented"
// 如果需要使用公众号 http 服务，messageHandler 应给被 wxOfficialAccount.SetMessageHandler 方法设置。
func NewWxOfficialAccount(appID, appSecret, token string) *WxOfficialAccount {
	wc := wechat.NewWechat()
	cfg := &config.Config{
        AppID:     appID,
        AppSecret: appSecret,
        Token:     token,
        // EncodingAESKey: "xxxx",
        Cache: cache.NewMemory(), // 这里本地内存保存access_token，也可选择redis，memcache或者自定cache
	}
	officialAccount := wc.GetOfficialAccount(cfg)

	return &WxOfficialAccount {
		OfficialAccount: *officialAccount,
		messageHandler: func(msg message.MixMessage) *message.Reply {
			log.Println("WxOfficialAccount: messageHandler not yet implemented")
			text := message.NewText("Server Error: Not Yet Implemented")
			return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
		},
	}
}

// SetMessageHandler 方法用来设置处理公众号被动消息回复接收到的消息的具体逻辑
// 如果需要使用公众号 http 服务，就应该调用该方法，设置具体的消息处理逻辑。
// Example:
//     w := NewWxOfficialAccount(...)
//     w.SetMessageHandler(func(msg message.MixMessage) *message.Reply {
//         // 回复消息
//         text := message.NewText(msg.Content)
//         return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
//    })
// message.MixMessage 即 "github.com/silenceper/wechat/v2/officialaccount/message".MixMessage，存放所有微信发送过来的消息和事件
// For more low-level details: https://github.com/silenceper/wechat/blob/release-2.0/officialaccount/server/server.go
func (w *WxOfficialAccount) SetMessageHandler(handler func(msg message.MixMessage) *message.Reply) {
	w.messageHandler = handler
}

// ServeHTTP 方法是一个 http.HandlerFunc，用来处理接收到的消息，提供公众号的 http 服务（具体逻辑由 messageHandler 提供）
// 使用 ServeHTTP 前，必须先调用 SetMessageHandler 实现具体逻辑，否则默认返回错误响应("Server Error: Not Yet Implemented")
// 这个方法也让 WxOfficialAccount 实现了 http.Handler 接口。
// Example:
//     w := NewWxOfficialAccount(...)
//     w.SetMessageHandler(...)
//     http.HandleFunc("/", w.ServerHandle)
//     http.ListenAndServe(":8001", nil)
// Or，use WxOfficialAccount as a http.Handler:
//     w := NewWxOfficialAccount(...)
//     w.SetMessageHandler(...)
//     http.ListenAndServe(":8001", w)
// For more low-level details: https://github.com/silenceper/wechat/blob/release-2.0/officialaccount/server/server.go
func (w *WxOfficialAccount) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	// 传入请求
	server := w.GetServer(r, writer)
	//设置接收消息的处理方法
	server.SetMessageHandler(w.messageHandler)
	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		log.Println(err)
		return
	}
	//发送回复的消息
	server.Send()
}

// SendCustomTextMessage 方法发送纯文本的客服消息。
// 该方法为 `wxOfficialAccount.GetCustomerMessageManager().Send(message.NewCustomerTextMessage(toUser, text))` 的 shorthand
// 传入参数 toUser 为要发送给的用户的 id，text 为消息内容。
// Example:
//     w := NewWxOfficialAccount(...)
//     w.SendCustomTextMessage("w110w9sXYZ2BxyzABCDa2_WSLNMD", "hello")
// For more low-level details: https://github.com/silenceper/wechat/blob/release-2.0/officialaccount/message/customer_message.go
func (w *WxOfficialAccount) SendCustomTextMessage(toUser, text string) {
	// manager := message.NewMessageManager(w.GetContext())
	manager := w.GetCustomerMessageManager()
	msg := message.NewCustomerTextMessage(toUser, text)
	if err := manager.Send(msg); err != nil {
		log.Println("send CustomMessage error:", err)
	}
}

// SendTemplateMessage 方法发送模版消息。
// 传入参数为 *"github.com/silenceper/wechat/v2/officialaccount/message".TemplateMessage, 是要发送的模板消息内容
// See https://github.com/silenceper/wechat/blob/release-2.0/officialaccount/message/template.go for more details
// 该方法为 `wxOfficialAccount.GetTemplate().Send(msg)` 的 shorthand
// Example:
//     msg := &message.TemplateMessage{
//         ToUser:     "w110w9sXYZ2BxyzABCDa2_WSLNMD",
//         TemplateID: "uRUXCN_s4Dn27rxSINVL_YX3sWoIomMz2HvSjO3p87e",
//         URL:        "https://www.baidu.com",
//         Data:       map[string]*message.TemplateDataItem{
//             "first": &message.TemplateDataItem {Value: "value"},
//             "course": &message.TemplateDataItem {Value: "value"},
//             "result": &message.TemplateDataItem {Value: "value"},
//             "remark": &message.TemplateDataItem {Value: "value"},
//         },
//     }
//     wx.SendTemplateMessage(msg)
func (w *WxOfficialAccount) SendTemplateMessage(msg *message.TemplateMessage) {
	// templateMsg := message.NewTemplate(w.GetContext())
	templateMsg := w.GetTemplate()
	if msgId, err := templateMsg.Send(msg); err != nil {
		log.Printf("send TemplateMessage (msgId=%v) error: %v", msgId, err)
	}
}
