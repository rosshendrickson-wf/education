package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"net/http"
	"runtime"
)

func main() {

	//data := cluster.CSVtoKNNData("load.csv")

	var classify = func(w http.ResponseWriter, r *http.Request) {
		//cls := cluster.NewKNNClasifierData(data)

		// read of URL params
		//Sepal length, Sepal width,Petal length, Petal width, Species

		values := r.URL.Query()
		sl := values.Get("sl")
		sw := values.Get("sw")
		pl := values.Get("pl")
		pw := values.Get("pw")

		// Make them into grid
		// Create some Attributes
		//Add a row
		//		values := []string{"5.1", "3.5", "1.4", "0.2", "Iris-setosa"}
		//		// Create a csv file
		//		hackName := "./hack.csv"
		//		f, _ := os.Create(hackName)
		//		csvw := csv.NewWriter(f)
		//		csvw.Write(values)
		//		csvw.Flush()
		//		f.Close()
		//		newInst := cluster.CSVtoKNNData(hackName)
		//		// Predict against them
		//		fmt.Println(newInst)
		//		predictions := cls.Predict(newInst)
		//		fmt.Println(predictions)
		//
		// Process prediction into something we can write
		fmt.Fprintf(w, "Hi there, I love %+v!", sl)
		fmt.Fprintf(w, "Hi there, I love %+v!", sw)
		fmt.Fprintf(w, "Hi there, I love %+v!", pl)
		fmt.Fprintf(w, "Hi there, I love %+v!", pw)

	}

	http.HandleFunc("/predict/", classify)
	panic(http.ListenAndServe(":8080", nil))
}

func classifyUser(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func float64bytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
