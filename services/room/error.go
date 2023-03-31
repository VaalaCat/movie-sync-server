package room

import (
	"log"

	socketio "github.com/googollee/go-socket.io"
)

func ErrorEndpoint(s socketio.Conn, e error) {
	log.Println("meet error:", e)
}
