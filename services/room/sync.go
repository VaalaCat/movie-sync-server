package room

import (
	"fmt"
	"log"
	"movie-sync-server/entities"
	"strings"

	socketio "github.com/googollee/go-socket.io"
)

func SyncEndpoint(s socketio.Conn, msg string) {
	server := entities.GetServer()
	log.Println("sync:", msg)
	Splitted := strings.Split(msg, ":::")
	room, showName := Splitted[0], Splitted[1]
	for _, r := range entities.Cinema {
		if (*r).Name() == room {
			(*r).Broadcast("sync", fmt.Sprintf("%d:::%d:::%s", (*r).GetMinTime(), (*r).GetMaxTime(), showName), server)
			break
		}
	}
}
