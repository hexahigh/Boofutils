package modules

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func saveToJSONFile(data interface{}, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}

func AskUserQuestions() {
	var userInput struct {
		Question1 string `json:"name"`
	}

	fmt.Println("Please answer the following questions:")
	fmt.Println("This data will never leave your device")
	fmt.Print("What is your name?: ")
	userInput.Question1 = AskInput()

	appdataDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Error getting user's appdata directory:", err)
		return
	}

	filePath := filepath.Join(appdataDir, "Boofdev", "boofutils.conf")
	if err := saveToJSONFile(userInput, filePath); err != nil {
		fmt.Println("Error saving data to JSON file:", err)
		return
	}

	fmt.Println("Data saved to", filePath)
}

func GetOptionFromConfig(optionName string) (string, error) {
	appdataDir, err := os.UserConfigDir()
	configFile, err := os.Open(filepath.Join(appdataDir, "Boofdev", "boofutils.conf"))
	if err != nil {
		return "", err
	}
	defer configFile.Close()

	// Parse the config file
	var configData map[string]string
	err = json.NewDecoder(configFile).Decode(&configData)
	if err != nil {
		return "", err
	}

	// Get the option from the config data
	option, ok := configData[optionName]
	if !ok {
		return "", fmt.Errorf("Option '%s' not found in config file", optionName)
	}

	return option, nil
}

func GenerateDefaultConfig() {
	var userInput struct {
		Question1 string `json:"name"`
	}

	userInput.Question1 = "Anon"

	appdataDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Error getting user's appdata directory:", err)
		return
	}

	filePath := filepath.Join(appdataDir, "Boofdev", "boofutils.conf")
	if err := saveToJSONFile(userInput, filePath); err != nil {
		fmt.Println("Error saving data to JSON file:", err)
		return
	}

	fmt.Println("Data saved to", filePath)
}

func CheckConfigFileExists() bool {
	appdataDir, err := os.UserConfigDir()
	_, err = os.Stat(filepath.Join(appdataDir, "Boofdev", "boofutils.conf"))
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	// handle other errors
	return false
}
