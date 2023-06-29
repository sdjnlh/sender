package senders

import (
	"encoding/json"
	"github.com/sdjnlh/communal/errors"
	"github.com/sdjnlh/communal/log"
	"github.com/sdjnlh/sender/model"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
)

const FUNGO_SMS = "fungo"

type FungoSMSSender struct {
	Address  string
	Username string
	Password string
}

func (sender *FungoSMSSender) TemplateContent(msg *model.Message) string {
	log.Logger.Debug("send email", zap.Any("recipient", msg.Recipient))
	return Sender.TemplateContent(sender, msg)
}

func (sender *FungoSMSSender) Send(msg *model.Message) (err error) {
	log.Logger.Info("send sms", zap.Any("req", msg))

	if msg.Recipient == "" || msg.Content == "" {
		return errors.InvalidParams()
	}

	resp, err := http.Post(sender.Address,
		"application/x-www-form-urlencoded",
		strings.NewReader("CpName="+sender.Username+"&CpPassword="+sender.Password+"&DesMobile="+msg.Recipient+"&Content="+msg.Content))
	if err != nil {
		log.Logger.Error("fail send mobile code{}", zap.Any("msg", msg), zap.String("error", err.Error()))
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var response map[string]interface{}
	json.Unmarshal(body, &response)
	respCode, _ := response["code"]
	if "0" != respCode.(string) {
		log.Logger.Error("send handler error{}", zap.Any("error code", respCode))
		return &errors.SimpleBizError{
			Code: model.GATEWAY_ERROR,
			Msg:  "failed to send sms by " + FUNGO_SMS + ", " + respCode.(string),
		}
	}

	return nil
}

func init() {
	builders[FUNGO_SMS] = func(conf *viper.Viper) (sd Sender, err error) {
		log.Slog.Debug("sender config", conf)
		sd = &FungoSMSSender{}
		err = conf.Unmarshal(sd)

		return
	}
}
