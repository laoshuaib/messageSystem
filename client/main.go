package main

import (
	"flag"
	"fmt"
)

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置目标服务器的IP地址，默认是127.0.0.1")
	flag.IntVar(&serverPort, "port", 8000, "设置目标服务器的端口，默认8000")
}
func main() {
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("服务器链接失败")
		return
	}
	fmt.Println("服务器链接成功")
	//业务
	client.Run()
}
