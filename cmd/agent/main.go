// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Main package for the agent binary
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/DataDog/datadog-agent/cmd/internal/runcmd"
)

func main() {
	executable := filepath.Base(os.Args[0])
	process := strings.TrimSuffix(executable, filepath.Ext(executable))

	if agentCmdBuilder := agents[process]; agentCmdBuilder != nil {
		if rootCmd := agentCmdBuilder(); rootCmd != nil {
			os.Exit(runcmd.Run(rootCmd))
		}

		// if not command is returned, main was already handled by the callback
		return
	}

	fmt.Fprintf(os.Stderr, "'%s' is an incorrect invocation of the Datadog Agent\n", process)
	os.Exit(1)
}
