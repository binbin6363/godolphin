package main

import (
	"fmt"
	"log"
	"net"
	"time"
	"strconv"

	"github.com/gansidui/gotcp/examples/echo"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "139.196.42.222:8989")
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	echoProtocol := &echo.EchoProtocol{}

	now := time.Now().UnixNano()
	// ping <--> pong
	for i := 0; i < 30000; i++ {
		// write
		conn.Write(echo.NewEchoPacket([]byte("hello - "+strconv.Itoa(i)), false).Serialize())

		// read
		p, err := echoProtocol.ReadPacket(conn)
		if err == nil {
			echoPacket := p.(*echo.EchoPacket)
			fmt.Printf("Server reply:[%v] [%v]\n", echoPacket.GetLength(), string(echoPacket.GetBody()))
		}

//		time.Sleep(2 * time.Second)
	}
    now2 := time.Now().UnixNano()
	fmt.Printf("cost time:%dms\n", (now2-now)/1000000)
	conn.Close()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
