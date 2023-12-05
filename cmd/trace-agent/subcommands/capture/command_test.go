// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package capture

import (
	"testing"

	"github.com/spf13/cobra"
	"gotest.tools/assert"

	"github.com/DataDog/datadog-agent/cmd/trace-agent/subcommands"
	"github.com/DataDog/datadog-agent/comp/trace/config"
	"github.com/DataDog/datadog-agent/pkg/util/fxutil"
)

func TestCaptureCommand(t *testing.T) {
	fxutil.TestOneShotSubcommand(t,
		[]*cobra.Command{MakeCommand(func() *subcommands.GlobalParams {
			return &subcommands.GlobalParams{}
		})},
		[]string{"capture"},
		capture,
		func(config config.Component) {
			assert.Assert(t, config.Object() != nil)
		})
}