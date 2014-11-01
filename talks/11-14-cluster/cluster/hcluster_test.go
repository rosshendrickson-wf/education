package cluster

import (
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var hTest struct {
	p *Point
}

func genPoints(num, max int) []*Point {

	if max == 0 {
		max = 1
	}
	numV := rand.Intn(max)
	//log.Println(numV)
	result := make([]*Point, num)
	for i := 0; i < num; i++ {
		values := make([]float64, numV)
		for j, _ := range values {
			values[j] = float64(j) * rand.Float64()
		}
		log.Printf("len %d, %+v", len(values), values)
		result[i] = &Point{Key: strconv.Itoa(i), Values: values}
	}

	return result
}

func TestHCluster(t *testing.T) {
	clusters := HCluster(genPoints(10, 30))
	//output, _ := json.Marshal(clusters)
	//output, _ := json.MarshalIndent(clusters, "", "    ")

	clusterKeys := ClusterKeys(clusters...)
	assert.Equal(t, 9, len(clusters), string(clusterKeys))
}

func BenchmarkHCluster(b *testing.B) {
	data := genPoints(100, 20)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HCluster(data)
	}
}

func init() {
	log.SetOutput(ioutil.Discard)
}
