package awsdocker

import (
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/e2e"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/environments"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/runner"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/optional"

	"github.com/DataDog/test-infra-definitions/components/datadog/agent"
	"github.com/DataDog/test-infra-definitions/components/docker"
	"github.com/DataDog/test-infra-definitions/resources/aws"
	"github.com/DataDog/test-infra-definitions/scenarios/aws/ec2"
	"github.com/DataDog/test-infra-definitions/scenarios/aws/fakeintake"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	provisionerBaseID = "aws-ec2docker-"
	defaultVMName     = "dockervm"
)

// ProvisionerParams contains all the parameters needed to create the environment
type ProvisionerParams struct {
	name string

	vmOptions         []ec2.VMOption
	agentOptions      []agent.DockerOption
	fakeintakeOptions []fakeintake.Option
	extraConfigParams runner.ConfigMap
}

func newProvisionerParams() *ProvisionerParams {
	// We use nil arrays to decide if we should create or not
	return &ProvisionerParams{
		name:              defaultVMName,
		vmOptions:         []ec2.VMOption{},
		agentOptions:      []agent.DockerOption{},
		fakeintakeOptions: []fakeintake.Option{},
		extraConfigParams: runner.ConfigMap{},
	}
}

// ProvisionerOption is a function that modifies the ProvisionerParams
type ProvisionerOption func(*ProvisionerParams) error

// WithName sets the name of the provisioner
func WithName(name string) ProvisionerOption {
	return func(params *ProvisionerParams) error {
		params.name = name
		return nil
	}
}

// WithEC2VMOptions sets the options for the EC2 VM
func WithEC2VMOptions(opts ...ec2.VMOption) ProvisionerOption {
	return func(params *ProvisionerParams) error {
		params.vmOptions = append(params.vmOptions, opts...)
		return nil
	}
}

// WithAgentOptions sets the options for the Docker Agent
func WithAgentOptions(opts ...agent.DockerOption) ProvisionerOption {
	return func(params *ProvisionerParams) error {
		params.agentOptions = append(params.agentOptions, opts...)
		return nil
	}
}

// WithFakeIntakeOptions sets the options for the FakeIntake
func WithFakeIntakeOptions(opts ...fakeintake.Option) ProvisionerOption {
	return func(params *ProvisionerParams) error {
		params.fakeintakeOptions = append(params.fakeintakeOptions, opts...)
		return nil
	}
}

// WithExtraConfigParams sets the extra config params for the environment
func WithExtraConfigParams(configMap runner.ConfigMap) ProvisionerOption {
	return func(params *ProvisionerParams) error {
		params.extraConfigParams = configMap
		return nil
	}
}

// WithoutFakeIntake deactivates the creation of the FakeIntake
func WithoutFakeIntake() ProvisionerOption {
	return func(params *ProvisionerParams) error {
		params.fakeintakeOptions = nil
		return nil
	}
}

// WithoutAgent deactivates the creation of the Docker Agent
func WithoutAgent() ProvisionerOption {
	return func(params *ProvisionerParams) error {
		params.agentOptions = nil
		params.fakeintakeOptions = nil
		return nil
	}
}

// Provisioner creates a VM environment with an EC2 VM with Docker, an ECS Fargate FakeIntake and a Docker Agent configured to talk to each other.
// FakeIntake and Agent creation can be deactivated by using [WithoutFakeIntake] and [WithoutAgent] options.
func Provisioner(opts ...ProvisionerOption) e2e.TypedProvisioner[environments.DockerVM] {
	params := newProvisionerParams()
	err := optional.ApplyOptions(params, opts)

	provisioner := e2e.NewPulumiTypedProvisioner(provisionerBaseID+params.name, func(ctx *pulumi.Context, env *environments.DockerVM) error {
		// We are abusing Pulumi RunFunc error to return our parameter parsing error, in the sake of the slightly simpler API.
		if err != nil {
			return err
		}

		awsEnv, err := aws.NewEnvironment(ctx)
		if err != nil {
			return err
		}

		host, err := ec2.NewVM(awsEnv, params.name, params.vmOptions...)
		if err != nil {
			return err
		}
		host.Export(ctx, &env.Host.HostOutput)

		manager, _, err := docker.NewManager(*awsEnv.CommonEnvironment, host, true)
		if err != nil {
			return err
		}

		// Create FakeIntake if required
		if params.fakeintakeOptions != nil {
			fakeIntake, err := fakeintake.NewECSFargateInstance(awsEnv, params.name, params.fakeintakeOptions...)
			if err != nil {
				return err
			}
			fakeIntake.Export(ctx, &env.FakeIntake.FakeintakeOutput)

			// TODO: Pending PR, but currently Docker Agent does not support fakeintake
		} else {
			// Suite inits all fields by default, so we need to explicitly set it to nil
			env.FakeIntake = nil
		}

		// Create Agent if required
		if params.agentOptions != nil {
			agent, err := agent.NewDockerAgent(*awsEnv.CommonEnvironment, host, manager, params.agentOptions...)
			if err != nil {
				return err
			}
			agent.Export(ctx, &env.Agent.DockerAgentOutput)
		} else {
			// Suite inits all fields by default, so we need to explicitly set it to nil
			env.Agent = nil
		}

		return nil
	}, params.extraConfigParams)

	return provisioner
}