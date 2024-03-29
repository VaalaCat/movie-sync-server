package room

import (
	"movie-sync-server/entities"
	"strings"

	socketio "github.com/googollee/go-socket.io"
	"github.com/sirupsen/logrus"
)

func GetUsersEndpoint(s socketio.Conn, msg string) {
	logrus.Println("getUsers:", msg)
	Splitted := strings.Split(msg, ":::")
	room, _ := Splitted[0], Splitted[1]
	hasRoom := false
	for _, r := range entities.Cinema {
		if (*r).Name() == room && len((*r).GetUsers()) > 0 {
			usernames := ""
			for _, u := range (*r).GetUsers() {
				usernames += (*u).GetUserName() + ","
			}
			s.Emit("allUsers", usernames)
			hasRoom = true
			break
		}
	}
	if !hasRoom {
		s.Emit("Null", "no room")
	}
}
