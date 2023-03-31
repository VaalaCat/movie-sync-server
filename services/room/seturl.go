package room

import (
	"fmt"
	"log"
	"movie-sync-server/entities"
	"strings"

	socketio "github.com/googollee/go-socket.io"
)

func SetUrlEndpoint(s socketio.Conn, msg string) {
	server := entities.GetServer()
	log.Println("setUrl:", msg)
	Splitted := strings.Split(msg, ":::")
	room, showName, url := Splitted[0], Splitted[1], Splitted[2]
	for _, r := range entities.Cinema {
		if (*r).Name() == room {
			(*r).SetUrl(url)
			(*r).Broadcast("setUrl", fmt.Sprintf("%s:::%s", showName, url), server)
			break
		}
	}
}
