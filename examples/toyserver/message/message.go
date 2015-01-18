package message

import (
	"bytes"
	"encoding/json"
	"log"
)

const MaxVector = 30

type Packet []byte

type Vector struct {
	X int
	Y int
}

type Message struct {
	Name    int
	Vectors []*Vector
	Value   string
}

var delim = byte(0)

func MessageToPacket(m *Message) Packet {

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.Encode(m)
	println("len m", buf.Len())
	buf.WriteByte(delim)
	packet := make([]byte, 512)
	println("len p", len(packet))
	n, err := buf.Read(packet)
	log.Println("n err", n, err)
	return packet
}

func PacketToMessage(p []byte) *Message {
	println("len p", len(p))
	var in bytes.Buffer
	in.Write(p)

	b, e := in.ReadBytes(delim)
	if e != nil {
		println(e)
	}
	println("len b", len(b))
	var m Message
	json.Unmarshal(b[:len(b)-1], &m)

	return &m
}
