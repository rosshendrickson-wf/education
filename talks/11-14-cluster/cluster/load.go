package cluster

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func CSVtoPoints(filename string) []*Point {

	result := make([]*Point, 0)

	csvfile, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer csvfile.Close()

	reader := csv.NewReader(csvfile)

	reader.FieldsPerRecord = -1 // see the Reader struct information below

	rawCSVdata, err := reader.ReadAll()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// sanity check, display to standard output
	for i, each := range rawCSVdata {
		// Header
		if i == 0 {
			continue
		}
		values := make([]float64, len(each)-1)
		class := each[len(each)-1]
		for j, v := range each[:len(each)-1] {
			f, e := strconv.ParseFloat(v, 64)
			if e == nil {
				values[j] = f
			}
		}
		result = append(result, &Point{Class: class, Key: strconv.Itoa(i), Values: values})
	}

	return result
}

func GetValue(v string) float64 {
	f, e := strconv.ParseFloat(v, 64)
	if e == nil {
		return f
	}
	return float64(0)
}
