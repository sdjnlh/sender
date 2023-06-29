package senders

import (
	"bytes"
	"errors"
	"github.com/sdjnlh/sender/model"
	"github.com/spf13/viper"
	"html/template"
	ttp "text/template"
)

//type SenderBuilder func (conf *viper.Viper) Sender

var builders = make(map[string]func(conf *viper.Viper) (Sender, error))
var Senders = make(map[string]Sender)
var Groups = make(map[string]Sender)
var templates = make(map[string]*MessageTemplate)
var TemplateOptions = make(map[string]*TemplateOption)
var TextTemplates = make(map[string]*ttp.Template)
var HtmlTemplate *template.Template

//var strategies = make(map[string][]Strategy)

type Sender interface {
	Send(msg *model.Message) error
	TemplateContent(msg *model.Message) string
}

func Send(senderGroup string, msg *model.Message) (err error) {
	sender := Groups[msg.Type+"."+senderGroup]

	if sender == nil {
		return errors.New("invalid sender group " + senderGroup)
	}

	return sender.Send(msg)
}
func EmailContent(senderGroup string, msg *model.Message) (content string, err error) {
	sender := Groups[msg.Type+"."+senderGroup]

	if sender == nil {
		return "", errors.New("invalid sender group " + senderGroup)
	}

	return sender.TemplateContent(msg), nil
}

func SendMessage(msg *model.Message) (err error) {
	sender := Groups[msg.Type+"."+msg.Group]

	if sender == nil {
		return errors.New("invalid sender group " + msg.Group)
	}

	return sender.Send(msg)
}

type TemplateOption struct {
	Name    string
	Type    string
	Dynamic bool
	Title   string
}

const (
	EmailSenderGroupSystem = "system"
)

type MessageTemplate struct {
	Name           string
	RemoteId       string //use RemoteId prior to text if it not empty
	Text           string
	SenderName     string
	Sender         Sender
	Params         []string
	RegexpTemplate *ttp.Template
}

func (template *MessageTemplate) Render(params map[string]interface{}) (content string, err error) {
	if template.Text == "" {
		return "", errors.New("no content of template " + template.Name)
	}

	sb := bytes.NewBufferString("")
	if err = template.RegexpTemplate.Execute(sb, params); err != nil {
		return content, err
	}

	content = sb.String()
	return
}

func GetTemplate(msg *model.Message) *MessageTemplate {
	if msg == nil || msg.Type == "" || msg.TemplateId == "" {
		return nil
	}

	return templates[msg.Type+"-"+msg.TemplateId]
}
