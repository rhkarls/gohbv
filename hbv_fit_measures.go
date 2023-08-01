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
	"math"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"
)

func Q_obs_to_array(inData []InputData) []float64 {
	var Q_values []float64
	for _, value := range inData {
		Q_values = append(Q_values, value.Discharge)
	}
	return Q_values
}

func Q_sim_to_array(mState []ModelState) []float64 {
	var Q_values []float64
	for _, value := range mState {
		Q_values = append(Q_values, value.Q_sim)
	}
	return Q_values
}

// Element wise power of items in array
// arr**y
func PowArray(arr []float64, y float64) []float64 {
	arr_pow := make([]float64, len(arr))

	for _, value := range arr {
		arr_pow = append(arr_pow, math.Pow(value, y))
	}

	return arr_pow
}

// Calculate Nash-Sutcliffe efficiency
// 1 - ( sum((obs-sim)**2) / sum((obs-mean_obs)**2) )
func NashSutcliffeEfficiency(sim []float64, obs []float64) float64 {

	obs_mean := stat.Mean(obs, nil)
	obs_sub_mean := obs
	floats.AddConst(-obs_mean, obs_sub_mean)

	obs_sub_sim := make([]float64, len(obs))
	floats.SubTo(obs_sub_sim, obs, sim)

	nse := 1 - (floats.Sum(PowArray(obs_sub_sim, 2)) / floats.Sum(PowArray(obs_sub_mean, 2)))

	return nse
}

func RSquared(sim []float64, obs []float64) float64 {

	r_sq := stat.RSquaredFrom(sim, obs, nil)

	return r_sq
}

// Calculate the relative volume error (i.e. bias)
// This evalulates systematic volume errors over longer periods
// See Lindström, G. (1997). A Simple Automatic Calibration Routine for the HBV Model.
// Nordic Hydrology, 28(3), 153–168. https://doi.org/10.2166/nh.1997.009
// Note: returns single float64, alternativelty return array of dv ?
func VolumeError(sim []float64, obs []float64) float64 {
	// Note: here we take the cumulative difference, not accumulating over an array
	sim_sub_obs := make([]float64, len(obs))
	floats.SubTo(sim_sub_obs, sim, obs)

	cum_diff := floats.Sum(sim_sub_obs)

	dv := cum_diff / floats.Sum(obs)

	return dv
}

// Penalize the NSE measure with remaining weighted volume error
func LindstromMeasure(sim []float64, obs []float64, w float64) float64 {
	// NSE - w*abs(dv)
	nse := NashSutcliffeEfficiency(sim, obs)
	dv := VolumeError(sim, obs)

	lm := nse - w*math.Abs(dv)

	return lm

}
