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

// BoundPort represents a port that is bound to a process
type BoundPort struct {
	localAddress string
	localPort    int
	processName  string
	pid          int
}

// LocalAddress returns the local address of the bound port
func (b *BoundPort) LocalAddress() string {
	return b.localAddress
}

// LocalPort returns the local port of the bound port
func (b *BoundPort) LocalPort() int {
	return b.localPort
}

// Process returns the process name of the bound port
func (b *BoundPort) Process() string {
	return b.processName
}

// PID returns the PID of the bound port
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

// ListBoundPorts returns a list of bound ports
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

// PutOrDownloadFile creates a file on the VM from a file/http URL
//
// If the URL is a local file, it will be uploaded to the VM.
// If the URL is a remote file, it will be downloaded from the VM
func PutOrDownloadFile(client client.VM, url string, destination string) error {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		_, err := client.ExecuteWithError(fmt.Sprintf("Invoke-WebRequest -Uri '%s' -OutFile '%s'", url, destination))
		if err != nil {
			return err
		}

		return nil
	}

	if strings.HasPrefix(url, "file://") {
		// URL is a local file
		localPath := strings.TrimPrefix(url, "file://")
		client.CopyFile(localPath, destination)
		return nil
	}

	// just assume it's a local file
	client.CopyFile(url, destination)
	return nil
}
