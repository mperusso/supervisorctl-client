package supervisorctl

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// ProgramInfo represents the information about a supervisor program.
// It contains details like name, state, PID, and uptime.
type ProgramInfo struct {
	Name        string `json:"name"`        // Program name
	Description string `json:"description"` // Program description
	State       int    `json:"state"`       // Program state code
	StateName   string `json:"state_name"`  // Program state name
	PID         int    `json:"pid"`         // Process ID
	Uptime      string `json:"uptime"`      // Program uptime
}

// Cmd interface for executing commands
type Cmd interface {
	Run() error
	Output() ([]byte, error)
}

// CommandExecutor interface for executing commands
type CommandExecutor interface {
	Command(name string, arg ...string) Cmd
}

// DefaultCommandExecutor implements CommandExecutor using os/exec
type DefaultCommandExecutor struct{}

func (d *DefaultCommandExecutor) Command(name string, arg ...string) Cmd {
	return exec.Command(name, arg...)
}

// Client represents a supervisorctl client.
// It maintains the configuration file path and provides methods to interact with supervisor.
type Client struct {
	configFile      string
	commandExecutor CommandExecutor
}

// NewClient creates a new supervisorctl client.
// If configFile is empty, it will use the default supervisor configuration.
func NewClient(configFile string) *Client {
	return &Client{
		configFile:      configFile,
		commandExecutor: &DefaultCommandExecutor{},
	}
}

// StatusOptions represents the options for the Status method.
type StatusOptions struct {
	// Names is a list of program names to get status for.
	// If empty, returns status for all programs.
	Names []string
}

// Status returns the status of supervisor programs.
// If no names are provided in opts, returns status for all programs.
func (c *Client) Status(opts StatusOptions) ([]ProgramInfo, error) {
	args := []string{"status"}
	if len(opts.Names) > 0 {
		args = append(args, opts.Names...)
	}
	if c.configFile != "" {
		args = append([]string{"-c", c.configFile}, args...)
	}
	cmd := c.commandExecutor.Command("supervisorctl", args...)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 3 || exitErr.ExitCode() == 4 {
			} else {
				return nil, fmt.Errorf("failed to get status: %w - output %s ", err, string(output))
			}
		} else {
			return nil, fmt.Errorf("failed to get status: %w", err)
		}
	}

	var programs []ProgramInfo
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		program, err := parseStatusLine(line)
		if err != nil {
			return nil, fmt.Errorf("failed to parse status line: %w", err)
		}
		programs = append(programs, program)
	}

	return programs, nil
}

// Start starts a supervisor program.
func (c *Client) Start(programName string) error {
	args := []string{"start", programName}
	if c.configFile != "" {
		args = append([]string{"-c", c.configFile}, args...)
	}
	cmd := c.commandExecutor.Command("supervisorctl", args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start program: %w", err)
	}
	return nil
}

// Stop stops a supervisor program.
func (c *Client) Stop(programName string) error {
	args := []string{"stop", programName}
	if c.configFile != "" {
		args = append([]string{"-c", c.configFile}, args...)
	}
	cmd := c.commandExecutor.Command("supervisorctl", args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop program: %w", err)
	}
	return nil
}

// Restart restarts a supervisor program.
func (c *Client) Restart(programName string) error {
	args := []string{"restart", programName}
	if c.configFile != "" {
		args = append([]string{"-c", c.configFile}, args...)
	}
	cmd := c.commandExecutor.Command("supervisorctl", args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restart program: %w", err)
	}
	return nil
}

// parseStatusLine parses a single line of supervisorctl status output.
func parseStatusLine(line string) (ProgramInfo, error) {
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return ProgramInfo{}, fmt.Errorf("invalid status line: %s", line)
	}

	name := parts[0]
	state := parts[1]
	description := strings.Join(parts[2:], " ")

	// Handle "no such process" error
	if strings.Contains(description, "no such process") {
		return ProgramInfo{
			Name:        name,
			State:       1000, // UNKNOWN
			StateName:   "UNKNOWN",
			Description: "no such process",
		}, nil
	}

	stateMap := map[string]int{
		"STOPPED":  0,
		"STARTING": 10,
		"RUNNING":  20,
		"BACKOFF":  30,
		"STOPPING": 40,
		"EXITED":   100,
		"FATAL":    200,
		"UNKNOWN":  1000,
	}

	stateCode, ok := stateMap[state]
	if !ok {
		return ProgramInfo{}, fmt.Errorf("unknown state: %s", state)
	}

	var pid int
	var uptime string
	if strings.Contains(description, "pid") {
		pidStr := strings.Split(description, "pid")[1]
		pidStr = strings.TrimSpace(strings.Split(pidStr, ",")[0])
		pid, _ = strconv.Atoi(pidStr)

		if strings.Contains(description, "uptime") {
			uptimeParts := strings.Split(description, "uptime")
			if len(uptimeParts) > 1 {
				uptime = strings.TrimSpace(strings.Split(uptimeParts[1], ",")[0])
			}
		}
	}

	return ProgramInfo{
		Name:        name,
		State:       stateCode,
		StateName:   state,
		Description: description,
		PID:         pid,
		Uptime:      uptime,
	}, nil
}
