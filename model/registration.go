package model

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"gopkg.in/redis.v5"

	"github.com/ianmcmahon/compliancebuddy/faaservices"
	"github.com/ianmcmahon/compliancebuddy/model/faaref"
)

const (
	keyspace_regdata string = "cb:faa:regdata"
)

var (
	ErrNonUSRegistrationNumber = errors.New("Not a US registration number")
	ErrRegistrationNotFound    = errors.New("Registration not found")
)

type RegistrationService struct {
	redisClient *redis.Client
	faaClient   *faaservices.Client
}

func NewRegistrationService(redisClient *redis.Client, faaClient *faaservices.Client) *RegistrationService {
	client := &RegistrationService{
		redisClient: redisClient,
		faaClient:   faaClient,
	}

	return client
}

func (s *RegistrationService) GetRegistration(tailNum string) (*faaref.RegistrationData, error) {
	if []byte(tailNum)[0] != 'N' {
		return nil, ErrNonUSRegistrationNumber
	} else {
		tailNum = tailNum[1:]
	}

	var regData faaref.RegistrationData

	if err := s.getObject(keyspace_regdata, "reg", tailNum, &regData); err != nil {
		return nil, ErrRegistrationNotFound
	}

	regData.AirframeData = &faaref.AirframeData{}
	regData.EngineData = &faaref.EngineData{}

	if err := s.getObject(keyspace_regdata, "acft", regData.AirframeCode, regData.AirframeData); err != nil {
		log.Printf("Couldn't find airframe by code %s: %v\n", regData.AirframeCode, err)
	}

	if err := s.getObject(keyspace_regdata, "eng", regData.EngineCode, regData.EngineData); err != nil {
		log.Printf("Couldn't find engine by code %s: %v\n", regData.EngineCode, err)
	}

	return &regData, nil
}

func (s *RegistrationService) getObject(keyspace string, objType string, key string, v interface{}) error {
	return s.redisClient.Get(fmt.Sprintf("%s:%s:%s", keyspace, objType, key)).Scan(v)
}

func (s *RegistrationService) setObject(keyspace string, objType string, key string, v interface{}, exp time.Duration) error {
	return s.redisClient.Set(fmt.Sprintf("%s:%s:%s", keyspace, objType, key), v, exp).Err()
}

func (s *RegistrationService) UpdateRegistrationData() error {
	c := make(chan interface{}, 10)
	wg := &sync.WaitGroup{}

	go s.receiveRegistrationData(c, wg)

	err := s.faaClient.ParseRegistryData(c)

	close(c)

	wg.Wait()

	return err
}

func (s *RegistrationService) receiveRegistrationData(c chan interface{}, wg *sync.WaitGroup) {
	wg.Add(1)
	for {
		v, more := <-c
		if more {
			switch ref := v.(type) {
			case faaservices.EngineReference:
				eng := faaref.EngineDataFromRef(ref)
				s.setObject(keyspace_regdata, "eng", eng.Code, eng, 0)
			case faaservices.AircraftReference:
				acft := faaref.AirframeDataFromRef(ref)
				s.setObject(keyspace_regdata, "acft", acft.Code, acft, 0)
			case faaservices.AircraftRegistration:
				reg := faaref.RegistrationDataFromRef(ref)
				s.setObject(keyspace_regdata, "reg", reg.RegistrationNumber, reg, 0)
			}
		} else {
			wg.Done()
			return
		}
	}
}
