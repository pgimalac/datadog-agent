// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package helper implement interfaces to get some information that can be OS specific
package helper

import (
	"github.com/DataDog/datadog-agent/test/new-e2e/tests/agent-platform/common"
)

// Unix implement helper function for Unix distributions
type Windows struct{}

var _ common.Helper = (*Windows)(nil)

// NewWindowsHelper create a new instance of Windows helper
func NewWindowsHelper() *Windows { return &Windows{} }

// GetInstallFolder return the install folder path
func (u *Windows) GetInstallFolder() string { return `C:\Program Files\Datadog\Datadog Agent` }

// GetConfigFolder return the config folder path
func (u *Windows) GetConfigFolder() string { return `C:\ProgramData\Datadog` }

// GetBinaryPath return the datadog-agent binary path
func (u *Windows) GetBinaryPath() string { return u.GetInstallFolder() + `\` + `bin\agent.exe` }

// GetConfigFileName return the config file name
func (u *Windows) GetConfigFileName() string { return "datadog.yaml" }

// GetServiceName return the service name
func (u *Windows) GetServiceName() string { return "datadogagent" }
