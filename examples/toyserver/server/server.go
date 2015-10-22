package toyserver

import (
	"net"

	"github.com/vova616/chipmunk"
)

// Server
// connect
// spin up shards
//	determine size
// shards connect over TCP

// Master
//  Main loop
//
//  listen to messages -> Push into shard queue based on pos
//  every 3ms
//    compact and send to shard for processing
//
////////////////////////////////////////
//
// Shard
//  Connect TCP then
//  Connect UDP
//
//      // Go routine listening to UDP
//		shard sorts messages into frame buckets that are >= than current
//		current frame revision
//
//  Main Loop
//
//      Loop
//		Check if we are supposed to calc new frame (Read TCP)
//      Check if new frame has shared entities
//			Loop
//				Send Ack that we know to bump revision (Write TCP)
//              Bump Revision
//
//
//		Frame Process Stage
//			applies updates to user movements for the current frame revision
//          calculate collisions and state
//          based on current POS of items do shard share phase
//
//			Frame Share phase
//				Max item size x or y -> Max Heap eventually
//				for this frame -> N number of entities need to exist in multiple
//				Shards
//
//      Post Frame Process Stage
//			Remove from physics stuff that is out of frame & buffers
//
//
//		Sends update revision finished (Write TCP)
//
//      Loop
//		Check frame calc was acked
//
//
// Cluster Consensus
//  M) Calc frame N -> A) B) C) D)
//  All) Calc frame N ack
//  A) I've calced N
//  M) You've calced frame N?
//  A) Yes I've calced N
//  M) A has calced frame N
//  B)  ....
//  M) Calc frame N+1

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
