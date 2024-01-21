package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	developmentFileName := ".env"
	exampleFileName := ".env.example"

	developmentContent, err := readFile(developmentFileName)
	if err != nil {
		fmt.Println("Error reading .env:", err)
		return
	}

	exampleContent := processDevelopmentContent(developmentContent)

	err = writeFile(exampleFileName, exampleContent)
	if err != nil {
		fmt.Println("Error writing .env.example:", err)
		return
	}

	fmt.Println(".env.example generated successfully.")
}

func readFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var content strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line contains an equal sign
		if strings.Contains(line, "=") {
			// If yes, add the variable name and placeholder value to .env.example
			parts := strings.SplitN(line, "=", 2)
			line = parts[0] + "=<placeholder>"
		}

		content.WriteString(line + "\n")
	}

	return content.String(), scanner.Err()
}

func writeFile(filename string, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

func processDevelopmentContent(developmentContent string) string {
	// You can add custom logic here to modify the development content
	// For example, replacing actual values with placeholders or adding comments

	// In this example, we return the development content as is
	return developmentContent
}
