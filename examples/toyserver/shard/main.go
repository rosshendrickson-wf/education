package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/rosshendrickson-wf/education/examples/toyserver/message"
)

var defaultPort = "8001"
var defaultAddr = "localhost"

var connected bool
var revision int
var proposed int
var name = 101

// states
var newRevision bool
var confirmedNew bool

func main() {
	var (
		addr = flag.String("address", defaultAddr, "Address to server")
		tcp  = flag.String("tcpport", defaultPort, "TCP port")
		//		udp  = flag.String("udpport", "8002", "UDP port")
	)

	servAddr := *addr + ":" + *tcp
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	tcpconn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	// Connection
	m := message.ConnectMessage(name, 1)
	packet := message.MessageToPacket(m)
	tcpconn.Write(packet)
	// Start Listening for connection
	for !connected {
		handleTCP(tcpconn)
	}
	//      Loop
	//		Check if we are supposed to calc new frame (Read TCP)
	//      Check if new frame has shared entities
	//			Loop
	//				Send Ack that we know to bump revision (Write TCP)
	//              Bump Revision

	for {
		for !confirmedNew {
			for !newRevision {
				handleTCP(tcpconn)
			}
			handleTCP(tcpconn)
		}
		println("Revision incremented to ", revision)
		// Reset Loop
		confirmedNew = false
		newRevision = false
	}

}

func handleTCP(conn net.Conn) {
	var buf []byte = make([]byte, 512)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		os.Exit(1)
	}

	m := message.PacketToMessage(buf)
	if m != nil && m.Revision > 0 {
		println("Got something", m)
	}

	switch m.Type {
	case message.Connect:
		connected = true
	case message.FrameUpdate:
		newRevision = true
		proposed = m.Revision
		ack := &message.Message{
			Name: name, Revision: revision, Type: message.FrameUpdateAck}
		packet := message.MessageToPacket(ack)
		conn.Write(packet)
		println("Proposed update to revision", proposed)
	case message.FrameUpdateAck:
		confirmedNew = true
		revision = proposed
	default:
	}

}
