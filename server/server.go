package server

import (
	"fmt"
	"net"
	"time"
)

/*
	Responsible for setting up a listener on a certain IP address
*/

func Test() {
	fmt.Println("Setting up link")
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening on port :8080")
	defer ln.Close()

	go func() {
		fmt.Println("Waiting one second")
		time.Sleep(1 * time.Second)
		TestClient()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go func(conn net.Conn) {
			defer func() {
				fmt.Println("Closing Client:", conn.LocalAddr())
				conn.Close()
			}()
			i := 0
			for {
				fmt.Println("Processing connection:", conn.LocalAddr())
				buf := make([]byte, 128)
				n, err := conn.Read(buf)
				if err != nil {
					if err.Error() == "EOF" {
						return
					}
				}
				fmt.Println("found bytes:", string(buf[:n]))
				conn.Write([]byte(fmt.Sprintf("%d: World", i)))
				i += 1
			}
		}(conn)
	}
}
