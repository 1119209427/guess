package api

import (
	"github.com/gin-gonic/gin"
	"guess/websocker"
)

func Route() {
	engine := gin.Default()
	engine.GET("guess", websocker.ServerWs)
	err := engine.Run()
	if err != nil {
		panic(err)
	}

}
