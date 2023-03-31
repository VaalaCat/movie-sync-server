package entities

import (
	socketio "github.com/googollee/go-socket.io"
)

type User interface {
	Name() string
	Send(event string, message string)
	GetSocket() *socketio.Conn
	GetTime() int
	SetTime(time int)
	SetName(name string)
	SetSocket(socket *socketio.Conn)
	SetUsername(username string)
	GetUserName() string
}

type UserImpl struct {
	name     string
	time     int
	username string
	socket   *socketio.Conn
}

func (u *UserImpl) Name() string {
	return u.name
}

func (u *UserImpl) Send(event string, message string) {
	(*(u.socket)).Emit(event, message)
}

func (u *UserImpl) SetTime(time int) {
	u.time = time
}

func (u *UserImpl) GetTime() int {
	return u.time
}

func (u *UserImpl) GetSocket() *socketio.Conn {
	return u.socket
}

func (u *UserImpl) SetName(name string) {
	u.name = name
}

func (u *UserImpl) SetSocket(socket *socketio.Conn) {
	u.socket = socket
}

func (u *UserImpl) SetUsername(username string) {
	u.username = username
}

func (u *UserImpl) GetUserName() string {
	return u.username
}
