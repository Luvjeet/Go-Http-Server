package main

import (
	"fmt"
	"net"
    "os"
	"strings"
    "bytes"
)

type HttpRequest struct {
    method string
    path string
    version string
    pathParts []string
    headers map[string]string
    body string
}

type HttpResponse struct {
    status string
    contentType string
    content string
}

func sendResponse(request HttpRequest)bytes.Buffer{
    fmt.Println(request)
    response := bytes.NewBufferString(fmt.Sprintf("HTTP/1.1 200 OK\r\n\r\nContent-Type: text/plain\r\nContent-Length:%s\r\n",request.headers["Content-Length"]))
    return *response
}

func handleConnection(conn *net.Conn) {
	defer (*conn).Close()
    input := make([]byte, 256);

    _, err := (*conn).Read(input)
	if err != nil {
		fmt.Println("Failed to read request", err.Error())
		os.Exit(1)
	}
    req := strings.Split(string(input), "\r\n\r\n")
    headers := req[0]
    body := req[1]
    rows := strings.Split(string(headers), "\r\n")

    var request HttpRequest;
    request.method = strings.Split(rows[0], " ")[0]
    request.path = strings.Split(rows[0], " ")[1]
    request.pathParts = strings.Split(strings.Split(rows[0], " ")[1],"/")
    request.headers = make(map[string]string) // Initialize the headers map
    for i:=1; i<len(rows);i++ {
        content := strings.Split(rows[i], ":")
        request.headers[content[0]] = content[1]
    }
    request.body = body

    response := sendResponse(request)
    bytesWritten, err := (*conn).Write(response.Bytes())
    if err != nil {
        fmt.Println("Failed to read request", err.Error())
    }
    fmt.Println(bytesWritten)
}

func main() {
	listener, err := net.Listen("tcp", ":6969")

	if err != nil {
		fmt.Println("Error occurred", err)
	}
	defer listener.Close()
	fmt.Println("Server listening on port 6969...")

	// Accept and handle incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(&conn)
	}

}
