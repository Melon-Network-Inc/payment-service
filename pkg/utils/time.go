package utils

import (
	"fmt"
	"time"
)

func GetPSTLocation() (*time.Location, error) {
	return GetTimeLocation("America/Los_Angeles")
}

func GetPSTTime() (time.Time, error) {
	return GetTimeByTimezone("America/Los_Angeles")
}

func GetESTLocation() (*time.Location, error) {
	return GetTimeLocation("America/New_York")
}

func GetESTTime() (time.Time, error) {
	return GetTimeByTimezone("America/New_York")
}

func GetCDTLocation() (*time.Location, error) {
	return GetTimeLocation("America/Chicago")
}

func GetCDTTime() (time.Time, error) {
	return GetTimeByTimezone("America/Chicago")
}

func GetTimeLocation(timezone string) (*time.Location, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		_ = fmt.Errorf("cannot parse location %s due to %s", timezone, err)
		return &time.Location{}, err
	}
	return loc, nil
}

func GetTimeByTimezone(timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		_ = fmt.Errorf("cannot get the time from location %s due to %s", timezone, err)
		return time.Now().UTC(), err
	}
	return time.Now().UTC().In(loc), nil
}
