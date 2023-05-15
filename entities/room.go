package entities

import (
	"time"

	socketio "github.com/googollee/go-socket.io"
)

type Room interface {
	Name() string
	SetName(name string)
	SetUrl(url string)
	GetUrl() string
	GetMinTime() int
	GetMaxTime() int
	AddUser(user *User, socket *socketio.Conn, server *socketio.Server)
	RemoveUser(username string, server *socketio.Server)
	Broadcast(event string, message string, server *socketio.Server)
	GetUsers() []*User
	GetUser(name string) *User
	InitUsers()
	Refresh(server *socketio.Server)
}

type RoomImpl struct {
	name     string
	url      string
	users    map[string]*User
	lastPlay time.Time
	lastStop time.Time
}

func (r *RoomImpl) Name() string {
	return r.name
}

func (r *RoomImpl) SetName(name string) {
	r.name = name
	r.lastPlay = time.Now()
	r.lastStop = time.Now()
}

func (r *RoomImpl) SetUrl(url string) {
	r.url = url
}

func (r *RoomImpl) GetUrl() string {
	return r.url
}

func (r *RoomImpl) GetMinTime() int {
	var times []int
	for _, user := range r.users {
		times = append(times, (*user).GetTime())
	}
	min := times[0]
	for _, time := range times {
		if time < min {
			min = time
		}
	}
	return min
}

func (r *RoomImpl) GetMaxTime() int {
	var times []int
	for _, user := range r.users {
		times = append(times, (*user).GetTime())
	}
	max := times[0]
	for _, time := range times {
		if time > max {
			max = time
		}
	}
	return max
}

func (r *RoomImpl) AddUser(user *User, socket *socketio.Conn, server *socketio.Server) {
	r.users[(*socket).ID()] = user
	server.JoinRoom("/", r.name, *socket)
}

func (r *RoomImpl) RemoveUser(username string, server *socketio.Server) {
	tmpUser, ok := r.users[username]
	if ok {
		server.LeaveRoom("/", r.name, *(*tmpUser).GetSocket())
		delete(r.users, username)
	}
}

func (r *RoomImpl) Broadcast(event string, message string, server *socketio.Server) {
	if event == "play" {
		r.lastPlay = time.Now()
	}
	if event == "pause" {
		r.lastStop = time.Now()
	}
	if time.Duration(r.lastStop.Sub(r.lastPlay)).Abs() < 300*time.Microsecond {
		return
	}
	server.BroadcastToRoom("/", r.name, event, message)
}

func (r *RoomImpl) InitUsers() {
	r.users = make(map[string]*User)
}

func (r *RoomImpl) Refresh(server *socketio.Server) {
	server.BroadcastToRoom("/", r.name, "refresh", "")
}

func (r *RoomImpl) GetUsers() []*User {
	var users []*User
	for _, user := range r.users {
		users = append(users, user)
	}
	return users
}

func (r *RoomImpl) GetUser(name string) *User {
	user, ok := r.users[name]
	if ok {
		return user
	}
	return nil
}
