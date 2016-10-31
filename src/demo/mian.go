package main

import (
	"fmt"
	"time"
	"runtime"
	"net"
)

func startRoutine(i int) {
	send := make(chan int)
	fmt.Println("start routine, ", i)
	send <- i

	// send data
	conn := net.DialTimeout("tcp", "127.0.0.1:9000", 10*time.Second)
	conn.Write(([]byte)("hi, nice to meet you!"))
	data := make([]byte, 20)

	conn.SetReadDeadline(200*time.Second)
	conn.Read(data)

	conn.Close()

	//time.Sleep(1 * time.Second)
	fmt.Println("done start routine, ", <- send)
}


func main() {

	runtime.GOMAXPROCS(2)
	fmt.Println("start demo pipe test.")

	cnt := 0
	sendlist := make(chan int, 2)

	for i := 0; i < 10; i++ {
		cnt++
		sendlist <- i
		go startRoutine(i)
	}

	for i := 0; i < 10; i++ {
		fmt.Println("recv routine, ",  <- sendlist)
	}

	fmt.Println("done demo test routine, ")
	return
}
