// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package boundport

import (
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/e2e/client"
	"github.com/DataDog/datadog-agent/test/new-e2e/tests/windows"
)

func boundPortsWindows(client client.VM) ([]BoundPort, error) {
	ports, err := windows.ListBoundPorts(client)
	if err != nil {
		return nil, err
	}
	// convert to BoundPort interface
	boundPorts := make([]BoundPort, 0, len(ports))
	for _, port := range ports {
		boundPorts = append(boundPorts, port)
	}
	return boundPorts, nil
}
