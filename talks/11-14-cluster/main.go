package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"net/http"
	"runtime"

	"github.com/rosshendrickson-wf/education/talks/11-14-cluster/cluster"
)

func main() {

	data := cluster.CSVtoPoints("load.csv")
	var classify = func(w http.ResponseWriter, r *http.Request) {
		cls := cluster.NewKNNClassifier()
		cls.Train(data...)
		// read of URL params
		//Sepal length, Sepal width,Petal length, Petal width, Species

		values := r.URL.Query()
		sl := values.Get("sl")
		sw := values.Get("sw")
		pl := values.Get("pl")
		pw := values.Get("pw")

		fvalues := []float64{cluster.GetValue(sl), cluster.GetValue(sw),
			cluster.GetValue(pl), cluster.GetValue(pw)}
		point := cluster.NewPoint("", r.URL.String(), fvalues)
		fmt.Fprintf(w, "%+v!\n", sl)
		fmt.Fprintf(w, "%+v!\n", sw)
		fmt.Fprintf(w, "%+v!\n", pl)
		fmt.Fprintf(w, "%+v!\n", pw)
		fmt.Fprintf(w, "Predicted Class %s", cls.Predict(point, 5))
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
