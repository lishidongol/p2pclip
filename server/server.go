package main

import (
	"bufio"
	"log"
	"net"
	"p2pclip/common"
	"strconv"
	"strings"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

// 一个点对点文件、文字即时同步GUI工具-服务端
func main() {
	// 创建socket服务端
	listen, err := net.Listen("tcp", "127.0.0.1:9001")
	if err != nil {
		panic(err)
	}
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
		go Process(c)
	}
}

func Process(client common.P2pclient) {
	log.Println(client.Name + " 连接!")
	// 处理完关闭连接
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(client.Conn)
	// 针对当前连接做发送和接受操作
	// 标示是否正在传输文件
	//isTransFile := false
	for {
		reader := bufio.NewReader(client.Conn)
		var buf [1024]byte
		n, err := reader.Read(buf[:])
		if err != nil {
			if "EOF" == err.Error() {
				log.Println("客户端 " + client.Name + " 断开连接！")
				return
			}
			log.Printf("read from conn failed, err:%v\n\n", err)
		}
		recv := string(buf[:n])
		if recv == "ok" {
			log.Println("数据接收完毕！")
		} else if strings.HasPrefix(recv, "str_") {
			split := strings.Split(recv, "_")
			if len(split) != 2 {
				log.Printf("格式错误！")
				continue
			}
			// 字符长度
			length, _ := strconv.Atoi(split[1])
			// 将接受到的数据返回给客户端
			_, err := client.Conn.Write([]byte("str_ack"))
			if err != nil {
				log.Printf("发送字符错误: %v\n", err)
				continue
			}
			if length <= 0 {
				log.Println("空数据，丢弃")
				continue
			}
			var buf = make([]byte, length)
			reader := bufio.NewReader(client.Conn)
			_, err = reader.Read(buf[:])
			if err != nil {
				log.Printf("接收数据错误：%v\n", err)
				continue
			}
			log.Printf(client.Name + " --> " + string(buf[:length]))

			// 返回响应
			_, err = client.Conn.Write([]byte("str_fin"))
			if err != nil {
				log.Printf("接收数据错误：%v\n", err)
				continue
			}
		}
	}
}
