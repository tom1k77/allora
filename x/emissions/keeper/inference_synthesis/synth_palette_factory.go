package inference_synthesis

import (
	alloraMath "github.com/allora-network/allora-chain/math"
)

// Could use Builder pattern in the future to make this cleaner
func (f *SynthPaletteFactory) BuildPaletteFromRequest(req SynthRequest) (SynthPalette, error) {
	inferenceByWorker := MakeMapFromInfererToTheirInference(req.Inferences.Inferences)
	forecastByWorker := MakeMapFromForecasterToTheirForecast(req.Forecasts.Forecasts)
	sortedInferers := alloraMath.GetSortedKeys(inferenceByWorker)
	sortedForecasters := alloraMath.GetSortedKeys(forecastByWorker)

	// Those values not from req are to be considered defaults
	palette := SynthPalette{
		Ctx:                              req.Ctx,
		K:                                req.K,
		TopicId:                          req.TopicId,
		Inferers:                         sortedInferers,
		InferenceByWorker:                inferenceByWorker,
		InfererRegrets:                   make(map[string]*StatefulRegret), // Populated below
		Forecasters:                      sortedForecasters,
		ForecastByWorker:                 forecastByWorker,
		ForecastImpliedInferenceByWorker: nil,                              // Populated below
		ForecasterRegrets:                make(map[string]*StatefulRegret), // Populated below
		AllInferersAreNew:                true,                             // Populated below
		SingleInfererNotNew:              "",                               // Populated below
		NetworkCombinedLoss:              req.NetworkCombinedLoss,
		Epsilon:                          req.Epsilon,
		FTolerance:                       req.FTolerance,
		PNorm:                            req.PNorm,
		CNorm:                            req.CNorm,
	}

	// Populates: infererRegrets, forecasterRegrets, allInferersAreNew
	palette.BootstrapRegretData()

	paletteCopy := palette.Clone()
	// Populates: forecastImpliedInferenceByWorker,
	err := paletteCopy.UpdateForecastImpliedInferences()
	if err != nil {
		return SynthPalette{}, err
	}
	palette.ForecastImpliedInferenceByWorker = paletteCopy.ForecastImpliedInferenceByWorker

	return palette, nil
}