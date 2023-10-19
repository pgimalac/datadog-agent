// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package e2e

import (
	"context"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/runner"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/infra"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type InfraDefinition[Env any] interface {
	GetInfraNoDeleteOnFailure(ctx context.Context, name string, failOnMissing bool) (*Env, map[string]interface{}, error)
	Delete(ctx context.Context, name string) error
}

var _ InfraDefinition[int] = (*PulumiStackDefinition[int])(nil)

type PulumiStackDefinition[Env any] struct {
	envFactory   func(ctx *pulumi.Context) (*Env, error)
	stackManager *infra.StackManager
	configMap    runner.ConfigMap
}

func (ps *PulumiStackDefinition[Env]) GetInfraNoDeleteOnFailure(ctx context.Context, name string, failOnMissing bool) (*Env, map[string]interface{}, error) {
	var env *Env

	deployFunc := func(ctx *pulumi.Context) error {
		var err error
		env, err = ps.envFactory(ctx)
		return err
	}
	_, _, err := ps.stackManager.GetStackNoDeleteOnFailure(ctx, name, ps.configMap, deployFunc, failOnMissing)

	result := make(map[string]interface{})
	/*
		type UpResult struct {
			StdOut  string
			StdErr  string
			Outputs OutputMap
			Summary UpdateSummary
		}
	*/
	return env, result, err
}

func (ps *PulumiStackDefinition[Env]) Delete(ctx context.Context, name string) error {
	return ps.stackManager.DeleteStack(ctx, name)
}

var _ InfraDefinition[int] = (*LocalInfraDefinition[int])(nil)

type LocalInfraDefinition[Env any] struct {
	envFactory func(jsonPath string) (*Env, error)
	configMap  runner.ConfigMap
}

func (li *LocalInfraDefinition[Env]) GetInfraNoDeleteOnFailure(ctx context.Context, name string, failOnMissing bool) (*Env, map[string]interface{}, error) {

	// TODO

	return nil, map[string]interface{}{}, nil
}

func (li *LocalInfraDefinition[Env]) Delete(ctx context.Context, name string) error {

	// TODO

	return nil
}
