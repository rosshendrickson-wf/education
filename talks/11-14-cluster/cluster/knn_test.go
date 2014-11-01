package cluster

import (
	"log"
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

func TestSimpleKNN(t *testing.T) {

	result := CSVtoPoints("iris_headers.csv")
	assert.NotNil(t, result)
	assert.Equal(t, 150, len(result))

	cls := NewKNNClassifier()
	assert.NotNil(t, cls)

	cls.Train(result...)
	assert.Equal(t, len(result), len(cls.points))

	prediction := cls.Predict(result[0], 5)
	assert.NotNil(t, prediction)
	assert.NotEqual(t, "", prediction)
	assert.Equal(t, result[0].Class, prediction)
}

func TestTrainSplitSimple(t *testing.T) {

	result := CSVtoPoints("iris_headers.csv")
	assert.NotNil(t, result)
	assert.Equal(t, 150, len(result))

	cls := NewKNNClassifier()
	assert.NotNil(t, cls)

	cls.Train(result[:120]...)
	assert.Equal(t, 120, len(cls.points))

	wrong := 0
	for i, point := range result[121:] {

		prediction := cls.Predict(point, 5)
		if prediction != result[i+121].Class {
			wrong++
		}
	}
	percent := float64(wrong) / float64(30)
	log.Printf("Error %f", percent)
	log.Printf("Wrong %d", wrong)

}
