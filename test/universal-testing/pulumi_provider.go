// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package universal_testing

import (
	"context"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/runner"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/e2e/client"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/infra"
)

var _ InfraProvider[int] = (*PulumiProvider[int])(nil)

// PulumiProvider uses a pulumi stack manager to initialize a pulumi stack & pass the resulting UpResult to
// any clients in the environment which implement the pulumiStackInitializer interface
type PulumiProvider[Env any] struct {
	envFactory func(ctx *pulumi.Context) (*Env, error)
	configMap  runner.ConfigMap

	stackManager *infra.StackManager
}

func NewPulumiProvider[Env any](envFactory func(ctx *pulumi.Context) (*Env, error)) *PulumiProvider[Env] {
	return &PulumiProvider[Env]{
		envFactory:   envFactory,
		configMap:    runner.ConfigMap{},
		stackManager: infra.GetStackManager(),
	}
}

func (ps *PulumiProvider[Env]) ProvisionInfraAndInitializeEnv(t *testing.T, ctx context.Context, name string, failOnMissing bool) (*Env, error) {
	var env *Env

	deployFunc := func(ctx *pulumi.Context) error {
		var err error
		env, err = ps.envFactory(ctx)
		return err
	}
	_, stackResult, err := ps.stackManager.GetStackNoDeleteOnFailure(ctx, name, ps.configMap, deployFunc, failOnMissing)

	err = client.CallStackInitializers(t, env, stackResult)
	return env, err
}

func (ps *PulumiProvider[Env]) DeleteInfra(ctx context.Context, name string) error {
	return ps.stackManager.DeleteStack(ctx, name)
}
