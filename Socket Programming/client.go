/*****************************************************************************
 * client.go
 * Name: Nicholas Nguyen, Lane Bryant
 * NetId: nicholas.nguyen, lcbyrant
 *****************************************************************************/

package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
)

const SEND_BUFFER_SIZE = 2048

/* The client function requires a server IP and port as arguments. Using these arguments,
 * client uses a TCP socket connection to send a payload to the server from Stdin.
 */
func client(server_ip string, server_port string) {
	serverAddr := server_ip + ":" + server_port
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	buf := make([]byte, SEND_BUFFER_SIZE)

	for {
		n, err := reader.Read(buf)

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		_, err = conn.Write(buf[:n])
		if err != nil {
			log.Fatal(err)
		}
	}
}

/* Main parses command-line arguments and calls client function
 * Two ways to run the client:
 * if there are 2 args, the client should accept them as the ip + port
 * then wait for input
 * 3 args: (server IP, server Port, file)
 */
func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: ./client [server IP] [server port] < [message file]")
	}

	server_ip := os.Args[1]
	server_port := os.Args[2]
	client(server_ip, server_port)
}
