// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package infra implements utilities to interact with testing infrastructure
package infra

import (
	"fmt"
	"sync"

	"github.com/DataDog/test-infra-definitions/common/utils"
)

var (
	vmManager     *LocalVMManager
	initVMManager sync.Once
)

type LocalVMManager struct {
	vms []*LocalVM
}

// GetLocalVMManager returns a vm manager, initialising on first call
func GetLocalVMManager() *LocalVMManager {
	initVMManager.Do(func() {
		vmManager = newLocalVMManager()
	})

	return vmManager
}

func newLocalVMManager() *LocalVMManager {
	return &LocalVMManager{}
}

func (m *LocalVMManager) AddVM(vm *LocalVM) {
	m.vms = append(m.vms, vm)
}

// Provision brings up all VMs which have been added to the manager
func (m *LocalVMManager) Provision() (map[string]*utils.Connection, error) {
	connections := map[string]*utils.Connection{}

	for _, vm := range m.vms {
		conn, err := m.create(vm)
		if err != nil {
			return map[string]*utils.Connection{}, err
		}

		connections[vm.name] = conn
	}

	return connections, nil
}

// Delete deletes all VMs which have been added to the manager
func (m *LocalVMManager) Delete() []error {
	errors := []error{}
	for _, vm := range m.vms {
		err := m.delete(vm)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

func (m *LocalVMManager) create(vm *LocalVM) (*utils.Connection, error) {
	if vm.isUp {
		return vm.getConnection()
	}

	// TODO
	return nil, fmt.Errorf("provisioning of new local VMs is not yet implemented")
}

func (m *LocalVMManager) delete(vm *LocalVM) error {
	// TODO
	return fmt.Errorf("deleting of local VMs is not yet implemented")
}
