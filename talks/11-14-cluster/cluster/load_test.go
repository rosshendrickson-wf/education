package cluster

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadCSV(t *testing.T) {

	result := CSVtoPoints("iris_headers.csv")
	assert.NotNil(t, result)
	assert.Equal(t, 150, len(result))

	point := result[0]

	assert.NotNil(t, point.Class)
	assert.Equal(t, 4, len(point.Values), fmt.Sprintf("%+v", point))
}
