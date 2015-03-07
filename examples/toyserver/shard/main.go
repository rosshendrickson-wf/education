package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/rosshendrickson-wf/education/examples/toyserver/message"
)

var defaultPort = "8001"
var defaultAddr = "localhost"

type Shard struct {
	revision int
	proposed int
	name     int

	// states
	newRevision  bool
	confirmedNew bool
	connected    bool
}

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
	s := &Shard{name: 101}
	m := message.ConnectMessage(s.name, 1)
	packet := message.MessageToPacket(m)
	tcpconn.Write(packet)
	// Start Listening for connection
	for !s.connected {
		handleTCP(tcpconn, s)
	}

	defer tcpconn.Close()
	//go shardState(s)
	for {
		handleTCP(tcpconn, s)
	}

}

func shardState(s *Shard) {

	//      Loop
	//		Check if we are supposed to calc new frame (Read TCP)
	//      Check if new frame has shared entities
	//			Loop
	//				Send Ack that we know to bump revision (Write TCP)
	//              Bump Revision

	for {
		for !s.confirmedNew {
			for !s.newRevision {

			}
		}

		// Reset Loop
		s.confirmedNew = false
		s.newRevision = false
	}
}

func handleTCP(conn net.Conn, s *Shard) {

	var buf []byte = make([]byte, 512)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		os.Exit(1)
	}

	m := message.PacketToMessage(buf)
	//	if m != nil && m.Revision >= 0 {
	//		log.Println("Got something", m.Type, m.Revision)
	//	}
	//
	switch m.Type {
	case message.Connect:
		s.connected = true
	case message.FrameUpdate:
		s.newRevision = true
		s.proposed = m.Revision
		ack := &message.Message{
			Name: s.name, Revision: s.proposed, Type: message.FrameUpdateAck}
		packet := message.MessageToPacket(ack)
		conn.Write(packet)
		log.Printf("Proposed update %d to revision %d", s.revision, s.proposed)
	case message.FrameUpdateAck:
		if s.proposed == m.Revision {
			s.confirmedNew = true
			s.revision = s.proposed
			log.Println("Confirmed update to revision", s.revision)
		} else {
			log.Println("Ack something else ", s.revision)
		}

	default:
		println("default")
	}
}
