package email

import (
	"bytes"
	"errors"
	"github.com/sdjnlh/communal/log"
	"github.com/sdjnlh/sender/model"
	"github.com/sdjnlh/sender/senders"
	"go.uber.org/zap"
	"net/smtp"
	"strings"
)

type SmtpSender struct {
	Name     string
	Address  string
	Password string
	Host     string
	Endpoint string
}

func (sender *SmtpSender) Send(msg *model.Message) error {
	auth := smtp.PlainAuth("", sender.Address, sender.Password, sender.Host)
	//sendTo := strings.Split(msg.Recipient, ";")
	//done := make(chan error, 1024)
	var contentType string
	log.Logger.Debug("send email", zap.Any("recipient", msg.Recipient))
	templateOption := senders.TemplateOptions[msg.Type+"."+msg.TemplateId]
	if templateOption == nil {
		return errors.New("invalid template id " + msg.Type + "." + msg.TemplateId)
	}
	title := "RSMC - Retrieve password"
	if msg.Language == "zh" {
		title = "RSMC - 找回密码"
	}
	templateOption.Title = title

	var contentBuffer = &bytes.Buffer{}
	if templateOption.Type == "html" {
		contentType = "Content-Type: text/html; charset=UTF-8"
		err := senders.HtmlTemplate.ExecuteTemplate(contentBuffer, msg.TemplateId+"-"+msg.Language+".html", msg.Params)
		if err != nil {
			log.Logger.Error("fail to render message template", zap.Error(err))
			return err
		}
	} else {
		contentType = "Content-Type: text/plain" + "; charset=UTF-8"
		//TODO execute text template
	}

	subject := msg.GetString("subject")
	if subject == "" {
		subject = templateOption.Title
	}
	recips := strings.Split(msg.Recipient, ";")
	senderEmail := msg.Params["from"].(string)
	str := "From: " + sender.Name + "<" + senderEmail + ">\r\nTo: " + msg.Recipient + "\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + contentBuffer.String()
	log.Logger.Debug("str", zap.Any("str", str))
	err := smtp.SendMail(
		sender.Endpoint,
		auth,
		sender.Address,
		recips,
		[]byte(str),
	)
	if err != nil {
		log.Logger.Debug("error", zap.Any("smtp error", err))
		return err
	}
	log.Logger.Debug("message sent ", zap.String("template", msg.TemplateId), zap.String("recipient", msg.Recipient))
	return nil
}

func (sender *SmtpSender) TemplateContent(msg *model.Message) string {
	log.Logger.Debug("send email", zap.Any("recipient", msg.Recipient))
	templateOption := senders.TemplateOptions[msg.Type+"."+msg.TemplateId]
	var contentBuffer = &bytes.Buffer{}
	if templateOption.Type == "html" {
		//contentType = "Content-Type: text/html; charset=UTF-8"
		err := senders.HtmlTemplate.ExecuteTemplate(contentBuffer, msg.TemplateId+"-"+msg.Language+".html", msg.Params)
		if err != nil {
			log.Logger.Error("fail to render message template", zap.Error(err))
		}
	} else {
		//contentType = "Content-Type: text/plain" + "; charset=UTF-8"
		//TODO execute text template
	}
	return contentBuffer.String()
}
