package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/rosshendrickson-wf/education/examples/toyserver/message"
)

// read from the connection ever 5ms and apply the updates
func runShip(address, port string) {

	commands := make(chan *message.Vector, 1000)
	ship := &Ship{commands: commands}

	ship.Connect(address, port)

	go func() {
		var buf []byte
		for {
			n, _, err := ship.conn.ReadFromUDP(buf)
			if err == nil {
				continue
			}
			if n > 0 {
				ship.updates <- buf
			}
		}
	}()

OuterLoop:
	for {
		select {
		case <-time.After(time.Millisecond * 5):
			ship.ApplyUpdates()
		case <-time.After(time.Millisecond * 1):
			ship.SendCommands()
		case <-time.After(time.Millisecond * 10):
			ship.Display()
		case <-time.After(time.Second * 1):
			ship.DisplayFrames()
		case <-time.After(time.Millisecond * 1):
			ship.commands <- RandomMove()
		case <-time.After(time.Second * 120):
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
	conn       net.UDPConn
	commands   chan *message.Vector
	updates    chan []byte
	shipTime   int
	serverTime int
	frames     int
}

func (s *Ship) Connect(address string, port string) {

	// spin off a goroutine to read from the connection
}

func (s *Ship) update(xdir, ydir int) {
	s.xp += xdir
	s.yp += ydir
}

func (s *Ship) DisplayFrames() {
	fmt.Printf("------------%s:%d", s.name, s.frames)
	s.frames = 0
}

func (s *Ship) Display() {
	fmt.Printf("%s:%d,%d", s.name, s.xp, s.yp)
}

func (s *Ship) gatherCommands() []*message.Vector {
	var commands []*message.Vector
OuterLoop:
	for {
		select {
		case v := <-s.commands:
			if v != nil {
				commands = append(commands, v)
			}
		case <-time.After(time.Millisecond * 1):
			break OuterLoop
		default:
		}
	}

	return commands
}

func (s *Ship) SendCommands() {

	commands := s.gatherCommands()
	for _, command := range s.gatherCommands() {
		fmt.Printf("command %+v", command)
		// optimistic move will switch to the right place after computation
		s.update(command.X, command.Y)
	}

	m := &message.Message{Vectors: commands}
	s.sendMessage(m)
}

func (s *Ship) sendMessage(m *message.Message) {
	// actuall send the vector over the wire as a command
	b, _ := json.Marshal(m)
	s.conn.Write(b)
}

// reads streaming in information
func (s *Ship) gatherUpdates() []*message.Vector {

	var updates []*message.Vector
OuterLoop:
	for {
		select {
		case b := <-s.updates:
			var m message.Message
			json.Unmarshal(b, m)
			if b != nil {
				updates = append(updates, m.Vectors...)
			}
		case <-time.After(time.Millisecond * 1):
			break OuterLoop
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
	for {
		newAddr := new(net.UDPAddr)
		*newAddr = *addr
		newAddr.IP = make(net.IP, len(addr.IP))
		copy(newAddr.IP, addr.IP)

		vectors := randVectors(100)

		ms := message.VectorsToMessages(vectors, 100)
		//conn.WriteToUDP(b, newAddr)
		for i, m := range ms {
			m.Value = strconv.Itoa(i)
			p := message.MessageToPacket(m)
			go conn.Write(p)
		}

		var buf []byte = make([]byte, 512)
		n, a, err := conn.ReadFromUDP(buf[0:])
		log.Printf("read %s %d", a, n)
		if err != nil {
			return
		}
		time.Sleep(time.Second * 10)
	}
}

func randVectors(num int) []*message.Vector {

	results := make([]*message.Vector, num)
	for i := range results {

		results[i] = RandomMove()
	}

	return results
}
