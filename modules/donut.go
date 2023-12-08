package modules

import (
	"fmt"
	"math"
	"time"
)

func Donut_main() {
	A := 0.0
	B := 0.0
	var b [1760]byte
	var z [1760]float64

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
				o := x + 80*y
				N := int(8 * ((f*e-c*d*g)*m - c*d*e - f*g - l*d*n))
				if 0 <= y && y < 22 && 0 <= x && x < 80 && D > z[o] {
					z[o] = D
					b[o] = ".,-~:;=!*#$@"[N]
				}
			}
		}

		fmt.Print("\x1b[H")
		for k := 0; k < 1761; k++ {
			if k%80 > 0 {
				fmt.Printf("%c", b[k])
			} else {
				fmt.Print("\n")
			}
		}

		time.Sleep(time.Second / 24)

		A += 0.04
		B += 0.02
	}
}
