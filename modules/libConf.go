package modules

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// saveToJSONFile saves the given data to a JSON file located at the specified file path.
//
// The function takes two parameters: data of type interface{} which represents the data to be saved,
// and filePath of type string which represents the path where the JSON file will be saved.
//
// The function returns an error in case there is an issue creating or writing to the file. If the operation is successful,
// it returns nil.
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

// AskUserQuestions prompts the user to answer a series of questions and saves the answers to a JSON file.
//
// No parameters.
// No return type.
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

// GetOptionFromConfig retrieves the value of the specified option from the configuration file.
//
// Parameters:
// - optionName: the name of the option to retrieve from the configuration file (type: string)
//
// Returns:
// - string: the value of the specified option (or an empty string if the option is not found)
// - error: an error if there was a problem reading the configuration file or parsing its contents
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

// GenerateDefaultConfig generates the default configuration.
//
// No parameters.
// No return type.
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

// CheckConfigFileExists checks if the configuration file exists.
//
// It does not take any parameters.
// It returns a boolean value indicating if the configuration file exists or not.
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
