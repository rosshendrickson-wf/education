package message

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVectorsToMessages(t *testing.T) {

	correct := 3
	v := randVectors(MaxVectors * correct)
	ms := VectorsToMessages(v, 1)
	assert.Equal(t, correct, len(ms))
	for _, m := range ms {
		vectors := PayloadToVectors(m.Payload)
		length := strconv.Itoa(len(vectors))
		assert.True(t, len(vectors) <= MaxVectors, length)
		for _, v := range vectors {
			assert.NotNil(t, v)
		}
	}
}

func TestVectorPayload(t *testing.T) {

	v := randVectors(MaxVectors)
	ms := VectorsToMessages(v, 1)

	for _, m := range ms {
		vectors := PayloadToVectors(m.Payload)
		assert.Equal(t, MaxVectors, len(vectors))
		for i, vector := range vectors {
			assert.NotNil(t, vector)
			assert.Equal(t, v[i], vector)
			assert.Equal(t, v[i].X, vector.X)
		}
	}
}

func TestVectorsRoundTrip(t *testing.T) {

	correct := 3
	v := randVectors(MaxVectors * correct)
	ms := VectorsToMessages(v, 1)
	assert.Equal(t, correct, len(ms))
	var packets []Packet
	for i, m := range ms {
		m.Revision = i
		packets = append(packets, MessageToPacket(m))
	}
	for i, p := range packets {
		expected := ms[i]
		m := PacketToMessage(p)
		vectors := PayloadToVectors(m.Payload)
		length := strconv.Itoa(len(vectors))
		assert.True(t, len(vectors) <= MaxVectors, length)

		eVectors := PayloadToVectors(expected.Payload)
		assert.Equal(t, len(eVectors), len(vectors))
		for j, v := range vectors {
			assert.NotNil(t, v)
			assert.Equal(t, eVectors[j].X, v.X)
		}

		assert.Equal(t, i, m.Revision)
	}
}
