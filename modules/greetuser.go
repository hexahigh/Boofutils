package modules

import (
	"time"
)

// Greet returns a string greeting based on the current time.
func Greet() string {
	now := time.Now()
	hour := now.Hour()

	switch {
	case hour < 3:
		return "It's late, go to bed"
	case hour < 12:
		return "Good morning"
	case hour < 17:
		return "Good afternoon"
	case hour < 20:
		return "Good evening"
	default:
		return "Good night"
	}
}
