package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int //当前状态
}

func NewClient(serverIp string, serverPort int) *Client {

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}

	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		conn:       conn,
		flag:       999,
		//Name  不初始化，服务器端将初始化为地址
	}
	return client
}

func (client *Client) dealRes() {
	_, err := io.Copy(os.Stdout, client.conn)
	if err != nil {
		fmt.Println("输出错误")
		return
	}
}

func (client *Client) Run() {
	go client.dealRes()
	//等待初始反馈
	time.Sleep(time.Second)
	for client.flag != 0 {
		for client.menu() != true {
			//阻塞
		}
		switch client.flag {
		case 1:
			client.publicChat()
			break
		case 2:
			client.privateChat()
			break
		case 3:
			client.rename()
			break
		}
	}
}

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	_, err := fmt.Scanln(&flag)
	if err != nil {
		fmt.Println("获取输出错误")
		return false
	}

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("请输入合法范围内的数字")
		return false
	}
}

func (client *Client) rename() bool {
	fmt.Println("请输入新用户名")
	_, err := fmt.Scanln(&client.Name)
	if err != nil {
		fmt.Println("获取输入错误（rename）")
		return false
	}
	sendMsg := "rename|" + client.Name + "\n"
	_, err = client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("发送消息错误（rename）")
		return false
	}
	//等待结果
	time.Sleep(time.Second)
	return true
}

func (client *Client) publicChat() {
	fmt.Println(">>>已进入公聊频道,输入exit退出")

	var chatMsg string
	_, err := fmt.Scanln(&chatMsg)
	if err != nil {
		fmt.Println("获取输入错误(publicChat)")
		return
	}
	for chatMsg != "exit" {
		//发给服务器
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("消息发送错误（publicChat）")
				return
			}
		}
		chatMsg = ""
		_, err := fmt.Scanln(&chatMsg)
		if err != nil {
			fmt.Println("获取输入错误(publicChat)")
			return
		}
	}
}

func (client *Client) privateChat() {
	fmt.Println(">>>已进入私聊频道,输入exit退出")
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("消息发送错误（privateChat）:", err)
		return
	}
	//等待返回
	time.Sleep(time.Second)

	fmt.Println(">>>请输入聊天对象用户名")
	var remoteName string
	var chatMsg string
	_, err = fmt.Scanln(&remoteName)
	if err != nil {
		fmt.Println("获取输入错误（privateChat）:", err)
		return
	}
	for remoteName != "exit" {
		fmt.Println(">>>请输入聊天内容")
		for chatMsg != "exit" {
			//发给服务器
			if len(chatMsg) != 0 {
				sendMsg = "to|" + remoteName + "|" + chatMsg + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("消息发送错误（privateChat）")
					return
				}
			}
			chatMsg = ""
			_, err := fmt.Scanln(&chatMsg)
			if err != nil {
				fmt.Println("获取输入错误(privateChat)")
				return
			}
		}
	}

}
