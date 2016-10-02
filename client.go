package main

import (
	"net"
	"fmt"
	"bufio"
	"time"
	"os"
	"strings"
)

func send_heartbeat(conn net.Conn){
	for{
		conn.Write([]byte("heartbeat_ping#"))
		time.Sleep(8*time.Second)
	}

}
func main() {
	name := os.Args[1]
	conn, err := net.Dial("tcp4", "127.0.0.1:8002")
	if err != nil {
		// handle error
		fmt.Println(err)
	}
	conn.Write([]byte(name))
	go send_heartbeat(conn)
	for{
		data, err := bufio.NewReader(conn).ReadString('#')
		if err != nil {
			fmt.Println("read err",err)
			return
		}
		if ! strings.EqualFold(data,"heartbeat_pong#") {
			fmt.Printf("data1:%s\n",string(data)) 
		}
	}
}
