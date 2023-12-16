// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package boundport

import (
	"encoding/json"

	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/e2e/client"
)

func boundPortsWindows(client client.VM) ([]BoundPort, error) {
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
