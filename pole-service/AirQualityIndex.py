from enum import Enum

class AQICategory(Enum):
    GOOD="GOOD"
    MODERATE="MODERATE"
    UNHEALTHY_FOR_SENSITIVE_GROUPS="UNHEALTHY_FOR_SENSITIVE_GROUPS"
    UNHEALTHY="UNHEALTHY"
    VERY_UNHEALTHY="VERY_UNHEALTHY"
    HAZARDOUS="HAZARDOUS"

class AQIBreakpoints:
    index_breakpoint_low: float
    index_breakpoint_high: float
    concentration_breakpoint_low: float
    concentration_breakpoint_high: float

class PollutantAQIIndexes:
    range: float # the difference range between levels (1 or 0.1)
    good_low: float
    good_high: float
    moderate_low: float
    moderate_high: float
    unhealthy_for_sensitive_groups_low: float
    unhealthy_for_sensitive_groups_high: float
    unhealthy_low: float
    unhealthy_high: float
    very_unhealthy_low: float
    very_unhealthy_high: float
    hazardous_high: float
    hazardous_low: float
    very_hazardous_low: float
    very_hazardous_high: float



# Default AQI Indexes
AQI_GOOD_LOW = 0
AQI_GOOD_HIGH = 50
AQI_MODERATE_LOW = 51
AQI_MODERATE_HIGH = 100
AQI_UNHEALTHY_FOR_SENSITIVE_GROUPS_LOW = 101
AQI_UNHEALTHY_FOR_SENSITIVE_GROUPS_HIGH = 150
AQI_UNHEALTHY_LOW = 151
AQI_UNHEALTHY_HIGH = 200
AQI_VERY_UNHEALTHY_LOW = 201
AQI_VERY_UNHEALTHY_HIGH = 300
AQI_HAZARDOUS_LOW = 301
AQI_HAZARDOUS_HIGH = 400
AQI_VERY_HAZARDOUS_LOW = 401
AQI_VERY_HAZARDOUS_HIGH = 500

class AirQualityIndex:
    def __init__(self, pollutant_aqi_index: PollutantAQIIndexes, concentration: float) -> None:
        self.pollutant_aqi_index = pollutant_aqi_index
        self.concentration = concentration

    def get_aqi_break_points(self) -> AQIBreakpoints:
        aqi_breakpoints = AQIBreakpoints()

        # Invalid (Negative concentration)
        if self.concentration < self.pollutant_aqi_index.good_low:
            return aqi_breakpoints

        # Good Level
        if self.concentration >= self.pollutant_aqi_index.good_low and self.concentration < self.pollutant_aqi_index.good_high + self.pollutant_aqi_index.range:
            aqi_breakpoints.index_breakpoint_low = AQI_GOOD_LOW
            aqi_breakpoints.index_breakpoint_high = AQI_GOOD_HIGH
            aqi_breakpoints.concentration_breakpoint_low = self.pollutant_aqi_index.good_low
            aqi_breakpoints.concentration_breakpoint_high = self.pollutant_aqi_index.good_high
            return aqi_breakpoints

        # Moderate Level
        if self.concentration >= self.pollutant_aqi_index.moderate_low and self.concentration < self.pollutant_aqi_index.moderate_high + self.pollutant_aqi_index.range:
            aqi_breakpoints.index_breakpoint_low = AQI_MODERATE_LOW
            aqi_breakpoints.index_breakpoint_high = AQI_MODERATE_HIGH
            aqi_breakpoints.concentration_breakpoint_low = self.pollutant_aqi_index.moderate_low
            aqi_breakpoints.concentration_breakpoint_high = self.pollutant_aqi_index.moderate_high
            return aqi_breakpoints

        # Unhealthy for sensitive groups Level
        if self.concentration >= self.pollutant_aqi_index.unhealthy_for_sensitive_groups_low and self.concentration < self.pollutant_aqi_index.unhealthy_for_sensitive_groups_high + self.pollutant_aqi_index.range:
            aqi_breakpoints.index_breakpoint_low = AQI_UNHEALTHY_FOR_SENSITIVE_GROUPS_LOW
            aqi_breakpoints.index_breakpoint_high = AQI_UNHEALTHY_FOR_SENSITIVE_GROUPS_HIGH
            aqi_breakpoints.concentration_breakpoint_low = self.pollutant_aqi_index.unhealthy_for_sensitive_groups_low
            aqi_breakpoints.concentration_breakpoint_high = self.pollutant_aqi_index.unhealthy_for_sensitive_groups_high
            return aqi_breakpoints

        # Unhealthy Level
        if self.concentration >= self.pollutant_aqi_index.unhealthy_low and self.concentration < self.pollutant_aqi_index.unhealthy_high + self.pollutant_aqi_index.range:
            aqi_breakpoints.index_breakpoint_low = AQI_UNHEALTHY_LOW
            aqi_breakpoints.index_breakpoint_high = AQI_UNHEALTHY_HIGH
            aqi_breakpoints.concentration_breakpoint_low = self.pollutant_aqi_index.unhealthy_low
            aqi_breakpoints.concentration_breakpoint_high = self.pollutant_aqi_index.unhealthy_high
            return aqi_breakpoints
        
        # Very Unhealthy Level
        if self.concentration >= self.pollutant_aqi_index.very_unhealthy_low and self.concentration < self.pollutant_aqi_index.very_unhealthy_high + self.pollutant_aqi_index.range:
            aqi_breakpoints.index_breakpoint_low = AQI_VERY_UNHEALTHY_LOW
            aqi_breakpoints.index_breakpoint_high = AQI_VERY_UNHEALTHY_HIGH
            aqi_breakpoints.concentration_breakpoint_low = self.pollutant_aqi_index.very_unhealthy_low
            aqi_breakpoints.concentration_breakpoint_high = self.pollutant_aqi_index.very_unhealthy_high
            return aqi_breakpoints
        
        # Hazardous Level
        if self.concentration >= self.pollutant_aqi_index.hazardous_low and self.concentration < self.pollutant_aqi_index.hazardous_high + self.pollutant_aqi_index.range:
            aqi_breakpoints.index_breakpoint_low = AQI_HAZARDOUS_LOW
            aqi_breakpoints.index_breakpoint_high = AQI_HAZARDOUS_HIGH
            aqi_breakpoints.concentration_breakpoint_low = self.pollutant_aqi_index.hazardous_low
            aqi_breakpoints.concentration_breakpoint_high = self.pollutant_aqi_index.hazardous_high
            return aqi_breakpoints

        # Very Hazardous Level
        aqi_breakpoints.index_breakpoint_low = AQI_VERY_HAZARDOUS_LOW
        aqi_breakpoints.index_breakpoint_high = AQI_VERY_HAZARDOUS_HIGH
        aqi_breakpoints.concentration_breakpoint_low = self.pollutant_aqi_index.very_hazardous_low
        aqi_breakpoints.concentration_breakpoint_high = self.pollutant_aqi_index.very_hazardous_high
        return aqi_breakpoints

    def calculate_air_quality_index(self) -> float:
        aqi_breakpoints = self.get_aqi_break_points()

        index_breakpoint_low = aqi_breakpoints.index_breakpoint_low
        index_breakpoint_high = aqi_breakpoints.index_breakpoint_high
        concentration_breakpoint_low = aqi_breakpoints.concentration_breakpoint_low
        concentration_breakpoint_high = aqi_breakpoints.index_breakpoint_high

        air_quality_index = ((index_breakpoint_high-index_breakpoint_low)/(concentration_breakpoint_high-concentration_breakpoint_low))*(self.concentration-concentration_breakpoint_low) + index_breakpoint_low
        return air_quality_index
    
    def get_air_quality_category(self) -> AQICategory:
        air_index = self.calculate_air_quality_index()

        # Good Level
        if air_index >= AQI_GOOD_LOW and air_index < AQI_GOOD_HIGH + 1:
            return AQICategory.GOOD

        # Moderate Level
        if air_index >= AQI_MODERATE_LOW and air_index < AQI_MODERATE_HIGH + 1:
            return AQICategory.MODERATE

        # Unhealthy for sensitive groups Level
        if air_index >= AQI_UNHEALTHY_FOR_SENSITIVE_GROUPS_LOW and air_index < AQI_UNHEALTHY_FOR_SENSITIVE_GROUPS_HIGH + 1:
            return AQICategory.UNHEALTHY_FOR_SENSITIVE_GROUPS

        # Unhealthy Level
        if air_index >= AQI_UNHEALTHY_LOW and air_index < AQI_UNHEALTHY_HIGH + 1:
            return AQICategory.UNHEALTHY
        
        # Very Unhealthy Level
        if air_index >= AQI_VERY_UNHEALTHY_LOW and air_index < AQI_VERY_UNHEALTHY_HIGH + 1:
            return AQICategory.VERY_UNHEALTHY
        
        # Very Unhealthy Level
        if air_index >= AQI_HAZARDOUS_LOW:
            return AQICategory.HAZARDOUS

# Level definitions: https://en.wikipedia.org/wiki/Air_quality_index

# 8 hour
class O3AirQualityIndex(AirQualityIndex):
    def __init__(self, concentration: float) -> None:
        pollutant_aqi_index = PollutantAQIIndexes()
        pollutant_aqi_index.range = 1
        pollutant_aqi_index.good_low = 0
        pollutant_aqi_index.good_high = 54
        pollutant_aqi_index.moderate_low = 55
        pollutant_aqi_index.moderate_high = 70
        pollutant_aqi_index.unhealthy_for_sensitive_groups_low = 71
        pollutant_aqi_index.unhealthy_for_sensitive_groups_high = 85
        pollutant_aqi_index.unhealthy_low = 86
        pollutant_aqi_index.unhealthy_high = 105
        pollutant_aqi_index.very_unhealthy_low = 106
        pollutant_aqi_index.very_unhealthy_high = 200
        pollutant_aqi_index.hazardous_high = 200 # Not defined because everything above 200 is hazardous
        pollutant_aqi_index.hazardous_low = 300
        pollutant_aqi_index.very_hazardous_low = 300
        pollutant_aqi_index.very_hazardous_high = 400
        super().__init__(pollutant_aqi_index=pollutant_aqi_index, concentration=concentration)

# 24 hour
class PM25AirQualityIndex(AirQualityIndex):
    def __init__(self, concentration: float) -> None:
        pollutant_aqi_index = PollutantAQIIndexes()
        pollutant_aqi_index.range = 0.1
        pollutant_aqi_index.good_low = 0
        pollutant_aqi_index.good_high = 12
        pollutant_aqi_index.moderate_low = 12.1
        pollutant_aqi_index.moderate_high = 35.4
        pollutant_aqi_index.unhealthy_for_sensitive_groups_low = 35.5
        pollutant_aqi_index.unhealthy_for_sensitive_groups_high = 55.4
        pollutant_aqi_index.unhealthy_low = 55.5
        pollutant_aqi_index.unhealthy_high = 150.4
        pollutant_aqi_index.very_unhealthy_low = 150.4
        pollutant_aqi_index.very_unhealthy_high = 250.4
        pollutant_aqi_index.hazardous_high = 250.5
        pollutant_aqi_index.hazardous_low = 350.4
        pollutant_aqi_index.very_hazardous_low = 350.5
        pollutant_aqi_index.very_hazardous_high = 500
        super().__init__(pollutant_aqi_index=pollutant_aqi_index, concentration=concentration)

# 24 hour
class PM10AirQualityIndex(AirQualityIndex):
    def __init__(self, concentration: float) -> None:
        pollutant_aqi_index = PollutantAQIIndexes()
        pollutant_aqi_index.range = 1
        pollutant_aqi_index.good_low = 0
        pollutant_aqi_index.good_high = 54
        pollutant_aqi_index.moderate_low = 55
        pollutant_aqi_index.moderate_high = 154
        pollutant_aqi_index.unhealthy_for_sensitive_groups_low = 155
        pollutant_aqi_index.unhealthy_for_sensitive_groups_high = 254
        pollutant_aqi_index.unhealthy_low = 255
        pollutant_aqi_index.unhealthy_high = 354
        pollutant_aqi_index.very_unhealthy_low = 355
        pollutant_aqi_index.very_unhealthy_high = 424
        pollutant_aqi_index.hazardous_high = 425
        pollutant_aqi_index.hazardous_low = 504
        pollutant_aqi_index.very_hazardous_low = 505
        pollutant_aqi_index.very_hazardous_high = 604
        super().__init__(pollutant_aqi_index=pollutant_aqi_index, concentration=concentration)

# 8h hour
class COAirQualityIndex(AirQualityIndex):
    def __init__(self, concentration: float) -> None:
        pollutant_aqi_index = PollutantAQIIndexes()
        pollutant_aqi_index.range = 0.1
        pollutant_aqi_index.good_low = 0
        pollutant_aqi_index.good_high = 4.4
        pollutant_aqi_index.moderate_low = 4.5
        pollutant_aqi_index.moderate_high = 9.4
        pollutant_aqi_index.unhealthy_for_sensitive_groups_low = 9.5
        pollutant_aqi_index.unhealthy_for_sensitive_groups_high = 12.4
        pollutant_aqi_index.unhealthy_low = 12.5
        pollutant_aqi_index.unhealthy_high = 15.4
        pollutant_aqi_index.very_unhealthy_low = 15.5
        pollutant_aqi_index.very_unhealthy_high = 30.4
        pollutant_aqi_index.hazardous_high = 30.5
        pollutant_aqi_index.hazardous_low = 40.4
        pollutant_aqi_index.very_hazardous_low = 40.5
        pollutant_aqi_index.very_hazardous_high = 50.4
        super().__init__(pollutant_aqi_index=pollutant_aqi_index, concentration=concentration)

class SO2AirQualityIndex(AirQualityIndex):
    def __init__(self, concentration: float) -> None:
        pollutant_aqi_index = PollutantAQIIndexes()
        pollutant_aqi_index.range = 1
        pollutant_aqi_index.good_low = 0
        pollutant_aqi_index.good_high = 35
        pollutant_aqi_index.moderate_low = 36
        pollutant_aqi_index.moderate_high = 75
        pollutant_aqi_index.unhealthy_for_sensitive_groups_low = 76
        pollutant_aqi_index.unhealthy_for_sensitive_groups_high = 185
        pollutant_aqi_index.unhealthy_low = 186
        pollutant_aqi_index.unhealthy_high = 304
        pollutant_aqi_index.very_unhealthy_low = 305
        pollutant_aqi_index.very_unhealthy_high = 604
        pollutant_aqi_index.hazardous_high = 605
        pollutant_aqi_index.hazardous_low = 804
        pollutant_aqi_index.very_hazardous_low = 805
        pollutant_aqi_index.very_hazardous_high = 1004
        super().__init__(pollutant_aqi_index=pollutant_aqi_index, concentration=concentration)
