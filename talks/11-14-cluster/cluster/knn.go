package cluster

import (
	"fmt"

	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/knn"
)

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
