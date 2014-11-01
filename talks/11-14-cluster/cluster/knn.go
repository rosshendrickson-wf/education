package cluster

import (
	"fmt"
	"sort"

	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/knn"
)

// Calc distance for all points
// find N nearest to new "point"
// Have the N vote, naive

type KNNClassifier struct {
	cache     map[string]*Distance
	points    []*Point
	distances []*Distance
}

func NewKNNClassifier() *KNNClassifier {

	cache := make(map[string]*Distance)
	distances := make([]*Distance, 0)

	return &KNNClassifier{cache: cache, distances: distances}

}

func (k *KNNClassifier) Predict(point *Point, kn int) string {

	ns := make([]*Distance, kn)
	ds := make([]*Distance, len(k.points))

	// Compute distances against existing points
	for i, p := range k.points {
		key := mergedKey(point.Key, p.Key)
		distance := NewDistance(point, p)
		k.cache[key] = distance
		ds[i] = distance
	}

	// Sort to find top k
	sort.Sort(ByDistance(ds))

	for i := 0; i < kn; i++ {
		ns[i] = ds[i]
	}

	// Vote
	return k.vote(ns...)
}

func (k *KNNClassifier) vote(distances ...*Distance) string {

	classes := make([]string, len(distances))
	classCount := make(map[string]int, 1)
	// Count the votes from the neighbors
	for i, distance := range distances {
		classes[i] = distance.P2.Class

		count := classCount[distance.P2.Class]
		if count == 0 {
			classCount[distance.P2.Class] = 1
		} else if count > 0 {
			classCount[distance.P2.Class] += 1
		}

	}

	// Return the most voted class
	max := 0
	class := ""
	for k, c := range classCount {
		if c > max {
			max = c
			class = k
		}
	}
	return class
}

func (k *KNNClassifier) Train(points ...*Point) {

	keyIndex := make(map[string]*Point, len(points))
	ps := make([]*Point, len(points))
	for i, point := range points {
		keyIndex[point.Key] = point
		ps[i] = point
	}

	for len(keyIndex) > 1 {
		for _, point := range keyIndex {
			for _, p := range keyIndex {
				key := mergedKey(point.Key, p.Key)
				cached := k.cache[key]
				if cached != nil {
					continue
				}
				if point == p {
					continue
				}
				distance := NewDistance(point, p)
				k.cache[key] = distance
				k.distances = append(k.distances, distance)
			}
			delete(keyIndex, point.Key)
		}
	}
	k.points = ps
}

type ByDistance []*Distance

func (a ByDistance) Len() int           { return len(a) }
func (a ByDistance) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDistance) Less(i, j int) bool { return a[i].Distance < a[j].Distance }

func CSVtoKNNData(filename string) base.FixedDataGrid {
	rawData, err := base.ParseCSVToInstances(filename, true)
	if err != nil {
		panic(err)
	}
	return rawData
}

func SplitData(data base.FixedDataGrid, split float64) (base.FixedDataGrid, base.FixedDataGrid) {
	return base.InstancesTrainTestSplit(data, split)
}

func NewTestTrial(filename string, split float64) bool {
	cls := knn.NewKnnClassifier("euclidean", 2)
	data := CSVtoKNNData(filename)
	train, test := base.InstancesTrainTestSplit(data, split)

	cls.Fit(train)
	//Calculates the Euclidean distance and returns the most popular label
	predictions := cls.Predict(test)
	fmt.Println(predictions)

	confusionMat, err := evaluation.GetConfusionMatrix(test, predictions)
	if err != nil {
		panic(fmt.Sprintf("Unable to get confusion matrix: %s", err.Error()))
	}
	fmt.Println(evaluation.GetSummary(confusionMat))

	return true
}

func NewKNNClasifierFile(filename string) *knn.KNNClassifier {
	//Initialises a new KNN classifier
	data := CSVtoKNNData(filename)
	return NewKNNClasifierData(data)
}

func NewKNNClasifierData(data base.FixedDataGrid) *knn.KNNClassifier {
	//Initialises a new KNN classifier
	cls := knn.NewKnnClassifier("euclidean", 2)
	cls.Fit(data)
	return cls
}
