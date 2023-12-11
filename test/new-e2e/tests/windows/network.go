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

// IsPortBound returns true if the port is bound
func IsPortBound(vmClient client.VM, port int) (bool, error) {
	out, err := vmClient.ExecuteWithError(fmt.Sprintf("(Get-NetTCPConnection -LocalPort %d -State Listen -ErrorAction SilentlyContinue) -ne $null", port))
	if err != nil {
		return false, err
	}
	return !strings.EqualFold(strings.TrimSpace(out), "False"), nil
}
