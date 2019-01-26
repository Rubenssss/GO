package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type client struct {
	channel chan<- string 
	name    string
}

var (
	entering = make(chan *client)
	leaving  = make(chan *client)
	messages = make(chan string) 
)

func broadcaster() {
	clients := make(map[string]*client) 
	for {
		select {
		case msg := <-messages:
			

			for _, cli := range clients {
				cli.channel <- msg
			}

		case cli := <-entering:
			go giveAllClients(cli.channel, clients)
			clients[cli.name] = cli

		case cli := <-leaving:
			delete(clients, cli.name)
			close(cli.channel)
		}
	}
}

func giveAllClients(channel chan<- string, clients map[string]*client) {
	if len(clients) > 1 {
		channel <- "All clients:"
		for _, cli := range clients {
			channel <- cli.name
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string) 
	me := &client{channel: ch, name: conn.RemoteAddr().String()}
	go clientWriter(conn, ch)

	me.channel <- "You are " + me.name
	messages <- me.name + " has arrived"
	entering <- me

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- me.name + ": " + input.Text()
	}
	


	leaving <- me
	messages <- me.name + " has left"
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




