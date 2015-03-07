package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/rosshendrickson-wf/education/examples/toyserver/message"
)

var count = 0
var vcount = 0

var revision = 0
var connections = make(map[net.Addr]*Connection, 0)

type Connection struct {
	name     int
	revision int
	ready    bool
	mu       sync.RWMutex
}

func (c *Connection) SetRevisionReady(revision int, ready bool) {

	c.mu.Lock()
	c.revision = revision
	c.ready = ready
	c.mu.Unlock()
	println("Udate connection rev & ready", revision, ready)
}

func (c *Connection) GetReady() bool {

	c.mu.RLock()
	r := c.ready
	c.mu.RUnlock()
	return r
}

func main() {

	// Listen for incoming connections.
	l, err := net.Listen("tcp", "localhost:8001")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application
	// closes.
	defer l.Close()
	fmt.Println("Listening on " + "local 8001")
	for {
		// Listen for an
		// incoming
		// connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go incrementRevision(conn)
		go handleRequest(conn)
	}
}

func incrementRevision(conn net.Conn) {
	for {
		time.Sleep(time.Second * 1)
		ready := true
		for _, v := range connections {
			if !v.GetReady() {
				ready = false
				break
			}
		}
		if !ready {
			continue
		}
		revision++
		for _, v := range connections {
			update := &message.Message{
				Name: v.name, Revision: revision, Type: message.FrameUpdate}
			pong := message.MessageToPacket(update)
			conn.Write(pong)
		}
	}
}

func handleRequest(conn net.Conn) {
	println("New Connection")

	//	update := &message.Message{
	//		Name: 0, Revision: revision + 1, Type: message.FrameUpdate}
	//	pong := message.MessageToPacket(update)
	//	conn.Write(pong)
	//
	println("Move to revision", revision)
	var connection *Connection
	for {

		var buf []byte = make([]byte, 512)
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Connection Error", err.Error())
			delete(connections, conn.LocalAddr())
			println("Shard Disconnected")
			println("Currently Connected", len(connections))
			return
		}

		m := message.PacketToMessage(buf)
		//		if m != nil && m.Revision > 0 {
		//			println("Got something")
		//		}

		switch m.Type {
		case message.Connect:
			pong := message.MessageToPacket(m)
			connection = &Connection{name: m.Name,
				revision: revision, ready: true}
			connections[conn.LocalAddr()] = connection
			conn.Write(pong)
			log.Printf("%d connectd", m.Name)
			println("Currently Connected", len(connections))
		case message.FrameUpdateAck:
			println("Server got Ack Frame", m.Name, m.Revision)
			pong := message.MessageToPacket(m)
			conn.Write(pong)
			connection.SetRevisionReady(m.Revision, true)
		}
	}
}
