package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/utils"
)

type HttpServer struct {
	Port    uint
	Version string
	dir     string
}

func NewServer(port uint, dir string) *HttpServer {
	return &HttpServer{
		Port:    port,
		Version: "1.1",
		dir:     dir,
	}
}

func (s *HttpServer) HandleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Step 1: Read the request line
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("[ERROR] Failed to read request line:", err)
		return
	}
	fmt.Println("[DEBUG] Request line:", requestLine)

	verb, urlPath, err := utils.BreakRequestData(requestLine)
	if err != nil {
		fmt.Println("[ERROR] Malformed request line:", err)
		return
	}

	// Step 2: Read headers
	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		fmt.Println("[DEBUG] Header line:", line)
		if err != nil {
			fmt.Println("[ERROR] Failed to read header:", err)
			return
		}
		if line == "\r\n" {
			break // End of headers
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
		fmt.Println("[DEBUG] Headers:", headers)
	}

	// Step 3: Read body (if present)
	var body []byte
	if cl, ok := headers["Content-Length"]; ok {
		length, err := strconv.Atoi(cl)
		if err != nil {
			fmt.Println("[ERROR] Invalid Content-Length:", err)
			return
		}
		body = make([]byte, length)
		_, err = io.ReadFull(reader, body)
		if err != nil {
			fmt.Println("[ERROR] Failed to read body:", err)
			return
		}
		fmt.Println("[DEBUG] Body:", string(body))
	}

	// Step 4: Handle the request
	var resp []byte
	switch verb {
	case "GET":
		resp = s.getHandler(urlPath, headers, body)
	case "POST":
		resp = s.postHandler(urlPath, headers, body)
	default:
		resp = []byte("HTTP/1.1 405 Method Not Allowed\r\n\r\n")
	}

	conn.Write(resp)
}

func (s *HttpServer) getHandler(urlPath string, headers map[string]string, body []byte) []byte {

	var r string

	fmt.Println("urlPath :: ", urlPath)
	if urlPath == "/" {
		r = utils.RespBody(headers, 200, "", "text/plain")
		// r = fmt.Sprintf("HTTP/%s 200 OK\r\n\r\n", s.Version)

	} else if strings.HasPrefix(urlPath, "/echo") {
		echo := strings.Split(urlPath, "/echo")

		var val strings.Builder

		for _, i := range echo {
			if len(strings.Trim(i, "/ ")) > 0 {
				val.WriteString(strings.Trim(i, "/ "))
			}
		}

		r = utils.RespBody(headers, 200, val.String(), "text/plain")
		// r = fmt.Sprintf("HTTP/%s 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", s.Version, len(val.String()), val.String())

	} else if strings.HasPrefix(urlPath, "/user-agent") {
		agent, ok := headers["User-Agent"]
		if !ok {
			r = utils.RespBody(headers, 400, "", "text/plain")
			// r = "HTTP/1.1 400 Bad Request\r\n\r\n%"
		} else {
			r = utils.RespBody(headers, 200, agent, "text/plain")
			// r = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(agent), agent)
		}

	} else if strings.HasPrefix(urlPath, "/files") {
		file := strings.Split(urlPath, "/files")
		fmt.Println("file :: ", file)

		var path string

		for _, i := range file {
			if len(strings.Trim(i, "/ ")) > 0 {
				path += strings.Trim(i, "/ ")
			}

			d, err := os.ReadFile(fmt.Sprintf("%s/%s", s.dir, path))
			if err != nil {
				fmt.Println("file error")
				r = utils.RespBody(headers, 404, "", "text/plain")
				// r = fmt.Sprintf("HTTP/%s 404 Not Found\r\n\r\n", s.Version)
			} else {
				r = utils.RespBody(headers, 200, string(d), "application/octet-stream")
				// r = fmt.Sprintf("HTTP/%s 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", s.Version, len(d), d)
				// r = fmt.Sprintf("HTTP/%s 404 %s\r\n\r\n", s.Version, d)
			}

		}

	} else {
		r = utils.RespBody(headers, 404, "", "text/plain")
		// r = fmt.Sprintf("HTTP/%s 404 Not Found\r\n\r\n", s.Version)
	}

	return []byte(r)
}

func (s *HttpServer) postHandler(urlPath string, headers map[string]string, body []byte) []byte {

	var r string

	if strings.HasPrefix(urlPath, "/files") {
		f := strings.Split(urlPath, "/files")

		var path string

		for _, i := range f {
			if len(strings.Trim(i, "/ ")) > 0 {
				path += strings.Trim(i, "/ ")
			}

			err := os.WriteFile(fmt.Sprintf("%s/%s", s.dir, path), body, 0644)
			if err != nil {
				fmt.Println("file error")
				r = fmt.Sprintf("HTTP/%s 400 Bad Request\r\n\r\n", s.Version)
			} else {
				r = fmt.Sprintf("HTTP/%s 201 Created\r\n\r\n", s.Version)
			}

		}
	} else {
		r = fmt.Sprintf("HTTP/%s 400 Bad Request\r\n\r\n", s.Version)
	}

	return []byte(r)
}
