// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package process

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/e2e/client"
)

func isProcessRunningUnix(client client.VM, process string) (bool, error) {
	_, err := client.ExecuteWithError(fmt.Sprintf("pgrep -f %s", process))
	if err != nil {
		return false, err
	}
	return true, nil
}

func findPIDUnix(client client.VM, process string) ([]int, error) {
	out, err := client.ExecuteWithError(fmt.Sprintf("pgrep -f '%s'", process))
	if err != nil {
		return nil, err
	}

	pids := []int{}
	for _, strPid := range strings.Split(out, "\n") {
		strPid = strings.TrimSpace(strPid)
		if strPid == "" {
			continue
		}
		pid, err := strconv.Atoi(strPid)
		if err != nil {
			return nil, err
		}
		pids = append(pids, pid)
	}

	return pids, nil
}
