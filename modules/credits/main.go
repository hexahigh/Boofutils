package credits

import (
	"context"
	"fmt"
	"time"

	m "github.com/hexahigh/boofutils/modules"
	con "github.com/hexahigh/boofutils/modules/constants"
)

func Main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go m.PlayAudioAdvanced(ctx, "credits.flac", 16, 2, 48000)

	fmt.Print(con.Art)

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
