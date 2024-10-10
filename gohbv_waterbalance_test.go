package main_test

import (
	"github.com/rhkarls/gohbv/hbv"
	"math"
	"testing"
)

func TestWaterBalanceAccounting(t *testing.T) {
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

	// Account for the initial states that are not included in the input data
	smAddedState := mPars.FC * mPars.LP
	lzAddedState := mPars.PERC / mPars.K2

	// Account for the snow routine SFCF loss
	// TODO write this loss as a column to output??
	// snowRoutineLoss is the difference between input precipitation and modelled snowfall, if snowfall >0
	snowRoutineLoss := 0.0
	for i, data := range inData {
		if hbvResult[i].Snowfall > 0 {
			snowRoutineLoss += data.Precipitation - hbvResult[i].Snowfall
		}
	}

	// Calculate the sum of the input data Precipitation field
	var sumPrecipitation float64
	for _, data := range inData {
		sumPrecipitation += data.Precipitation
	}

	// Calculate the sum of the outputs and storages
	var sumOutputs, sumStorages float64
	for _, state := range hbvResult {
		sumOutputs += state.Q_sim + state.AET
		sumStorages = state.S_snow + state.S_soil + state.S_gw_suz + state.S_gw_slz
	}

	sumOutputs += snowRoutineLoss
	sumStorages -= smAddedState + lzAddedState
	// Compare the two sums
	// Should not be > 1 mm difference
	if math.Abs(sumPrecipitation-(sumOutputs+sumStorages)) > 1.0 {
		t.Errorf("Failed: Water balance mismatch: sum of inputs = %f, sum of outputs and storages = %f",
			sumPrecipitation, sumOutputs+sumStorages)
	} else {
		t.Logf("Passed: Water balance accounted for: sum of inputs = %f, sum of outputs and storages = %f",
			sumPrecipitation, sumOutputs+sumStorages)
	}
}
