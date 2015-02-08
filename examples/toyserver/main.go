package main

import (
	"bufio"
	"log"
	"net"
	"runtime"
	"time"

	"github.com/rosshendrickson-wf/education/examples/toyserver/message"
)

// Global - Bad - testing
var count int = 0
var old int = 0

var vcount int = 0
var vold int = 0

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	addr, err := net.ResolveUDPAddr("udp4", ":10234")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(time.Second * 2)
	go func() {
		for _ = range ticker.C {
			log.Printf("Processed %d messages", count)
			log.Printf("Processed %d vectors", vcount)
			old += count
			count = 0
			vcount = 0
		}
	}()

	defer ticker.Stop()
	defer conn.Close()

	log.Printf("Read loop Start %+v", addr)
	for {
		handleClient(conn)
	}

}

func handleClient(conn *net.UDPConn) {

	var buf []byte = make([]byte, 512)
	//	conn.ReadFromUDP(buf[0:])

	_, a, err := conn.ReadFromUDP(buf[0:])
	//	log.Printf("read %s %d", a, n)
	if err != nil {
		return
	}

	m := message.PacketToMessage(buf)
	if m != nil && m.Revision > 0 {
		println("Got something")
	}

	switch m.Type {
	case message.Connect:
		pong := message.MessageToPacket(m)
		conn.WriteTo(pong, a)
		log.Printf("CONNECTED %+v: %+v", a, m)

	case message.VectorUpdate:
		count++
		if count%100 == 0 {
			log.Printf("PONG %+v", a)
			pong := message.MessageToPacket(m)
			conn.WriteTo(pong, a)
		}
		vectors := message.PayloadToVectors(m.Payload)
		vcount += len(vectors)
	default:
		log.Printf("DEFAULT %+v", m)
	}
}
