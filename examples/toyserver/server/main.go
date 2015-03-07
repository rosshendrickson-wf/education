package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/rosshendrickson-wf/education/examples/toyserver/message"
)

var count = 0
var vcount = 0

var connections = make([]*Connection, 1)

type Connection struct {
	name     int
	revision int
	ready    bool
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
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {

	revision := 0
	for {
		var buf []byte = make([]byte, 512)
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			return
		}

		m := message.PacketToMessage(buf)
		if m != nil && m.Revision > 0 {
			println("Got something")
		}

		switch m.Type {
		case message.Connect:
			pong := message.MessageToPacket(m)
			connection := &Connection{m.Name, m.Revision, false}
			connections = append(connections, connection)
			conn.Write(pong)
			log.Printf("CONNECTED %+v", m.Name)
		case message.FrameUpdateAck:
			println("Server got Ack Frame")
			pong := message.MessageToPacket(m)
			conn.Write(pong)
		case message.InputUpdate:
			log.Printf("Got input")
		case message.VectorUpdate:
			count++
			if count%100 == 0 {
				log.Printf("PONG")
				pong := message.MessageToPacket(m)
				conn.Write(pong)
			}
			vectors := message.PayloadToVectors(m.Payload)
			vcount += len(vectors)
		default:
			log.Printf("DEFAULT %+v", m)
		}

		revision++
		update := &message.Message{
			Name: m.Name, Revision: revision, Type: message.FrameUpdate}

		pong := message.MessageToPacket(update)
		println("Move to revision", revision)
		conn.Write(pong)

	}
}
