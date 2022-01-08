/*
This file contains HBV model routines for
- Snow accumulation, storage and melt
- Soil moisture and evaporation
- Groundwater storage and runoff
- Routing
*/

package gohbv

import (
	"math"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/integrate"
)

// Helper function to calculate MAXBAS triangular weights
// This function outputs the values that are integrated in RoutingMaxbasWeights
func routingMaxbas(x []float64, p_maxbas float64) []float64 {
	a := 2 / p_maxbas
	c := 4 / (math.Pow(p_maxbas, 2))

	var maxbas_x = make([]float64, len(x))

	for i := 0; i < len(x); i++ {
		maxbas_x[i] = a - math.Abs(float64(x[i])-p_maxbas/2.0)*c
	}

	return maxbas_x
}

// Calculate MAXBAS triangular weights using trapezoidal integration
func RoutingMaxbasWeights(mPars Parameters) []float64 {
	dx := 0.1
	var maxbas = make([]float64, int(math.Ceil(mPars.MAXBAS)))
	for i := 0; i < int(math.Ceil(mPars.MAXBAS)); i++ {
		x := floats.Span(make([]float64, int(1.0/dx)+1), float64(i), float64(i+1))
		y := routingMaxbas(x, mPars.MAXBAS)
		x_int := integrate.Trapezoidal(x, y)
		maxbas[i] = x_int
	}

	// Adjust the triangular weights for inaccuracies in integratation so that sum equals 1.0
	maxbas_sum := 0.0
	for i := 0; i < len(maxbas); i++ {
		maxbas_sum += maxbas[i]
	}
	adj_mb_factor := 1.0 / maxbas_sum
	for i := 0; i < len(maxbas); i++ {
		maxbas[i] = maxbas[i] * adj_mb_factor
	}

	return maxbas
}

// Snow routine (also called precipitation routine) for calculating snow- or rainfall
// accumulation, melt and refreezing of snow storage.
func SnowRoutine(mState []ModelState, mPars Parameters, inData []InputData, i int) {
	mState[i].Snow_solid = mState[i-1].Snow_solid

	// Snow cover beginning of day
	if mState[i].Snow_solid > 0 {
		mState[i].Snow_cover = 1
	} else {
		mState[i].Snow_cover = 0
	}

	if inData[i].Temperature <= mPars.TT { // Temperature below threshold
		// If air temp bellow threshold (p_TT) then calculate snowfall and refreezing
		// Snowfall added to snow storage
		mState[i].Snowfall = inData[i].Precipitation * mPars.SFCF
		mState[i].Rainfall = 0
		mState[i].Snow_solid = mState[i].Snow_solid + mState[i].Snowfall
		// Refreezing, moved from liquid to solid snow storage
		var pot_refreeze float64 = mPars.CFMAX * mPars.CFR * (mPars.TT - inData[i].Temperature)
		refreezing := math.Min(pot_refreeze, mState[i-1].Snow_liquid)
		mState[i].Snow_solid = mState[i].Snow_solid + refreezing
		mState[i].Snow_liquid = mState[i-1].Snow_liquid - refreezing // free water content in snowpack
		// No snowmelt or liquid water infiltrating
		mState[i].Snow_melt = 0
		mState[i].Liquid_in = 0
	} else { // Precipitation as rain and snow can melt
		mState[i].Rainfall = inData[i].Precipitation
		mState[i].Snowfall = 0

		snowmelt_potential := math.Max(mPars.CFMAX*(inData[i].Temperature-mPars.TT), 0.0)
		// Snow melt is limited to frozen solid part of the snow pack
		mState[i].Snow_melt = math.Min(snowmelt_potential, mState[i].Snow_solid)
		// Remove snow melt from the solid part of the snow pack
		mState[i].Snow_solid = math.Max(mState[i].Snow_solid-mState[i].Snow_melt, 0.0)

		// Snowpack can retain CWH fraction of meltwater, which can later refreeze
		// Water holding capacity is updated after subtracting melt from solid part of snow pack
		// Max liquid water the snowpack can hold
		pot_liqwater_snow := mState[i].Snow_solid * mPars.CWH
		// Calculate liquid water in the snowpack, snowmelt and rainfall can be held
		// Liquid water in snow pack from previousstep + snowmelt + preciptiation
		mState[i].Snow_liquid = mState[i-1].Snow_liquid + inData[i].Precipitation + mState[i].Snow_melt

		// pot_liqwater_snow is held in remaining snowpack, rest infiltrates
		// Excess meltwater and rainfall goes to infiltration (liquid_in)
		// snow_liquid is not "melted" but will be released here when snowpack can no longer hold it
		mState[i].Liquid_in = math.Max(mState[i].Snow_liquid-pot_liqwater_snow, 0)
		mState[i].Snow_liquid = mState[i].Snow_liquid - mState[i].Liquid_in // Update snowpack liquid water
	}

	// Update total snow storage, combined solid and liquid part
	mState[i].S_snow = mState[i].Snow_solid + mState[i].Snow_liquid

}

func SoilRoutine(mState []ModelState, mPars Parameters, inData []InputData, i int) {
	// Soil routine. Recharge and Evapotranspiration
	// Split input to soil moisture and upper groundwater recharge
	// 1 mm at the time to avoid numerical issues
	soil_s_current := mState[i-1].S_soil
	soil_s_in := 0.0
	recharge_gw_in := 0.0
	recharge_gw_in_total := 0.0

	if mState[i].Liquid_in > 0 {
		liquid_in_last := mState[i].Liquid_in - math.Floor(mState[i].Liquid_in) // last remaining non-whole 1 mm
		liquid_in_int := int(math.Floor(mState[i].Liquid_in))
		for i := 1; i <= liquid_in_int; i++ { // Note i not used
			recharge_gw_in = 1 * math.Pow((soil_s_current/mPars.FC), mPars.BETA) // 1 mm each step
			soil_s_in = 1 - recharge_gw_in                                       // 1 mm each step
			soil_s_current += soil_s_in
			recharge_gw_in_total += recharge_gw_in
			// fmt.Printf("%+v\n", recharge_gw_in)
			// fmt.Printf("%+v\n", soil_s_current)
			// fmt.Printf("%+v\n", math.Pow((soil_s_current/mPars.FC), mPars.BETA))
		}

		recharge_gw_in = liquid_in_last * math.Pow((soil_s_current/mPars.FC), mPars.BETA)
		soil_s_in = liquid_in_last - recharge_gw_in
		soil_s_current += soil_s_in
		recharge_gw_in_total += recharge_gw_in

		mState[i].Recharge_sm = soil_s_current - mState[i-1].S_soil
		mState[i].Recharge_gwuz = recharge_gw_in_total
	} else {
		mState[i].Recharge_gwuz = 0
		mState[i].Recharge_sm = 0
	}

	// ET only if no snow on the ground (as in HBV-light) using mean soils moisture over the recharge day
	sm_aet := (soil_s_current-mState[i-1].S_soil)/2 + mState[i-1].S_soil
	if mState[i].Snow_cover == 1 {
		mState[i].AET = 0
	} else {
		mState[i].AET = inData[i].PotentialET * math.Min(1, (sm_aet*(1/(mPars.LP*mPars.FC))))
	}

	mState[i].S_soil = mState[i-1].S_soil - mState[i].AET + mState[i].Recharge_sm
}

func ResponseRoutine(mState []ModelState, mPars Parameters, i int) {
	// Groundwater recharge and percolation
	mState[i].S_gw_suz = mState[i-1].S_gw_suz + mState[i].Recharge_gwuz
	percolation := math.Min(mPars.PERC, mState[i].S_gw_suz)
	mState[i].S_gw_suz = mState[i].S_gw_suz - percolation
	mState[i].S_gw_slz = mState[i-1].S_gw_slz + percolation

	// Groundwater discharge
	q_lz := mPars.K2 * mState[i].S_gw_slz
	q_uz := mPars.K1 * mState[i].S_gw_suz
	q_uzt := mPars.K0 * math.Max(mState[i].S_gw_suz-mPars.UZL, 0)
	mState[i].Q_gw = q_lz + q_uz + q_uzt

	// Update groundwater storages
	// TODO can they go negative? should not be possible but double think it
	mState[i].S_gw_slz = mState[i].S_gw_slz - q_lz
	mState[i].S_gw_suz = mState[i].S_gw_suz - q_uz - q_uzt
}

// The RoutingRoutine applies maxbas weights to the groundwater response
// to calculate the simulated runoff
func RoutingRoutine(mState []ModelState, mPars Parameters, inData []InputData, i int, maxbas []float64) {

	for j := 0; j < len(maxbas); j++ {
		ij := i + j
		if ij >= len(inData) {
			break
		}
		mState[ij].Q_sim = mState[ij].Q_sim + mState[i].Q_gw*maxbas[j]
	}
}
