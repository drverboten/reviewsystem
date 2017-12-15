package models

import (
	"time"
)

// GetSubmitTime returns the time on the specific time zone when the link was submitted default: America/Mexico_City
func GetSubmitTime() string {
	nowUTC := time.Now().UTC()

	locationTime, err := time.LoadLocation("America/Mexico_City")
	if err != nil {
		panic(err)
	}

	submitTime := nowUTC.In(locationTime)
	return submitTime.Format(time.RFC3339)
}
