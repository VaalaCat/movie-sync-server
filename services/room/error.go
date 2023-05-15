package room

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/sirupsen/logrus"
)

func ErrorEndpoint(s socketio.Conn, e error) {
	logrus.Println("meet error:", e)
}
