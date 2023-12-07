// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package processmanager implement interfaces to manage processes
package processmanager

import (
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/e2e/client"
	"github.com/DataDog/datadog-agent/test/new-e2e/tests/windows"
)

// Windows struct for Windows process manager
type Windows struct {
	vmClient client.VM
}

// NewServiceSvcManager return service service manager
func NewWindowsProcessManager(vmClient client.VM) *Windows {
	return &Windows{vmClient}
}

// IsProcessRunning returns true if process is running
func (s *Windows) IsProcessRunning(process string) (bool, error) {
	return windows.IsProcessRunning(s.vmClient, process)
}
