package main

import "flag"

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器的IP地址，默认是127.0.0.1")
	flag.IntVar(&serverPort, "port", 8000, "设置服务器的端口，默认8000")
}
func main() {
	//命令行解析
	flag.Parse()

	server := NewServer(serverIp, serverPort)
	server.Start()
}
