package email

import (
	"errors"
	"github.com/sdjnlh/communal"
	"github.com/sdjnlh/communal/app"
	"github.com/sdjnlh/communal/log"
	"github.com/sdjnlh/sender/app"
	"github.com/sdjnlh/sender/senders"
	"github.com/spf13/viper"
	"html/template"
	"strings"
)

var builders = make(map[string]func(conf *viper.Viper) (senders.Sender, error))

type Starter struct {
	*app.BaseStarter
}

func (star *Starter) Start(ctx *communal.Context) (err error) {
	//load senders
	senderConfigs := sender.Service.RawConfig.Sub("email.sender")
	configMap := sender.Service.RawConfig.GetStringMap("email.sender")

	for name, _ := range configMap {
		if strings.Index(name, ".") > 0 {
			continue
		}
		config := senderConfigs.Sub(name)
		tp := config.GetString("type")
		builder := builders[tp]
		if builder == nil {
			return errors.New("no builder found for email sender of type " + tp)
		}
		sd, err := builder(config)
		if err != nil {
			return err
		}
		senders.Senders["email."+name] = sd
	}

	//load html templates
	templateDir := sender.Service.RawConfig.GetString("sender.htmlTemplateDir")
	senders.HtmlTemplate, err = template.ParseGlob(templateDir + "/*.html")
	if err != nil {
		return err
	}

	//load template options
	templateConfigs := sender.Service.RawConfig.Sub("email.template")
	templateConfigMap := sender.Service.RawConfig.GetStringMap("email.template")

	for name, _ := range templateConfigMap {
		if strings.Index(name, ".") > 0 {
			continue
		}
		config := templateConfigs.Sub(name)
		option := senders.TemplateOption{}
		err = config.Unmarshal(&option)
		if err != nil {
			return err
		}

		option.Name = name
		senders.TemplateOptions["email."+name] = &option
	}

	//load group sender mapping
	groupMap := sender.Service.RawConfig.GetStringMapString("email.group")
	for key, value := range groupMap {
		sd := senders.Senders["email."+value]
		if sd == nil {
			return errors.New("invalid email group, " + key + ":" + value + ", sender not found")
		}

		senders.Groups["email."+key] = sd
	}

	return nil
}

func NewStarter() *Starter {
	ss := &Starter{
		BaseStarter: app.NewBaseStarter("email_sender_starter", app.PriorityLowest),
	}

	return ss
}

func init() {
	builders["smtp"] = func(conf *viper.Viper) (sd senders.Sender, err error) {
		log.Slog.Debug("smtp sender config", conf)
		sd = &SmtpSender{}
		err = conf.Unmarshal(sd)

		return
	}
}
