package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	addr := net.UDPAddr{
		Port: 6000,
		IP:   net.ParseIP("127.0.0.1"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	conn.Write([]byte("start\n"))

	log.Printf("Listening")
	// Do something with `conn`
	connbuf := bufio.NewReader(conn)

	count := 0
	var buf []byte
	for {

		n, remote_addr, err := conn.ReadFromUDP(buf)
		switch {
		case n != 0:
			fmt.Printf("from %v got message %q\n", remote_addr, string(buf[:n]))
		case err != nil:
			log.Fatal(err)
		}
		count++
		if count%1000 == 0 {
			log.Printf("%d", count)
		}

		str, err := connbuf.ReadString('\n')
		//var buf [1024]byte
		//n, err := conn.Read(buf[:])

		if err != nil {
			log.Printf("err %+v", err)

		}

		log.Printf("ping %+v", str)
		_, e := conn.Write([]byte("hello\n"))

		if e != nil {
			log.Printf("e%+v", e)

		}

	}
}
