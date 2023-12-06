package modules

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/go-git/go-git/v5"
)

func Upd_main(binary bool) {
	if binary {
		Upd_main_binary()
	} else {
		Upd_main_source()
	}
}

func Upd_main_source() {
	if runtime.GOOS == "windows" {
		fmt.Println("This autoupdater does not work on Windows. Please use the -b flag when updating to use a precompiled binary.")
		os.Exit(0)
	}
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the update process...")

	fmt.Println("Cloning the latest version from:", "https://github.com/hexahigh/boofutils")
	_, err = git.PlainClone("/tmp/bu", false, &git.CloneOptions{
		URL:      "https://github.com/hexahigh/boofutils",
		Progress: os.Stdout,
	})
	CheckIfError(err)

	fmt.Println("Compiling the new version...")

	outputFile := exePath + "_new"

	cmd := exec.Command("go", "build", "-o", outputFile)
	cmd.Dir = "/tmp/bu"
	err = cmd.Run()
	CheckIfError(err)

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

	fmt.Println("Fixing permissions")
	err = os.Chmod(exePath, 0755)
	if err != nil {
		panic(err)
	}

	fmt.Println("Update completed successfully!")
	os.Exit(0)
}

func Upd_main_binary() {
	fmt.Println("Starting the update process...")

	if runtime.GOOS == "windows" {
		fmt.Println("Warning: Windows does not like it when you rename running executables.\nIt may work or it might not.")
	}

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

	if resp.StatusCode != 200 {
		panic(fmt.Errorf("unexpected status: %d %s", resp.StatusCode, resp.Status))
	}

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

	fmt.Println("Fixing permissions")
	err = os.Chmod(exePath, 0755)
	if err != nil {
		panic(err)
	}

	fmt.Println("Update completed successfully!")
	os.Exit(0)
}
