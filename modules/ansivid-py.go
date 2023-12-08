package modules

import (
	"fmt"
	"os/exec"
)

func Ansivid_py_main(filename string, strat string, install bool) {
	if install {
		// Install video-to-ascii
		cmd := exec.Command("pip3", "install", "video-to-ascii", "--install-option=\"--with-audio\"")
		out, err := cmd.Output()

		if err != nil {
			fmt.Println("Error installing video-to-ascii:", err)
			fmt.Println("Command output:", string(out))
			return
		}
	}
	// Run video-to-ascii
	cmd := exec.Command("video-to-ascii", "-f", filename, "-s", strat, "-a")
	out, err := cmd.Output()

	if err != nil {
		fmt.Println("Error running video-to-ascii:", err)
		return
	}

	fmt.Println(string(out))
}
