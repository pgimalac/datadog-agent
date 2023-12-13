// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package processmanager

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/e2e/client"
)

// Unix struct for Unix process manager
type Unix struct {
	vmClient client.VM
}

// NewUnixProcessManager return process manager for Unix systems
func NewUnixProcessManager(vmClient client.VM) *Unix {
	return &Unix{vmClient}
}

// IsProcessRunning returns true if process is running
func (u *Unix) IsProcessRunning(process string) (bool, error) {
	_, err := u.vmClient.ExecuteWithError(fmt.Sprintf("pgrep -f %s", process))
	if err != nil {
		return false, err
	}
	return true, nil
}

// FindPID returns the PID of a process
func (u *Unix) FindPID(process string) ([]int, error) {
	out, err := u.vmClient.ExecuteWithError(fmt.Sprintf("pgrep -f '%s'", process))
	if err != nil {
		return nil, err
	}

	pids := []int{}
	for _, strPid := range strings.Split(out, "\n") {
		strPid = strings.TrimSpace(strPid)
		if strPid == "" {
			continue
		}
		pid, err := strconv.Atoi(strPid)
		if err != nil {
			return nil, err
		}
		pids = append(pids, pid)
	}

	return pids, nil
}
