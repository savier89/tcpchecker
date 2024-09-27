package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
)

var (
	server_addr string
	server_port int
)

func main() {

	flag.StringVar(&server_addr, "s", "0.0.0.0", "Server ip address or hostname for listen.")
	flag.IntVar(&server_port, "p", 8080, "Server TCP port for listen.")

	flag.Parse()

	server, err := net.Listen("tcp", *&server_addr+":"+strconv.Itoa(*&server_port))
	if err != nil {
		log.Fatalln(err)
	}
	defer server.Close()

	log.Println("Server is running on:" + *&server_addr + ":" + strconv.Itoa(*&server_port))

	for {
		conn, err := server.Accept()
		fmt.Printf("New connection from addr: %s\n", conn.RemoteAddr().String())
		if err != nil {
			log.Println("Failed to accept conn.", err)
			continue
		}

		go func(conn net.Conn) {
			defer func() {
				conn.Close()
			}()

			clientAddr := conn.RemoteAddr().String()
			buf := make([]byte, 1024)
			for {
				n, _ := conn.Read(buf)
				if io.EOF == err {
					fmt.Println("Connecton EOF")
					break
				}
				if n > 0 {
					fmt.Printf("Client %s Data: %s \n", clientAddr, buf[:n])
				}
			}
			io.Copy(conn, conn)
		}(conn)
	}
}
