package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

type vector struct {
	xdir int
	ydir int
}

// read from the connection ever 5ms and apply the updates
func runShip(address, port string) {

	commands := make(chan *vector, 1000)
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

type Message struct {
	Vectors []*vector
}

type Ship struct {
	xp         int
	yp         int
	health     int
	name       string
	conn       net.UDPConn
	commands   chan *vector
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

func (s *Ship) gatherCommands() []*vector {
	var commands []*vector
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
		s.update(command.xdir, command.ydir)
	}

	m := &Message{commands}
	s.sendMessage(m)
}

func (s *Ship) sendMessage(m *Message) {
	// actuall send the vector over the wire as a command
	b, _ := json.Marshal(m)
	b = append(b, []byte("\n")...)
	s.conn.Write(b)
}

// reads streaming in information
func (s *Ship) gatherUpdates() []*vector {

	var updates []*vector
OuterLoop:
	for {
		select {
		case b := <-s.updates:
			var m Message
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
		s.xp = update.xdir
		s.yp = update.ydir
	}
	fmt.Printf("%s:%d,%d", s.name, s.xp, s.yp)
	s.frames++
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func RandomMove() *vector {

	xdir := random(1, 6)
	ydir := random(1, 10)

	return &vector{xdir, ydir}
}

func main() {

	conn, err := net.Dial("udp", "127.0.0.1:6000")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	connbuf := bufio.NewReader(conn)
	for {
		str, err := connbuf.ReadString('\n')
		if len(str) > 0 {
			fmt.Println("pong", str)
		}
		if err != nil {
			break
		}

		_, e := conn.Write([]byte("hello\n"))

		if e != nil {
			log.Printf("e%+v", e)

		}

	}
}
