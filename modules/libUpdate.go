package modules

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
)

func Upd_main() {
	fmt.Println("Starting the update process...")

	url := "https://github.com/hexahigh/Boofutils/releases/download/latest_auto/boofutils-" + runtime.GOOS + "-" + runtime.GOARCH
	outputFile := "boofutils_new"
	if runtime.GOOS == "windows" {
		url += ".exe"
		outputFile += ".exe"
	}

	fmt.Println("Downloading the latest version from:", url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Saving the new version as:", outputFile)
	out, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}

	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	oldExePath := exePath + "_old"
	fmt.Println("Renaming the current version to:", oldExePath)
	err = os.Rename(exePath, oldExePath)
	if err != nil {
		panic(err)
	}

	fmt.Println("Replacing the current version with the new version")
	err = os.Rename(outputFile, exePath)
	if err != nil {
		panic(err)
	}

	fmt.Println("Update completed successfully!")
}
