/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/rhkarls/gohbv/hbv"
	"github.com/spf13/cobra"
	"os"
)

var outputFile string

var RunCmd = &cobra.Command{
	Use:   "run [input data file] [parameter file]",
	Short: "Run the HBV model",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		inputDataFile := args[0]
		parameterFile := args[1]

		// Read input forcing data
		inData := hbv.ReadInputData(inputDataFile)
		fmt.Printf("Read input data from %s\n", inputDataFile)

		// Read parameters from json file
		mPars := hbv.ReadParameters(parameterFile)
		fmt.Printf("Read parameters from %s\n", parameterFile)

		// Run model
		hbvResult, err := hbv.RunModel(mPars, inData)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Model run completed\n")

		// Write model output (timeseries)
		hbv.WriteCSVResults(hbvResult, inData, outputFile)
		fmt.Printf("Wrote model results to %s\n", outputFile)
	},
}

func init() {
	rootCmd.AddCommand(RunCmd)

	RunCmd.Flags().StringVarP(&outputFile, "output", "o",
		"gohbv_result.csv", "Output file for the results")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// RunCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// RunCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
