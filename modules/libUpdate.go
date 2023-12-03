package modules

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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

	if runtime.GOOS == "windows" {
		fmt.Println("Creating update script")
		batch := []byte(fmt.Sprintf(`@echo off
:wait
tasklist /fi "imagename eq %s" 2>nul | find /i /n "%s">nul
if "%errorlevel%"=="0" goto wait
move /Y "%s" "%s"
del "%s"
`, filepath.Base(exePath), filepath.Base(exePath), outputFile, exePath, oldExePath))

		script := filepath.Join(filepath.Dir(exePath), "bu_update.bat")
		err = os.WriteFile(script, batch, 0644)
		if err != nil {
			panic(err)
		}

		fmt.Println("Running update script")
		cmd := exec.Command("cmd.exe", "/C", "start", "/B", script)
		err = cmd.Start()
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Replacing the current version with the new version")
		err = os.Rename(outputFile, exePath)
		if err != nil {
			panic(err)
		}

		fmt.Println("Update completed successfully!")
	}
}
