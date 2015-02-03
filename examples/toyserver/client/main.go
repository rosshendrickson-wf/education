package main

import (
	"bufio"
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

	//ApplyUpdates := time.NewTicker(time.Millisecond * 5).C
	SendCommands := time.NewTicker(time.Millisecond * 500).C
	//Display := time.NewTicker(time.Millisecond * 120).C
	DisplayFrames := time.NewTicker(time.Second * 1).C
	Random := time.NewTicker(time.Millisecond * 5).C
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

	// spin off a goroutine to read from the connection
	addr, err := net.ResolveUDPAddr("udp4", ":"+clientPort)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	saddr, err := net.ResolveUDPAddr("udp4", address+":"+serverPort)
	if err != nil {
		log.Fatal(err)
	}
	s.serverAddr = saddr

	sconn, err := net.DialUDP("udp", nil, saddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	s.sconn = sconn

	go func() {
		for {
			defer conn.Close()
			connbuf := bufio.NewReader(conn)
			s.handleUpdate(conn, connbuf)
			if s.Stop() {
				return
			}
		}
	}()

	s.conn = conn
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

func (s *Ship) handleUpdate(conn *net.UDPConn, reader *bufio.Reader) {
	var buf []byte = make([]byte, 512)
	conn.ReadFromUDP(buf[0:])
	m := message.PacketToMessage(buf)
	log.Printf("ship got message %+v", m)
	s.frames++
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
	s.revision++
	log.Printf("Sending Commands at Rev %d", s.revision)
	ms := message.VectorsToMessages(s.gatherCommands(), s.revision)
	s.sendMessages(ms...)
}

func (s *Ship) sendMessages(ms ...*message.Message) {
	// actuall send the vector over the wire as a command
	for _, m := range ms {
		b := message.MessageToPacket(m)
		s.sconn.Write(b)
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

	runShip("", "10234", "55102")
}

func test() {
	//////////////////////////////////// Sends Random traffic
	addr, err := net.ResolveUDPAddr("udp", ":10234")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer conn.Close()
	log.Println("Connected to ", addr)
	num := 100000
	vectors := randVectors(num)

	ms := message.VectorsToMessages(vectors, 101)

	for {
		//		newAddr := new(net.UDPAddr)
		//		*newAddr = *addr
		//		newAddr.IP = make(net.IP, len(addr.IP))
		//		copy(newAddr.IP, addr.IP)
		//
		//conn.WriteToUDP(b, newAddr)
		for i, m := range ms {
			m.Revision = i
			p := message.MessageToPacket(m)
			conn.Write(p)
		}

		//	var buf []byte = make([]byte, 512)
		//	n, a, err := conn.ReadFromUDP(buf[0:])
		//	log.Printf("read %s %d", a, n)
		//	if err != nil {
		//		return
		//	}
		//log.Printf("Sent %d vectors in %d", num, len(ms))
		//time.Sleep(time.Second * 1)
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
