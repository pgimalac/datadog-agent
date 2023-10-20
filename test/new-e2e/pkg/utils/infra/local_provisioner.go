// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package infra implements utilities to interact with testing infrastructure
package infra

import (
	"fmt"
	"github.com/DataDog/test-infra-definitions/common/utils"
	commonvm "github.com/DataDog/test-infra-definitions/components/vm"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

var _ utils.RemoteServiceDeserializer[commonvm.ClientData] = (*LocalVM)(nil)

type LocalVM struct{}

func (l LocalVM) Deserialize(connResult auto.UpResult) (*commonvm.ClientData, error) {
	result, ok := connResult.Outputs["vm-connection"]
	if !ok {
		return nil, fmt.Errorf("connection result did not contain vm connection information")
	}

	connInfo, ok := result.Value.(map[string]string)
	if !ok {
		return nil, fmt.Errorf("vm connection information was malformed")
	}

	host, ok := connInfo["host"]
	if !ok {
		return nil, fmt.Errorf("vm connection information was malformed: missing host field")
	}

	user, ok := connInfo["user"]
	if !ok {
		return nil, fmt.Errorf("vm connection information was malformed: missing user field")
	}

	return &commonvm.ClientData{Connection: utils.Connection{Host: host, User: user}}, nil
}

type LocalVMManager struct {
	// TODO
}

func (m *LocalVMManager) ProvisionLocalVM(jsonPath string) (auto.UpResult, error) {

	// TODO: provision local VM & return SSH information as an auto.UpResult

	return auto.UpResult{}, nil
}

func (m *LocalVMManager) Delete() error {

	// TODO: tear down local VMs

	return nil
}
