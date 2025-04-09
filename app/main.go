package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/http-server-starter-go/server"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	dir := flag.String("directory", "", "directrory name")

	flag.Parse()

	// Uncomment this block to pass the first stage
	httpServer := server.NewServer(4221, *dir)

	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", httpServer.Port))

	if err != nil {
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			os.Exit(1)
		}

		if err != nil {
			fmt.Println("Err : ", err)
			os.Exit(1)
		}

		fmt.Println("Ok handling connection")
		go httpServer.HandleConnection(conn)
	}

}
