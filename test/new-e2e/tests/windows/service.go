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

func GetServiceAccountName(vm client.VM, service string) (string, error) {
	cmd := fmt.Sprintf("(Get-WmiObject Win32_Service -Filter \"Name=`'%s`'\").StartName", service)
	return vm.ExecuteWithError(cmd)
}

func GetServiceInfo(vm client.VM, service string) (map[string]any, error) {
	cmd := fmt.Sprintf("Get-Service -Name '%s' | ConvertTo-Json", service)
	output, err := vm.ExecuteWithError(cmd)
	if err != nil {
		fmt.Println(output)
		return nil, err
	}

	var result map[string]any
	err = json.Unmarshal([]byte(output), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetServiceSDDL(vm client.VM, service string) (string, error) {
	cmd := fmt.Sprintf("sc.exe sdshow '%s'", service)
	out, err := vm.ExecuteWithError(cmd)
	return strings.TrimSpace(out), err
}

func GetServiceStatus(vm client.VM, service string) (string, error) {
	cmd := fmt.Sprintf("(Get-Service -Name '%s').Status", service)
	out, err := vm.ExecuteWithError(cmd)
	return strings.TrimSpace(out), err
}

func StopService(vm client.VM, service string) error {
	cmd := fmt.Sprintf("Stop-Service -Force -Name '%s'", service)
	_, err := vm.ExecuteWithError(cmd)
	return err
}

func StartService(vm client.VM, service string) error {
	cmd := fmt.Sprintf("Start-Service -Name '%s'", service)
	_, err := vm.ExecuteWithError(cmd)
	return err
}

func RestartService(vm client.VM, service string) error {
	cmd := fmt.Sprintf("Restart-Service -Force -Name '%s'", service)
	_, err := vm.ExecuteWithError(cmd)
	return err
}

// GetServicePID returns the PID running the service
func GetServicePID(vm client.VM, service string) (int, error) {
	info, err := GetServiceInfo(vm, service)
	if err != nil {
		return 0, err
	}
	pid := info["ProcessId"].(float64)
	return int(pid), nil
}
