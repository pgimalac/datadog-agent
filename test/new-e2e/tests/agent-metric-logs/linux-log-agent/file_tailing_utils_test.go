// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package logagent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/DataDog/datadog-agent/test/fakeintake/aggregator"
	fi "github.com/DataDog/datadog-agent/test/fakeintake/client"
	componentos "github.com/DataDog/test-infra-definitions/components/os"
	"github.com/stretchr/testify/assert"
)

// appendLog appen log with 'content', which is then repeated 'reccurrence' times and verifies log contents.
func appendLog(s *LinuxFakeintakeSuite, content string, recurrence int) {
	// Determine the OS and set the appropriate log path and command.
	var logPath, cmd, checkCmd string
	t := s.T()
	t.Helper()

	var os string
	switch s.Env().RemoteHost.OSFamily {
	case componentos.WindowsFamily:
		os = "windows"
		t.Log("Generating Windows log.")
		logPath = "C:\\logs\\hello-world.log"
		cmd = fmt.Sprintf("echo %s > %s", strings.Repeat(content+" ", recurrence), logPath)
		checkCmd = fmt.Sprintf("Get-Content %s", logPath)
	default: // Assuming Linux if not Windows.
		os = "linux"
		t.Log("Generating Linux log.")
		logPath = "/var/log/hello-world.log"
		cmd = fmt.Sprintf("echo '%s' | sudo tee -a %s", strings.Repeat(content+" ", recurrence), logPath)
		checkCmd = fmt.Sprintf("sudo cat %s", logPath)
	}

	s.EventuallyWithT(func(c *assert.CollectT) {
		// Generate the log content
		output, err := s.Env().RemoteHost.Execute(cmd)
		if err != nil {
			assert.FailNowf(c, "Having issue generating %s log with error: %s", os, output)
		}
		// Verify the log content locally
		output, err = s.Env().RemoteHost.Execute(checkCmd)
		if err != nil {
			assert.FailNowf(c, "Log content %s not found, instead received:: %s", content, output)
		}
		if strings.Contains(output, content) {
			t.Logf("Finished generating %s log with content: '%s' \n", os, content)
		}
	}, 2*time.Minute, 10*time.Second)
}

// checkLogFile verifies the presence or absence of a log file path
func checkLogFilePresence(s *LinuxFakeintakeSuite, logPath string) {
	t := s.T()

	switch s.Env().RemoteHost.OSFamily {
	case componentos.WindowsFamily:
		checkCmd := fmt.Sprintf("Get-Content %s", logPath)
		_, err := s.Env().RemoteHost.Execute(checkCmd)
		if err != nil {
			assert.FailNow(t, "Log File not found")
		}
	default: // Assuming Linux if not Windows.
		checkCmd := fmt.Sprintf("sudo cat %s", logPath)
		_, err := s.Env().RemoteHost.Execute(checkCmd)
		if err != nil {
			assert.FailNow(t, "Log File not found")
		}
	}
}

// checkLogs verifies the presence or absence of logs in the intake based on the expectLogs flag.
func checkLogs(s *LinuxFakeintakeSuite, service, content string, expectLogs bool) {
	client := s.Env().FakeIntake.Client()
	t := s.T()
	t.Helper()
	s.EventuallyWithT(func(c *assert.CollectT) {
		names, err := client.GetLogServiceNames()
		if !assert.NoErrorf(c, err, "Error found: %s", err) {
			return
		}

		if len(names) > 0 {
			logs, err := client.FilterLogs(service)
			if !assert.NoErrorf(c, err, "Error found: %s", err) {
				return
			}
			if !assert.NotEmpty(c, logs, "No logs with service matching '%s' found, instead got '%s'", service, names) {
				return
			}

			logs, err = client.FilterLogs(service, fi.WithMessageMatching(content))
			intakeLogs := logsToString(logs)
			assert.NoErrorf(c, err, "Error found: %s", err)

			if expectLogs {
				t.Logf("Logs from service: '%s' with content: '%s' collected", names, content)
				assert.NotEmpty(c, logs, "Expected at least 1 log with content: '%s', from service: %s but received %s logs.", content, names, intakeLogs)
			} else {
				t.Logf("No logs from service: '%s' with content: '%s' collected as expected", names, content)
				assert.Empty(c, logs, "No logs with content: '%s' is expected to be found from service: %s instead found: %s", content, names, intakeLogs)
			}
		}
	}, 2*time.Minute, 10*time.Second)
}

// cleanUp cleans up any existing log files (only useful when running dev mode/local runs).
func (s *LinuxFakeintakeSuite) cleanUp() {
	t := s.T()

	s.EventuallyWithT(func(c *assert.CollectT) {
		err := s.Env().FakeIntake.Client().FlushServerAndResetAggregators()
		if assert.NoErrorf(c, err, "Having issue flushing server and resetting aggregators, retrying...") {
			t.Log("Successfully flushed server and reset aggregators.")
		}
	}, 1*time.Minute, 10*time.Second)
}

// prettyPrintLog pretty prints a log entry.
func prettyPrintLog(log *aggregator.Log) string {
	// Unmarshal and re-marshal the message field for pretty printing
	var messageObj map[string]interface{}
	if err := json.Unmarshal([]byte(log.Message), &messageObj); err == nil {
		prettyMessage, _ := json.MarshalIndent(messageObj, "", "  ")
		log.Message = string(prettyMessage)
	}

	// Marshal the entire log entry
	logStr, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		// Handle the error appropriately
		return fmt.Sprintf("Error marshaling log: %v", err)
	}
	return string(logStr)
}

// logsToString converts a slice of logs to a string.
func logsToString(logs []*aggregator.Log) string {
	var logsStrings []string
	for _, log := range logs {
		logsStrings = append(logsStrings, prettyPrintLog(log))
	}
	return fmt.Sprintf("[%s]", strings.Join(logsStrings, ",\n"))
}
