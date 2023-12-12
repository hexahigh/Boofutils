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

func Upd_main(binary bool, allow_win bool) {
	if binary {
		Upd_main_binary()
	} else {
		Upd_main_source(allow_win)
	}
}

func Upd_main_source(allow_win bool) {

	// Check if we are on Windows
	if runtime.GOOS == "windows" && !allow_win {
		fmt.Println("This autoupdater does not work on Windows. Please use the -b flag when updating to use a precompiled binary.")
		os.Exit(0)
	}

	/*if !checkGo() {
		fmt.Println("It seems like Go is not installed, would you like to install it? (Y/N)")
		if YNtoBool(AskInput()) {
			installGo()
		}
	}*/

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

	fmt.Println("Cleaning up...")
	err = os.RemoveAll("/tmp/bu")
	if err != nil {
		panic(err)
	}
	err = os.Remove(oldExePath)
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

func checkGo() bool {
	// Check if Go is installed
	_, err := exec.LookPath("go")
	if err != nil {
		return false
	}
	return true
}

func installGo() {
	fmt.Println("Installing Bison")
	cmd := exec.Command("sudo", "apt", "install", "-y", "bison")
	err := cmd.Run()
	CheckIfError(err)

	fmt.Println("Installing GVM")
	cmd = exec.Command("bash", "-c", "curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer | bash")
	err = cmd.Run()
	CheckIfError(err)

	fmt.Println("Installing Go 1.4")
	cmd = exec.Command("gvm", "install", "go1.4")
	err = cmd.Run()
	CheckIfError(err)
	exec.Command("gvm", "use", "go1.4").Run()

	fmt.Println("Installing Go 1.5")
	cmd = exec.Command("gvm", "install", "go1.5")
	err = cmd.Run()
	CheckIfError(err)
	exec.Command("gvm", "use", "go1.5").Run()

	fmt.Println("Installing Go 1.9")
	cmd = exec.Command("gvm", "install", "go1.9")
	err = cmd.Run()
	CheckIfError(err)
	exec.Command("gvm", "use", "go1.9").Run()

	fmt.Println("Installing Go 1.18")
	cmd = exec.Command("gvm", "install", "go1.18")
	err = cmd.Run()
	CheckIfError(err)
	exec.Command("gvm", "use", "go1.18").Run()

	fmt.Println("Installing Go 1.21.4")
	cmd = exec.Command("gvm", "install", "go1.21.4")
	err = cmd.Run()
	CheckIfError(err)
	exec.Command("gvm", "use", "go1.21.4", "--default").Run()
}
