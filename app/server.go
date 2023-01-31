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

func main() {
    // You can use print statements as follows for debugging, they'll be visible when running tests.
    buffSize := 1024
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
                fmt.Printf("[WORKER %d] waiting connection\n", tNo)
                conn, err := l.Accept()
                if err != nil {
                    fmt.Println("Error accepting connection: ", err.Error())
                    os.Exit(1)
                }

                fmt.Printf("[WORKER %d] accepting connection\n", tNo)

                for {
                    buff := make([]byte, buffSize)
                    n, err := conn.Read(buff)
                    if err == io.EOF {
                        break
                    }
                    if err != nil {
                        fmt.Printf("error when reading client request, err: %s\n", err)
                    }

                    if n > buffSize {
                        fmt.Println("got buffer bigger than expected")
                    }

                    fmt.Printf("[WORKER %d] reading request, got: %#v\n", tNo, string(buff[:n]))
                    _, err = conn.Write([]byte("+PONG\r\n"))
                    if err != nil {
                        fmt.Printf("error: %s", err)
                    }
                }

                fmt.Printf("[WORKER %d] closing connection\n", tNo)
                conn.Close()
            }
        }(i)
    }

    <-stopChan
}
