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
		length := strconv.Itoa(len(m.Vectors))
		assert.True(t, len(m.Vectors) <= MaxVectors, length)
		for _, v := range m.Vectors {
			assert.NotNil(t, v)
		}
	}
}
