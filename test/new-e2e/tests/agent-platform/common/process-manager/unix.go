// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package processmanager

import (
	"fmt"

	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/e2e/client"
)

// Unix struct for Unix process manager
type Unix struct {
	vmClient client.VM
}

// NewProcessManager return process manager
func NewUnixProcessManager(vmClient client.VM) *Unix {
	return &Unix{vmClient}
}

// IsProcessRunning returns true if process is running
func (u *Unix) IsProcessRunning(process string) (bool, error) {
	_, err := u.vmClient.ExecuteWithError(fmt.Sprintf("pgrep -f %s", process))
	if err != nil {
		return false, err
	}
	return true, nil
}
