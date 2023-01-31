package main

import (
	"fmt"
	"io"
	"net"
	"os"
	// Uncomment this block to pass the first stage
	// "net"
	// "os"
)

var (
	buffSize = 1024
)

func handleConnection(tNo int, conn net.Conn) {
	for {
		buff := make([]byte, buffSize)
		n, err := conn.Read(buff)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("[WORKER - CHILD %d] Error when reading client request, err: %s\n", tNo, err)
		}

		if n > buffSize {
			fmt.Printf("[WORKER - CHILD %d] Got buffer bigger than expected\n", tNo)
		}

		fmt.Printf("[WORKER - CHILD %d] Reading request, got: %#v\n", tNo, string(buff[:n]))
		_, err = conn.Write([]byte("+PONG\r\n"))
		if err != nil {
			fmt.Printf("[WORKER - CHILD %d] Error: %s\n", tNo, err)
		}
	}

	fmt.Printf("[WORKER - CHILD %d] Closing connection\n", tNo)
	conn.Close()
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	totalThread := 5
	fmt.Println("Logs from your program will appear here!")
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()

	stopChan := make(chan bool)
	for i := 0; i < totalThread; i++ {
		go func(tNo int) {
			for {
				fmt.Printf("[WORKER %d] Waiting connection\n", tNo)
				conn, err := l.Accept()
				if err != nil {
					fmt.Printf("[WORKER %d] Error accepting connection: %s", tNo, err)
					os.Exit(1)
				}

				fmt.Printf("[WORKER %d] Accepting connection\n", tNo)

				go handleConnection(tNo, conn)
			}
		}(i)
	}

	<-stopChan
}
