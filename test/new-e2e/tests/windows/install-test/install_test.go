// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package installtest

import (
	"flag"
	"fmt"

	"github.com/DataDog/test-infra-definitions/scenarios/aws/vm/ec2os"
	"github.com/DataDog/test-infra-definitions/scenarios/aws/vm/ec2params"

	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/e2e"
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

	msiURL       string
	majorVersion string
}

func TestMSI(t *testing.T) {
	var opts []func(*params.Params)

	if *devMode {
		opts = append(opts, params.WithDevMode())
	}

	majorVersion := windowsAgent.GetMajorVersionFromEnv()
	msiURL, err := windowsAgent.GetMSIURLFromEnv()
	if err != nil {
		t.Fatalf("failed to get MSI URL from env: %v", err)
	}
	t.Logf("Using MSI URL: %v", msiURL)

	s := &agentMSISuite{
		msiURL:       msiURL,
		majorVersion: majorVersion,
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

	// TODO: Add apikey option
	apikey := "00000000000000000000000000000000"
	is.Run("install the agent", func() {
		args := fmt.Sprintf(`APIKEY="%s"`, apikey)
		err := windows.InstallMSI(vm, is.msiURL, args, "install.log")
		is.Require().NoError(err, "should install the agent")
	})

	client := common.NewWindowsTestClient(is.T(), vm)

	is.Run("agent runtime behavior", func() {
		common.CheckInstallation(is.T(), client)
		common.CheckAgentBehaviour(is.T(), client)
		common.CheckAgentStops(is.T(), client)
		common.CheckAgentRestarts(is.T(), client)
		common.CheckIntegrationInstall(is.T(), client)
		if is.majorVersion == "6" {
			common.CheckAgentPython(is.T(), client, "2")
		}
		common.CheckAgentPython(is.T(), client, "3")
		common.CheckApmEnabled(is.T(), client)
		common.CheckApmDisabled(is.T(), client)
		// TODO: common.CheckCWSBehaviour(is.T(), client)
	})

	is.Run("uninstall the agent", func() {
		err := windowsAgent.UninstallAgent(vm, "uninstall.log")
		is.Require().NoError(err, "should uninstall the agent")

		common.CheckUninstallation(is.T(), client)
	})

}
