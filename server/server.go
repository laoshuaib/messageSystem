package main

import (
	"fmt"
	"io"
	"net"
	"runtime"
	"sync"
	"time"
)

const BUFFER = 4096
const TIME = 60 * 5 // 活跃时间

type Server struct {
	Ip   string
	Port int

	//在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//消息广播channel
	Message chan string
}

func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

func (server *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println("server-Start net.Listen err:", err)
		return
	}
	//close listen socket
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("server-Start net.Close err:", err)
			return
		}
	}(listener)

	//启动消息监听
	go server.ListenMessage()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("server-Start net.Accept err:", err)
			continue
		}

		//do handler
		go server.Handler(conn)
	}

}

func (server *Server) Handler(conn net.Conn) {
	//业务
	user := NewUser(conn, server)
	user.Online()
	//是否活跃
	isLive := make(chan bool)
	//接受客户端消息
	go func() {
		buf := make([]byte, BUFFER)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err: ", err)
				return
			}
			//处理消息
			msg := string(buf[:n-1])
			user.DoMessage(msg)
			//记录活跃
			isLive <- true
		}
	}()
	//handler阻塞
	for {
		select {
		case <-isLive:
			//活跃，重置定时器
			//DoNothing, 让select重新更新下面的计时器
		case <-time.After(time.Second * TIME):
			//超时
			user.SendMsg("长时间不活跃，已被踢出")
			time.Sleep(time.Second)
			//销毁资源
			user.Offline()
			err := conn.Close()
			if err != nil {
				fmt.Println(user.Name, "已超时但未正确退出")
			}
			runtime.Goexit()
		}
	}

}

func (server *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	server.Message <- sendMsg
}

func (server *Server) ListenMessage() {
	for {
		msg := <-server.Message
		server.mapLock.Lock()
		for _, user := range server.OnlineMap {
			user.SendMsg(msg)
		}
		server.mapLock.Unlock()
	}
}
