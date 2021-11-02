package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strconv"
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
	sync.Add(1)
	go ListenLocalInput(conn)
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
		// 添加字符串标识位
		flag := "str_" + strconv.Itoa(len(s))
		// 发送字符串长度
		_, err := conn.Write([]byte(flag))
		if err != nil {
			log.Printf("发送错误：%v\n", err)
			continue
		}
		// 读取服务器返回确认数据
		var buf [1024]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			log.Printf("发送错误：%v\n", err)
			continue
		}
		str := string(buf[:n])
		if str == "str_ack" {
			// 开始发送字符
			_, err := conn.Write([]byte(s))
			if err != nil {
				log.Printf("发送错误：%v\n", err)
				continue
			}
			buf := make([]byte, 1024)
			n, err := conn.Read(buf[:])
			if err != nil {
				log.Printf("接收错误：%v\n", err)
				continue
			}
			str = string(buf[:n])
			if str == "str_fin" {
				continue
			}
		}
		log.Printf("服务器错误！")
	}
}
