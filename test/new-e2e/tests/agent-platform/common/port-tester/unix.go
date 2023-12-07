// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package svcmanager

import (
	"fmt"

	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/e2e/client"
)

// Unix struct for Unix port testing
type Unix struct {
	vmClient client.VM
}

// NewUnixPortTester return port tester
func NewUnixPortTester(vmClient client.VM) *Unix {
	return &Unix{vmClient}
}

// IsPortBound returns true if the port is bound
func (u *Unix) IsPortBound(port int) (bool, error) {
	netstatCmd := "sudo netstat -lntp | grep %v"
	if _, err := u.vmClient.ExecuteWithError("command -v netstat"); err != nil {
		netstatCmd = "sudo ss -lntp | grep %v"
	}

	_, err := u.vmClient.ExecuteWithError(fmt.Sprintf(netstatCmd, port))
	// TODO: distinguish grep not matching vs some other error
	if err != nil {
		return false, nil
	}

	return true, nil
}
