package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
	sync2 "sync"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

var sync sync2.WaitGroup

// 一个点对点文件、文字即时同步GUI工具-客户端
func main() {
	// 1、与服务端建立连接
	conn, err := net.Dial("tcp", "127.0.0.1:9001")
	if err != nil {
		log.Printf("conn server failed, err:%v\n", err)
		return
	}
	defer func(Conn net.Conn) {
		err := Conn.Close()
		if err != nil {
			log.Println("server 连接关闭!")
			return
		}
	}(conn)
	sync.Add(1)
	go ListenLocalInput(conn)
	go ListenRemoteInput(conn)
	sync.Wait()
}

func ListenLocalInput(conn net.Conn) {
	// 2、使用 conn 连接进行数据的发送和接收
	for {
		input := bufio.NewReader(os.Stdin)
		s, _ := input.ReadString('\n')
		s = strings.TrimSpace(s)
		if strings.ToUpper(s) == "Q" {
			err := conn.Close()
			if err != nil {
				return
			}
			os.Exit(0)
		}
		if len(s) <= 0 {
			continue
		}
		// 发送字符串长度
		_, err := conn.Write([]byte(s))
		if err != nil {
			log.Printf("发送错误：%v\n", err)
			continue
		}
	}
}

func ListenRemoteInput(conn net.Conn) {
	for {
		var buf = make([]byte, 10240)
		input := bufio.NewReader(conn)
		s, err := input.Read(buf[:])
		if err != nil {
			if err.Error() == "EOF" {
				log.Printf("server 退出!")
				os.Exit(0)
			}
			log.Println("接收失败! err: " + err.Error())
		}
		str := string(buf[:s])
		log.Printf("server --> " + str)
	}
}
