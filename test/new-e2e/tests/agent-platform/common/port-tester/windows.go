// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package svcmanager

import (
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/e2e/client"
	"github.com/DataDog/datadog-agent/test/new-e2e/tests/windows"
)

// Windows struct for Windows port testing
type Windows struct {
	vmClient client.VM
}

// NewWindowsPortTester return port tester
func NewWindowsPortTester(vmClient client.VM) *Windows {
	return &Windows{vmClient}
}

// IsPortBound returns true if the port is bound
func (s *Windows) IsPortBound(port int) (bool, error) {
	return windows.IsPortBound(s.vmClient, port)
}
