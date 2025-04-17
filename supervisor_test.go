package supervisorctl

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCommandExecutor struct {
	mock.Mock
}

func (m *MockCommandExecutor) Command(name string, arg ...string) Cmd {
	args := m.Called(name, arg)
	return args.Get(0).(Cmd)
}

type MockCmd struct {
	mock.Mock
	MockOutput []byte
	Err        error
}

func (m *MockCmd) Run() error {
	return m.Err
}

func (m *MockCmd) Output() ([]byte, error) {
	return m.MockOutput, m.Err
}

func TestNewClient(t *testing.T) {
	client := NewClient("")
	assert.NotNil(t, client)
	assert.Empty(t, client.configFile)

	client = NewClient("/path/to/config")
	assert.NotNil(t, client)
	assert.Equal(t, "/path/to/config", client.configFile)
}

func TestStatus(t *testing.T) {
	tests := []struct {
		name          string
		opts          StatusOptions
		output        string
		expectedError error
		expectedLen   int
	}{
		{
			name: "successful status all",
			opts: StatusOptions{},
			output: `program1 RUNNING pid 123, uptime 1:23:45
program2 STOPPED not started
program3 STARTING start in progress`,
			expectedError: nil,
			expectedLen:   3,
		},
		{
			name: "successful status specific",
			opts: StatusOptions{
				Names: []string{"program1", "program2"},
			},
			output: `program1 RUNNING pid 123, uptime 1:23:45
program2 STOPPED not started`,
			expectedError: nil,
			expectedLen:   2,
		},
		{
			name:          "command error",
			opts:          StatusOptions{},
			output:        "",
			expectedError: errors.New("command failed"),
			expectedLen:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExecutor := new(MockCommandExecutor)
			mockCmd := &MockCmd{
				MockOutput: []byte(tt.output),
				Err:        tt.expectedError,
			}

			args := []string{"status"}
			if len(tt.opts.Names) > 0 {
				args = append(args, tt.opts.Names...)
			}

			mockExecutor.On("Command", "supervisorctl", args).Return(mockCmd)

			client := &Client{
				configFile:      "",
				commandExecutor: mockExecutor,
			}

			programs, err := client.Status(tt.opts)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, programs)
			} else {
				assert.NoError(t, err)
				assert.Len(t, programs, tt.expectedLen)
			}

			mockExecutor.AssertExpectations(t)
		})
	}
}

func TestStart(t *testing.T) {
	tests := []struct {
		name          string
		programName   string
		expectedError error
	}{
		{
			name:          "successful start",
			programName:   "program1",
			expectedError: nil,
		},
		{
			name:          "command error",
			programName:   "program1",
			expectedError: errors.New("command failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExecutor := new(MockCommandExecutor)
			mockCmd := &MockCmd{
				Err: tt.expectedError,
			}

			mockExecutor.On("Command", "supervisorctl", []string{"start", tt.programName}).Return(mockCmd)

			client := &Client{
				configFile:      "",
				commandExecutor: mockExecutor,
			}

			err := client.Start(tt.programName)

			if tt.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockExecutor.AssertExpectations(t)
		})
	}
}

func TestStop(t *testing.T) {
	tests := []struct {
		name          string
		programName   string
		expectedError error
	}{
		{
			name:          "successful stop",
			programName:   "program1",
			expectedError: nil,
		},
		{
			name:          "command error",
			programName:   "program1",
			expectedError: errors.New("command failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExecutor := new(MockCommandExecutor)
			mockCmd := &MockCmd{
				Err: tt.expectedError,
			}

			mockExecutor.On("Command", "supervisorctl", []string{"stop", tt.programName}).Return(mockCmd)

			client := &Client{
				configFile:      "",
				commandExecutor: mockExecutor,
			}

			err := client.Stop(tt.programName)

			if tt.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockExecutor.AssertExpectations(t)
		})
	}
}

func TestRestart(t *testing.T) {
	tests := []struct {
		name          string
		programName   string
		expectedError error
	}{
		{
			name:          "successful restart",
			programName:   "program1",
			expectedError: nil,
		},
		{
			name:          "command error",
			programName:   "program1",
			expectedError: errors.New("command failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExecutor := new(MockCommandExecutor)
			mockCmd := &MockCmd{
				Err: tt.expectedError,
			}

			mockExecutor.On("Command", "supervisorctl", []string{"restart", tt.programName}).Return(mockCmd)

			client := &Client{
				configFile:      "",
				commandExecutor: mockExecutor,
			}

			err := client.Restart(tt.programName)

			if tt.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockExecutor.AssertExpectations(t)
		})
	}
}

func TestParseStatusLine(t *testing.T) {
	tests := []struct {
		name          string
		line          string
		expectedInfo  ProgramInfo
		expectedError error
	}{
		{
			name: "running program",
			line: "program1 RUNNING pid 123, uptime 1:23:45",
			expectedInfo: ProgramInfo{
				Name:        "program1",
				State:       20,
				StateName:   "RUNNING",
				Description: "pid 123, uptime 1:23:45",
				PID:         123,
				Uptime:      "1:23:45",
			},
			expectedError: nil,
		},
		{
			name: "stopped program",
			line: "program2 STOPPED not started",
			expectedInfo: ProgramInfo{
				Name:        "program2",
				State:       0,
				StateName:   "STOPPED",
				Description: "not started",
			},
			expectedError: nil,
		},
		{
			name: "no such process",
			line: "program3 ERROR (no such process)",
			expectedInfo: ProgramInfo{
				Name:        "program3",
				State:       1000,
				StateName:   "UNKNOWN",
				Description: "no such process",
			},
			expectedError: nil,
		},
		{
			name:          "invalid line",
			line:          "invalid",
			expectedInfo:  ProgramInfo{},
			expectedError: errors.New("invalid status line: invalid"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := parseStatusLine(tt.line)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedInfo, info)
			}
		})
	}
}
