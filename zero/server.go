package main

import "net"
import "strconv"
import "bufio"
import "sync"
import "fmt"

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:3333")
	if err != nil {
		return
	}
	defer listener.Close()
	var value int
	valueLock := &sync.Mutex{}

	connectionRoutine := func(conn net.Conn) {
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

		valueLock.Lock()
		defer valueLock.Unlock()
		value += incr
		fmt.Fprintf(conn, "After incrementing, new value is %d\n", value)

	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		go connectionRoutine(conn)

	}
}
