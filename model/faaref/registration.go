package faaref

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/ianmcmahon/compliancebuddy/faaservices"
)

func CategoryMap() map[string]string {
	return map[string]string{
		"airplane":   "Airplane",
		"rotorcraft": "Rotorcraft",
		"glider":     "Glider",
		"lta":        "Lighter Than Air",
		"pp":         "Powered Parachute",
		"ws":         "Weight Shift",
	}
}

func ClassMap() map[string]map[string]string {
	return map[string]map[string]string{
		"airplane": map[string]string{
			"SEL": "Single-Engine Land",
			"MEL": "Multi-Engine Land",
			"SES": "Single-Engine Sea",
			"MES": "Multi-Engine Sea",
			"SEA": "Single-Engine Amphib",
			"MEA": "Multi-Engine Amphib",
		},
		"rotorcraft": map[string]string{
			"helicopter": "Helicopter",
			"gyroplane":  "Gyroplane",
		},
		"glider": map[string]string{
			"glider": "Glider",
		},
		"lta": map[string]string{
			"airship": "Airship",
			"balloon": "Balloon",
		},
		"pp": map[string]string{
			"pp": "Powered Parachute",
		},
		"ws": map[string]string{
			"ws": "Weight Shift",
		},
	}
}

func mustParseInt(s string) int32 {
	s = strings.TrimSpace(s)
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0
	}
	return int32(i)
}

func mustParseDate(s string) time.Time {
	s = strings.TrimSpace(s)
	t, err := time.Parse("20060102", s)
	if err != nil {
		return time.Time{}
	}
	return t
}

func mustParseBool(s string) bool {
	s = strings.TrimSpace(s)
	return s == "Y"
}

type EngineType int32

const (
	EngineType_None       EngineType = iota // 0
	EngineType_Recip                        // 1
	EngineType_Turboprop                    // 2
	EngineType_Turboshaft                   // 3
	EngineType_Turbojet                     // 4
	EngineType_Turbofan                     // 5
	EngineType_Ramjet                       // 6
	EngineType_2Cycle                       // 7
	EngineType_4Cycle                       // 8
	EngineType_Unknown                      // 9
	EngineType_Electric                     // 10
	EngineType_Rotary                       // 11
)

func (t EngineType) Power() string {
	if t == EngineType_Turbojet || t == EngineType_Turbofan || t == EngineType_Ramjet {
		return "lb"
	}
	return "hp"
}

type EnginePower struct {
	Unit  string `json:"unit"`
	Value int32  `json:"val"`
}

type EngineData struct {
	Code         string      `json:"code"`
	Manufacturer string      `json:"mfg"`
	Model        string      `json:"model"`
	Type         EngineType  `json:"type"`
	Power        EnginePower `json:"power"`
}

func (d *EngineData) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (d *EngineData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

func EngineDataFromRef(ref faaservices.EngineReference) *EngineData {
	d := &EngineData{
		Code:         strings.TrimSpace(ref.Code),
		Manufacturer: strings.TrimSpace(ref.Manufacturer),
		Model:        strings.TrimSpace(ref.Model),
		Type:         EngineType(mustParseInt(ref.Type)),
	}

	d.Power = EnginePower{
		Unit: d.Type.Power(),
	}

	switch d.Power.Unit {
	case "hp":
		d.Power.Value = mustParseInt(ref.Horsepower)
	case "lb":
		d.Power.Value = mustParseInt(ref.Thrust)
	}

	return d
}

type AircraftType int32

const (
	AircraftType_Unknown          AircraftType = iota // 0
	AircraftType_Glider                               // 1
	AircraftType_Balloon                              // 2
	AircraftType_Blimp                                // 3
	AircraftType_SingleEngine                         // 4
	AircraftType_MultiEngine                          // 5
	AircraftType_Rotorcraft                           // 6
	AircraftType_WeightShift                          // 7
	AircraftType_PoweredParachute                     // 8
	AircraftType_Gyroplane                            // 9
)

type AircraftCategory int32

const (
	AircraftCategory_Unknown AircraftCategory = iota // 0
	AircraftCategory_Land                            // 1
	AircraftCategory_Sea                             // 2
	AircraftCategory_Amphib                          // 3
)

type CertificationType int32

const (
	CertificationType_Unknown             CertificationType = iota // 0
	CertificationType_TypeCertificated                             // 1
	CertificationType_NonTypeCertificated                          // 2
	CertificationType_LightSport                                   // 3
)

type WeightClass int32

const (
	WeightClass_Unknown WeightClass = iota // 0
	WeightClass_Light                      // 1 - below 12,500lb
	WeightClass_Medium                     // 2 - 12,500 - 19,999lb
	WeightClass_Heavy                      // 3 - 20,000lb+
	WeightClass_UAV                        // 4 - UAV up to 55lb
)

type AirframeData struct {
	Code          string            `json:"code"`
	Manufacturer  string            `json:"mfg"`
	Model         string            `json:"model"`
	Type          AircraftType      `json:"type"`
	EngineType    EngineType        `json:"eng_type"`
	Category      AircraftCategory  `json:category"`
	Certification CertificationType `json:"cert_type"`
	NumEngines    int32             `json:"num_engines"`
	NumSeats      int32             `json:"num_seats"`
	Weight        WeightClass       `json:"weight"`
	Speed         int32             `json:"speed"`
}

func (d *AirframeData) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (d *AirframeData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

func AirframeDataFromRef(ref faaservices.AircraftReference) *AirframeData {
	return &AirframeData{
		Code:          strings.TrimSpace(ref.Code),
		Manufacturer:  strings.TrimSpace(ref.Manufacturer),
		Model:         strings.TrimSpace(ref.Model),
		Type:          AircraftType(mustParseInt(ref.TypeAircraft)),
		EngineType:    EngineType(mustParseInt(ref.TypeEngine)),
		Category:      AircraftCategory(mustParseInt(ref.Category)),
		Certification: CertificationType(mustParseInt(ref.BuilderCert)),
		NumEngines:    mustParseInt(ref.NumEngines),
		NumSeats:      mustParseInt(ref.NumSeats),
		Weight:        WeightClass(mustParseInt(ref.Weight)),
		Speed:         mustParseInt(ref.Speed),
	}
}

type RegistrantType int32

const (
	RegistrantType_Unknown            RegistrantType = iota // 0
	RegistrantType_Individual                               // 1
	RegistrantType_Partnership                              // 2
	RegistrantType_Corporation                              // 3
	RegistrantType_CoOwned                                  // 4
	RegistrantType_Government                               // 5
	RegistrantType_NonCitizenCorp                           // 6
	RegistrantType_NonCitizenCoOwnerd                       // 7
)

type Region string

const (
	Region_Unknown           Region = ""
	Region_Eastern           Region = "1"
	Region_Southwestern      Region = "2"
	Region_Central           Region = "3"
	Region_WesternPacific    Region = "4"
	Region_Alaskan           Region = "5"
	Region_Southern          Region = "7"
	Region_European          Region = "8"
	Region_GreatLakes        Region = "C"
	Region_NewEngland        Region = "E"
	Region_NorthwestMountain Region = "S"
)

type RegistrationData struct {
	RegistrationNumber   string         `json:"n_num"`
	SerialNumber         string         `json:"ser_num"`
	AirframeCode         string         `json:"airframe_code"`
	AirframeData         *AirframeData  `json:"airframe_data"`
	EngineCode           string         `json:"engine_code"`
	EngineData           *EngineData    `json:"engine_data"`
	YearManufactured     int32          `json:"year"`
	RegistrantType       RegistrantType `json:"registrant_type"`
	Name                 string         `json:"name"`
	Street               string         `json:"street"`
	Street2              string         `json:"street2"`
	City                 string         `json:"city"`
	State                string         `json:"state"`
	Zip                  string         `json:"zip"`
	Region               Region         `json:"region"`
	County               string         `json:"county"`
	Country              string         `json:"country"`
	LastActivityDate     time.Time      `json:"last_activity_date"`
	CertificateIssueDate time.Time      `json:"cert_issue_date"`
	Certification        string         `json:"certification"` // todo: this is a complicated field that can be parsed if i care
	StatusCode           string         `json:"status_code"`   // todo: this can be parsed
	ModeSCode            string         `json:"mode_s"`
	Fractional           bool           `json:"fractional"`
	AirworthinessDate    time.Time      `json:"airworthiness_date"`
	OtherNames           []string       `json:"other_names"`
	ExpirationDate       time.Time      `json:"exp_date"`
	UniqueID             string         `json:"unique_id"`
	KitManufacturer      string         `json:"kit_mfg"`
	KitModel             string         `json:"kit_model"`
	ModeSCodeHex         string         `json:"mode_s_hex"`
}

func (d *RegistrationData) MarshalBinary() ([]byte, error) {
	d.AirframeData = nil
	d.EngineData = nil
	return json.Marshal(d)
}

func (d *RegistrationData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

func RegistrationDataFromRef(ref faaservices.AircraftRegistration) *RegistrationData {
	d := &RegistrationData{
		RegistrationNumber:   strings.TrimSpace(ref.RegNumber),
		SerialNumber:         strings.TrimSpace(ref.SerialNumber),
		AirframeCode:         strings.TrimSpace(ref.AircraftCode),
		EngineCode:           strings.TrimSpace(ref.EngineCode),
		YearManufactured:     mustParseInt(ref.MfgYear),
		RegistrantType:       RegistrantType(mustParseInt(ref.RegistrantType)),
		Name:                 strings.TrimSpace(ref.Name),
		Street:               strings.TrimSpace(ref.Street),
		Street2:              strings.TrimSpace(ref.Street2),
		City:                 strings.TrimSpace(ref.City),
		State:                strings.TrimSpace(ref.State),
		Zip:                  strings.TrimSpace(ref.Zip),
		Region:               Region(strings.TrimSpace(ref.Region)),
		County:               strings.TrimSpace(ref.County),
		Country:              strings.TrimSpace(ref.Country),
		LastActivityDate:     mustParseDate(ref.LastActionDate),
		CertificateIssueDate: mustParseDate(ref.CertIssueDate),
		Certification:        strings.TrimSpace(ref.Certification),
		StatusCode:           strings.TrimSpace(ref.StatusCode),
		ModeSCode:            strings.TrimSpace(ref.ModeSCode),
		ModeSCodeHex:         strings.TrimSpace(ref.ModeSCodeHex),
		Fractional:           mustParseBool(ref.FractOwner),
		AirworthinessDate:    mustParseDate(ref.AirworthinessDate),
		ExpirationDate:       mustParseDate(ref.ExpirationDate),
		UniqueID:             strings.TrimSpace(ref.UniqueID),
		KitManufacturer:      strings.TrimSpace(ref.KitMfg),
		KitModel:             strings.TrimSpace(ref.KitModel),
		OtherNames:           []string{},
	}

	for _, s := range []string{ref.OtherNames1, ref.OtherNames2, ref.Othernames3, ref.Othernames4, ref.Othernames5} {
		s = strings.TrimSpace(s)
		if s != "" {
			d.OtherNames = append(d.OtherNames, s)
		}
	}

	return d
}
