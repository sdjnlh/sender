package service

//func Register(server *server.Server) {
//	server.RegisterName("Sender", &SendService{}, "")
//}

//type SenderStarter struct {
//	*starter.BaseStarter
//	Config *viper.Viper
//}
//
//func (star *SenderStarter) Start() error {
//	if star.Config == nil {
//		return errors.New("trying to build db connection, but config is not loaded")
//	}
//
//	defaultSenderName := star.Config.GetString("sender.default")
//
//	var err error
//	senderConf := star.Config.Sub("sender.senders." + defaultSenderName)
//	switch defaultSenderName {
//	case "fungo":
//		Sender, err = sender.NewFungo(senderConf)
//		break
//	case "mockup":
//		Sender, err = sender.NewMockup(senderConf)
//		break
//	default:
//		err = errors.New("default sender not found")
//		break
//	}
//
//	return err
//}
