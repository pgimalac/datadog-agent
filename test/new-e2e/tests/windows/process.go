// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-present Datadog, Inc.

package windows

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/e2e/client"
)

// ListProcess returns list of process
func ListProcess(vm client.VM) ([]string, error) {
	//cmd := "Get-CimInstance Win32_Process | Select-Object -ExpandProperty CommandLine"
	//cmd := "Get-Process | Select-Object -ExpandProperty ProcessName"
	cmd := "Get-CimInstance Win32_Process | Select-Object -ExpandProperty Name"
	output, err := vm.ExecuteWithError(cmd)
	if err != nil {
		return nil, err
	}

	return strings.Split(output, "\n"), nil
}

// IsProcessRunning returns true if process is running
func IsProcessRunning(vm client.VM, imageName string) (bool, error) {
	cmd := fmt.Sprintf(`tasklist /fi "ImageName -eq '%s'"`, imageName)
	out, err := vm.ExecuteWithError(cmd)
	if err != nil {
		return false, err
	}
	return strings.Contains(out, imageName), nil
}

// KillProcess kill process
func KillProcess(vm client.VM, pid int) error {
	cmd := fmt.Sprintf("Stop-Process -Id '%d'", pid)
	_, err := vm.ExecuteWithError(cmd)
	return err
}

// PkillProcess kill process
func PkillProcess(vm client.VM, pattern string) error {
	cmd := fmt.Sprintf("Stop-Process -Name '%s'", pattern)
	_, err := vm.ExecuteWithError(cmd)
	return err
}

func FindPID(vm client.VM, pattern string) ([]int, error) {
	cmd := fmt.Sprintf("Get-Process -Name '%s' | Select-Object -ExpandProperty Id", pattern)
	out, err := vm.ExecuteWithError(cmd)
	if err != nil {
		return nil, err
	}
	var pids []int
	for _, strPID := range strings.Split(out, "\n") {
		strPID = strings.TrimSpace(strPID)
		if strPID == "" {
			continue
		}
		pid, err := strconv.Atoi(strPID)
		if err != nil {
			return nil, err
		}
		pids = append(pids, pid)
	}
	return pids, nil
}
