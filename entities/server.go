package entities

import (
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

var server *socketio.Server
var router *gin.Engine

func init() {
	server = socketio.NewServer(nil)
	router = gin.New()
}

func GetServer() *socketio.Server {
	return server
}

func GetRouter() *gin.Engine {
	return router
}
