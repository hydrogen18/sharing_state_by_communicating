package main

import "net"
import "strconv"
import "bufio"
import "fmt"

type increment struct {
	amount   int
	response chan<- int
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:3333")
	if err != nil {
		return
	}
	defer listener.Close()
	var value int

	connectionRoutine := func(conn net.Conn, outgoing chan<- increment) {
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
		incoming := make(chan int, 1)
		outgoing <- increment{incr, incoming}

		fmt.Fprintf(conn, "After incrementing, new value is %d\n", <-incoming)

	}

	increments := make(chan increment)
	conns := make(chan net.Conn, 16)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				break
			}
			conns <- conn
			go connectionRoutine(conn, increments)
		}
	}()

	for {
		select {
		case incr := <-increments:
			value += incr.amount
			incr.response <- value

		case conn := <-conns:
			fmt.Printf("New connection from %v\n", conn.RemoteAddr())
		}
	}
}
