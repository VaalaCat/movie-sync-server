package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func GinMiddleware(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Request.Header.Del("Origin")
		c.Next()
	}
}

type Room interface {
	Name() string
	SetName(name string)
	SetUrl(url string)
	GetUrl() string
	GetMinTime() int
	GetMaxTime() int
	AddUser(user *User, socket *socketio.Conn)
	RemoveUser(username string)
	Broadcast(event string, message string)
	GetUsers() []*User
	GetUser(name string) *User
	InitUsers()
	Refresh()
}

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

type RoomImpl struct {
	name  string
	url   string
	users map[string]*User
}

type UserImpl struct {
	name     string
	time     int
	username string
	socket   *socketio.Conn
}

func (r *RoomImpl) Name() string {
	return r.name
}

func (r *RoomImpl) SetName(name string) {
	r.name = name
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

func (r *RoomImpl) AddUser(user *User, socket *socketio.Conn) {
	r.users[(*user).Name()] = user
	server.JoinRoom("/", r.name, *socket)
}

func (r *RoomImpl) RemoveUser(username string) {
	tmpUser, ok := r.users[username]
	if ok {
		server.LeaveRoom("/", r.name, *(*tmpUser).GetSocket())
		delete(r.users, username)
	}
}

func (r *RoomImpl) Broadcast(event string, message string) {
	server.BroadcastToRoom("/", r.name, event, message)
}

func (r *RoomImpl) InitUsers() {
	r.users = make(map[string]*User)
}

func (r *RoomImpl) Refresh() {
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

var server = socketio.NewServer(nil)

func main() {
	var Cinema []*Room

	router := gin.New()
	//如果有一个客户端连接上
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())
		return nil
	})

	//连接后必须加入房间
	server.OnEvent("/", "join", func(s socketio.Conn, msg string) {
		log.Println("join:", msg)
		Spited := strings.Split(msg, ":")
		room, showName := Spited[0], Spited[1]
		username := s.ID()
		//首先判断当前用户是否想要加入已有的房间，如果房间不存在则新建房间
		joined := false
		var joinedRoom Room
		for _, r := range Cinema {
			if (*r).Name() == room {
				joinedRoom = *r
				var newUser User
				newUser = new(UserImpl)
				newUser.SetName(username)
				newUser.SetSocket(&s)
				newUser.SetUsername(showName)
				(*r).AddUser(&newUser, &s)
				joined = true
				break
			}
		}
		if joined == false {
			var newRoom Room
			newRoom = new(RoomImpl)
			newRoom.SetName(room)
			newRoom.InitUsers()
			var newUser User
			newUser = new(UserImpl)
			newUser.SetName(username)
			newUser.SetUsername(showName)
			newUser.SetSocket(&s)
			newRoom.AddUser(&newUser, &s)
			joinedRoom = newRoom
			Cinema = append(Cinema, &newRoom)
		}
		joinedRoom.Broadcast("join", username)
	})

	// 设置房间的视频地址
	server.OnEvent("/", "setUrl", func(s socketio.Conn, msg string) {
		log.Println("setUrl:", msg)
		Spited := strings.Split(msg, ":::")
		room, url := Spited[0], Spited[1]
		for _, r := range Cinema {
			if (*r).Name() == room {
				(*r).SetUrl(url)
				break
			}
		}
	})

	// 如果用户回复了时间则将每个用户的时间更新到对应的用户并广播最大值和最小值
	server.OnEvent("/", "time", func(s socketio.Conn, msg string) {
		log.Println("time:", msg)
		Spited := strings.Split(msg, ":")
		room, _, time := Spited[0], Spited[1], Spited[2]
		username := s.ID()
		for _, r := range Cinema {
			if (*r).Name() == room {
				for _, u := range (*r).GetUsers() {
					if (*u).Name() == username {
						tmpTime := strings.Split(time, ".")
						timeNum, err := strconv.Atoi(tmpTime[0])
						if err != nil {
							log.Println("time is not a number")
						} else {
							(*u).SetTime(timeNum)
						}
					}
				}
				(*r).Broadcast("sync", fmt.Sprintf("%d:%d", (*r).GetMinTime(), (*r).GetMaxTime()))
				break
			}
		}
	})

	// 如果有用户发出getTime请求则开始广播getTime请求
	server.OnEvent("/", "getTime", func(s socketio.Conn, msg string) {
		log.Println("getTime:", msg)
		Spited := strings.Split(msg, ":")
		room, _ := Spited[0], Spited[1]
		username := s.ID()
		for _, r := range Cinema {
			if (*r).Name() == room {
				tmpUser := (*r).GetUser(username)
				if tmpUser != nil {
					(*r).Broadcast("getTime", (*tmpUser).GetUserName())
					break
				}
			}
		}
	})

	// 如果有用户发出setTime请求则使用sync消息将所有用户的时间同步
	server.OnEvent("/", "setTime", func(s socketio.Conn, msg string) {
		log.Println("setTime:", msg)
		Spited := strings.Split(msg, ":")
		room, _, time := Spited[0], Spited[1], Spited[2]
		for _, r := range Cinema {
			if (*r).Name() == room {
				(*r).Broadcast("setTime", time)
				for _, u := range (*r).GetUsers() {
					tmpTime := strings.Split(time, ".")
					timeNum, err := strconv.Atoi(tmpTime[0])
					if err != nil {
						log.Println("time is not a number")
					} else {
						(*u).SetTime(timeNum)
					}
				}
				break
			}
		}
	})

	// 用户发出getUsers请求则返回所有的用户名
	server.OnEvent("/", "getUsers", func(s socketio.Conn, msg string) {
		log.Println("getUsers:", msg)
		Spited := strings.Split(msg, ":")
		room, _ := Spited[0], Spited[1]
		hasRoom := false
		for _, r := range Cinema {
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
	})

	// 用户发出getUrl请求则返回url给对应用户
	server.OnEvent("/", "getUrl", func(s socketio.Conn, msg string) {
		log.Println("getUrl:", msg)
		room := msg
		for _, r := range Cinema {
			if (*r).Name() == room {
				tmpUrl := (*r).GetUrl()
				s.Emit("setUrl", tmpUrl)
				break
			}
		}
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		log.Println("closed", msg)
		if len(Cinema) > 0 {
			for i, r := range Cinema {
				(*r).RemoveUser(s.ID())
				if len((*r).GetUsers()) == 0 {
					Cinema = append(Cinema[:i], Cinema[i+1:]...)
				}
			}
		}
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer server.Close()

	// 从.env文件中读取端口号
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	allowOrigin := os.Getenv("ALLOW_ORIGIN")
	router.Use(GinMiddleware(allowOrigin))
	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/movie")
	})
	router.StaticFile("/movie", "../asset/index.html")
	router.StaticFile("/movie/login", "../asset/index.html")
	router.StaticFS("/movie/css", http.Dir("../asset/css"))
	router.StaticFS("/movie/js", http.Dir("../asset/js"))
	router.GET("/movie/room/*any", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/movie")
	})
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal("failed run app: ", err)
	}
}
