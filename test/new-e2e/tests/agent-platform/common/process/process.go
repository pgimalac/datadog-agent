// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package process provides utilities for testing processes
package process

import (
	"fmt"

	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/e2e/client"
	"github.com/DataDog/datadog-agent/test/new-e2e/tests/windows"
	componentos "github.com/DataDog/test-infra-definitions/components/os"
)

// IsProcessRunning returns true if process is running
func IsProcessRunning(client client.VM, process string) (bool, error) {
	os := client.GetOSType()
	if os == componentos.UnixType {
		return isProcessRunningUnix(client, process)
	} else if os == componentos.WindowsType {
		return windows.IsProcessRunning(client, process)
	}
	return false, fmt.Errorf("unsupported OS type: %v", os)
}

// FindPID returns list of PIDs that match process
func FindPID(client client.VM, process string) ([]int, error) {
	os := client.GetOSType()
	if os == componentos.UnixType {
		return findPIDUnix(client, process)
	} else if os == componentos.WindowsType {
		return windows.FindPID(client, process)
	}
	return nil, fmt.Errorf("unsupported OS type: %v", os)
}
