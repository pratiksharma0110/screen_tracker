package tracker

import (
	"fmt"
	"time"
)

type AppData struct {
	LastClass  string
	LastSwitch time.Time
	TimeSpent  map[string]time.Duration
}

var currentClass string = ""

func (a *AppData) TrackTimer(newClass string) map[string]float64 {

	if a.TimeSpent == nil {
		a.TimeSpent = make(map[string]time.Duration)
	}
	now := time.Now()

	if currentClass != "" && currentClass != newClass {

		elapsed := now.Sub(a.LastSwitch)

		a.TimeSpent[currentClass] += elapsed
	}

	currentClass = newClass
	a.LastSwitch = now
	fmt.Printf("%v", a.LastSwitch)

	secondsMap := make(map[string]float64)
	for k, v := range a.TimeSpent {
		secondsMap[k] = v.Seconds()
	}

	return secondsMap //send time duration in second for easiness

}
