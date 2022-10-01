package main

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	c    chan string
	conn net.Conn

	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String() //string是里面的地址
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		c:      make(chan string),
		conn:   conn,
		server: server,
	}
	//启动监听
	go user.ListenMessage()

	return user
}

func (user *User) ListenMessage() {
	for {
		msg := <-user.c
		_, err := user.conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("user :", user.Name, "conn.Write eer :", err)
			return
		}
	}
}

func (user *User) Online() {
	//onlineMap
	user.server.mapLock.Lock()
	user.server.OnlineMap[user.Name] = user
	user.server.mapLock.Unlock()
	//广播
	user.server.BroadCast(user, "已上线")
}

func (user *User) Offline() {
	//onlineMap
	user.server.mapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.mapLock.Unlock()
	//广播
	user.server.BroadCast(user, "已下线")
}

func (user *User) DoMessage(msg string) {
	if msg == "who" {
		//查询在线用户
		user.server.mapLock.Lock()
		for _, u := range user.server.OnlineMap {
			onlineMsg := u.Name + ":" + "在线.."
			user.c <- onlineMsg
		}
		user.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		//定义更名格式为rename|
		newName := msg[7:]
		if _, ok := user.server.OnlineMap[newName]; ok {
			user.c <- "当前用户名已被使用"
		} else {
			//改map
			user.server.mapLock.Lock()
			delete(user.server.OnlineMap, user.Name)
			user.server.OnlineMap[newName] = user
			user.server.mapLock.Unlock()
			//改自己
			user.Name = newName

			//end
			user.c <- "用户名更正成功,您现在的用户名是 : " + user.Name
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		//消息格式：to|user.Name|内容
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			user.SendMsg("消息格式不正确，消息格式：to|Name|内容")
			return
		}
		remoteUser, ok := user.server.OnlineMap[remoteName]
		if !ok {
			user.SendMsg("对方用户不存在")
			return
		}
		content := strings.Split(msg, "|")[2]
		if content == "" {
			user.SendMsg("空消息，请重发")
			return
		}
		remoteUser.SendMsg("[" + user.Name + "]向你私聊:" + content)
		user.SendMsg("向" + "[" + remoteUser.Name + "]私聊:" + content)
	} else {
		user.server.BroadCast(user, msg)
	}
}
func (user *User) SendMsg(msg string) {
	user.c <- msg
}
