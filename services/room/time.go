package room

import (
	"log"
	"movie-sync-server/entities"
	"strconv"
	"strings"

	socketio "github.com/googollee/go-socket.io"
)

func TimeEndpoint(s socketio.Conn, msg string) {
	log.Println("time:", msg)
	Splitted := strings.Split(msg, ":::")
	room, _, time := Splitted[0], Splitted[1], Splitted[2]
	username := s.ID()
	for _, r := range entities.Cinema {
		if (*r).Name() == room {
			u := (*r).GetUser(username)
			if u != nil {
				tmpTime := strings.Split(time, ".")
				timeNum, err := strconv.Atoi(tmpTime[0])
				if err != nil {
					log.Println("time is not a number", err)
				} else {
					(*u).SetTime(timeNum)
				}
			}
			break
		}
	}
}
