package message

import (
	"bytes"
	"encoding/json"
	"io"
)

// Very VERY naive protocol - totally can do a ton here to get more data
// compressions, go to a byte specific protocol
// 22 will keep the message under 512, our byte limit (based on router issues)
const MaxVectors = 22
const PacketSize = 512

// Message types
const (
	VectorUpdate = 1
	Connect      = 2
	InputUpdate  = 3
)

type Packet []byte

type Message struct {
	Name     int
	Revision int
	Type     int
	Payload  []byte
}

var delim = byte(0)

func ConnectMessage(name, revision int) *Message {
	return &Message{Name: name, Revision: revision, Type: Connect}
}

func MessageToPacket(m *Message) Packet {

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.Encode(m)

	buf.WriteByte(delim)

	return buf.Bytes()
}

func PacketToMessage(p []byte) *Message {
	var in bytes.Buffer
	in.Write(p)

	b, e := in.ReadBytes(delim)
	if e != nil && e != io.EOF {
		if e == io.EOF {
			println("ERR", e.Error())
		}
		println("ERR", e.Error())
	}
	var m Message
	json.Unmarshal(b[:len(b)-1], &m)

	return &m
}
