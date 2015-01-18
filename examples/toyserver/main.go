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
	connbuf := bufio.NewReader(conn)
	for {
		handleClient(conn, connbuf)
	}

}

func handleClient(conn *net.UDPConn, reader *bufio.Reader) {

	var buf []byte = make([]byte, 512)
	conn.ReadFromUDP(buf[0:])

	//n, a, err := conn.ReadFromUDP(buf[0:])
	//log.Printf("read %s %d", a, n)
	//	if err != nil {
	//		return
	//	}

	m := message.PacketToMessage(buf)
	if m.Value != "" {
		count++
		vcount += len(m.Vectors)
	}
	//log.Printf("deserialized: %s", m.Value)
	//	conn.WriteToUDP([]byte("hello"), a)
}
