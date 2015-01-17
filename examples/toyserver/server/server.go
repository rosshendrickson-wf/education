package toyserver

import (
	"net"

	"github.com/vova616/chipmunk"
)

type Client struct {
	inbound  chan []byte
	outbound chan []byte
	socket   net.Conn
	state    *widget
}

type widget struct {
	shape  *chipmunk.Shape
	health int
}

// client observer - reads from the UDP socket and shoves the updates into the
// transaction log. timestamps when the server recieved them, calculates client
// latency through pulses of information (round trip to-from) not through
// timestamps. Shards are used to track which clients are in which areas

func Observe(Client) {
	for {

	}
}

type Server struct {
	Shards  []*Shard
	clients map[string]*Client
}

type Shard struct {
	clients  []*Client
	inbound  chan []byte
	outbound chan []byte
	space    *chipmunk.Space
	xw       int
	yw       int
	xpos     int
	ypos     int
}

func (s *Shard) step() {

	// compute updates from incoming stream of information vs clients
}

// move a client between shards
