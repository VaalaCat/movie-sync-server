package room

import (
	"movie-sync-server/entities"

	socketio "github.com/googollee/go-socket.io"
	"github.com/sirupsen/logrus"
)

func GetUrlEndpoint(s socketio.Conn, msg string) {
	logrus.Println("getUrl:", msg)
	room := msg
	for _, r := range entities.Cinema {
		if (*r).Name() == room {
			tmpUrl := (*r).GetUrl()
			s.Emit("setUrl", tmpUrl)
			break
		}
	}
}
