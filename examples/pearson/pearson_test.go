package examples

import (
	"log"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testData struct {
	a []float64
	b []float64
}

func genData(num, length int) []*testData {

	result := make([]*testData, num)
	for i := 0; i < num; i++ {
		result[i] = genDatum(length)
	}

	return result

}

func genDatum(length int) *testData {
	test := &testData{a: make([]float64, length), b: make([]float64, length)}
	for j := 0; j < length; j++ {
		test.a[j] = float64(j) * rand.Float64()
		test.b[j] = float64(j) * rand.Float64()
	}
	return test
}

func TestPearson(t *testing.T) {
	length := 10
	test := genDatum(length)
	result := Pearson(test.a, test.b)
	log.Println("result", result)
	assert.NotNil(t, result)

	if result > 0.5 {
		t.Errorf("Random binary correlation was greater then 50%")
	}
}

var bench float64

func BenchmarkPearson(b *testing.B) {

	data := genDatum(100)
	result := float64(0.0)
	for i := 0; i < b.N; i++ {
		result = Pearson(data.a, data.b)
	}
	bench = result
}

func BenchmarkFasterPearson(b *testing.B) {

	data := genDatum(100)
	result := float64(0.0)
	for i := 0; i < b.N; i++ {
		result = FasterPearson(data.a, data.b)
	}
	bench = result
}
