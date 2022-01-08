/*
Function to run the gohbv model.

To be implemented:
cli support
*/

package gohbv

// Model state and fluxes - available across module
type ModelState struct {
	AET           float64 // Actual evapotranspiration [mm]
	Q_gw          float64 // Groundwater discharge [mm]
	Q_sim         float64 // Simulated runoff [mm]
	Snowfall      float64 // Snowfall [mm]
	Rainfall      float64 // Rainfall [mm]
	S_snow        float64 // Storage snow, Snow_solid + Snow_liquid [mm]
	Snow_solid    float64 // Solid water content in snowpack [mm]
	Snow_liquid   float64 // Liquid water content in snowpack [mm]
	Snow_cover    int     // Snow cover [0/1]
	Snow_melt     float64 // Snow melt [mm]
	Liquid_in     float64 // Snow melt + liquid precipitation [mm]
	S_soil        float64 // Soil water storage [mm]
	S_gw_suz      float64 // Groundwater storage upper zone [mm]
	S_gw_slz      float64 // Groundwater storage lower zone [mm]
	Recharge_sm   float64 // Recharge/infiltration to soil moisture storage [mm]
	Recharge_gwuz float64 // Recharge to upper groundwater storage [mm]
}

// GoHBV function is used to run the HBV model based on an InputData and Parameters structs.
// The function call returns an array of ModelState structs
func RunModel(mPars Parameters, inData []InputData) ([]ModelState, error) {
	// Set model states
	// All the uninitialized fields are set to zero value (0)
	var mState = make([]ModelState, len(inData))

	// Set initial soil moisture to FP * LP
	mState[0].S_soil = mPars.FC * mPars.LP
	// Initial Lower Zone groundwater storage
	mState[0].S_gw_slz = mPars.PERC / mPars.K2
	// Routing with MAXBAS (only need to get the maxbas array once)
	maxbas := RoutingMaxbasWeights(mPars)

	// Forward Euler loop
	for i := 1; i < len(inData); i++ { // Loop starts on second index (1), first is initial state (zeros)
		SnowRoutine(mState, mPars, inData, i)
		SoilRoutine(mState, mPars, inData, i)
		ResponseRoutine(mState, mPars, i)
		RoutingRoutine(mState, mPars, inData, i, maxbas)
	}

	return mState, nil
}
