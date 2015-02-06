package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"github.com/rosshendrickson-wf/education/examples/toyserver/message"
)

//1000ms/sec / 900FPS = 1.111.. ms per frame
//1000ms/sec / 450FPS = 2.222.. ms per frame
//Increase in execution time: 1.111.. ms
//
//1000ms/sec / 60FPS = 16.666.. ms per frame
//1000ms/sec / 56.25FPS = 17.777.. ms per frame

// read from the connection ever 5ms and apply the updates
func runShip(address, serverPort, clientPort string) {

	commands := make(chan *message.Vector, 1000)
	ship := &Ship{commands: commands}

	ship.Connect(address, serverPort, clientPort)

	//time.Sleep(time.Millisecond * 1000)

	//ApplyUpdates := time.NewTicker(time.Millisecond * 5).C
	SendCommands := time.NewTicker(time.Millisecond * 10).C
	//Display := time.NewTicker(time.Millisecond * 120).C
	DisplayFrames := time.NewTicker(time.Second * 1).C
	Random := time.NewTicker(time.Millisecond * 1).C
	DieTime := time.NewTicker(time.Second * 10).C

OuterLoop:
	for {
		select {
		//	case <-ApplyUpdates:
		//		ship.ApplyUpdates()
		case <-SendCommands:
			ship.SendCommands()
		//	case <-Display:
		//		ship.Display()
		case <-DisplayFrames:
			ship.DisplayFrames()
		case <-Random:
			ship.commands <- RandomMove()
		case <-DieTime:
			ship.Close()
			log.Printf("Ship %+v dead", ship)
			break OuterLoop
		default:
		}
	}
}

type Ship struct {
	xp         int
	yp         int
	health     int
	name       [8]byte
	conn       *net.UDPConn
	sconn      *net.UDPConn
	commands   chan *message.Vector
	updates    chan []byte
	shipTime   int
	serverTime int
	frames     int
	lock       sync.Mutex
	stop       bool
	revision   int
	serverAddr net.Addr
}

func (s *Ship) Connect(address, serverPort, clientPort string) {

	addr, err := net.ResolveUDPAddr("udp", ":"+clientPort)
	if err != nil {
		log.Fatal(err)
	}

	saddr, err := net.ResolveUDPAddr("udp", address+":"+serverPort)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", addr, saddr)
	log.Printf("Listening on %+v", addr)
	log.Printf("Sending on %+v", saddr)
	//conn, err := net.ListenUDP("udp", addr)
	//	conn, err := net.DialUDP("udp", addr, saddr)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	s.conn = conn

	//	sconn, err := net.DialUDP("udp", nil, saddr)
	//	if err != nil {
	//		fmt.Println(err)
	//		os.Exit(1)
	//	}

	go func() {
		//		defer conn.Close()
		for {
			//			println("HEY")
			var buf []byte = make([]byte, 512)
			conn.ReadFromUDP(buf[0:])
			message.PacketToMessage(buf)
			s.frames++
			if s.Stop() {
				println("NO")
				return
			}
			//			println("Looped")
		}
	}()

}

func (s *Ship) Close() {
	s.lock.Lock()
	s.stop = true
	s.lock.Unlock()
}

func (s *Ship) Stop() bool {
	s.lock.Lock()
	result := s.stop
	s.lock.Unlock()
	return result
}

func (s *Ship) handleUpdate(conn *net.UDPConn) {
	//	s.updates <- buf
}

func (s *Ship) update(xdir, ydir int) {
	s.xp += xdir
	s.yp += ydir
}

func (s *Ship) DisplayFrames() {
	log.Printf("------------%s:%d", s.name, s.frames)
	s.frames = 0
}

func (s *Ship) Display() {
	fmt.Printf("%s:%d,%d", s.name, s.xp, s.yp)
}

func (s *Ship) gatherCommands() []*message.Vector {

	var stop bool
	time.AfterFunc(time.Millisecond*1, func() {
		stop = true
	})

	var commands []*message.Vector
OuterLoop:
	for {
		select {
		case v := <-s.commands:
			if v != nil {
				commands = append(commands, v)
			}
			if stop {
				break OuterLoop
			}
		default:
			if stop {
				break OuterLoop
			}
		}
	}

	return commands
}

func (s *Ship) SendCommands() {

	num := 1000
	vectors := randVectors(num)
	s.revision++

	ms := message.VectorsToMessages(vectors, s.revision)
	s.sendMessages(ms...)
	if s.revision%10 == 0 {
		log.Printf("Sent %d Commands at Rev %d", len(ms), s.revision)
	}

}

func (s *Ship) sendMessages(ms ...*message.Message) {
	// actuall send the vector over the wire as a command
	for _, m := range ms {
		b := message.MessageToPacket(m)
		s.conn.Write(b)
	}
}

// reads in as many updates as possible in 1 millisecond
func (s *Ship) gatherUpdates() []*message.Vector {

	var stop bool
	time.AfterFunc(time.Millisecond*1, func() {
		stop = true
	})

	var updates []*message.Vector
OuterLoop:
	for {
		select {
		case b := <-s.updates:
			var m message.Message
			json.Unmarshal(b, m)
			if m.Type != "" {
				println("Ship got message")
				//				updates = append(updates, m.Vectors...)
			}
			if stop {
				break OuterLoop
			}
		default:
		}
	}
	return updates
}

func (s *Ship) ApplyUpdates() {

	// Updates override the position as they are state updates
	for _, update := range s.gatherUpdates() {
		s.xp = update.X
		s.yp = update.Y
	}
	fmt.Printf("%s:%d,%d", s.name, s.xp, s.yp)
	s.frames++
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func RandomMove() *message.Vector {

	xdir := random(1, 6)
	ydir := random(1, 10)

	return &message.Vector{xdir, ydir}
}

func main() {

	address := ""
	clientPort := "36503"
	serverPort := "10234"
	//	clientPort := "55102"
	runShip(address, serverPort, clientPort)
	//test(address, serverPort, clientPort)
}

func test(address, serverPort, clientPort string) {

	//////////////////////////////////// Sends Random traffic
	//	addr, err := net.ResolveUDPAddr("udp", ":10234")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	maddr, err := net.ResolveUDPAddr("udp", ":"+clientPort)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	addr, err := net.ResolveUDPAddr("udp", ":"+clientPort)
	if err != nil {
		log.Fatal(err)
	}

	saddr, err := net.ResolveUDPAddr("udp", address+":"+serverPort)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", addr, saddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer conn.Close()
	log.Println("Connected to ", addr)
	num := 1000
	vectors := randVectors(num)

	ms := message.VectorsToMessages(vectors, 101)

	go func() {
		for {
			var buf []byte = make([]byte, 512)
			n, a, err := conn.ReadFromUDP(buf[0:])
			log.Printf("read %s %d", a, n)
			if err != nil {
				//		return
			}
		}
	}()

	for {

		for i, m := range ms {
			m.Revision = i
			p := message.MessageToPacket(m)
			conn.Write(p)
		}

		log.Printf("Sent %d vectors in %d messages", num, len(ms))
		println("Sent Messages")
		time.Sleep(time.Millisecond * 1000)
	}
}

func randVectors(num int) []*message.Vector {

	results := make([]*message.Vector, num)
	for i := range results {

		results[i] = RandomMove()
	}

	return results
}
