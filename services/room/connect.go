package room

import (
	"log"

	socketio "github.com/googollee/go-socket.io"
)

func ConnectEndpoint(s socketio.Conn) error {
	s.SetContext("")
	log.Println("connected:", s.ID())
	return nil
}
