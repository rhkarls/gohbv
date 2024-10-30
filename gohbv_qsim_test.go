package main_test

import (
	"encoding/csv"
	"github.com/rhkarls/gohbv/hbv"
	"math"
	"os"
	"strconv"
	"testing"
)

func TestRunModelQSim(t *testing.T) {
	// Read input data
	inputDataFile := "test_data/test_case_sp0_go_hbv_input.csv"
	inData := hbv.ReadInputData(inputDataFile)

	// Read parameters
	parameterFile := "test_data/test_case_sp0_hbv_parameters.json"
	mPars := hbv.ReadParameters(parameterFile)

	// Run model
	hbvResult, err := hbv.RunModel(mPars, inData)
	if err != nil {
		t.Fatalf("Error running model: %v", err)
	}

	// Read benchmark data
	benchmarkFile := "test_data/test_case_sp0_results_HBV_Light_2014_2020_Benchmark.txt"
	file, err := os.Open(benchmarkFile)
	if err != nil {
		t.Fatalf("Error opening benchmark file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Error reading benchmark file: %v", err)
	}

	// Compare Q_sim values
	var sumModel, sumBenchmark float64
	n_failed := 0
	n_passed := 0

	for i, record := range records {
		if i <= 1 {
			continue // Skip header, kip first row sensitive to initial conditions
		}

		qSimBenchmark, err := strconv.ParseFloat(record[1], 64) // Assuming Qsim is the 2nd column
		if err != nil {
			t.Fatalf("Error parsing benchmark Qsim value: %v", err)
		}
		qSimModel := hbvResult[i-1].Q_sim

		// each individual time step value should not have an absolute error above 0.001 mm
		// *and* a relative error above 0.5 % (might be more reasonable to set to 1%?)
		// i.e. very small values can get large relative error, but is very close absolute
		// very large values can get large absolute errors, but these are small relative errors
		if (math.Abs(qSimModel-qSimBenchmark) > 0.001) && ((qSimModel-qSimBenchmark)/qSimBenchmark*100 > 0.5) {
			t.Errorf("Mismatch at row %d: model Q_sim = %f, benchmark Q_sim = %f", i, qSimModel, qSimBenchmark)
			n_failed++
		} else {
			n_passed++
		}

		sumModel += qSimModel
		sumBenchmark += qSimBenchmark
	}

	t.Logf("Passed Q simulated timestep tests: %d, Failed tests: %d\n", n_passed, n_failed)

	if math.Abs(sumModel-sumBenchmark) > 0.1 {
		t.Errorf("Failed Q sum test: model sum = %f, benchmark sum = %f", sumModel, sumBenchmark)
	} else {
		t.Logf("Passed Q sum test: model sum = %f, benchmark sum = %f", sumModel, sumBenchmark)
	}

}
