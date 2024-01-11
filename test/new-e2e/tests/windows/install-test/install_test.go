// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package installtest contains e2e tests for the Windows agent installer
package installtest

import (
	"flag"
	"fmt"
	"strings"

	"github.com/DataDog/test-infra-definitions/scenarios/aws/vm/ec2os"
	"github.com/DataDog/test-infra-definitions/scenarios/aws/vm/ec2params"

	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/e2e"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/e2e/params"
	windows "github.com/DataDog/datadog-agent/test/new-e2e/tests/windows"
	windowsAgent "github.com/DataDog/datadog-agent/test/new-e2e/tests/windows/agent"

	"testing"
)

var (
	devMode = flag.Bool("devmode", false, "enable devmode")
)

type agentMSISuite struct {
	e2e.Suite[e2e.VMEnv]

	agentPackage *windowsAgent.Package
	majorVersion string
}

func TestMSI(t *testing.T) {
	var opts []func(*params.Params)

	if *devMode {
		opts = append(opts, params.WithDevMode())
	}

	agentPackage, err := windowsAgent.GetPackageFromEnv()
	if err != nil {
		t.Fatalf("failed to get MSI URL from env: %v", err)
	}
	t.Logf("Using Agent: %#v", agentPackage)

	// Set stack name to avoid conflicts with other tests
	// Include channel if we're not running in a CI pipeline.
	// E2E auto includes the pipeline ID in the stack name, so we don't need to do that here.
	stackNameChannelPart := ""
	if agentPackage.PipelineID == "" && agentPackage.Channel != "" {
		stackNameChannelPart = fmt.Sprintf("-%s", agentPackage.Channel)
	}
	opts = append(opts, params.WithStackName(fmt.Sprintf("windows-msi-test-v%s-%s%s", agentPackage.Version, agentPackage.Arch, stackNameChannelPart)))

	s := &agentMSISuite{
		agentPackage: agentPackage,
		majorVersion: strings.Split(agentPackage.Version, ".")[0],
	}

	e2e.Run(t,
		s,
		e2e.EC2VMStackDef(ec2params.WithOS(ec2os.WindowsOS)),
		opts...)
}

func (is *agentMSISuite) prepareVM() {
	vm := is.Env().VM

	if !is.Run("prepare VM", func() {
		is.Run("disable defender", func() {
			err := windows.DisableDefender(vm)
			is.Require().NoError(err, "should disable defender")
		})
	}) {
		is.T().Fatal("failed to prepare VM")
	}
}

func (is *agentMSISuite) TestInstall() {
	vm := is.Env().VM
	is.prepareVM()

	t, err := NewTester(is.T(), vm, WithExpectedAgentVersion(is.agentPackage.AgentVersion()))
	is.Require().NoError(err, "should create tester")

	if !t.TestInstallAgentPackage(is.T(), is.agentPackage, "", "install.log") {
		is.T().Fatal("failed to install agent")
	}
	t.TestRuntimeExpectations(is.T())
	t.TestUninstall(is.T())
}

func (is *agentMSISuite) TestUpgrade() {
	vm := is.Env().VM
	is.prepareVM()

	t, err := NewTester(is.T(), vm, WithExpectedAgentVersion(is.agentPackage.AgentVersion()))
	is.Require().NoError(err, "should create tester")

	// install old agent
	lastStableAgentPackage := is.installLastStable(t)

	// upgrade to new agent
	if !t.TestInstallAgentPackage(is.T(), is.agentPackage, "", "upgrade.log") {
		is.T().Fatal("failed to upgrade agent")
	}

	// Check that the agent was upgraded
	newVersion, err := t.InstallTestClient.GetAgentVersion()
	is.Require().NoError(err, "should get agent version")
	is.Assert().NotEqual(lastStableAgentPackage.AgentVersion(), newVersion, "new version should be installed")
	is.Assert().Equal(is.agentPackage.AgentVersion(), newVersion, "new version should be installed")

	t.TestRuntimeExpectations(is.T())
	t.TestUninstall(is.T())
}

func (is *agentMSISuite) installLastStable(t *Tester) *windowsAgent.Package {
	var agentPackage *windowsAgent.Package

	if !is.Run("install prev stable agent", func() {
		var err error

		agentPackage, err = windowsAgent.GetLastStablePackageFromEnv()
		is.Require().NoError(err, "should get last stable agent package from env")

		t.InstallAgentPackage(is.T(), agentPackage, "", "install.log")

		agentVersion, err := t.InstallTestClient.GetAgentVersion()
		is.Require().NoError(err, "should get agent version")
		is.Assert().Equal(agentPackage.AgentVersion(), agentVersion, "installed agent version should match expected agent version")
	}) {
		is.T().Fatal("failed to install last stable agent")
	}

	return agentPackage
}
