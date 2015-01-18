package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/rosshendrickson-wf/education/examples/toyserver/message"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp4", ":10234")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	log.Printf("Read loop Start %+v", addr)
	connbuf := bufio.NewReader(conn)
	for {
		handleClient(conn, connbuf)
	}
}

func handleClient(conn *net.UDPConn, reader *bufio.Reader) {

	var buf []byte = make([]byte, 512)
	n, a, err := conn.ReadFromUDP(buf[0:])
	log.Printf("read %s %d", a, n)
	if err != nil {
		return
	}

	m := message.PacketToMessage(buf)
	//	log.Printf("bytes: %+v", buf)
	log.Printf("deserialized: %+v", m)
	conn.WriteToUDP([]byte("hello"), a)
}
