package message

import (
	"bytes"
	"io"
	"log"
	"math/rand"
	"time"

	"encoding/json"
)

// Very VERY naive protocol - totally can do a ton here to get more data
// compressions, go to a byte specific protocol
// 30 will keep the message under 512, our byte limit (based on router issues)
const MaxVectors = 22
const PacketSize = 512
const VectorUpdate = "Vector"

type Packet []byte

type Vector struct {
	X int
	Y int
}

type Message struct {
	Name     int
	Revision int
	Type     string
	Payload  []byte
}

type VectorPayload struct {
	Vectors []*Vector
}

var delim = byte(0)

// VectorsToMessages will take a list of vectors and split them up into
// as many messages as are needed. Room for speeding this up
func VectorsToMessages(vectors []*Vector, name int) []*Message {

	if len(vectors) <= MaxVectors {

		log.Printf("Vec: %+v", vectors)
		log.Printf("Vec: %+v", vectors[0])

		results := make([]*Message, 1)
		vec := VectorPayload{vectors}
		payload, _ := json.Marshal(vec)
		results[0] = &Message{Name: name, Payload: payload}
		return results
	}

	numMessages := len(vectors) / MaxVectors
	results := make([]*Message, numMessages)

	j := 0
	for i := range results {
		m := &Message{}
		chunk := make([]*Vector, 0)
		for len(chunk) < MaxVectors && j < len(vectors)-1 {
			chunk = append(chunk, vectors[j])
			j++
		}
		vec := VectorPayload{chunk}
		payload, _ := json.Marshal(vec)
		m.Payload = payload
		m.Type = VectorUpdate
		results[i] = m
	}

	return results
}

func PayloadToVectors(payload []byte) []*Vector {

	var vPayload VectorPayload
	json.Unmarshal(payload, &vPayload)

	return vPayload.Vectors
}

func MessageToPacket(m *Message) Packet {

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.Encode(m)
	//println("len m1", buf.Len())

	//	packet := make([]byte, PacketSize)
	//	buf.Read(packet)
	//	n, err := buf.Read(packet)
	buf.WriteByte(delim)

	//println("len m2", buf.Len())

	//	println("n err", n, err)
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
