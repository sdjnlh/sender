package senders

import (
	"errors"
	"github.com/sdjnlh/communal/app"
	"github.com/sdjnlh/communal/config"
	"regexp"
)

const STARTER_SENDER = "SENDER"

type SenderStarter struct {
	*app.BaseStarter
}

func (star *SenderStarter) Start() (err error) {
	if config.Config == nil {
		return errors.New("trying to start senders, but config is not loaded")
	}

	//conf := config.Config.Sub("sender")
	//types := conf.AllKeys()

	//senders := conf.GetStringMap("senders")

	types := config.Config.GetStringMap("sender")

	for typ, typConf := range types {
		if typConf == nil {
			continue
		}

		tc := typConf.(map[string]interface{})
		defaultSenderName := tc["default"]

		sendersConf := tc["senders"]

		if sendersConf == nil {
			continue
		}

		senders := sendersConf.(map[string]interface{})
		for senderName, senderConf := range senders {
			if senderConf == nil {
				return errors.New("no configuration for sender " + senderName)
			}

			sc := senderConf.(map[string]interface{})
			senderTypeIntf := sc["type"]

			if senderTypeIntf == nil {
				return errors.New("type is required for sender " + senderName)
			}

			senderType := senderTypeIntf.(string)
			senderBuilder := builders[senderType]

			if senderBuilder == nil {
				return errors.New("unsupported sender type " + senderType)
			}

			sender, err := senderBuilder(config.Config.Sub("sender." + typ + ".senders." + senderName))

			if err != nil {
				return err
			}

			Senders[typ+"-"+senderName] = sender
		}

		var defaultSender Sender
		if defaultSenderName != nil {
			defaultSender = Senders[defaultSenderName.(string)]
		}

		templatesConf := tc["templates"]

		if templatesConf == nil {
			continue
		}

		templates := templatesConf.(map[string]interface{})

		for templateName, _ := range templates {
			template := &MessageTemplate{Name: typ + "-" + templateName}
			err = config.Config.UnmarshalKey("sender."+typ+".templates."+templateName, template)
			if err != nil {
				return err
			}

			if template.SenderName != "" {
				template.Sender = Senders[typ+"-"+template.SenderName]
			} else {
				template.Sender = defaultSender
			}

			if err = analyseTemplate(template); err != nil {
				return err
			}

			templates[typ+"-"+templateName] = template
		}
	}

	//log.Logger.Debug("senders", zap.Any("types", types))

	//for name, senderConf := range senders {
	//	//sender, err := builders[name]()
	//	log.Logger.Debug("sender ", zap.String("name", name), zap.Any("config", senderConf))
	//}
	return nil
}

func analyseTemplate(template *MessageTemplate) (err error) {
	if template.Sender == nil {
		return errors.New("no sender for template " + template.Name)
	}

	//TODO check if sender has corresponding template configuration, and use it if it has.
	if template.RemoteId == "" {
		if template.Text == "" {
			return errors.New("no content of template " + template.Name)
		} else {
			subMatches := paramRegexp.FindAllStringSubmatch(template.Text, -1)

			template.Params = []string{}
			for _, matches := range subMatches {
				template.Params = append(template.Params, matches[1])
			}

			//template.RegexpTemplate, err = ttp.New("test").Parse(template.Text)
			//if err != nil {
			//	return err
			//}
		}
	}

	return nil
}

func NewSenderStarter() *SenderStarter {
	ss := &SenderStarter{
		BaseStarter: app.NewBaseStarter(STARTER_SENDER, app.PriorityMiddle),
	}

	return ss
}

var paramRegexp *regexp.Regexp

func init() {
	paramRegexp, _ = regexp.Compile("{{.([a-z]+)}}")
}
