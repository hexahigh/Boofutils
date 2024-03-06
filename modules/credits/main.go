package credits

import (
	"context"
	"fmt"
	"strings"
	"time"

	m "github.com/hexahigh/boofutils/modules"
	color "github.com/hexahigh/boofutils/modules/color"
	con "github.com/hexahigh/boofutils/modules/constants"
)

func Main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go m.PlayAudioAdvanced(ctx, "credits.flac", 16, 2, 48000)

	lines := strings.Split(con.Art2, "\n")
	fmt.Print("\033[38;2;0;255;0;48;2;0;0;0m")
	for _, line := range lines {
		fmt.Println(line)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Print(color.ColorReset)

	for _, line := range con.Credits {
		// Print each character of the line one by one
		for _, char := range line {
			fmt.Print(string(char))
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Println()
		time.Sleep(1000 * time.Millisecond) // Pause before starting the next line
	}

	fmt.Println("Press enter to exit...")
	fmt.Scanln()
	cancel()

}
