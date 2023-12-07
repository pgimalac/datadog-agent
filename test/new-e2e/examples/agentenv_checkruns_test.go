// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package examples

import (
	"testing"
	"time"

	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/e2e"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/environments"
	awsvm "github.com/DataDog/datadog-agent/test/new-e2e/pkg/environments/aws/vm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type agentSuiteEx5 struct {
	e2e.BaseSuite[environments.VM]
}

func TestAgentSuiteEx5(t *testing.T) {
	e2e.Run(t, &agentSuiteEx5{}, e2e.WithProvisioner(awsvm.Provisioner()))
}

func (s *agentSuiteEx5) TestCheckRuns() {
	s.EventuallyWithT(func(c *assert.CollectT) {
		checkRuns, err := s.Env().FakeIntake.Client().GetCheckRun("datadog.agent.up")
		require.NoError(c, err)
		assert.Greater(c, len(checkRuns), 0)
	}, 30*time.Second, 500*time.Millisecond)
}
