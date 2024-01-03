// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build !linux || !cgo

//nolint:revive // TODO(CINT) Fix revive linter
package ebpf

import "github.com/DataDog/datadog-agent/pkg/collector/check"

const (
	Enabled   = false
	CheckName = "ebpf"
)

func Factory() check.Check {
	return nil
}
