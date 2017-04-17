package faaservices

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/spkg/bom"
)

func (c *Client) ParseRegistryData(outChan chan interface{}) error {
	dir, err := ioutil.TempDir("", "faadata")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	c.RegistryFetcher.FetchRegistryData(dir)

	wg := &sync.WaitGroup{}

	engChan := make(chan EngineReference)
	acftChan := make(chan AircraftReference)
	regChan := make(chan AircraftRegistration)

	go parseEngineReferenceFile(dir, engChan, wg)
	go parseAircraftReferenceFile(dir, acftChan, wg)
	go parseRegistrationFile(dir, regChan, wg)

	done := make(chan bool)

	go func() {
		for {
			select {
			case e := <-engChan:
				if e != (EngineReference{}) {
					outChan <- e
				}
			case a := <-acftChan:
				if a != (AircraftReference{}) {
					outChan <- a
				}
			case r := <-regChan:
				if r != (AircraftRegistration{}) {
					outChan <- r
				}
			case <-done:
				return
			}
		}
	}()

	time.Sleep(1 * time.Millisecond) // block long enough for goroutines to spin up.

	wg.Wait()

	done <- true

	return nil
}

type EngineReference struct {
	Code         string `csv:"CODE"`
	Manufacturer string `csv:"MFR"`
	Model        string `csv:"MODEL"`
	Type         string `csv:"TYPE"`
	Horsepower   string `csv:"HORSEPOWER"`
	Thrust       string `csv:"THRUST"`
}

type AircraftReference struct {
	Code         string `csv:"CODE"`
	Manufacturer string `csv:"MFR"`
	Model        string `csv:"MODEL"`
	TypeAircraft string `csv:"TYPE-ACFT"`
	TypeEngine   string `csv:"TYPE-ENG"`
	Category     string `csv:"AC-CAT"`
	BuilderCert  string `csv:"BUILD-CERT-IND"`
	NumEngines   string `csv:"NO-ENG"`
	NumSeats     string `csv:"NO-SEATS"`
	Weight       string `csv:"AC-WEIGHT"`
	Speed        string `csv:"SPEED"`
}

type AircraftRegistration struct {
	RegNumber         string `csv:"N NUMBER"`
	SerialNumber      string `csv:"SERIAL NUMBER"`
	AircraftCode      string `csv:"MFR MDL CODE"`
	EngineCode        string `csv:"ENG MFR MDL"`
	MfgYear           string `csv:"YEAR MFR"`
	RegistrantType    string `csv:"TYPE REGISTRANT"`
	Name              string `csv:"NAME"`
	Street            string `csv:"STREET"`
	Street2           string `csv:"STREET2"`
	City              string `csv:"CITY"`
	State             string `csv:"STATE"`
	Zip               string `csv:"ZIP CODE"`
	Region            string `csv:"REGION"`
	County            string `csv:"COUNTY"`
	Country           string `csv:"COUNTRY"`
	LastActionDate    string `csv:"LAST ACTION DATE"`
	CertIssueDate     string `csv:"CERT ISSUE DATE"`
	Certification     string `csv:"CERTIFICATION"`
	AircraftType      string `csv:"TYPE AIRCRAFT"`
	EngineType        string `csv:"TYPE ENGINE"`
	StatusCode        string `csv:"STATUS CODE"`
	ModeSCode         string `csv:"MODE S CODE"`
	FractOwner        string `csv:"FRACT OWNER"`
	AirworthinessDate string `csv:"AIR WORTH DATE"`
	OtherNames1       string `csv:"OTHER NAMES(1)"`
	OtherNames2       string `csv:"OTHER NAMES(2)"`
	Othernames3       string `csv:"OTHER NAMES(3)"`
	Othernames4       string `csv:"OTHER NAMES(4)"`
	Othernames5       string `csv:"OTHER NAMES(5)"`
	ExpirationDate    string `csv:"EXPIRATION DATE"`
	UniqueID          string `csv:"UNIQUE ID"`
	KitMfg            string `csv:"KIT MFR"`
	KitModel          string `csv:"KIT MODEL"`
	ModeSCodeHex      string `csv:"MODE S CODE HEX"`
}

func parseEngineReferenceFile(srcdir string, ch chan EngineReference, wg *sync.WaitGroup) {
	wg.Add(1)
	filename := path.Join(srcdir, "ENGINE.txt")
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening %s: %v\n", filename, err)
		wg.Done()
		return
	}
	defer file.Close()

	time.Sleep(2 * time.Second)

	gocsv.SetCSVReader(gocsv.LazyCSVReader)
	gocsv.FailIfDoubleHeaderNames = true
	gocsv.FailIfUnmatchedStructTags = true

	if err := gocsv.UnmarshalToChan(bom.NewReader(file), ch); err != nil {
		log.Printf("Error unmarshaling engine data: %v\n", err)
	}
	wg.Done()
}

func parseAircraftReferenceFile(srcdir string, ch chan AircraftReference, wg *sync.WaitGroup) {
	wg.Add(1)
	filename := path.Join(srcdir, "ACFTREF.txt")
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening %s: %v\n", filename, err)
		wg.Done()
		return
	}
	defer file.Close()

	time.Sleep(2 * time.Second)

	gocsv.SetCSVReader(gocsv.LazyCSVReader)
	gocsv.FailIfDoubleHeaderNames = true
	gocsv.FailIfUnmatchedStructTags = true

	if err := gocsv.UnmarshalToChan(bom.NewReader(file), ch); err != nil {
		log.Printf("Error unmarshaling aircraft data: %v\n", err)
	}
	wg.Done()
}

func parseRegistrationFile(srcdir string, ch chan AircraftRegistration, wg *sync.WaitGroup) {
	wg.Add(1)
	filename := path.Join(srcdir, "MASTER.txt")
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening %s: %v\n", filename, err)
		wg.Done()
		return
	}
	defer file.Close()

	time.Sleep(2 * time.Second)

	gocsv.SetCSVReader(gocsv.LazyCSVReader)
	gocsv.FailIfDoubleHeaderNames = true
	gocsv.FailIfUnmatchedStructTags = true

	if err := gocsv.UnmarshalToChan(bom.NewReader(file), ch); err != nil {
		log.Printf("Error unmarshaling registration data: %v\n", err)
	}
	wg.Done()
}
