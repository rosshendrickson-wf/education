package message

import (
	"encoding/json"
	"math/rand"
	"time"
)

type Vector struct {
	X int
	Y int
}

type Vec struct {
	X, Y float32
}

type VectorPayload struct {
	Vectors []*Vector
}

// VectorsToMessages will take a list of vectors and split them up into
// as many messages as are needed. Room for speeding this up
func VectorsToMessages(vectors []*Vector, name int) []*Message {

	if len(vectors) <= MaxVectors {

		//	log.Printf("Vec: %+v", vectors)
		//	log.Printf("Vec: %+v", vectors[0])

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

func randVecs(num int) []*Vec {

	results := make([]*Vec, 0)
	for i := 0; i < num; i++ {
		xdir := randomFloat(1.0)
		ydir := randomFloat(1.0)

		v := &Vec{xdir, ydir}

		results = append(results, v)
	}
	return results
}

func randomFloat(min float32) float32 {
	rand.Seed(time.Now().Unix())
	return rand.Float32() * min

}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
