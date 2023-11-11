// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-present Datadog, Inc.

// Package config implements the configuration type for the traps server and
// a component that provides it.
package config

// team: network-device-monitoring

// Component is the component type.
type Component interface {
	// Enabled returns true if and only if the traps server is enabled.
	Enabled() bool
	// Get returns the traps config, or nil if the traps server is disabled.
	Get() *TrapsConfig
}
