package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/sdjnlh/communal"
	"github.com/sdjnlh/communal/log"
	"github.com/sdjnlh/communal/util"
	"github.com/sdjnlh/sender/model"
	"github.com/sdjnlh/sender/senders"
	"go.uber.org/zap"
	"math/rand"
	"strings"
	"time"
)

const CodeExpireMinutes = 10

type SendService struct {
}

func (service *SendService) VerifyCode(ctx context.Context, msg *model.Message, result *communal.Result) (err error) {
	//todo generate code and set it into msg.params
	ut := uuid.NewV4()

	if msg.Params == nil {
		msg.Params = make(map[string]interface{})
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	msg.Params["code"] = fmt.Sprintf("%06v", rnd.Int31n(1000000))
	msg.Params["expire"] = CodeExpireMinutes

	//todo generate token and cache it in redis as key with code as value
	err = send(msg)
	if err != nil {
		log.Logger.Debug("failed to send verify code ", zap.String("error", err.Error()))
		return
	}
	result.Data = ut.String()
	return
}

func SendDenialLetter(msg *model.Message) (err error) {
	return senders.Send(senders.EmailSenderGroupSystem, msg)
}
func EmailTemplateContent(msg *model.Message) (content string, err error) {
	return senders.EmailContent(senders.EmailSenderGroupSystem, msg)
}

func send(msg *model.Message) (err error) {
	template := senders.GetTemplate(msg)
	if template == nil {
		log.Logger.Warn("unsupported template ", zap.Any("message", msg))
		return errors.New("unsupported template ")
	}

	for _, param := range template.Params {
		if util.IsEmpty(msg.Params[param]) {
			return errors.New("empty param " + param)
		}
	}

	if msg.Content, err = template.Render(msg.Params); err != nil {
		return err
	}
	if err = template.Sender.Send(msg); err != nil {
		return err
	}

	return nil
}

func generateCode(length int) string {
	//seed := "1234556789"
	rnd := rand.New(rand.NewSource(9))

	sb := strings.Builder{}
	for i := 0; i < length; i++ {
		sb.WriteRune(rune(rnd.Int() + 1))
	}

	return sb.String()
}
