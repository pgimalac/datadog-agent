// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package host contains the definition of the AWS Host environment.
package awshost

import (
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/e2e"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/environments"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/runner"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/optional"

	"github.com/DataDog/test-infra-definitions/components/datadog/agent"
	"github.com/DataDog/test-infra-definitions/components/datadog/agentparams"
	"github.com/DataDog/test-infra-definitions/resources/aws"
	"github.com/DataDog/test-infra-definitions/scenarios/aws/ec2"
	"github.com/DataDog/test-infra-definitions/scenarios/aws/fakeintake"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	provisionerBaseID = "aws-ec2vm-"
	defaultVMName     = "vm"
)

// ProvisionerParams is a set of parameters for the Provisioner.
type ProvisionerParams struct {
	name string

	instanceOptions   []ec2.VMOption
	agentOptions      []agentparams.Option
	fakeintakeOptions []fakeintake.Option
	extraConfigParams runner.ConfigMap
}

func newProvisionerParams() *ProvisionerParams {
	// We use nil arrays to decide if we should create or not
	return &ProvisionerParams{
		name:              defaultVMName,
		instanceOptions:   []ec2.VMOption{},
		agentOptions:      []agentparams.Option{},
		fakeintakeOptions: []fakeintake.Option{},
		extraConfigParams: runner.ConfigMap{},
	}
}

// ProvisionerOption is a provisioner option.
type ProvisionerOption func(*ProvisionerParams) error

// WithName sets the name of the provisioner.
func WithName(name string) ProvisionerOption {
	return func(params *ProvisionerParams) error {
		params.name = name
		return nil
	}
}

// WithEC2InstanceOptions adds options to the EC2 VM.
func WithEC2InstanceOptions(opts ...ec2.VMOption) ProvisionerOption {
	return func(params *ProvisionerParams) error {
		params.instanceOptions = append(params.instanceOptions, opts...)
		return nil
	}
}

// WithAgentOptions adds options to the Agent.
func WithAgentOptions(opts ...agentparams.Option) ProvisionerOption {
	return func(params *ProvisionerParams) error {
		params.agentOptions = append(params.agentOptions, opts...)
		return nil
	}
}

// WithFakeIntakeOptions adds options to the FakeIntake.
func WithFakeIntakeOptions(opts ...fakeintake.Option) ProvisionerOption {
	return func(params *ProvisionerParams) error {
		params.fakeintakeOptions = append(params.fakeintakeOptions, opts...)
		return nil
	}
}

// WithExtraConfigParams adds extra config parameters to the ConfigMap.
func WithExtraConfigParams(configMap runner.ConfigMap) ProvisionerOption {
	return func(params *ProvisionerParams) error {
		params.extraConfigParams = configMap
		return nil
	}
}

// WithoutFakeIntake disables the creation of the FakeIntake.
func WithoutFakeIntake() ProvisionerOption {
	return func(params *ProvisionerParams) error {
		params.fakeintakeOptions = nil
		return nil
	}
}

// WithoutAgent disables the creation of the Agent.
func WithoutAgent() ProvisionerOption {
	return func(params *ProvisionerParams) error {
		params.agentOptions = nil
		return nil
	}
}

// ProvisionerNoAgentNoFakeIntake wraps Provisioner with hardcoded WithoutAgent and WithoutFakeIntake options.
func ProvisionerNoAgentNoFakeIntake(opts ...ProvisionerOption) e2e.TypedProvisioner[environments.Host] {
	mergedOpts := make([]ProvisionerOption, 0, len(opts)+2)
	mergedOpts = append(mergedOpts, opts...)
	mergedOpts = append(mergedOpts, WithoutAgent(), WithoutFakeIntake())

	return Provisioner(mergedOpts...)
}

// ProvisionerNoFakeIntake wraps Provisioner with hardcoded WithoutFakeIntake option.
func ProvisionerNoFakeIntake(opts ...ProvisionerOption) e2e.TypedProvisioner[environments.Host] {
	mergedOpts := make([]ProvisionerOption, 0, len(opts)+1)
	mergedOpts = append(mergedOpts, opts...)
	mergedOpts = append(mergedOpts, WithoutFakeIntake())

	return Provisioner(mergedOpts...)
}

// Provisioner creates a VM environment with an EC2 VM, an ECS Fargate FakeIntake and a Host Agent configured to talk to each other.
// FakeIntake and Agent creation can be deactivated by using [WithoutFakeIntake] and [WithoutAgent] options.
func Provisioner(opts ...ProvisionerOption) e2e.TypedProvisioner[environments.Host] {
	params := newProvisionerParams()
	err := optional.ApplyOptions(params, opts)

	provisioner := e2e.NewTypedPulumiProvisioner(provisionerBaseID+params.name, func(ctx *pulumi.Context, env *environments.Host) error {
		// We are abusing Pulumi RunFunc error to return our parameter parsing error, in the sake of the slightly simpler API.
		if err != nil {
			return err
		}

		awsEnv, err := aws.NewEnvironment(ctx)
		if err != nil {
			return err
		}

		host, err := ec2.NewVM(awsEnv, params.name, params.instanceOptions...)
		if err != nil {
			return err
		}
		err = host.Export(ctx, &env.RemoteHost.HostOutput)
		if err != nil {
			return err
		}

		// Create FakeIntake if required
		if params.fakeintakeOptions != nil {
			fakeIntake, err := fakeintake.NewECSFargateInstance(awsEnv, params.name, params.fakeintakeOptions...)
			if err != nil {
				return err
			}
			err = fakeIntake.Export(ctx, &env.FakeIntake.FakeintakeOutput)
			if err != nil {
				return err
			}

			// Normally if FakeIntake is enabled, Agent is enabled, but just in case
			if params.agentOptions != nil {
				// Prepend in case it's overridden by the user
				newOpts := []agentparams.Option{agentparams.WithFakeintake(fakeIntake)}
				params.agentOptions = append(newOpts, params.agentOptions...)
			}
		} else {
			// Suite inits all fields by default, so we need to explicitly set it to nil
			env.FakeIntake = nil
		}

		// Create Agent if required
		if params.agentOptions != nil {
			agent, err := agent.NewHostAgent(awsEnv.CommonEnvironment, host, params.agentOptions...)
			if err != nil {
				return err
			}

			err = agent.Export(ctx, &env.Agent.HostAgentOutput)
			if err != nil {
				return err
			}
		} else {
			// Suite inits all fields by default, so we need to explicitly set it to nil
			env.Agent = nil
		}

		return nil
	}, params.extraConfigParams)

	return provisioner
}
