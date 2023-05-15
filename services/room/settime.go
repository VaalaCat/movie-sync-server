package room

import (
	"fmt"
	"movie-sync-server/entities"
	"strconv"
	"strings"

	socketio "github.com/googollee/go-socket.io"
	"github.com/sirupsen/logrus"
)

func SetTimeEndpoint(s socketio.Conn, msg string) {
	server := entities.GetServer()
	logrus.Println("setTime:", msg)
	Splitted := strings.Split(msg, ":::")
	room, showName, time := Splitted[0], Splitted[1], Splitted[2]
	for _, r := range entities.Cinema {
		if (*r).Name() == room {
			(*r).Broadcast("setTime", fmt.Sprintf("%s:::%s", showName, time), server)
			for _, u := range (*r).GetUsers() {
				tmpTime := strings.Split(time, ".")
				timeNum, err := strconv.Atoi(tmpTime[0])
				if err != nil {
					logrus.Println("time is not a number")
				} else {
					(*u).SetTime(timeNum)
				}
			}
			break
		}
	}
}
