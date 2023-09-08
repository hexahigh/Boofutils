package main

import (
	"time"
)

func Greet() string {
	currentTime := time.Now()
	hour := currentTime.Hour()

	if hour < 12 && hour >= 5 {
		return "Good morning"
	} else if hour < 17 {
		return "Good afternoon"
	} else if hour < 20 {
		return "Good evening"
	} else if hour < 24 {
		return "Good night"
	} else {
		return "You should go to bed"
	}
}
