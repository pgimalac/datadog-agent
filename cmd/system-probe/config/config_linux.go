// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build linux

package config

import (
	"fmt"
	"path/filepath"
	ebpfkernel "github.com/DataDog/datadog-agent/pkg/security/ebpf/kernel"
	"github.com/DataDog/datadog-agent/pkg/util/log"
)

const (
	// defaultSystemProbeAddress is the default unix socket path to be used for connecting to the system probe
	defaultSystemProbeAddress = "/opt/datadog-agent/run/sysprobe.sock"

	defaultConfigDir = "/etc/datadog-agent"
)

// ValidateSocketAddress validates that the sysprobe socket config option is of the correct format.
func ValidateSocketAddress(sockPath string) error {
	if !filepath.IsAbs(sockPath) {
		return fmt.Errorf("socket path must be an absolute file path: `%s`", sockPath)
	}
	return nil
}

func canEnablePES() bool {
	kernelVersion, err := ebpfkernel.NewKernelVersion()
	if err != nil {
		log.Errorf("unable to detect the kernel version: %s", err)
		return true
	}
	if !kernelVersion.IsRH7Kernel() && !kernelVersion.IsRH8Kernel() && kernelVersion.Code < ebpfkernel.Kernel4_15 {
		log.Warn("disabling process event monitoring as it is not supported for this kernel version")
		return false
	}

	return true
}
