package model

import (
	"encoding/json"
	"testing"

	"gopkg.in/redis.v5"

	"github.com/alicebob/miniredis"
	"github.com/ianmcmahon/compliancebuddy/faaservices"
)

func TestRegistrationFetch(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	redisClient := redis.NewClient(&redis.Options{Addr: s.Addr()})
	faaClient := faaservices.NewClient()
	faaClient.RegistryFetcher = faaservices.LocalRegistryFetcher{"../faaservices/testdata/abbrev"}

	service := NewRegistrationService(redisClient, faaClient)

	if err := service.UpdateRegistrationData(); err != nil {
		t.Fatal(err)
	}

	reg, err := service.GetRegistration("N7706Y")
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"n_num":"7706Y","ser_num":"30-794","airframe_code":"7103002","airframe_data":{"code":"7103002","mfg":"PIPER","model":"PA-30","type":5,"eng_type":1,"Category":1,"cert_type":0,"num_engines":2,"num_seats":4,"weight":0,"speed":139},"engine_code":"41509","engine_data":{"code":"41509","mfg":"LYCOMING","model":"IO-320 SERIES","type":1,"power":{"unit":"hp","val":150}},"year":0,"registrant_type":3,"name":"ANIKIN AVIATION INC","street":"1201 N ORANGE ST STE 600","street2":"","city":"WILMINGTON","state":"DE","zip":"198011171","region":"1","county":"003","country":"US","last_activity_date":"2014-03-27T00:00:00Z","cert_issue_date":"2014-03-27T00:00:00Z","certification":"1N","status_code":"27","mode_s":"52466110","fractional":false,"airworthiness_date":"1965-05-04T00:00:00Z","other_names":[],"exp_date":"2017-03-31T00:00:00Z","unique_id":"00273477","kit_mfg":"","kit_model":"","mode_s_hex":"AA6C48"}`
	got, err := json.Marshal(reg)
	if err != nil {
		t.Fatal(err)
	}

	if expected != string(got) {
		t.Errorf("expected: %s\ngot:      %s\n", expected, got)
	}
}
