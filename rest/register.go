package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sdjnlh/sender/rest/handler"
)

func RegisterAPIs(router *gin.Engine) {
	handler.NewSenderAPI().Register(router)
}
