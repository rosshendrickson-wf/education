package cluster

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKNN(t *testing.T) {

	result := NewTestTrial("iris_headers.csv", 0.7)
	assert.True(t, result)
}

func TestKNNFunc(t *testing.T) {
	filename := "iris_headers.csv"
	data := CSVtoKNNData(filename)
	assert.NotNil(t, data)
	train, test := SplitData(data, 0.7)
	assert.NotNil(t, train)
	assert.NotNil(t, test)
	cls := NewKNNClasifierData(train)
	assert.NotNil(t, cls)
	cls.Predict(test)
}
