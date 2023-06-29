package senders

import (
	"github.com/sdjnlh/communal/log"
	"github.com/sdjnlh/sender/model"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const MOCKUP = "mockup"

type MockupSender struct{}

func (sender *MockupSender) TemplateContent(msg *model.Message) string {
	log.Logger.Info("send mockup sms", zap.Any("req", msg))

	return ""
}

func (sender *MockupSender) Send(msg *model.Message) (err error) {
	log.Logger.Info("send mockup sms", zap.Any("req", msg))

	return nil
}

func init() {
	builders[MOCKUP] = func(conf *viper.Viper) (sd Sender, err error) {
		log.Slog.Debug("sender config", conf)
		return &MockupSender{}, nil
	}
}
