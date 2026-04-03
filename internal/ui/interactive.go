package ui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Confirm prompts the user for yes/no confirmation.
func Confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [Y/n] ", prompt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "" || input == "y" || input == "yes"
}

// PromptChoice shows options and returns the user's choice.
func PromptChoice(message string, options []string) int {
	fmt.Println(message)
	for i, opt := range options {
		fmt.Printf("  %d) %s\n", i+1, opt)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Choice: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		for i := range options {
			if input == fmt.Sprintf("%d", i+1) {
				return i
			}
		}
		fmt.Println("Invalid choice, try again.")
	}
}

// EditInEditor opens the given text in the user's $EDITOR and returns the edited result.
func EditInEditor(initial string) (string, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		editor = "vi"
	}

	tmpFile, err := os.CreateTemp("", "gitwise-*.txt")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(initial); err != nil {
		tmpFile.Close()
		return "", err
	}
	tmpFile.Close()

	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editor exited with error: %w", err)
	}

	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}
