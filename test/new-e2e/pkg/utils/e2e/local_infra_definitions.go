// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package e2e

import (
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/runner"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/e2e/client"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/infra"
)

// NewLocalInfraDef creates a custom local infra definition
func NewLocalInfraDef[Env any](envFactory func() (*Env, error), configMap runner.ConfigMap) InfraDefinition[Env] {
	return &LocalInfraDefinition[Env]{envFactory: envFactory, infraProvider: &infra.LocalVMManager{}, configMap: configMap}
}

func LocalVMInfraDef(jsonPath string) InfraDefinition[VMEnv] {
	cm := runner.ConfigMap{}
	cm.Set("jsonPath", jsonPath, false)

	return NewLocalInfraDef(
		func() (*VMEnv, error) {
			return &VMEnv{
				VM: client.NewLocalVM(infra.LocalVM{}),
			}, nil
		}, cm)
}
