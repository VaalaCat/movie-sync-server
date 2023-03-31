package room

import (
	"log"
	"movie-sync-server/entities"
	"strings"

	socketio "github.com/googollee/go-socket.io"
)

func GetTimeEndpoint(s socketio.Conn, msg string) {
	server := entities.GetServer()
	log.Println("getTime:", msg)
	Splitted := strings.Split(msg, ":::")
	room, showName := Splitted[0], Splitted[1]
	username := s.ID()
	for _, r := range entities.Cinema {
		if (*r).Name() == room {
			tmpUser := (*r).GetUser(username)
			if tmpUser != nil {
				(*r).Broadcast("getTime", showName, server)
				break
			}
		}
	}
}
