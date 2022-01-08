// Example running gohbv from Go
package main

import (
	"fmt"
	"time"

	"github.com/rhkarls/gohbv"
)

const input_data_file = "hbv_data/go_hbv_input.csv"
const parameter_file = "hbv_data/hbv_parameters.json"

func main() {

	// Read input forcing data
	inData := gohbv.ReadInputData(input_data_file)

	// Read parameters from json file
	mPars := gohbv.ReadParameters(parameter_file)

	// Run model 10000 times for timing it (~0.5 s for 10 000, 0.0005 s/run)
	// fmt printing and IO will slow this down significantly unless using goroutines
	t_s := time.Now()
	for i := 1; i < 10000; i++ {
		_, _ = gohbv.RunModel(mPars, inData)
	}
	t_e := time.Now()
	fmt.Printf("\nThe call took %v ns to run 10 000 HBV calculation.\n", 1*(t_e.UnixNano()-t_s.UnixNano()))

	// Run model
	hbv_result, err := gohbv.RunModel(mPars, inData)
	if err != nil {
		fmt.Println(err)
	}

	// Write model output (timeseries)
	gohbv.WriteCSVResults(hbv_result, "gohbv_model_results.csv")

	// Write a summary file of the simulation
}
