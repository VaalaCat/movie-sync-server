package room

import (
	"log"
	"movie-sync-server/entities"

	socketio "github.com/googollee/go-socket.io"
)

func DisconnectEndpoint(s socketio.Conn, msg string) {
	server := entities.GetServer()
	log.Println("closed", msg)
	if len(entities.Cinema) > 0 {
		for i, r := range entities.Cinema {
			(*r).RemoveUser(s.ID(), server)
			if len((*r).GetUsers()) == 0 {
				entities.Cinema = append(entities.Cinema[:i], entities.Cinema[i+1:]...)
			} else {
				(*r).Broadcast("leaveRoom", s.ID(), server)
			}
		}
	}
}
