package cluster

import "math"

// Heirachal clustering
// values grouped to some id

type Point struct {
	Class  string
	Key    string
	Values []float64
}

type Cluster struct {
	Key    string
	Points []*Point
}

type Distance struct {
	P1       *Point
	P2       *Point
	Distance float64
}

// for every point in the set of data compute a distance to every other set
// group the two closest into a new point with averaged values
// repeat until all points are within a single cluster
func HCluster(points []*Point) []*Cluster {

	cache := make(map[string]*Distance, len(points))
	//eIndex := make(map[string]*list.Element, len(points))
	keyIndex := make(map[string]*Point, len(points))
	for _, point := range points {
		keyIndex[point.Key] = point
	}

	results := make([]*Cluster, 0)
	// As things are merged we need to remove them from the points slice
	for len(keyIndex) > 1 {
		var closest *Distance
		for _, point := range keyIndex {
			for _, p := range keyIndex {
				cached := cache[mergedKey(point.Key, p.Key)]
				if cached != nil {
					if closest != nil && cached.Distance > closest.Distance {
						println("A")
						closest = cached
					}
					continue
				}

				if point == p {
					continue
				}
				//log.Println("CL", closest)

				distance := NewDistance(point, p)
				//log.Printf("new Distance %+v", distance)
				//cache[mergedKey(point.key, p.key)] = distance
				if closest == nil {
					//log.Println("B")
					closest = distance
					continue
				}
				if closest != nil && distance.Distance > closest.Distance {
					//log.Println("C")
					closest = distance
				}
			}
		}
		if closest == nil {
			//log.Println("NIL closest")
			//return nil
			continue
		}
		// Remove merged points
		delete(keyIndex, closest.P1.Key)
		delete(keyIndex, closest.P2.Key)
		// Add new merged node
		newPoint := AveragedPoints(closest)
		keyIndex[newPoint.Key] = newPoint

		//log.Printf("New Cluster (%s)(%s)", closest.P1.Key, closest.P2.Key)
		results = append(results, NewCluster(closest.P1, closest.P2))
	}
	return results
}

func NewPoint(class string, key string, values []float64) *Point {
	return &Point{Class: class, Key: key, Values: values}
}

func NewCluster(points ...*Point) *Cluster {
	key := ""
	c := &Cluster{}
	c.Points = make([]*Point, len(points))
	for i, point := range points {
		key += point.Key
		key += ":"
		c.Points[i] = point
	}
	c.Key = key

	return c
}

func AveragedPoints(closest *Distance) *Point {
	if closest == nil {
		return nil
	}
	values := make([]float64, len(closest.P1.Values))
	result := &Point{Key: mergedKey(closest.P1.Key, closest.P2.Key), Values: values}
	for i, value := range closest.P1.Values {
		v2 := closest.P2.Values[i]
		values[i] = (value + v2) / 2
	}

	return result
}

func NewDistance(P1, P2 *Point) *Distance {
	//log.Printf("calc distance for (%s)(%s)", P1.Key, P2.Key)
	return &Distance{P1: P1, P2: P2, Distance: Pearson(P1.Values, P2.Values)}
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
		//log.Println(den)
		return 0
	}
	//log.Println("distance was ", 1.0-num/den)
	return 1.0 - num/den
}

func ClusterKeys(clusters ...*Cluster) string {

	all := "\n"
	for _, cluster := range clusters {
		all += cluster.Key
		all += "\n"
	}
	return all
}
