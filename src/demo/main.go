package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"net"
	"runtime"
	"time"
)

func startRoutine(i int) {
	send := make(chan int)
	fmt.Println("start routine, ", i)
	send <- i

	// send data
	conn, err := net.DialTimeout("tcp", "127.0.0.1:9000", 10*time.Second)
	if err != nil {
		return
	}
	conn.Write(([]byte)("hi, nice to meet you!"))
	data := make([]byte, 20)

	conn.SetReadDeadline(time.Now().Add(200 * time.Second))
	conn.Read(data)

	conn.Close()

	//time.Sleep(1 * time.Second)
	fmt.Println("done start routine, ", <-send)
}

func makeList() chan int {
	chanLish := make(chan int, 2)
	go func() {
		for i := 0; i < 20; i++ {
			chanLish <- i * 10
		}
		close(chanLish)
	}()
	return chanLish
}

func main() {

	runtime.GOMAXPROCS(2)
	log := logs.NewLogger(1000)

	log.SetLogger("file", `{"filename":"test.log"}`)
	log.Info("start demo pipe test, this is a log file.")
	log.Warn("start demo pipe test, this is a log file.")
	log.Trace("start demo pipe test, this is a log file.")
	log.Debug("start demo pipe test, this is a log file.")
	log.Critical("start demo pipe test, this is a log file.")
	for i := 0; i < 10000; i++ {
		nowTime := time.Now()
		log.Info("time now:%v-----ns:%d-------index:%d", nowTime, nowTime.Nanosecond(), i)
	}
	fmt.Println("start demo pipe test.")

	//cnt := 0
	//sendlist := make(chan int, 2)

	for item := range makeList() {
		fmt.Println("show list item: ", item)
	}

	fmt.Println("done demo test routine, ")
	return
}
