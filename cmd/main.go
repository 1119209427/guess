package main

import (
	"github.com/gin-gonic/gin"
	"guess/websocker"
)

func main() {
	route := gin.Default()
	route.GET("guess", websocker.ServerWs)
	route.Run()
	hub := websocker.NewHub()
	go hub.Run()

}
