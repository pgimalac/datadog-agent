// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package porttester

import (
	"encoding/json"

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

// BoundPorts returns a map of ports to the process name they are bound to
func (s *Windows) BoundPorts() ([]BoundPort, error) {
	out, err := s.vmClient.ExecuteWithError(`Get-NetTCPConnection -State Listen | Foreach-Object {
		@{
			LocalAddress=$_.LocalAddress
			LocalPort = $_.LocalPort
			Process = (Get-Process -Id $_.OwningProcess).Name
			PID = $_.OwningProcess
		}} | ConvertTo-JSON`)
	if err != nil {
		return nil, err
	}

	// unmarshal out as JSON
	var ports []map[string]any
	err = json.Unmarshal([]byte(out), &ports)
	if err != nil {
		return nil, err
	}

	// process JSON to BoundPort
	boundPorts := make([]BoundPort, 0, len(ports))
	for _, port := range ports {
		boundPorts = append(boundPorts, &boundPort{
			localAddress: port["LocalAddress"].(string),
			localPort:    int(port["LocalPort"].(float64)),
			processName:  port["Process"].(string),
			pid:          int(port["PID"].(float64)),
		})
	}

	return boundPorts, nil
}
