// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package e2e

import (
	"context"
	"fmt"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/runner"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/infra"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type InfraDefinition[Env any] interface {
	GetInfraNoDeleteOnFailure(ctx context.Context, name string, failOnMissing bool) (*Env, auto.UpResult, error)
	Delete(ctx context.Context, name string) error
}

var _ InfraDefinition[int] = (*PulumiStackDefinition[int])(nil)

type PulumiStackDefinition[Env any] struct {
	envFactory   func(ctx *pulumi.Context) (*Env, error)
	stackManager *infra.StackManager
	configMap    runner.ConfigMap
}

func (ps *PulumiStackDefinition[Env]) GetInfraNoDeleteOnFailure(ctx context.Context, name string, failOnMissing bool) (*Env, auto.UpResult, error) {
	var env *Env

	deployFunc := func(ctx *pulumi.Context) error {
		var err error
		env, err = ps.envFactory(ctx)
		return err
	}
	_, stackResult, err := ps.stackManager.GetStackNoDeleteOnFailure(ctx, name, ps.configMap, deployFunc, failOnMissing)
	return env, stackResult, err
}

func (ps *PulumiStackDefinition[Env]) Delete(ctx context.Context, name string) error {
	return ps.stackManager.DeleteStack(ctx, name)
}

var _ InfraDefinition[int] = (*LocalInfraDefinition[int])(nil)

type LocalInfraDefinition[Env any] struct {
	envFactory    func() (*Env, error)
	infraProvider *infra.LocalVMManager
	configMap     runner.ConfigMap
}

func (li *LocalInfraDefinition[Env]) GetInfraNoDeleteOnFailure(_ context.Context, _ string, failOnMissing bool) (*Env, auto.UpResult, error) {
	jsonPath, ok := li.configMap["jsonPath"]
	if !ok {
		return nil, auto.UpResult{}, fmt.Errorf("'jsonPath' key must be provided in the config map for locally provisioned VMs")
	}

	connResult, err := li.infraProvider.ProvisionLocalVM(jsonPath.Value)
	if err != nil {
		return nil, auto.UpResult{}, fmt.Errorf("error provisioning local VM: %s", err)
	}

	env, err := li.envFactory()
	return env, connResult, err
}

func (li *LocalInfraDefinition[Env]) Delete(ctx context.Context, name string) error {
	return li.infraProvider.Delete()
}
