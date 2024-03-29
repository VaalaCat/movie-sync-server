package room

import (
	"movie-sync-server/entities"
	"strings"

	socketio "github.com/googollee/go-socket.io"
	"github.com/sirupsen/logrus"
)

func PlayEndpoint(s socketio.Conn, msg string) {
	server := entities.GetServer()
	logrus.Println("play:", msg)
	Splitted := strings.Split(msg, ":::")
	room, showName := Splitted[0], Splitted[1]
	for _, r := range entities.Cinema {
		if (*r).Name() == room {
			(*r).Broadcast("play", showName, server)
		}
	}
}
