package infra

import (
	"context"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/runner"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type StackManagerInterface interface {
	GetStackNoDeleteOnFailure(ctx context.Context, name string, config runner.ConfigMap, deployFunc pulumi.RunFunc, failOnMissing bool) (*auto.Stack, auto.UpResult, error)
	DeleteStack(ctx context.Context, name string) error
}

func GetStackManagerInterface(fakeStack bool) StackManagerInterface {
	if fakeStack {
		return &FakeStackManager{}
	}
	return GetStackManager()
}

type FakeStackManager struct{}

func (fsm *FakeStackManager) GetStackNoDeleteOnFailure(_ context.Context, _ string, _ runner.ConfigMap, deployFunc pulumi.RunFunc, _ bool) (*auto.Stack, auto.UpResult, error) {
	// Fake stack doesn't need to do anything on creation except call deployFunc
	err := deployFunc(nil)
	return nil, auto.UpResult{}, err
}

func (fsm *FakeStackManager) DeleteStack(ctx context.Context, name string) error {
	return nil
}
