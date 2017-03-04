package redis_replica

import (
	"net"
	"fmt"
	"log"
)

type Server struct {
	Host string
	Port int
}

// NewServer создаёт новый сервер реплики
func NewServer() (*Server) {
	s := new(Server)
	return s
}

// Close
func (s *Server) Close() error {
	return nil
}

// ListenAndServe
func (s *Server) ListenAndServe(host string, port int) error {
	s.Host = host
	s.Port = port

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:d%", s.Host, s.Port))
	if err != nil {
		return err
	}
	go s.listener(ln)
	return nil
}

// Получаем соединения
func (s *Server) listener(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Unable to accept: %v\n", err)
			continue
		}
		go s.replica(conn)
	}
}

// Копируем данные клиенту
func (s *Server) replica(conn net.Conn) {

}