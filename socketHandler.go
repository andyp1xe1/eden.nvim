package main

import (
	"bufio"
	"log"
	"net"
	"os"
)


type SocketHandler struct {
	path string
	msgChan chan string
}

func NewSocketHandler(path string) SocketHandler {
	return SocketHandler{path, make(chan string, 24)}
}

func (s SocketHandler) Listen() {
	if err := os.RemoveAll(s.path); err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("unix", s.path)
	if err != nil {
		log.Fatal("Error creating Unix socket:", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		go func() {
			defer conn.Close()
			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				s.msgChan <- scanner.Text()
			}
			if err := scanner.Err(); err != nil {
				log.Println("Error reading from socket:", err)
			}
		}()
	}
}



