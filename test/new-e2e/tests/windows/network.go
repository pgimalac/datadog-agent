// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-present Datadog, Inc.

package windows

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/e2e/client"
)

type BoundPort struct {
	localAddress string
	localPort    int
	processName  string
	pid          int
}

func (b *BoundPort) LocalAddress() string {
	return b.localAddress
}

func (b *BoundPort) LocalPort() int {
	return b.localPort
}

func (b *BoundPort) Process() string {
	return b.processName
}

func (b *BoundPort) PID() int {
	return b.pid
}

// IsPortBound returns true if the port is bound
func IsPortBound(vmClient client.VM, port int) (bool, error) {
	out, err := vmClient.ExecuteWithError(fmt.Sprintf("(Get-NetTCPConnection -LocalPort %d -State Listen -ErrorAction SilentlyContinue) -ne $null", port))
	if err != nil {
		return false, err
	}
	return !strings.EqualFold(strings.TrimSpace(out), "False"), nil
}

func ListBoundPorts(client client.VM) ([]*BoundPort, error) {
	out, err := client.ExecuteWithError(`Get-NetTCPConnection -State Listen | Foreach-Object {
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
	boundPorts := make([]*BoundPort, 0, len(ports))
	for _, port := range ports {
		boundPorts = append(boundPorts, &BoundPort{
			localAddress: port["LocalAddress"].(string),
			localPort:    int(port["LocalPort"].(float64)),
			processName:  port["Process"].(string),
			pid:          int(port["PID"].(float64)),
		})
	}

	return boundPorts, nil
}
