package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func countIdleTime(conn net.Conn, notIdleCh <-chan bool) {
	ticker := time.NewTicker(time.Second)
	counter := 0
	max := 20 
	for {
		select {
		case <-ticker.C:
			counter++
			if counter == max {
				msg := conn.RemoteAddr().String() + " idle too long. Kicked out."
				fmt.Println(msg)
				fmt.Fprintln(conn, msg) 
				ticker.Stop()
				conn.Close()
				return
			}
		case <-notIdleCh:
			counter = 0
		}
	}
}

func handleConn(c net.Conn) {

	input := bufio.NewScanner(c)
	notIdleCh := make(chan bool)
	go countIdleTime(c, notIdleCh)

	for input.Scan() {
		notIdleCh <- true
		go echo(c, input.Text(), 1*time.Second)
	}

	c.Close()

}

func main() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err) 
			continue
		}
		go handleConn(conn)
	}
}


