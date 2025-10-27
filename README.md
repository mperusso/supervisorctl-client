# supervisorctl-client

[![Go Reference](https://pkg.go.dev/badge/github.com/mperusso/supervisorctl-client.svg)](https://pkg.go.dev/github.com/mperusso/supervisorctl-client)

A Go client for interacting with supervisor programs.

## Features

- Get status of all programs with structured data
- Start programs
- Stop programs
- Restart programs
- Start/Restart without waiting (async)
- JSON output support
- Error handling

## Installation

```bash
go get github.com/mperusso/supervisorctl-client
```

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/mperusso/supervisorctl-client"
)

func main() {
	// Create a new client
	client := supervisorctl.NewClient("")

	// Get status of all programs
	programs, err := client.Status(supervisorctl.StatusOptions{})
	if err != nil {
		log.Fatal(err)
	}

	// Print program information
	for _, program := range programs {
		fmt.Printf("Name: %s\n", program.Name)
		fmt.Printf("State: %s\n", program.StateName)
		fmt.Printf("PID: %d\n", program.PID)
		fmt.Printf("Uptime: %s\n", program.Uptime)
		fmt.Println("---")
	}

	// Get status of specific programs
	programs, err = client.Status(supervisorctl.StatusOptions{
		Names: []string{"program_name", "program_name2"},
	})

 // Control programs
	err = client.Start("program_name")
	err = client.Stop("program_name")
	err = client.Restart("program_name")

	// Start/Restart without waiting (async)
	_ = client.StartAsync("program_name")   // returns immediately
	_ = client.RestartAsync("program_name") // returns immediately
}
```

## ProgramInfo Structure

```go
type ProgramInfo struct {
	Name        string `json:"name"`        // Program name
	Description string `json:"description"` // Program description
	State       int    `json:"state"`       // Program state code
	StateName   string `json:"state_name"`  // Program state name
	PID         int    `json:"pid"`         // Process ID
	Uptime      string `json:"uptime"`      // Program uptime
}
```

## Program States

The program states are represented by the following codes:

- `0`: STOPPED
- `10`: STARTING
- `20`: RUNNING
- `30`: BACKOFF
- `40`: STOPPING
- `100`: EXITED
- `200`: FATAL
- `1000`: UNKNOWN

## License

MIT 