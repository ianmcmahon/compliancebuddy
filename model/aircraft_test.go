package model

import (
	"fmt"
	"testing"
)

func TestCreateAircraft(t *testing.T) {
	twinkie := twinkie()

	client := NewClient()

	err := client.UpsertAircraft(twinkie)
	if err != nil {
		t.Error(err)
	}
}

func TestListAircraft(t *testing.T) {
	client := NewClient()

	aircraft, err := client.ListAircraft()
	if err != nil {
		t.Error(err)
	}

	for _, a := range aircraft {
		fmt.Printf("%s (%s %s)\n", a.Registration, a.Model.Make, a.Model.Model)
	}
}

func twinkie() *Aircraft {
	aircraftHobbs := Meter{
		Type:        MASTER_HOBBS,
		Name:        "Airframe Hobbs",
		LastReading: 767.0,
	}

	heaterHobbs := Meter{
		Type:        ACCESSORY_HOBBS,
		Name:        "Heater Hobbs",
		LastReading: 73.2,
	}

	twinkie := &Aircraft{
		Registration: "N7706Y",
		Model: FAAModel{
			Make:  "Piper Aircraft, Inc.",
			Model: "PA-30",
		},
		Serial: "30-794",
		Meters: []Meter{aircraftHobbs, heaterHobbs},
		Time:   TimeInService{aircraftHobbs, 3269.7 - 767.0},
		Engines: []Engine{
			Engine{
				Position: LEFT,
				Model:    FAAModel{"Lycoming Engines", "IO-320-B1A"},
				Serial:   "L-2727-55A",
				Time:     TimeInService{aircraftHobbs, 3269.7 - 435.0},
			},
			Engine{
				Position: RIGHT,
				Model:    FAAModel{"Lycoming Engines", "IO-320-B1A"},
				Serial:   "L-1706-55A",
				Time:     TimeInService{aircraftHobbs, 3269.7 - 1335.0},
			},
		},
		Propellers: []Propeller{
			Propeller{
				Position: LEFT,
				Model:    FAAModel{"Hartzell Propeller", "HC-E2YL-2"},
				Serial:   "BG520",
				Blades:   []PropellerBlade{PropellerBlade{Serial: "C5350"}, PropellerBlade{Serial: "C4256"}},
			},
			Propeller{
				Position: RIGHT,
				Model:    FAAModel{"Hartzell Propeller", "HC-E2YL-2"},
				Serial:   "BG4391",
				Blades:   []PropellerBlade{PropellerBlade{Serial: "B83836"}, PropellerBlade{Serial: "B83801"}},
			},
		},
		Accessories: []Accessory{
			Accessory{
				Position: NOSE,
				Model:    FAAModel{"Janitrol", "B-Series Combustion Heater"},
				Serial:   "Unknown",
				Time:     TimeInService{heaterHobbs, 0},
			},
		},
	}

	return twinkie
}
