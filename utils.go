package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"

	"inputs"
	"multichoice"
	"unichoice"

	"github.com/fatih/color"
)

type FileResponse struct {
	Content string `json:"content"` // Base64 encoded content
}

func checkError(err error, message string) {
	if err != nil {
		color.Red("%s: %s", message, err)
		os.Exit(1)
	}
}

// This function clears the terminal screen.
func ClearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func updatePrompt() string {
	result, err := unichoice.Run([]string{"Yes", "No"}, "Do you agree to update the system packages")
	if err != nil {
		color.Red("Failed to prompt for update: %s", err)
		return "No"
	}
	return result
}

func userNamePrompt(defaultUsername string) string {
	name, err := inputs.Run(defaultUsername, "Enter your username")
	checkError(err, "Prompt failed")
	if name == "" {
		return defaultUsername
	}
	// Ensure the username is Unix compliant
	validUsernameRegex := regexp.MustCompile(`^[a-z_][a-z0-9_-]{0,30}[^-]$`)
	for !validUsernameRegex.MatchString(name) {
		color.Red("Invalid username! your username must be Unix compliant, please try again.")
		name, err = inputs.Run(defaultUsername, "Enter your username")
		if name == "" {
			return defaultUsername
		}
		checkError(err, "Prompt failed")
	}
	return name
}

func fetchJsonFiles(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch JSON file: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func parseJsonFile(jsonFile string) (jsonData, error) {
	var data jsonData
	err := json.Unmarshal([]byte(jsonFile), &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func profilePrompt(profiles []string) []string {
	title := "Choose your profiles (Profiles are used to select the packages that might interest you)"
	descriptions := make([]string, len(profiles)) // Empty descriptions for profiles
	result, err := multichoice.Run(profiles, descriptions, title)
	if err != nil {
		color.Red("Failed to select profile: %s", err)
		os.Exit(1)
	}
	return result
}

func shellPrompt(shells []string) string {
	title := "Choose your shell"
	result, err := unichoice.Run(shells, title)
	if err != nil {
		color.Red("Failed to select shell: %s", err)
		os.Exit(1)
	}
	return result
}

func packagesPrompt(packages, descriptions []string) []string {
	title := "Select the packages you want to install with Spacebar and confirm with Enter"
	result, err := multichoice.Run(packages, descriptions, title)
	if err != nil {
		color.Red("Failed to select packages: %s", err)
		os.Exit(1)
	}
	return result
}
