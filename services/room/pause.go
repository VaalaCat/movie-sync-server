package room

import (
	"log"
	"movie-sync-server/entities"
	"strings"

	socketio "github.com/googollee/go-socket.io"
)

func PauseEndpoint(s socketio.Conn, msg string) {
	server := entities.GetServer()
	log.Println("stop:", msg)
	Splitted := strings.Split(msg, ":::")
	room, showName := Splitted[0], Splitted[1]
	for _, r := range entities.Cinema {
		if (*r).Name() == room {
			(*r).Broadcast("pause", showName, server)
		}
	}
}
