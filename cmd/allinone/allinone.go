// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build !windows

// Main package for the allinone binary
package main

import (
	"fmt"
	"os"
	"path"
	"runtime/debug"
	"strings"
	"time"

	"github.com/DataDog/datadog-agent/cmd/internal/runcmd"
	"github.com/DataDog/datadog-agent/pkg/util/log"

	"github.com/spf13/cobra"
)

const temporaryGOGCDuration = 5 * time.Second
const temporaryGOGCPercent = 1

var agents = map[string]func() *cobra.Command{}

func registerAgent(getCommand func() *cobra.Command, names ...string) {
	for _, name := range names {
		agents[name] = getCommand
	}
}

func main() {
	previousGCPercent := debug.SetGCPercent(temporaryGOGCPercent)

	go func() {
		time.Sleep(temporaryGOGCDuration)

		log.Infof("Restoring GOGC to %d (previous GOGC: %d)\n", previousGCPercent, temporaryGOGCPercent)
		debug.SetGCPercent(previousGCPercent)
	}()

	executable := path.Base(os.Args[0])
	process := strings.TrimSuffix(executable, path.Ext(executable))

	if agentCmdBuilder := agents[process]; agentCmdBuilder != nil {
		rootCmd := agentCmdBuilder()
		os.Exit(runcmd.Run(rootCmd))
	}

	fmt.Fprintf(os.Stderr, "'%s' is an incorrect invocation of the Datadog Agent\n", process)
	os.Exit(1)
}
