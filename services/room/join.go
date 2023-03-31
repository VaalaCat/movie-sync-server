package room

import (
	"log"
	"movie-sync-server/entities"
	"strings"

	socketio "github.com/googollee/go-socket.io"
)

func JoinEndpoint(s socketio.Conn, msg string) {
	server := entities.GetServer()
	log.Println("join:", msg)
	Splitted := strings.Split(msg, ":::")
	room, showName := Splitted[0], Splitted[1]
	username := s.ID()
	//首先判断当前用户是否想要加入已有的房间，如果房间不存在则新建房间
	joined := false
	var joinedRoom entities.Room
	for _, r := range entities.Cinema {
		if (*r).Name() == room {
			joinedRoom = *r
			var newUser entities.User
			newUser = new(entities.UserImpl)
			newUser.SetName(username)
			newUser.SetSocket(&s)
			newUser.SetUsername(showName)
			(*r).AddUser(&newUser, &s, server)
			joined = true
			break
		}
	}
	if joined == false {
		var newRoom entities.Room
		newRoom = new(entities.RoomImpl)
		newRoom.SetName(room)
		newRoom.InitUsers()
		var newUser entities.User
		newUser = new(entities.UserImpl)
		newUser.SetName(username)
		newUser.SetUsername(showName)
		newUser.SetSocket(&s)
		newRoom.AddUser(&newUser, &s, server)
		joinedRoom = newRoom
		entities.Cinema = append(entities.Cinema, &newRoom)
	}
	joinedRoom.Broadcast("join", showName, server)
}
