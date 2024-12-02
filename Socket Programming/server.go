/*****************************************************************************
 * server.go
 * Name: Nicholas Nguyen, Lane Bryant
 * NetId: nicholas.nguyen, lcbyrant
 *****************************************************************************/

package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

const RECV_BUFFER_SIZE = 2048
const IP = "127.0.0.1"

/* The server function takes a port as an argument, then creates a listener socket
 * to receive client payloads indefinitely.
 */
func server(server_port string) {
	addr := IP + ":" + server_port

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
		}

		handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, RECV_BUFFER_SIZE)

	for {
		n, err := conn.Read(buf)

		if err == io.EOF {
			break
		}
		
		if err != nil {
			log.Print(err)
		}

		fmt.Print(string(buf[:n]))
	}
}

// Main parses command-line arguments and calls server function
func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: ./server [server port]")
	}
	server_port := os.Args[1]
	server(server_port)
}
