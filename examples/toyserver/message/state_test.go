package message

import (
	"log"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatesToMessages(t *testing.T) {

	correct := 3
	v := randState(maxState * correct)

	log.Printf("--%+v", v)

	for _, s := range v {

		log.Printf("++%v", s.Position)
	}
	ms := StatesToMessages(v)
	assert.Equal(t, correct, len(ms))
	for _, m := range ms {
		vectors := PayloadToStates(m.Payload)
		log.Printf("%+v", vectors)
		log.Printf("%+v", m.Payload)
		length := strconv.Itoa(len(vectors.States))
		assert.True(t, len(vectors.States) <= maxState, length)
		assert.True(t, len(vectors.States) > 0)
		for _, v := range vectors.States {
			assert.NotNil(t, v)
		}
	}

}

func randState(num int) []*State {
	results := make([]*State, 0)
	vectors := randVecs(10)
	for i, vec := range vectors {
		println(i)
		println("000000000000")
		results = append(results, &State{Kind: i, UUID: i, Position: vec})
	}

	return results
}
