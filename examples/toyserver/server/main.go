package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/Workiva/go-datastructures/queue"

	"github.com/rosshendrickson-wf/education/examples/toyserver/message"
)

var count = 0
var vcount = 0

var revision = 0
var connections = make(map[net.Addr]*Connection, 0)
var clients = make(map[net.Addr]*Connection, 0)
var bQueue *queue.Queue

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
	//println("Udate connection rev & ready", revision, ready)
}

func (c *Connection) GetReady() bool {

	c.mu.RLock()
	r := c.ready
	c.mu.RUnlock()
	return r
}

func main() {

	// Stat Ticker
	ticker := time.NewTicker(time.Second * 2)
	go func() {
		for _ = range ticker.C {
			log.Printf("Processed ~%d Frames", count/2)
			count = 0
		}
	}()

	// UDP - clients
	addr, err := net.ResolveUDPAddr("udp4", ":10234")
	if err != nil {
		log.Fatal(err)
	}

	udpconn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	go func(conn *net.UDPConn) {
		for {
			handleUDPClient(conn)
		}
	}(udpconn)
	log.Printf("Read UDP loop Start %+v", addr)

	// TCP
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

	// Setup Broadcaster
	bQueue = queue.New(1000)

	go func(queue *queue.Queue) {

		for {
			packet, err := queue.Get(1)
			if err != nil {
				log.Printf("Error accessing items from queue %s", err)
				return
			}

			Broadcast(udpconn, packet[0].([]byte))
		}
	}(bQueue)

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
		go handleRequest(conn, udpconn)
	}
}

func incrementRevision(conn net.Conn) {
	for {
		time.Sleep(time.Second * 1 / 120)
		ready := true
		for _, v := range connections {
			if !v.GetReady() {
				ready = false
				break
			}
		}
		if !ready || len(connections) == 0 {
			continue
		}
		revision++
		if revision > 100000 {
			os.Exit(1)
		}
		count++
		for _, v := range connections {
			update := &message.Message{
				Name: v.name, Revision: revision, Type: message.FrameUpdate}
			pong := message.MessageToPacket(update)
			conn.Write(pong)
		}
	}
}

func handleRequest(conn net.Conn, udpconn *net.UDPConn) {
	println("New Connection")
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
			//println("Server got Ack Frame", m.Name, m.Revision)
			pong := message.MessageToPacket(m)
			conn.Write(pong)
			connection.SetRevisionReady(m.Revision, true)
		case message.StateUpdate:
			println("got stat update")
			bQueue.Put(buf)
		}
	}
}

func handleUDPClient(conn *net.UDPConn) {
	var connection *Connection
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
		connection = &Connection{name: m.Name,
			revision: revision, ready: true}
		clients[a] = connection

	case message.InputUpdate:
		log.Printf("Got input")
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

func Broadcast(conn *net.UDPConn, packet []byte) {
	println("BroadCast")
	for a, _ := range clients {
		println("STUCK")
		conn.WriteTo(packet, a)
		println("HERE")

	}
}
