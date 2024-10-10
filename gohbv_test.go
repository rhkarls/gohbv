package main_test

import (
	"os"
	"testing"

	"github.com/rhkarls/gohbv/hbv/cmd"
)

func TestRunCommand(t *testing.T) {
	// Save the original command-line arguments
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	// Set up command-line arguments
	os.Args = []string{"gohbv", "run", "test_data/test_case_sp0_go_hbv_input.csv", "test_data/test_case_sp0_hbv_parameters.json", "-o", "test_output.csv"}

	// Execute the command
	err := cmd.RunCmd.Execute()
	if err != nil {
		t.Fatalf("Error executing run command: %v", err)
	}
}

func BenchmarkRunCommand(b *testing.B) {
	// Save the original command-line arguments
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	// Set up command-line arguments
	os.Args = []string{"gohbv", "run", "test_data/test_case_sp0_go_hbv_input.csv", "test_data/test_case_sp0_hbv_parameters.json", "-o", "test_output.csv"}

	// Reset the timer to exclude setup time from the benchmark
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Execute the command
		err := cmd.RunCmd.Execute()
		if err != nil {
			b.Fatalf("Error executing run command: %v", err)
		}
	}
}
