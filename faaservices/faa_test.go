package faaservices

import (
	"encoding/xml"
	"fmt"
	"strings"
	"testing"
)

func TestXMLResponseFormat(t *testing.T) {
	documents := documents{
		List: []document{
			document{
				DocumentNumber: "12345",
				Title:          "foobar",
				Type:           "Airworthiness Directives",
				Uri:            "http://google.com",
				Subject:        "A bogus document",
			},
			document{
				DocumentNumber: "54321",
				Title:          "barbaz",
				Type:           "Airworthiness Directives",
				Uri:            "http://google.com",
				Subject:        "A bogus document",
			},
		},
	}

	b, err := xml.Marshal(&documents)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", b)
}

// func TestADSearch(t *testing.T) {
// 	twinkie := twinkie()

// 	ads, err := ADSearch(twinkie)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("%+v\n", ads)
// }

func TestFetchRegistry(t *testing.T) {
	client := NewClient()
	client.RegistryFetcher = LocalRegistryFetcher{"testdata/abbrev"}

	recv := make(chan interface{})
	done := make(chan bool)

	engineMap := map[string]*EngineReference{}
	aircraftMap := map[string]*AircraftReference{}
	registrationMap := map[string]*AircraftRegistration{}

	go func() {
		for {
			select {
			case v := <-recv:
				switch ref := v.(type) {
				case EngineReference:
					engineMap[ref.Code] = &ref
				case AircraftReference:
					aircraftMap[ref.Code] = &ref
				case AircraftRegistration:
					registrationMap[ref.RegNumber] = &ref
				}
			case <-done:
				return
			}
		}
	}()

	if err := client.ParseRegistryData(recv); err != nil {
		t.Fatal(err)
	}
	done <- true

	aerostar, ok := registrationMap["601UK"]
	if !ok {
		t.Errorf("Couldn't find N601UK in registration map\n")
		t.Fail()
	}

	acft, ok := aircraftMap[aerostar.AircraftCode]
	if !ok {
		t.Errorf("Couldn't find %s in aircraft map\n", aerostar.AircraftCode)
		t.Fail()
	}

	eng, ok := engineMap[aerostar.EngineCode]
	if !ok {
		t.Errorf("Couldn't find %s in engine map\n", aerostar.EngineCode)
		t.Fail()
	}

	if strings.TrimSpace(acft.Model) != "AEROSTAR 601P" {
		t.Errorf("unexpected model: %s\n", strings.TrimSpace(acft.Model))
	}

	if strings.TrimSpace(eng.Model) != "IO-540 SER" {
		t.Errorf("unexpected model: %s\n", strings.TrimSpace(eng.Model))
	}
}
