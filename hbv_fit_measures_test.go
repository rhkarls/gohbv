/*
Functions to calculate goodness of fit measure
- Nash-Sutcliffe effciency
- R-squared
- Volume error
- Lindström measure

To be implemented
Tests (also for missing values in obs data)
KGE
log NSE
Spearman rank correlation

For ranking, see argsort for numpy
Go implementation: https://github.com/mkmik/argsort
SO: https://stackoverflow.com/questions/31141202/get-the-indices-of-the-array-after-sorting-in-golang
*/

package gohbv

import (
	"reflect"
	"testing"
)

func TestPowArray(t *testing.T) {
	type args struct {
		arr []float64
		y   float64
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PowArray(tt.args.arr, tt.args.y); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PowArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
