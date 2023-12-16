// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package boundport

import (
	"fmt"

	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/e2e/client"
)

func isPortBoundUnix(client client.VM, port int) (bool, error) {
	netstatCmd := "sudo netstat -lntp | grep %v"
	if _, err := client.ExecuteWithError("command -v netstat"); err != nil {
		netstatCmd = "sudo ss -lntp | grep %v"
	}

	_, err := client.ExecuteWithError(fmt.Sprintf(netstatCmd, port))
	// TODO: distinguish grep not matching vs some other error
	if err != nil {
		return false, nil
	}

	return true, nil
}

func boundPortsUnix(client client.VM) ([]BoundPort, error) {
	if _, err := client.ExecuteWithError("command -v netstat"); err == nil {
		out, err := client.ExecuteWithError("sudo netstat -lntp")
		if err != nil {
			return nil, err
		}
		return FromNetstat(out)
	}

	if _, err := client.ExecuteWithError("command -v ss"); err == nil {
		out, err := client.ExecuteWithError("sudo ss -lntp")
		if err != nil {
			return nil, err
		}
		return FromSs(out)
	}

	return nil, fmt.Errorf("no netstat or ss found")
}
