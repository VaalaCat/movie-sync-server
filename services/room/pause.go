package room

import (
	"movie-sync-server/entities"
	"strings"

	socketio "github.com/googollee/go-socket.io"
	"github.com/sirupsen/logrus"
)

func PauseEndpoint(s socketio.Conn, msg string) {
	server := entities.GetServer()
	logrus.Println("stop:", msg)
	Splitted := strings.Split(msg, ":::")
	room, showName := Splitted[0], Splitted[1]
	for _, r := range entities.Cinema {
		if (*r).Name() == room {
			(*r).Broadcast("pause", showName, server)
		}
	}
}
