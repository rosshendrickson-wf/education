package message

import (
	"bytes"
	"math/rand"
	"time"

	"encoding/json"
	"log"
)

// Very VERY naive protocol - totally can do a ton here to get more data
// compressions, go to a byte specific protocol
// 30 will keep the message under 512, our byte limit (based on router issues)
const MaxVectors = 32

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

// VectorsToMessages will take a list of vectors and split them up into
// as many messages as are needed
func VectorsToMessages(vectors []*Vector, name int) []*Message {

	if len(vectors) <= MaxVectors {

		results := make([]*Message, 0)
		results = append(results, &Message{Name: name, Vectors: vectors})
		return results
	}

	numMessages := len(vectors) / MaxVectors
	results := make([]*Message, numMessages)

	for i := range results {
		m := &Message{Vectors: make([]*Vector, 0)}
		j := 0
		for len(m.Vectors) < MaxVectors || j < len(vectors) {
			m.Vectors = append(m.Vectors, vectors[j])
			j++
		}
		results[i] = m
	}

	return results
}

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

func randVectors(num int) []*Vector {

	results := make([]*Vector, num)
	for i := range results {
		xdir := random(1, 6)
		ydir := random(1, 10)

		v := &Vector{xdir, ydir}

		results[i] = v
	}

	return results
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
