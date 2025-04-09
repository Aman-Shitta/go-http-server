package server

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/http-server-starter-go/utils"
)

type HttpServer struct {
	Port    uint
	Version string
}

func NewServer(port uint) *HttpServer {
	return &HttpServer{
		Port:    port,
		Version: "1.1",
	}
}

func (s *HttpServer) HandleConnection(conn net.Conn) {

	var reader = bufio.NewReader(conn)

	data, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println("Error in reeading data : ", err)
		os.Exit(1)
	}

	verb, urlPath, err := utils.BreakRequestData(data)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var resp []byte

	fmt.Println("Request recieved :: ", verb)
	switch verb {
	case "GET":
		resp = s.getHandler(urlPath)
	default:
		fmt.Printf("%s not supported\n", verb)
		return
	}

	conn.Write(resp)

}

func (s *HttpServer) getHandler(urlPath string) []byte {

	var r string

	if urlPath == "/" {
		r = fmt.Sprintf("HTTP/%s 200 OK\r\n\r\n", s.Version)
	} else {
		r = fmt.Sprintf("HTTP/%s 404 Not Found\r\n\r\n", s.Version)
	}

	return []byte(r)
}
