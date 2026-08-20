package server

import (
	"net"
	"personal-http-server/parser/request"
)

type Server struct {
	Request request.Request
	Port    int
	Host    string
}

// func NewServer(serverName, host, port string, lc net.ListenConfig) *Server {
// 	addr := net.JoinHostPort(host, port)
// 	// tcpAddr, err := net.ResolveIPAddr("tcp", addr)
// 	return
// }
