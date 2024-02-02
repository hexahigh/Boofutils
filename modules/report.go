package modules

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"
)

type SystemReport struct {
	CPUInfo      map[string]string `json:"cpu_info"`
	OS           map[string]string `json:"os"`
	MEMInfo      map[string]string `json:"mem_info"`
	Env          map[string]string `json:"environment_variables"`
	BlockDevices []BlockDevice     `json:"block_devices"`
	SwapInfo     map[string]string `json:"swap_info"`
	LscpuInfo    map[string]string `json:"lscpu_info"`
}

type BlockDevice struct {
	Name        string   `json:"name"`
	MajMin      string   `json:"maj:min"`
	Rm          bool     `json:"rm"`
	Size        string   `json:"size"`
	Ro          bool     `json:"ro"`
	Type        string   `json:"type"`
	MountPoints []string `json:"mountpoints"`
}

var lsblkOutput struct {
	BlockDevices []BlockDevice `json:"blockdevices"`
}

var memMap, cpuMap, osMap, swapMap, lscpuMap map[string]string

func Report(out_file string, stdout bool) {
	// Get CPU Info
	cpuInfo, err := os.Open("/proc/cpuinfo")
	if err != nil {
		fmt.Println("Error opening /proc/cpuinfo:", err)
	} else {
		defer cpuInfo.Close()

		scanner := bufio.NewScanner(cpuInfo)
		cpuMap = make(map[string]string)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				cpuMap[key] = val
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading /proc/cpuinfo:", err)
			return
		}
	}

	// Get OS Info
	osInfo, err := os.Open("/etc/os-release")
	if err != nil {
		fmt.Println("Error opening /etc/os-release:", err)
	} else {
		defer osInfo.Close()

		scanner := bufio.NewScanner(osInfo)
		osMap = make(map[string]string)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := parts[0]
				val := strings.Trim(parts[1], `"`)
				osMap[key] = val
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading /etc/os-release:", err)
			return
		}
	}

	// Get mem info
	memInfo, err := os.Open("/proc/meminfo")
	if err != nil {
		fmt.Println("Error opening /proc/meminfo:", err)
	} else {
		defer cpuInfo.Close()

		scanner := bufio.NewScanner(memInfo)
		memMap = make(map[string]string)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				memMap[key] = val
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading /proc/meminfo:", err)
			return
		}
	}

	// Get env
	envMap := make(map[string]string)
	for _, envVar := range os.Environ() {
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			val := parts[1]
			envMap[key] = val
		}
	}

	// Get swap info
	swapInfo, err := os.Open("/proc/swaps")
	if err != nil {
		fmt.Println("Error opening /proc/swaps:", err)
	} else {
		defer swapInfo.Close()

		scanner := bufio.NewScanner(swapInfo)
		swapMap = make(map[string]string)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				key := parts[0]
				val := strings.Join(parts[1:], " ")
				swapMap[key] = val
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading /proc/swaps:", err)
			return
		}
	}

	// Execute lsblk -J
	cmd := exec.Command("lsblk", "-J")
	out, err := executeCmd(cmd)
	if err != nil {
		fmt.Println("Error executing lsblk -J:", err)
	} else {

		// Parse lsblk output
		err = json.Unmarshal(out, &lsblkOutput)
		if err != nil {
			fmt.Println("Error parsing lsblk output:", err)
			return
		}
	}

	// Lscpu
	cmd = exec.Command("lscpu")
	out, err = executeCmd(cmd)
	if err != nil {
		fmt.Println("Error executing lscpu", err)
	} else {
		defer cpuInfo.Close()

		scanner := bufio.NewScanner(bytes.NewReader(out))
		lscpuMap = make(map[string]string)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				lscpuMap[key] = val
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error executing lscpu", err)
			return
		}
	}

	report := &SystemReport{
		CPUInfo:      cpuMap,
		OS:           osMap,
		MEMInfo:      memMap,
		Env:          envMap,
		BlockDevices: lsblkOutput.BlockDevices,
		SwapInfo:     swapMap,
		LscpuInfo:    lscpuMap,
	}

	// Marshal the report to JSON
	jsonData, err := json.MarshalIndent(report, "", " ")
	if err != nil {
		fmt.Println("Error marshalling report to JSON:", err)
		return
	}

	if stdout {
		fmt.Print(string(jsonData))
		return
	} else {
		// Write the JSON data to the output file
		err = os.WriteFile(out_file, jsonData, fs.ModePerm)
		if err != nil {
			fmt.Println("Error writing report to file:", err)
			return
		}

		fmt.Println("System report completed.")
	}
}

func executeCmd(cmd *exec.Cmd) ([]byte, error) {
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.Bytes(), err
}
