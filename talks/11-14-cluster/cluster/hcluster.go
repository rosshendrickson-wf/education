package cluster

import (
	"log"
	"math"
)

// Heirachal clustering
// values grouped to some id

type Point struct {
	key    string
	values []float64
}

type Cluster struct {
	Points []*Point
}

type Distance struct {
	p1       *Point
	p2       *Point
	distance float64
}

// for every point in the set of data compute a distance to every other set
// group the two closest into a new point with averaged values
// repeat until all points are within a single cluster
func HCluster(points []*Point) []*Cluster {

	cache := make(map[string]*Distance, len(points))
	keyIndex := make(map[string]*Point, len(points))
	for _, point := range points {
		keyIndex[point.key] = point
	}
	results := make([]*Cluster, 1)
	// As things are merged we need to remove them from the points slice
	remaining := len(points)
	for remaining > 1 {
		var closest *Distance
		for _, point := range points {
			for _, p := range points {
				cached := cache[mergedKey(point.key, p.key)]
				if cached != nil {
					if closest != nil && cached.distance > closest.distance {
						closest = cached
					}
					continue
				}
				distance := NewDistance(point, p)
				cache[mergedKey(point.key, p.key)] = distance
				if closest == nil {
					closest = distance
				}
				if closest != nil && distance.distance > closest.distance {
					closest = distance
				}
			}
		}
		if closest == nil {
		}
		points = append(points, AveragedPoints(closest))
		remaining--
		results = append(results, NewCluster(closest.p1, closest.p2))
	}
	return results
}

func NewCluster(points ...*Point) *Cluster {
	return &Cluster{Points: points}
}

func AveragedPoints(closest *Distance) *Point {

	values := make([]float64, len(closest.p1.values))
	result := &Point{key: closest.p1.key + closest.p2.key, values: values}
	for i, value := range closest.p1.values {
		v2 := closest.p2.values[i]
		values[i] = (value + v2) / 2
	}

	return result
}

func NewDistance(p1, p2 *Point) *Distance {
	return &Distance{p1: p1, p2: p2, distance: Pearson(p1.values, p2.values)}
}

func mergedKey(k1, k2 string) string {
	return k1 + ":" + k2
}

func Pearson(v1 []float64, v2 []float64) float64 {

	sum1 := float64(0.0)
	sum2 := float64(0.0)
	sum1Sq := float64(0.0)
	sum2Sq := float64(0.0)
	pSum := float64(0.0)
	var value2 float64
	for i, value := range v1 {
		value2 = v2[i]
		// Simple Sums
		sum1 += value
		sum2 += value2

		// Sum of the squares
		sum1Sq += math.Pow(value, 2)
		sum2Sq += math.Pow(value2, 2)

		// Sum of the products
		pSum += value * value2
	}

	// Calculate r (Pearson score)
	viLen := float64(len(v1))
	num := pSum - sum1*sum2/viLen
	den := math.Sqrt((sum1Sq - math.Pow(sum1, 2)/viLen) * (sum2Sq - math.Pow(sum2, 2)/viLen))

	if den == 0 {
		log.Println(den)
		return 0
	}
	return 1.0 - num/den
}
