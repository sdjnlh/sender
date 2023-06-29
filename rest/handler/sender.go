package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sdjnlh/communal"
	"github.com/sdjnlh/communal/log"
	"github.com/sdjnlh/communal/rpc"
	"github.com/sdjnlh/communal/web"
	"github.com/sdjnlh/sender/model"
	"github.com/sdjnlh/sender/service"
	"go.uber.org/zap"
)

type SenderAPI struct {
	*web.RestHandler
}

func NewSenderAPI() *SenderAPI {
	return &SenderAPI{
		RestHandler: web.DefaultRestHandler,
	}
}

func (api *SenderAPI) SendVerifyCode(c *gin.Context) {
	var msg = &model.Message{
		Type:       "sms",
		TemplateId: "verify-code",
		Params:     make(map[string]interface{}),
	}

	req := &model.SMSSendRequest{}
	err := api.BindAndValidate(c, req, "")

	//err := validator.ValidateMap("", c.Keys)
	if err != nil {
		api.BadRequestWithError(c, err)
		return
	}

	// TODO check verify code against redis or call rpc method

	msg.Recipient = req.Mobile

	var result = &communal.Result{}

	if err := rpc.Call(context.Background(), "Sender", "VerifyCode", msg, result); err != nil {
		log.Logger.Error("failed to call rpc", zap.Any("error", err))
		api.RPCError(c)
		return
	}

	//log.Fatal(r.Name)
	log.Logger.Debug("verify code send ", zap.Any("result", result))

	api.SuccessWithPair(c, "result", result.Data)
}

func (api *SenderAPI) SendEmail(c *gin.Context) {
	var msg = &model.Message{
		Type:   "email",
		Params: make(map[string]interface{}),
	}

	err := api.Bind(c, msg)
	if err != nil {
		api.BadRequestWithError(c, err)
		return
	}

	var result = &communal.Result{}

	if err := service.SendDenialLetter(msg); err != nil {
		log.Logger.Error("failed to send email", zap.Any("error", err))
		api.RPCError(c)
		return
	}

	//log.Fatal(r.Name)
	log.Logger.Debug("send static email ", zap.Any("result", result))

	api.SuccessWithPair(c, "result", result.Data)
}

func (api *SenderAPI) Register(router *gin.Engine) {
	v1 := router.Group("/v1")
	v1.POST("/verifycode", api.SendVerifyCode)
	v1.POST("/send/email", api.SendEmail)
}
