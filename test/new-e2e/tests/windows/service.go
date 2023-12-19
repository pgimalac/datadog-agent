// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-present Datadog, Inc.

package windows

import (
	"fmt"
	"strings"

	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/e2e/client"
)

// GetServiceStatus returns the status of the service
func GetServiceStatus(vm client.VM, service string) (string, error) {
	cmd := fmt.Sprintf("(Get-Service -Name '%s').Status", service)
	out, err := vm.ExecuteWithError(cmd)
	return strings.TrimSpace(out), err
}

// StopService stops the service
func StopService(vm client.VM, service string) error {
	cmd := fmt.Sprintf("Stop-Service -Force -Name '%s'", service)
	_, err := vm.ExecuteWithError(cmd)
	return err
}

// StartService starts the service
func StartService(vm client.VM, service string) error {
	cmd := fmt.Sprintf("Start-Service -Name '%s'", service)
	_, err := vm.ExecuteWithError(cmd)
	return err
}

// RestartService restarts the service
func RestartService(vm client.VM, service string) error {
	cmd := fmt.Sprintf("Restart-Service -Force -Name '%s'", service)
	_, err := vm.ExecuteWithError(cmd)
	return err
}
