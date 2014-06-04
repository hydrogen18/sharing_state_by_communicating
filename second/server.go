package main

import "net"
import "strconv"
import "bufio"
import "fmt"

func connectionRoutine(conn net.Conn, increment func(amount int) int) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	line = line[:len(line)-1]
	incr, err := strconv.Atoi(line)
	if err != nil {
		fmt.Fprintf(conn, "The value %q is not a valid number\n", line)
		return
	}

	fmt.Fprintf(conn, "After incrementing, new value is %d\n", increment(incr))
}

type increment struct {
	amount   int
	response chan<- int
}

type Server struct {
	listener   net.Listener
	conns      chan net.Conn
	increments chan increment
}

func (s *Server) Accept() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			break
		}
		s.conns <- conn
		go connectionRoutine(conn, s.Increment)
	}
}

func (s *Server) Increment(amount int) int {
	response := make(chan int, 1)
	increment := increment{}
	increment.amount = amount
	increment.response = response
	s.increments <- increment
	return <-response
}

func main() {
	var err error
	s := Server{}
	s.listener, err = net.Listen("tcp", "0.0.0.0:3333")
	if err != nil {
		return
	}
	defer s.listener.Close()
	var value int

	s.increments = make(chan increment)
	s.conns = make(chan net.Conn, 16)

	go s.Accept()

	for {
		select {
		case incr := <-s.increments:
			value += incr.amount
			incr.response <- value

		case conn := <-s.conns:
			fmt.Printf("New connection from %v\n", conn.RemoteAddr())
		}
	}
}
