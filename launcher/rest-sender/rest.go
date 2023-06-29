package main

import (
	"flag"

	"github.com/json-iterator/go/extra"
	"github.com/sdjnlh/communal/log"
	"github.com/sdjnlh/communal/util"
	"github.com/sdjnlh/communal/web"
	sender "github.com/sdjnlh/sender/app"
	"github.com/sdjnlh/sender/senders/email"
	"golang.org/x/sync/errgroup"

	"github.com/gin-gonic/gin"
	"github.com/sdjnlh/communal/app"
	"github.com/sdjnlh/sender/rest"
)

func main() {
	sender.Api.Mount(sender.Service)
	app.RegisterStarter(sender.Api)
	app.RegisterStarter(email.NewStarter())
	if err := app.Start(); err != nil {
		panic(err)
	}

	flag.Parse()
	log.Slog.Info("starting rest")
	extra.SetNamingStrategy(util.LowerFirst)
	frontApiServer := gin.New()
	frontApiServer.Use(gin.Logger())
	frontApiServer.Use(gin.Recovery())
	frontApiServer.Use(web.CorsHandler(sender.Api))
	//sender.Api.ManageUserState(frontApiServer, userApp.UserBuilder)
	rest.RegisterAPIs(frontApiServer)
	var g errgroup.Group

	g.Go(func() error {
		return frontApiServer.Run(":" + sender.Api.Port)
	})

	if err := g.Wait(); err != nil {
		log.Slog.Fatal(err)
	}
}
