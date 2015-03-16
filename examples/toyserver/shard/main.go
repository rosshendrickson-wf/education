package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"

	"github.com/rosshendrickson-wf/education/examples/toyserver/message"
)

var (
	ballRadius = 25
	ballMass   = 1
	count      = 0
	fcount     = 0

	space       *chipmunk.Space
	balls       []*chipmunk.Shape
	staticLines []*chipmunk.Shape
	deg2rad     = math.Pi / 180
)

type State struct {
	Shapes []*chipmunk.Shape
}

// createBodies sets up the chipmunk space and static bodies
func createBodies() {
	space = chipmunk.NewSpace()
	space.Gravity = vect.Vect{0, -900}

	staticBody := chipmunk.NewBodyStatic()
	staticLines = []*chipmunk.Shape{
		chipmunk.NewSegment(vect.Vect{111.0, 280.0}, vect.Vect{407.0, 246.0}, 0),
		chipmunk.NewSegment(vect.Vect{407.0, 246.0}, vect.Vect{407.0, 343.0}, 0),
	}
	for _, segment := range staticLines {
		segment.SetElasticity(0.6)
		staticBody.AddShape(segment)
	}
	space.AddBody(staticBody)
}

func addBall() {
	x := rand.Intn(350-115) + 115
	ball := chipmunk.NewCircle(vect.Vector_Zero, float32(ballRadius))
	ball.SetElasticity(0.95)

	body := chipmunk.NewBody(vect.Float(ballMass), ball.Moment(float32(ballMass)))
	body.SetPosition(vect.Vect{vect.Float(x), 600.0})
	body.SetAngle(vect.Float(rand.Float32() * 2 * math.Pi))

	body.AddShape(ball)
	space.AddBody(body)
	balls = append(balls, ball)
}

// step advances the physics engine and cleans up any balls that are off-screen
func step(dt float32) []*message.State {
	space.Step(vect.Float(dt))
	states := make([]*message.State, len(balls))

	for i := 0; i < len(balls); i++ {

		ball := balls[i]
		rot := ball.Body.Angle() * chipmunk.DegreeConst
		p := ball.Body.Position()
		vec := message.Vec{X: float32(p.X), Y: float32(p.Y)}
		s := &message.State{Kind: 0, Position: vec, Rotation: float32(rot)}
		if p.Y < -100 {
			space.RemoveBody(balls[i].Body)
			balls[i] = nil
			balls = append(balls[:i], balls[i+1:]...)
			i-- // consider same index again
		} else {
			states[i] = s
		}
	}
	return states
}

var defaultPort = "8001"
var defaultAddr = "localhost"

type Shard struct {
	revision int
	proposed int
	name     int
	mu       sync.RWMutex

	// Connections
	udpconn *net.UDPConn

	// states
	newRevision  bool
	confirmedNew bool
	connected    bool
}

func (s *Shard) Connect(address, serverPort, clientPort string) bool {

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

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	s.udpconn = conn

	count := 0
Connect:
	for {

		// Send connection message
		m := message.ConnectMessage(s.name, 0)
		pong := message.MessageToPacket(m)
		conn.Write(pong)

		log.Println("Connecting . . . . ", count)
		s.mu.Lock()
		if s.connected {
			s.mu.Unlock()
			break Connect
		}
		s.mu.Unlock()
		count++

		s.handleUDPUpdate(conn)

		if count >= 10 {
			return false
		}

		time.Sleep(time.Second * 1)
	}

	return true
}

func (s *Shard) handleUDPUpdate(conn *net.UDPConn) {

	var buf []byte = make([]byte, 512)
	conn.ReadFromUDP(buf[0:])
	m := message.PacketToMessage(buf)

	if m == nil {
		println("a")
		return
	}
	switch m.Type {
	case message.Connect:
		s.Connected(true)
	default:
		println("DEFAULT")
	}
}

func (s *Shard) Connected(value bool) {
	s.mu.Lock()
	s.connected = value
	s.mu.Unlock()
}

func (s *Shard) CalcNextFrame() bool {
	s.mu.RLock()
	r := s.confirmedNew
	s.mu.RUnlock()
	return r
}

func (s *Shard) SetConfirmedNew(value bool, revision int) {
	s.mu.Lock()
	s.confirmedNew = value
	if value {
		s.revision = revision
	}
	s.mu.Unlock()
}

func main() {

	// Stat Ticker
	tickerSecond := 1
	ticker := time.NewTicker(time.Second * 1)
	go func() {
		for _ = range ticker.C {
			log.Printf("Received ~%d Frames Updates", count/tickerSecond)
			count = 0

			log.Printf("Calc ~%d Frames", fcount/tickerSecond)
			fcount = 0
		}
	}()

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

	// UDP STATE
	address := ""
	clientPort := "36504"
	serverPort := "10235"

	connected := s.Connect(address, serverPort, clientPort)

	if !connected {
		panic(errors.New("Unable to establish UDP State updates"))
	}

	m := message.ConnectMessage(s.name, 1)
	packet := message.MessageToPacket(m)
	tcpconn.Write(packet)
	// Start Listening for connection
	for !s.connected {
		handleTCP(tcpconn, s)
	}

	defer tcpconn.Close()

	// Set up Physics state
	createBodies()

	runtime.LockOSThread()

	go shardState(s)
	for {
		handleTCP(tcpconn, s)
	}

}

func shardState(s *Shard) {
	ticksToNextBall := 10
	for {
		if s.CalcNextFrame() {
			ticksToNextBall--
			if ticksToNextBall == 0 {
				ticksToNextBall = rand.Intn(100) + 1
				addBall()
			}
			states := step(1.0 / 60.0)
			s.SetConfirmedNew(false, 0)
			messages := message.StatesToMessages(states)
			for _, m := range messages {
				packet := message.MessageToPacket(m)
				s.udpconn.Write(packet)
				println("Wrote state")
				fcount++
			}
		}
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
		//log.Printf("Proposed update %d to revision %d", s.revision, s.proposed)
	case message.FrameUpdateAck:
		if s.proposed == m.Revision {
			s.SetConfirmedNew(true, s.proposed)
			s.revision = s.proposed
			count++
			//	log.Println("Confirmed update to revision", s.revision)
		} else {
			log.Println("Ack something else ", s.revision)
		}

		//	default:
		//		println("default")
	}
}
