// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package inventoryhost exposes the interface for the component to generate the 'host_metadata' metadata payload for inventory.
package inventoryhost

// team: agent-shared-components

// Component is the component type.
type Component interface {
	// GetAsJSON returns the payload as a JSON string. Useful to be displayed in the CLI or added to a flare.
	GetAsJSON() ([]byte, error)
	// Refresh trigger a new payload to be send while still respecting the minimal interval between two updates.
	Refresh()
}
