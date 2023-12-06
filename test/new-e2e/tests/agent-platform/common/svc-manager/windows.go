// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package svcmanager

import (
	"fmt"
	"strings"

	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/e2e/client"
	"github.com/DataDog/datadog-agent/test/new-e2e/tests/windows"
)

// WindowsSvcManager struct for Windows service manager (SCM)
type WindowsSvcManager struct {
	vmClient client.VM
}

// NewServiceSvcManager return service service manager
func NewWindowsSvcManager(vmClient client.VM) *WindowsSvcManager {
	return &WindowsSvcManager{vmClient}
}

// Status returns status from service
func (s *WindowsSvcManager) Status(service string) (string, error) {
	status, err := windows.GetServiceStatus(s.vmClient, service)
	if err != nil {
		return status, err
	}

	// TODO: The other service managers return an error if the service is not running,
	//       is that really what we want?
	if !strings.EqualFold(status, "Running") {
		return status, fmt.Errorf("service %s is not running", service)
	}

	return status, nil
}

// Stop executes stop command from service
func (s *WindowsSvcManager) Stop(service string) (string, error) {
	return "", windows.StopService(s.vmClient, service)
}

// Start executes start command from service
func (s *WindowsSvcManager) Start(service string) (string, error) {
	return "", windows.StartService(s.vmClient, service)
}

// Restart executes restart command from service
func (s *WindowsSvcManager) Restart(service string) (string, error) {
	return "", windows.RestartService(s.vmClient, service)
}
