package model

const (
	PARKING_TWIN_INTERFACE      = "ngsi-ld-city-offstreetparking"
	PARKING_SPOT_TWIN_INTERFACE = "ngsi-ld-city-offstreetparkingspot"
)

type Facility string
type Layout string
type UsageScenario string
type Security string
type SpecialLocation string

const (
	BikeParking             Facility = "bikeParking"
	CashMachine             Facility = "cashMachine"
	CopyMachineOrService    Facility = "copyMachineOrService"
	Defibrillator           Facility = "defibrillator"
	DumpingStation          Facility = "dumpingStation"
	ElectricChargingStation Facility = "electricChargingStation"
	Elevator                Facility = "elevator"
	FaxMachineOrService     Facility = "faxMachineOrService"
	FireHose                Facility = "fireHose"
	FireExtinguisher        Facility = "fireExtinguisher"
	FireHydrant             Facility = "fireHydrant"
	FirstAidEquipment       Facility = "firstAidEquipment"
	FreshWater              Facility = "freshWater"
	IceFreeScaffold         Facility = "iceFreeScaffold"
	InformationPoint        Facility = "informationPoint"
	InternetWireless        Facility = "internetWireless"
	LuggageLocker           Facility = "luggageLocker"
	PayDesk                 Facility = "payDesk"
	PaymentMachine          Facility = "paymentMachine"
	Playground              Facility = "playground"
	PublicPhone             Facility = "publicPhone"
	RefuseBin               Facility = "refuseBin"
	SafeDeposit             Facility = "safeDeposit"
	Shower                  Facility = "shower"
	Toilet                  Facility = "toilet"
	TollTerminal            Facility = "tollTerminal"
	VendingMachine          Facility = "vendingMachine"
	WasteDisposal           Facility = "wasteDisposal"

	AutomatedParkingGarage Layout = "automatedParkingGarage"
	Carports               Layout = "carports"
	Covered                Layout = "covered"
	Field                  Layout = "field"
	GarageBoxes            Layout = "garageBoxes"
	MultiLevel             Layout = "multiLevel"
	MultiStorey            Layout = "multiStorey"
	Nested                 Layout = "nested"
	OpenSpace              Layout = "openSpace"
	Rooftop                Layout = "rooftop"
	Sheds                  Layout = "sheds"
	SingleLevel            Layout = "singleLevel"
	Surface                Layout = "surface"
	Other                  Layout = "other"

	AutomaticParkingGuidance UsageScenario = "automaticParkingGuidance"
	CarSharing               UsageScenario = "carSharing"
	DropOffWithValet         UsageScenario = "dropOffWithValet"
	DropOffMechanical        UsageScenario = "dropOffMechanical"
	DropOff                  UsageScenario = "dropOff"
	EventParking             UsageScenario = "eventParking"
	KissAndRide              UsageScenario = "kissAndRide"
	LiftShare                UsageScenario = "liftShare"
	LoadingBay               UsageScenario = "loadingBay"
	OvernightParking         UsageScenario = "overnightParking"
	ParkAndCycle             UsageScenario = "parkAndCycle"
	ParkAndRide              UsageScenario = "parkAndRide"
	ParkAndWalk              UsageScenario = "parkAndWalk"
	RestArea                 UsageScenario = "restArea"
	ServiceArea              UsageScenario = "serviceArea"
	StaffGuidesToSpace       UsageScenario = "staffGuidesToSpace"
	TruckParking             UsageScenario = "truckParking"
	VehicleLift              UsageScenario = "vehicleLift"
)

type UpdateVehicleCountCommand struct {
	VehicleEntranceCount int `json:"vehicleEntranceCount,omitempty"`
	VehicleExitCount     int `json:"vehicleExitCount,omitempty"`
}

func NewOffStreetParking() OffStreetParking {
	return OffStreetParking{}
}

type OffStreetParking struct {
	AccessModified      string            `json:"accessModified,omitempty"`
	AggregateRating     float64           `json:"aggregateRating,omitempty"`
	HighestFloor        int               `json:"highestFloor,omitempty"`
	Images              string            `json:"images,omitempty"`
	LowestFloor         int               `json:"lowestFloor,omitempty"`
	OpeningHours        string            `json:"openingHours,omitempty"`
	PriceCurrency       string            `json:"priceCurrency,omitempty"`
	PriceRatePerMinute  float64           `json:"priceRatePerMinute,omitempty"`
	Provider            string            `json:"provider,omitempty"`
	Facilities          []Facility        `json:"facilities,omitempty"`
	Layout              Layout            `json:"layout,omitempty"`
	UsageScenario       UsageScenario     `json:"usageScenario,omitempty"`
	Security            []Security        `json:"security,omitempty"`
	SpecialLocation     []SpecialLocation `json:"specialLocation,omitempty"`
	Category            string            `json:"category,omitempty"`
	ExtCategory         string            `json:"extCategory,omitempty"`
	FirstAvailableFloor int               `json:"firstAvailableFloor,omitempty"`
	MeasuresPeriod      float64           `json:"measuresPeriod,omitempty"`
	MeasuresPeriodUnit  string            `json:"measuresPeriodUnit,omitempty"`
	Occupancy           float64           `json:"occupancy,omitempty"`
	OccupancyModified   float64           `json:"occupancyModified,omitempty"`
	OccupiedSpotNumber  int               `json:"occupiedSpotNumber,omitempty"`
	Status              string            `json:"status,omitempty"`
}

type Category string
type Status string

const (
	OffStreet Category = "offStreet"
	OnStreet  Category = "onStreet"
	Occupied  Status   = "occupied"
	Free      Status   = "free"
	Closed    Status   = "closed"
	Unknown   Status   = "unknown"
)

type ParkingSpot struct {
	DateObserved float64  `json:"dateObserved,omitempty"`
	Width        float64  `json:"width,omitempty"`
	Length       float64  `json:"length,omitempty"`
	TimeInstant  string   `json:"timeInstant,omitempty"`
	Image        string   `json:"image,omitempty"`
	Color        string   `json:"color,omitempty"`
	Category     Category `json:"category,omitempty"`
	Status       Status   `json:"status,omitempty"`
}
