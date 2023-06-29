package model

type SMSSendRequest struct {
	Type   string `json:"type" valid:"required"`
	Mcc    string `json:"mcc"`
	Mobile string `json:"mobile" valid:"required"`
	Token  string `json:"token" valid:"required"`
	Code   string `json:"code" valid:"required"`
}

type TokenReply struct {
	Token string
}

type Message struct {
	Group      string `json:"-"`
	Type       string
	TemplateId string
	Recipient  string
	Content    string
	Language   string
	Params     map[string]interface{}
}

func (message *Message) GetString(key string) string {
	v := message.Params[key]
	if v == nil {
		return ""
	}

	return v.(string)
}

const (
	GATEWAY_ERROR = "gateway_error"
)
