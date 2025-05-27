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

func (a *AppData) TrackTimer(newClass string) map[string]time.Duration {

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

	return a.TimeSpent

}

func (a *AppData) Print(details map[string]time.Duration) {
	for app, time := range details {
		if app != "" {
			fmt.Printf("You spent %v in %v\n", time, app)

		}
	}
}
