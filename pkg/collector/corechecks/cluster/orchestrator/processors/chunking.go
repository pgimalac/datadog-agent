// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build orchestrator

//nolint:revive // TODO(CAPP) Fix revive linter
package processors

import model "github.com/DataDog/agent-payload/v5/process"

func weightForOrchestratorPayload(payloads []interface{}, i int) int {
	if i >= len(payloads) {
		return 0
	}
	return len(payloads[i].(*model.Manifest).Content)
}
