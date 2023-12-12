package modules

import (
	"fmt"
	"math"
	"time"
)

func Donut_main(spinSpeed float64, isRainbow bool, trueColor bool) {
	A := 0.0
	B := 0.0
	var b [1760]byte
	var z [1760]float64
	/*var colors []string
	if trueColor {
		colors = []string{
			"\033[38;2;255;0;0m",     // Red
			"\033[38;2;255;165;0m",   // Orange
			"\033[38;2;255;255;0m",   // Yellow
			"\033[38;2;0;128;0m",     // Green
			"\033[38;2;0;0;255m",     // Blue
			"\033[38;2;75;0;130m",    // Indigo
			"\033[38;2;238;130;238m", // Violet
		}
	} else {
		colors = []string{
			"\033[38;5;196m", // Red
			"\033[38;5;202m", // Orange
			"\033[38;5;226m", // Yellow
			"\033[38;5;46m",  // Green
			"\033[38;5;21m",  // Blue
			"\033[38;5;129m", // Indigo
			"\033[38;5;201m", // Violet
		}
	}*/
	for {
		for i := range b {
			b[i] = ' '
		}
		for i := range z {
			z[i] = 0
		}

		for j := 0.0; j < 6.28; j += 0.07 {
			for i := 0.0; i < 6.28; i += 0.02 {
				c := math.Sin(i)
				d := math.Cos(j)
				e := math.Sin(A)
				f := math.Sin(j)
				g := math.Cos(A)
				h := d + 2
				D := 1 / (c*h*e + f*g + 5)
				l := math.Cos(i)
				m := math.Cos(B)
				n := math.Sin(B)
				t := c*h*g - f*e
				x := int(40 + 30*D*(l*h*m-t*n))
				y := int(12 + 15*D*(l*h*n+t*m))
				if 0 <= y && y < 22 && 0 <= x && x < 80 {
					o := x + 80*y
					if 0 <= o && o < len(z) {
						z[o] = D
						if N := int(8 * ((f*e-c*d*g)*m - c*d*e - f*g - l*d*n)); N >= 0 && N < len(".,-~:;=!*#$@") {
							b[o] = ".,-~:;=!*#$@"[N]
						}
					}
				}
			}
		}

		fmt.Print("\x1b[H")
		for k := 0; k < 1761; k++ {
			if k%80 > 0 {
				if isRainbow {
					fmt.Printf("%s%c", rainbow(k), b[k])
				} else {
					fmt.Printf("%c", b[k])
				}
			} else {
				fmt.Print("\n")
			}
		}

		time.Sleep(time.Second / 24)

		A += 0.04 * spinSpeed
		B += 0.02 * spinSpeed
	}
}

func rainbow(i int) string {
	freq := 0.1
	red := math.Sin(freq*float64(i)+0)*127 + 128
	green := math.Sin(freq*float64(i)+2*math.Pi/3)*127 + 128
	blue := math.Sin(freq*float64(i)+4*math.Pi/3)*127 + 128
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", int(red), int(green), int(blue))
}
