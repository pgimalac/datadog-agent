// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package installtest

import (
	"flag"
	"fmt"
	"strings"

	"github.com/DataDog/test-infra-definitions/scenarios/aws/vm/ec2os"
	"github.com/DataDog/test-infra-definitions/scenarios/aws/vm/ec2params"

	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/e2e"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/e2e/client"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/e2e/params"
	"github.com/DataDog/datadog-agent/test/new-e2e/tests/agent-platform/common"
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

func (is *agentMSISuite) TestInstallAgent() {
	vm := is.Env().VM

	err := windows.DisableDefender(vm)
	is.Require().NoError(err, "should disable defender")

	client := common.NewWindowsTestClient(is.T(), vm)

	// TODO: Add apikey option
	apikey := "00000000000000000000000000000000"
	is.Run("install the agent", func() {
		args := fmt.Sprintf(`APIKEY="%s"`, apikey)
		err := windows.InstallMSI(vm, is.agentPackage.URL, args, "install.log")
		is.Require().NoError(err, "should install the agent")

		common.CheckInstallation(is.T(), client)
		is.testCodeSignature(client.VMClient)
	})

	is.Run("agent runtime behavior", func() {
		common.CheckAgentBehaviour(is.T(), client)
		common.CheckAgentStops(is.T(), client)
		common.CheckAgentRestarts(is.T(), client)
		common.CheckIntegrationInstall(is.T(), client)
		if is.IsPython2Installed() {
			common.CheckAgentPython(is.T(), client, "2")
		}
		common.CheckAgentPython(is.T(), client, "3")
		common.CheckApmEnabled(is.T(), client)
		common.CheckApmDisabled(is.T(), client)
		// TODO: CWS on Windows isn't available yet
		// common.CheckCWSBehaviour(is.T(), client)
	})

	is.Run("uninstall the agent", func() {
		err := windowsAgent.UninstallAgent(vm, "uninstall.log")
		is.Require().NoError(err, "should uninstall the agent")

		common.CheckUninstallation(is.T(), client)
	})

}

func (is *agentMSISuite) IsPython2Installed() bool {
	return is.majorVersion == "6"
}

func (is *agentMSISuite) testCodeSignature(client client.VM) {
	root := `C:\Program Files\Datadog\Datadog Agent\`
	paths := []string{
		// user binaries
		root + `bin\agent.exe`,
		root + `bin\libdatadog-agent-three.dll`,
		root + `bin\agent\trace-agent.exe`,
		root + `bin\agent\process-agent.exe`,
		root + `bin\agent\system-probe.exe`,
		// drivers
		root + `bin\agent\driver\ddnpm.sys`,
	}
	// Python3 should be signed by Python, since we don't build our own anymore
	// We still build our own Python2, so we need to check that
	if is.IsPython2Installed() {
		paths = append(paths, []string{
			root + `bin\libdatadog-agent-three.dll`,
			root + `embedded2\python.exe`,
			root + `embedded2\pythonw.exe`,
			root + `embedded2\python27.dll`,
		}...)
	}

	windowsAgent.TestValidDatadogCodeSignatures(is.T(), client, paths)
}
