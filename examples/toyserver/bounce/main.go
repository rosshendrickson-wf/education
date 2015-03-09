package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"

	"github.com/rosshendrickson-wf/education/examples/toyserver/message"
)

//1000ms/sec / 900FPS = 1.111.. ms per frame
//1000ms/sec / 450FPS = 2.222.. ms per frame
//Increase in execution time: 1.111.. ms
//
//1000ms/sec / 60FPS = 16.666.. ms per frame
//1000ms/sec / 56.25FPS = 17.777.. ms per frame

var (
	ballRadius = 25
	ballMass   = 1

	space       *chipmunk.Space
	balls       []*chipmunk.Shape
	staticLines []*chipmunk.Shape
	deg2rad     = math.Pi / 180
	window      *glfw.Window
)

// drawCircle draws a circle for the specified radius, rotation angle, and the specified number of sides
func drawCircle(radius float64, sides int) {
	gl.Begin(gl.LINE_LOOP)
	for a := 0.0; a < 2*math.Pi; a += (2 * math.Pi / float64(sides)) {
		gl.Vertex2d(math.Sin(a)*radius, math.Cos(a)*radius)
	}
	gl.Vertex3f(0, 0, 0)
	gl.End()
}

// OpenGL draw function
func draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.POINT_SMOOTH)
	gl.Enable(gl.LINE_SMOOTH)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.LoadIdentity()

	gl.Begin(gl.LINES)
	gl.Color3f(.2, .2, .2)
	for i := range staticLines {
		x := staticLines[i].GetAsSegment().A.X
		y := staticLines[i].GetAsSegment().A.Y
		gl.Vertex3f(float32(x), float32(y), 0)
		x = staticLines[i].GetAsSegment().B.X
		y = staticLines[i].GetAsSegment().B.Y
		gl.Vertex3f(float32(x), float32(y), 0)
	}
	gl.End()

	gl.Color4f(.3, .3, 1, .8)
	// draw balls
	for _, ball := range balls {
		gl.PushMatrix()
		pos := ball.Body.Position()
		rot := ball.Body.Angle() * chipmunk.DegreeConst
		gl.Translatef(float32(pos.X), float32(pos.Y), 0.0)
		gl.Rotatef(float32(rot), 0, 0, 1)
		drawCircle(float64(ballRadius), 60)
		gl.PopMatrix()
	}
}

// onResize sets up a simple 2d ortho context based on the window size
func onResize(window *glfw.Window, w, h int) {
	w, h = window.GetSize() // query window to get screen pixels
	width, height := window.GetFramebufferSize()
	gl.Viewport(0, 0, width, height)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(w), 0, float64(h), -1, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.ClearColor(1, 1, 1, 1)
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

func addBall(x, y, rot float32) {

	ball := chipmunk.NewCircle(vect.Vector_Zero, float32(ballRadius))
	ball.SetElasticity(0.95)

	body := chipmunk.NewBody(vect.Float(ballMass), ball.Moment(float32(ballMass)))
	body.SetPosition(vect.Vect{vect.Float(x), vect.Float(y)})
	body.SetAngle(vect.Float(rot))

	body.AddShape(ball)
	space.AddBody(body)
	balls = append(balls, ball)
}

// read from the connection ever 5ms and apply the updates
func runShip(address, serverPort, clientPort string) {
	// set up physics
	createBodies()

	commands := make(chan *message.Vector, 1000)
	ship := &Ship{name: 100, commands: commands}

	connected := ship.Connect(address, serverPort, clientPort)

	if connected {
		log.Println("Ship is connected to server")
	} else {
		log.Println("Unable to start ship")
		return
	}

	//ApplyUpdates := time.NewTicker(time.Millisecond * 5).C
	//	SendCommands := time.NewTicker(time.Millisecond * 100).C
	//Display := time.NewTicker(time.Millisecond * 120).C
	DisplayFrames := time.NewTicker(time.Second * 1).C
	Random := time.NewTicker(time.Millisecond * 1).C
	DieTime := time.NewTicker(time.Second * 10).C

OuterLoop:
	for {
		select {
		//	case <-ApplyUpdates:
		//		ship.ApplyUpdates()
		//		case <-SendCommands:
		//			ship.SendCommands()
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
	name       int
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
	connected  bool
}

func (s *Ship) Connected(value bool) {
	s.lock.Lock()
	s.connected = value
	s.lock.Unlock()
}

func (s *Ship) Connect(address, serverPort, clientPort string) bool {

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
	s.conn = conn

	count := 0
Connect:
	for {

		// Send connection message
		m := message.ConnectMessage(s.name, 0)
		pong := message.MessageToPacket(m)
		conn.Write(pong)

		log.Println("Connecting . . . . ", count)
		s.lock.Lock()
		if s.connected {
			s.lock.Unlock()
			break Connect
		}
		s.lock.Unlock()
		count++

		s.handleUpdate(conn)

		if count >= 10 {
			return false
		}

		time.Sleep(time.Second * 1)
	}

	go func() {
		defer conn.Close()

		// Set up graphics
		// initialize glfw
		if !glfw.Init() {
			panic("Failed to initialize GLFW")
		}
		defer glfw.Terminate()

		// create window
		window, _ = glfw.CreateWindow(600, 600, os.Args[0], nil, nil)
		//	if err != nil {
		//		panic(err)
		//	}
		window.SetFramebufferSizeCallback(onResize)

		window.MakeContextCurrent()
		// set up opengl context
		onResize(window, 600, 600)

		// set up physics
		createBodies()

		runtime.LockOSThread()
		glfw.SwapInterval(1)

		for {
			s.handleUpdate(conn)
			if s.Stop() {
				return
			}
		}
	}()
	return true
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

	var buf []byte = make([]byte, 512)
	conn.ReadFromUDP(buf[0:])
	m := message.PacketToMessage(buf)
	s.frames++

	if m == nil {
		println("a")
		return
	}
	switch m.Type {
	case message.Connect:
		s.Connected(true)
	case message.VectorUpdate:
		println("got a correction")
	case message.StateUpdate:
		balls = make([]*chipmunk.Shape, 0)
		states := message.PayloadToStates(m.Payload)
		for _, state := range states {
			addBall(state.Position.X, state.Position.Y, state.Rotation)
		}
		if window != nil {
			draw()
			window.SwapBuffers()
			glfw.PollEvents()
		}

	default:
		log.Printf("DEFAULT %+v", m)
	}
}

func (s *Ship) update(xdir, ydir int) {
	s.xp += xdir
	s.yp += ydir
}

func (s *Ship) DisplayFrames() {
	log.Printf("------------:%d", s.frames)
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

	num := 100
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
			if m.Type > 0 {
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

	xdir := random(1, 10)
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
