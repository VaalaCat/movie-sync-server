package room

import (
	"log"
	"movie-sync-server/entities"

	socketio "github.com/googollee/go-socket.io"
)

func GetUrlEndpoint(s socketio.Conn, msg string) {
	log.Println("getUrl:", msg)
	room := msg
	for _, r := range entities.Cinema {
		if (*r).Name() == room {
			tmpUrl := (*r).GetUrl()
			s.Emit("setUrl", tmpUrl)
			break
		}
	}
}
