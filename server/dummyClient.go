package server

import (
	"fmt"
	"net"
)

// This is a dumb thing that will just send a data packet once every N second

type Client interface {
	Connect()
	Disconnect()
	Write(msg string) error
	Read() []byte
}

type DummyTcpClient struct {
	id   int
	addr string
	msg  string
	conn net.Conn
	buf  []byte
}

func (client *DummyTcpClient) Connect() {
	conn, err := net.Dial("tcp", client.addr)
	if err != nil {
		panic(err)
	}
	client.conn = conn
}

func (client *DummyTcpClient) Disconnect() {
	client.conn.Close()
}

func (client *DummyTcpClient) Write(msg []byte) {
	if _, err := client.conn.Write(msg); err != nil {
		panic(err)
	}
}

func (client *DummyTcpClient) Read() []byte {
	if numRead, err := client.conn.Read(client.buf); err != nil {
		panic(err)
	} else {
		fmt.Println("Client", client.id, "read bytes[", string(client.buf[:numRead]), "]")
		return client.buf[:numRead]
	}
}

var lastID int = 0

func NewDummyTcpClient(addr string, msg string) (client *DummyTcpClient) {
	client = &DummyTcpClient{
		id:   lastID,
		addr: addr,
		msg:  msg,
		buf:  make([]byte, 128),
	}
	lastID += 1
	return
}

func TestClient() {
	clients := make([]*DummyTcpClient, 0)
	messages := []string{"Hello", "World", "Other"}
	fmt.Println("Connecting Clients")
	for i := range 3 {
		clients = append(clients, NewDummyTcpClient(":8080", messages[i]))
		fmt.Println("Client [", i, "]")
		clients[i].Connect()
	}

	for i := range 5 {
		for _, client := range clients {
			client.Write([]byte(fmt.Sprintf("%d: %s", i, client.msg)))
			client.Read()
		}
	}

	for _, client := range clients {
		client.Disconnect()
	}
}

// Alrighty, so how do I get this going on the server side?
