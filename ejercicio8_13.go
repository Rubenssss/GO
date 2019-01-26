package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

type client chan<- string 

var (
	entering             = make(chan client)
	leaving              = make(chan client)
	messages             = make(chan string) 
	registeringAddress   = make(chan string)
	unregisteringAddress = make(chan string)
)

func broadcaster() {
	clients := make(map[client]bool) 
	addresses := make(map[string]bool)
	for {
		select {
		case msg := <-messages:
			

			for cli := range clients {
				cli <- msg
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)

		case address := <-registeringAddress:
			addresses[address] = true
			allClients := "All clients:"
			for addr := range addresses {
				allClients = fmt.Sprintf("%s\n%s", allClients, addr)
			}
			go func() { messages <- allClients }() 
		case address := <-unregisteringAddress:
			delete(addresses, address)
			allClients := "All clients:"
			for addr := range addresses {
				allClients = fmt.Sprintf("%s\n%s", allClients, addr)
			}
			go func() { messages <- allClients }() 
		}

	}
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
				messages <- msg
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

func handleConn(conn net.Conn) {
	ch := make(chan string) 
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- ch
	registeringAddress <- who

	notIdleCh := make(chan bool)
	go countIdleTime(conn, notIdleCh)

	input := bufio.NewScanner(conn)
	for input.Scan() {
		notIdleCh <- true
		messages <- who + ": " + input.Text()
	}
	

	leaving <- ch
	messages <- who + " has left"
	unregisteringAddress <- who
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) 
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}


