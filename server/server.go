package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"p2pclip/common"
	"strings"
	sync2 "sync"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

var sync sync2.WaitGroup
var clientMap = make(map[string]common.P2pclient)

// 一个点对点文件、文字即时同步GUI工具-服务端
func main() {
	// 创建socket服务端
	listen, err := net.Listen("tcp", "127.0.0.1:9001")
	if err != nil {
		panic(err)
	}
	go ListenLocalInput()
	for {
		accept, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		// 随机给定一个名称
		c := common.P2pclient{}
		c.Name = accept.RemoteAddr().String()
		c.Conn = accept
		clientMap[c.Name] = c
		go Process(c)
	}
}

func Process(client common.P2pclient) {
	log.Println(client.Name + " 连接!")
	// 处理完关闭连接
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("连接关闭! err: " + err.Error())
			return
		}
	}(client.Conn)

	ListenRemoteInput(client)
}

func ListenLocalInput() {
	// 2、使用 conn 连接进行数据的发送和接收
	for {
		input := bufio.NewReader(os.Stdin)
		s, _ := input.ReadString('\n')
		s = strings.TrimSpace(s)
		if strings.ToUpper(s) == "Q" {
			os.Exit(0)
		}
		if len(s) <= 0 {
			continue
		}
		if len(clientMap) > 0 {
			// 循环给客户端发送数据
			for _, pclient := range clientMap {
				_, err := pclient.Conn.Write([]byte(s))
				if err != nil {
					log.Println(pclient.Name + " 发送失败! err: " + err.Error())
				}
			}
		}
	}
}

func ListenRemoteInput(client common.P2pclient) {
	for {
		var buf = make([]byte, 10240)
		input := bufio.NewReader(client.Conn)
		s, err := input.Read(buf[:])
		if err != nil {
			if err.Error() == "EOF" {
				// 删除客户端
				log.Println("删除客户端:" + client.Name)
				delete(clientMap, client.Name)
				break
			}
			log.Println("接收失败! err: " + err.Error())
		}
		str := string(buf[:s])
		log.Printf(client.Name + " --> " + str)
		// 转发
		if len(clientMap) > 0 {
			// 循环给客户端发送数据
			for _, pclient := range clientMap {
				if pclient.Name != client.Name {
					_, err := pclient.Conn.Write([]byte(str))
					if err != nil {
						log.Println(pclient.Name + " 发送失败! err: " + err.Error())
					}
				}
			}
		}
	}
}
