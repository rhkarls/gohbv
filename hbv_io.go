/*
This file contains I/O functions to read and write model and data files,
as well as defitions of structs for data.

To be implemented/wip:
Struct for holding model results.
Read parameters csv (many parameter sets, output array of par struct)
Read parameters nested json (many parameter sets, output array of par struct)

Write timeseries output (wip, need to include dates and refine a bit, exclude warmup time)
Write summary output (input struct with summarized goodness of fit and file path (exclude warmup), output none)
Write batch output (Simulation based on many parameter sets)
*/

package gohbv

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"os"
	"strconv"
)

// Parameters struct
type Parameters struct {
	TT     float64
	CFMAX  float64
	SFCF   float64
	CFR    float64
	CWH    float64
	FC     float64
	LP     float64
	BETA   float64
	PERC   float64
	UZL    float64
	K0     float64
	K1     float64
	K2     float64
	MAXBAS float64
	PCALT  float64
	TCALT  float64
}

// Input forcing and observations
type InputData struct {
	Date          string
	Precipitation float64
	Temperature   float64
	Discharge     float64
	PotentialET   float64
}

// Not yet in use
// type OutputData struct {
// 	Date       string
// 	SimulatedQ float64
// }

// Parse input data csv rows to InputData struct
func ParseInputData(r [][]string) ([]InputData, error) {
	var indata []InputData = make([]InputData, len(r))
	var err error

	for i, row := range r {

		indata[i].Date = row[0]

		indata[i].Precipitation, err = strconv.ParseFloat(row[1], 64)
		if err != nil {
			fmt.Printf("Error converting Precipitation string: %v", err)
		}

		indata[i].Temperature, err = strconv.ParseFloat(row[2], 64)
		if err != nil {
			fmt.Printf("Error converting Temperature string: %v", err)
		}

		indata[i].Discharge, err = strconv.ParseFloat(row[3], 64)
		if err != nil {
			fmt.Printf("Error converting Discharge string: %v", err)
		}

		indata[i].PotentialET, err = strconv.ParseFloat(row[4], 64)
		if err != nil {
			fmt.Printf("Error converting PotentialET string: %v", err)
		}

	}
	return indata, nil
}

// Read csv input forcing data to InputData struct
// CSV data must be in order Date, Precipitation, Temperature, Discharge, Potential ET
func ReadInputData(inputDataFile string) []InputData {
	csvFile, err := os.Open(inputDataFile)
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	_, _ = csvReader.Read() // Read() reads first line which is just the header, not part of data. Needed to skip line for data
	csvLines, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	inData, err := ParseInputData(csvLines)
	if err != nil {
		fmt.Println(err)
	}

	return inData
}

// Parse parameter .json file
// Unmarshal to Parameters struct
func ParseParameters(jsonData []byte) (Parameters, error) {
	var pars Parameters

	err := json.Unmarshal(jsonData, &pars)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	return pars, nil
}

// Read parameter json file and parse it to Parameters struct
func ReadParameters(parameterFile string) Parameters {
	jsonFileObj, err := os.Open(parameterFile)
	if err != nil {
		fmt.Println(err)
	}
	// defer closing so we can parse it later
	defer jsonFileObj.Close()

	// read parameter json as byte
	parsByte, _ := ioutil.ReadAll(jsonFileObj)

	pars, _ := ParseParameters(parsByte)

	return pars
}

// Write model result (array of ModelState structs) to outputFile
func WriteCSVResults(mState []ModelState, outputFile string) {
	fmt.Println(outputFile)
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	err = writer.Write([]string{"S_snow", "Recharge_gwuz", "S_soil", "S_gw_suz", "Q"})
	if err != nil {
		fmt.Println(err)
	}

	for _, value := range mState {
		err := writer.Write([]string{
			strconv.FormatFloat(value.S_snow, 'f', 4, 64),
			strconv.FormatFloat(value.Recharge_gwuz, 'f', 4, 64),
			strconv.FormatFloat(value.S_soil, 'f', 4, 64),
			strconv.FormatFloat(value.S_gw_suz, 'f', 4, 64),
			strconv.FormatFloat(value.Q_sim, 'f', 4, 64)})
		if err != nil {
			fmt.Println(err)
		}
	}

}
