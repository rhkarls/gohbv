// Example running gohbv from Go
package main

import (
	"fmt"
	"github.com/rhkarls/gohbv/hbv"
	"time"
)

const input_data_file = "hbv_data/go_hbv_input.csv"
const parameter_file = "hbv_data/hbv_parameters.json"

func main() {

	// Read input forcing data
	inData := hbv.ReadInputData(input_data_file)

	// Read parameters from json file
	mPars := hbv.ReadParameters(parameter_file)

	// Run model 10000 times for timing it (~0.5 s for 10 000, 0.0005 s/run)
	// fmt printing and IO will slow this down significantly unless using goroutines
	t_s := time.Now()
	for i := 1; i < 10000; i++ {
		_, _ = hbv.RunModel(mPars, inData)
	}
	t_e := time.Now()
	fmt.Printf("\nThe call took %v ns to run 10 000 HBV calculation.\n", 1*(t_e.UnixNano()-t_s.UnixNano()))

	// Run model
	hbv_result, err := hbv.RunModel(mPars, inData)
	if err != nil {
		fmt.Println(err)
	}

	// Write model output (timeseries)
	hbv.WriteCSVResults(hbv_result, inData, "gohbv_model_results.csv")

	// Write a summary file of the simulation
}
