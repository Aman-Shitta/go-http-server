package server

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

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

	fmt.Println("urlPath :: ", urlPath)
	if urlPath == "/" {
		r = fmt.Sprintf("HTTP/%s 200 OK\r\n\r\n", s.Version)

	} else if strings.HasPrefix(urlPath, "/echo") {
		echo := strings.Split(urlPath, "/echo")
		fmt.Println("Echo :: ", echo)

		var val strings.Builder

		for _, i := range echo {
			if len(strings.Trim(i, "/ ")) > 0 {
				val.WriteString(strings.Trim(i, "/ "))
			}

		}

		r = fmt.Sprintf("HTTP/%s 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", s.Version, len(val.String()), val.String())
	} else {
		r = fmt.Sprintf("HTTP/%s 404 Not Found\r\n\r\n", s.Version)
	}

	return []byte(r)
}
