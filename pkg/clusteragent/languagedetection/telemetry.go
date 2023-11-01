// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-present Datadog, Inc.

package languagedetection

import "github.com/DataDog/datadog-agent/pkg/telemetry"

const subsystem = "language_detection_dca_handler"

var (
	commonOpts = telemetry.Options{NoDoubleUnderscoreSep: true}
)

var (
	// PatchRetries determines the number of times a patch request fails and is retried for a gived
	PatchRetries = telemetry.NewCounterWithOpts(
		subsystem,
		"retries",
		[]string{"owner_kind", "owner_name", "namespace"},
		"Tracks the number of retries while patching deployments with language annotations",
		commonOpts,
	)

	// SuccessPatches tracks the number of successful annotation patch operations
	SuccessPatches = telemetry.NewCounterWithOpts(
		subsystem,
		"success_patch",
		[]string{"owner_kind", "owner_name", "namespace"},
		"Tracks the number of successful annotation patch operations",
		commonOpts,
	)

	// FailedPatches tracks the number of failing annotation patch operations
	FailedPatches = telemetry.NewCounterWithOpts(
		subsystem,
		"fail_patch",
		[]string{"owner_kind", "owner_name", "namespace"},
		"Tracks the number of failing annotation patch operations",
		commonOpts,
	)

	// SkippedPatches tracks the number of times a patch was skipped because no new languages are detected
	SkippedPatches = telemetry.NewSimpleCounter(
		subsystem,
		"skipped_patch",
		"Tracks the number of times a patch was skipped because no new languages are detected",
	)

	// OkResponses tracks the number the request was processed successfully
	OkResponses = telemetry.NewSimpleCounter(
		subsystem,
		"ok_response",
		"Tracks the number the request was processed successfully",
	)

	// ErrorResponses tracks the number of times request processsing fails
	ErrorResponses = telemetry.NewSimpleCounter(
		subsystem,
		"fail_response",
		"Tracks the number of times request processsing fails",
	)
)
