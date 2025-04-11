package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mperusso/supervisorctl-client"
)

func main() {
	// Create a new client without passing a config file path in the command line
	client := supervisorctl.NewClient("")

	// Get status of all programs
	programs, err := client.Status(supervisorctl.StatusOptions{})
	if err != nil {
		log.Fatalf("Failed to get status: %v", err)
	}

	// Print all programs as JSON
	jsonData, err := json.MarshalIndent(programs, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Println("All programs:")
	fmt.Println(string(jsonData))

	// Get status of specific programs
	specificPrograms, err := client.Status(supervisorctl.StatusOptions{
		Names: []string{"program_name", "program_name2"},
	})
	if err != nil {
		log.Fatalf("Failed to get specific programs status: %v", err)
	}

	// Print specific programs as JSON
	jsonData, err = json.MarshalIndent(specificPrograms, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Println("\nSpecific programs:")
	fmt.Println(string(jsonData))

	// Start a program
	if err := client.Start("program_name"); err != nil {
		log.Printf("Failed to start program: %v", err)
	}

	// Stop a program
	if err := client.Stop("program_name"); err != nil {
		log.Printf("Failed to stop program: %v", err)
	}

	// Restart a program
	if err := client.Restart("program_name"); err != nil {
		log.Printf("Failed to restart program: %v", err)
	}
}
