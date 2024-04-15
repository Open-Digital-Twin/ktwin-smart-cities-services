package model

const (
	AIR_QUALITY_OBSERVED_TWIN_INTERFACE = "ngsi-ld-city-airqualityobserved"
)

type AQICategory string

const (
	Good                        AQICategory = "GOOD"
	Moderate                    AQICategory = "MODERATE"
	UnhealthyForSensitiveGroups AQICategory = "UNHEALTHY_FOR_SENSITIVE_GROUPS"
	Unhealthy                   AQICategory = "UNHEALTHY"
	VeryUnhealthy               AQICategory = "VERY_UNHEALTHY"
	Hazardous                   AQICategory = "HAZARDOUS"
)

// AirQualityEvent represents the structure for an air quality event
type AirQualityEvent struct {
	AirQualityIndex          float64     `json:"airQualityIndex,omitempty"`
	Reliability              float64     `json:"reliability,omitempty"`
	VolatileOrganicCompounds int         `json:"volatileOrganicCompoundsTotal,omitempty"`
	TypeOfLocation           string      `json:"typeOfLocation,omitempty"`
	CO2Density               float64     `json:"CO2Density,omitempty"`
	CODensity                float64     `json:"CODensity,omitempty"`
	PM1Density               float64     `json:"PM1Density,omitempty"`
	PM10Density              float64     `json:"PM10Density,omitempty"`
	PM25Density              float64     `json:"PM25Density,omitempty"`
	NODensity                float64     `json:"NODensity,omitempty"`
	SO2Density               float64     `json:"SO2Density,omitempty"`
	C6H6Density              float64     `json:"C6H6Density,omitempty"`
	NIDensity                float64     `json:"NIDensity,omitempty"`
	ASDensity                float64     `json:"ASDensity,omitempty"`
	CDDensity                float64     `json:"CDDensity,omitempty"`
	NO2Density               float64     `json:"NO2Density,omitempty"`
	O3Density                float64     `json:"O3Density,omitempty"`
	PBDensity                float64     `json:"PBDensity,omitempty"`
	SH2Density               float64     `json:"SH2Density,omitempty"`
	Precipitation            float64     `json:"precipitation,omitempty"`
	RelativeHumidity         float64     `json:"relativeHumidity,omitempty"`
	Temperature              float64     `json:"temperature,omitempty"`
	WindDirection            float64     `json:"WindDirection,omitempty"`
	WindSpeed                float64     `json:"WindSpeed,omitempty"`
	COAqiLevel               AQICategory `json:"COAqiLevel,omitempty"`
	PM10AqiLevel             AQICategory `json:"PM10AqiLevel,omitempty"`
	PM25AqiLevel             AQICategory `json:"PM25AqiLevel,omitempty"`
	SO2AqiLevel              AQICategory `json:"SO2AqiLevel,omitempty"`
	O3AqiLevel               AQICategory `json:"O3AqiLevel,omitempty"`
}

type UpdateAirQualityIndexCommand struct {
	AqiLevel AQICategory `json:"aqiLevel,omitempty"`
}

type AQIBreakpoints struct {
	IndexBreakpointLow          float64
	IndexBreakpointHigh         float64
	ConcentrationBreakpointLow  float64
	ConcentrationBreakpointHigh float64
}

type PollutantAQIIndexes struct {
	Range                           float64
	GoodLow                         float64
	GoodHigh                        float64
	ModerateLow                     float64
	ModerateHigh                    float64
	UnhealthyForSensitiveGroupsLow  float64
	UnhealthyForSensitiveGroupsHigh float64
	UnhealthyLow                    float64
	UnhealthyHigh                   float64
	VeryUnhealthyLow                float64
	VeryUnhealthyHigh               float64
	HazardousHigh                   float64
	HazardousLow                    float64
	VeryHazardousLow                float64
	VeryHazardousHigh               float64
}

var (
	AQIGoodLow                         = 0.0
	AQIGoodHigh                        = 50.0
	AQIModerateLow                     = 51.0
	AQIModerateHigh                    = 100.0
	AQIUnhealthyForSensitiveGroupsLow  = 101.0
	AQIUnhealthyForSensitiveGroupsHigh = 150.0
	AQIUnhealthyLow                    = 151.0
	AQIUnhealthyHigh                   = 200.0
	AQIVeryUnhealthyLow                = 201.0
	AQIVeryUnhealthyHigh               = 300.0
	AQIHazardousLow                    = 301.0
	AQIHazardousHigh                   = 400.0
	AQIVeryHazardousLow                = 401.0
	AQIVeryHazardousHigh               = 500.0
)

type AirQualityIndex struct {
	PollutantAQIIndex PollutantAQIIndexes
	Concentration     float64
}

func (aqi *AirQualityIndex) GetAQIBreakPoints() AQIBreakpoints {
	aqiBreakpoints := AQIBreakpoints{}

	// Invalid (Negative concentration)
	if aqi.Concentration < aqi.PollutantAQIIndex.GoodLow {
		return aqiBreakpoints
	}

	// Good Level
	if aqi.Concentration >= aqi.PollutantAQIIndex.GoodLow && aqi.Concentration < aqi.PollutantAQIIndex.GoodHigh+aqi.PollutantAQIIndex.Range {
		aqiBreakpoints.IndexBreakpointLow = AQIGoodLow
		aqiBreakpoints.IndexBreakpointHigh = AQIGoodHigh
		aqiBreakpoints.ConcentrationBreakpointLow = aqi.PollutantAQIIndex.GoodLow
		aqiBreakpoints.ConcentrationBreakpointHigh = aqi.PollutantAQIIndex.GoodHigh
		return aqiBreakpoints
	}

	// Moderate Level
	if aqi.Concentration >= aqi.PollutantAQIIndex.ModerateLow && aqi.Concentration < aqi.PollutantAQIIndex.ModerateHigh+aqi.PollutantAQIIndex.Range {
		aqiBreakpoints.IndexBreakpointLow = AQIModerateLow
		aqiBreakpoints.IndexBreakpointHigh = AQIModerateHigh
		aqiBreakpoints.ConcentrationBreakpointLow = aqi.PollutantAQIIndex.ModerateLow
		aqiBreakpoints.ConcentrationBreakpointHigh = aqi.PollutantAQIIndex.ModerateHigh
		return aqiBreakpoints
	}

	// Unhealthy for sensitive groups Level
	if aqi.Concentration >= aqi.PollutantAQIIndex.UnhealthyForSensitiveGroupsLow && aqi.Concentration < aqi.PollutantAQIIndex.UnhealthyForSensitiveGroupsHigh+aqi.PollutantAQIIndex.Range {
		aqiBreakpoints.IndexBreakpointLow = AQIUnhealthyForSensitiveGroupsLow
		aqiBreakpoints.IndexBreakpointHigh = AQIUnhealthyForSensitiveGroupsHigh
		aqiBreakpoints.ConcentrationBreakpointLow = aqi.PollutantAQIIndex.UnhealthyForSensitiveGroupsLow
		aqiBreakpoints.ConcentrationBreakpointHigh = aqi.PollutantAQIIndex.UnhealthyForSensitiveGroupsHigh
		return aqiBreakpoints
	}

	// Unhealthy Level
	if aqi.Concentration >= aqi.PollutantAQIIndex.UnhealthyLow && aqi.Concentration < aqi.PollutantAQIIndex.UnhealthyHigh+aqi.PollutantAQIIndex.Range {
		aqiBreakpoints.IndexBreakpointLow = AQIUnhealthyLow
		aqiBreakpoints.IndexBreakpointHigh = AQIUnhealthyHigh
		aqiBreakpoints.ConcentrationBreakpointLow = aqi.PollutantAQIIndex.UnhealthyLow
		aqiBreakpoints.ConcentrationBreakpointHigh = aqi.PollutantAQIIndex.UnhealthyHigh
		return aqiBreakpoints
	}

	// Very Unhealthy Level
	if aqi.Concentration >= aqi.PollutantAQIIndex.VeryUnhealthyLow && aqi.Concentration < aqi.PollutantAQIIndex.VeryUnhealthyHigh+aqi.PollutantAQIIndex.Range {
		aqiBreakpoints.IndexBreakpointLow = AQIVeryUnhealthyLow
		aqiBreakpoints.IndexBreakpointHigh = AQIVeryUnhealthyHigh
		aqiBreakpoints.ConcentrationBreakpointLow = aqi.PollutantAQIIndex.VeryUnhealthyLow
		aqiBreakpoints.ConcentrationBreakpointHigh = aqi.PollutantAQIIndex.VeryUnhealthyHigh
		return aqiBreakpoints
	}

	// Hazardous Level
	if aqi.Concentration >= aqi.PollutantAQIIndex.HazardousLow && aqi.Concentration < aqi.PollutantAQIIndex.HazardousHigh+aqi.PollutantAQIIndex.Range {
		aqiBreakpoints.IndexBreakpointLow = AQIHazardousLow
		aqiBreakpoints.IndexBreakpointHigh = AQIHazardousHigh
		aqiBreakpoints.ConcentrationBreakpointLow = aqi.PollutantAQIIndex.HazardousLow
		aqiBreakpoints.ConcentrationBreakpointHigh = aqi.PollutantAQIIndex.HazardousHigh
		return aqiBreakpoints
	}

	// Very Hazardous Level
	aqiBreakpoints.IndexBreakpointLow = AQIVeryHazardousLow
	aqiBreakpoints.IndexBreakpointHigh = AQIVeryHazardousHigh
	aqiBreakpoints.ConcentrationBreakpointLow = aqi.PollutantAQIIndex.VeryHazardousLow
	aqiBreakpoints.ConcentrationBreakpointHigh = aqi.PollutantAQIIndex.VeryHazardousHigh
	return aqiBreakpoints
}

func (aqi *AirQualityIndex) CalculateAirQualityIndex() float64 {
	aqiBreakpoints := aqi.GetAQIBreakPoints()

	indexBreakpointLow := aqiBreakpoints.IndexBreakpointLow
	indexBreakpointHigh := aqiBreakpoints.IndexBreakpointHigh
	concentrationBreakpointLow := aqiBreakpoints.ConcentrationBreakpointLow
	concentrationBreakpointHigh := aqiBreakpoints.IndexBreakpointHigh

	airQualityIndex := ((indexBreakpointHigh-indexBreakpointLow)/(concentrationBreakpointHigh-concentrationBreakpointLow))*(aqi.Concentration-concentrationBreakpointLow) + indexBreakpointLow
	return airQualityIndex
}

func (aqi *AirQualityIndex) GetAirQualityCategory() AQICategory {
	airIndex := aqi.CalculateAirQualityIndex()

	// Good Level
	if airIndex >= AQIGoodLow && airIndex < AQIGoodHigh+1 {
		return Good
	}

	// Moderate Level
	if airIndex >= AQIModerateLow && airIndex < AQIModerateHigh+1 {
		return Moderate
	}

	// Unhealthy for sensitive groups Level
	if airIndex >= AQIUnhealthyForSensitiveGroupsLow && airIndex < AQIUnhealthyForSensitiveGroupsHigh+1 {
		return UnhealthyForSensitiveGroups
	}

	// Unhealthy Level
	if airIndex >= AQIUnhealthyLow && airIndex < AQIUnhealthyHigh+1 {
		return Unhealthy
	}

	// Very Unhealthy Level
	if airIndex >= AQIVeryUnhealthyLow && airIndex < AQIVeryUnhealthyHigh+1 {
		return VeryUnhealthy
	}

	// Hazardous Level
	if airIndex >= AQIHazardousLow {
		return Hazardous
	}

	return ""
}

type O3AirQualityIndex struct {
	AirQualityIndex
}

func NewO3AirQualityIndex(concentration float64) *O3AirQualityIndex {
	pollutantAQIIndex := PollutantAQIIndexes{
		Range:                           1,
		GoodLow:                         0,
		GoodHigh:                        54,
		ModerateLow:                     55,
		ModerateHigh:                    70,
		UnhealthyForSensitiveGroupsLow:  71,
		UnhealthyForSensitiveGroupsHigh: 85,
		UnhealthyLow:                    86,
		UnhealthyHigh:                   105,
		VeryUnhealthyLow:                106,
		VeryUnhealthyHigh:               200,
		HazardousHigh:                   200,
		HazardousLow:                    300,
		VeryHazardousLow:                300,
		VeryHazardousHigh:               400,
	}

	return &O3AirQualityIndex{AirQualityIndex{pollutantAQIIndex, concentration}}
}

type PM25AirQualityIndex struct {
	AirQualityIndex
}

func NewPM25AirQualityIndex(concentration float64) *PM25AirQualityIndex {
	pollutantAQIIndex := PollutantAQIIndexes{
		Range:                           0.1,
		GoodLow:                         0,
		GoodHigh:                        12,
		ModerateLow:                     12.1,
		ModerateHigh:                    35.4,
		UnhealthyForSensitiveGroupsLow:  35.5,
		UnhealthyForSensitiveGroupsHigh: 55.4,
		UnhealthyLow:                    55.5,
		UnhealthyHigh:                   150.4,
		VeryUnhealthyLow:                150.4,
		VeryUnhealthyHigh:               250.4,
		HazardousHigh:                   250.5,
		HazardousLow:                    350.4,
		VeryHazardousLow:                350.5,
		VeryHazardousHigh:               500,
	}

	return &PM25AirQualityIndex{AirQualityIndex{pollutantAQIIndex, concentration}}
}

type PM10AirQualityIndex struct {
	AirQualityIndex
}

func NewPM10AirQualityIndex(concentration float64) *PM10AirQualityIndex {
	pollutantAQIIndex := PollutantAQIIndexes{
		Range:                           1,
		GoodLow:                         0,
		GoodHigh:                        54,
		ModerateLow:                     55,
		ModerateHigh:                    154,
		UnhealthyForSensitiveGroupsLow:  155,
		UnhealthyForSensitiveGroupsHigh: 254,
		UnhealthyLow:                    255,
		UnhealthyHigh:                   354,
		VeryUnhealthyLow:                355,
		VeryUnhealthyHigh:               424,
		HazardousHigh:                   425,
		HazardousLow:                    504,
		VeryHazardousLow:                505,
		VeryHazardousHigh:               604,
	}

	return &PM10AirQualityIndex{AirQualityIndex{pollutantAQIIndex, concentration}}
}

type COAirQualityIndex struct {
	AirQualityIndex
}

func NewCOAirQualityIndex(concentration float64) *COAirQualityIndex {
	pollutantAQIIndex := PollutantAQIIndexes{
		Range:                           0.1,
		GoodLow:                         0,
		GoodHigh:                        4.4,
		ModerateLow:                     4.5,
		ModerateHigh:                    9.4,
		UnhealthyForSensitiveGroupsLow:  9.5,
		UnhealthyForSensitiveGroupsHigh: 12.4,
		UnhealthyLow:                    12.5,
		UnhealthyHigh:                   15.4,
		VeryUnhealthyLow:                15.5,
		VeryUnhealthyHigh:               30.4,
		HazardousHigh:                   30.5,
		HazardousLow:                    40.4,
		VeryHazardousLow:                40.5,
		VeryHazardousHigh:               50.4,
	}

	return &COAirQualityIndex{AirQualityIndex{pollutantAQIIndex, concentration}}
}

type SO2AirQualityIndex struct {
	AirQualityIndex
}

func NewSO2AirQualityIndex(concentration float64) *SO2AirQualityIndex {
	pollutantAQIIndex := PollutantAQIIndexes{
		Range:                           1,
		GoodLow:                         0,
		GoodHigh:                        35,
		ModerateLow:                     36,
		ModerateHigh:                    75,
		UnhealthyForSensitiveGroupsLow:  76,
		UnhealthyForSensitiveGroupsHigh: 185,
		UnhealthyLow:                    186,
		UnhealthyHigh:                   304,
		VeryUnhealthyLow:                305,
		VeryUnhealthyHigh:               604,
		HazardousHigh:                   605,
		HazardousLow:                    804,
		VeryHazardousLow:                805,
		VeryHazardousHigh:               1004,
	}

	return &SO2AirQualityIndex{AirQualityIndex{pollutantAQIIndex, concentration}}
}
