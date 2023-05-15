package room

import (
	"fmt"
	"movie-sync-server/entities"
	"strings"

	socketio "github.com/googollee/go-socket.io"
	"github.com/sirupsen/logrus"
)

func SetUrlEndpoint(s socketio.Conn, msg string) {
	server := entities.GetServer()
	logrus.Println("setUrl:", msg)
	Splitted := strings.Split(msg, ":::")
	room, showName, url := Splitted[0], Splitted[1], Splitted[2]
	for _, r := range entities.Cinema {
		if (*r).Name() == room {
			// eurl := video.GetUrl(url)
			// if eurl != "" {
			// 	logrus.Infof("extract SetUrl: %s", url)
			// 	(*r).SetUrl(eurl)
			// 	(*r).Broadcast("setUrl", fmt.Sprintf("%s:::%s", showName+"admin", eurl), server)
			// 	break
			// }
			(*r).SetUrl(url)
			(*r).Broadcast("setUrl", fmt.Sprintf("%s:::%s", showName, url), server)
			break
		}
	}
}
