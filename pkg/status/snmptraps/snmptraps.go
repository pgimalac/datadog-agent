// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package snmptraps fetch information needed to render the 'snmptraps' section of the status page.
package snmptraps

import (
	"github.com/DataDog/datadog-agent/pkg/config"
	"github.com/DataDog/datadog-agent/pkg/snmp/traps"
)

// PopulateStatus populates the status stats
func PopulateStatus(stats map[string]interface{}) {
	if traps.IsEnabled(config.Datadog) {
		stats["snmpTrapsStats"] = traps.GetStatus()
	}
}
