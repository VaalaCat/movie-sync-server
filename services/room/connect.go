package room

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/sirupsen/logrus"
)

func ConnectEndpoint(s socketio.Conn) error {
	s.SetContext("")
	logrus.Println("connected:", s.ID())
	return nil
}
