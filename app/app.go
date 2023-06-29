package sender

import (
	"github.com/sdjnlh/communal/app"
)

var Service = app.NewService("sender-service")
var Api = app.NewWeb("sender-api")
