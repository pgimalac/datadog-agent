// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package trace

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/DataDog/datadog-agent/comp/core"
	"github.com/DataDog/datadog-agent/comp/core/tagger"
	"github.com/DataDog/datadog-agent/comp/core/workloadmeta"
	"github.com/DataDog/datadog-agent/comp/trace/agent"
	"github.com/DataDog/datadog-agent/comp/trace/config"
	"github.com/DataDog/datadog-agent/pkg/trace/telemetry"
	"github.com/DataDog/datadog-agent/pkg/util/fxutil"
)

// team: agent-apm

func TestBundleDependencies(t *testing.T) {
	fxutil.TestBundle(t, Bundle,
		fx.Provide(func() context.Context { return context.TODO() }), // fx.Supply(ctx) fails with a missing type error.
		fx.Supply(core.BundleParams{}),
		core.Bundle,

		fx.Supply(workloadmeta.NewParams()),
		workloadmeta.Module,
		fx.Provide(func(cfg config.Component) telemetry.TelemetryCollector { return telemetry.NewCollector(cfg.Object()) }),
		fx.Supply(tagger.NewFakeTaggerParams()),
		tagger.Module,
		fx.Supply(&agent.Params{}),
	)
}

func TestMockBundleDependencies(t *testing.T) {
	os.Setenv("DD_APP_KEY", "abc1234")
	defer func() { os.Unsetenv("DD_APP_KEY") }()

	os.Setenv("DD_DD_URL", "https://example.com")
	defer func() { os.Unsetenv("DD_DD_URL") }()

	cfg := fxutil.Test[config.Component](t, fx.Options(
		fx.Provide(func() context.Context { return context.TODO() }), // fx.Supply(ctx) fails with a missing type error.
		fx.Supply(core.BundleParams{}),
		core.MockBundle,
		fx.Supply(workloadmeta.NewParams()),
		workloadmeta.Module,
		fx.Invoke(func(_ config.Component) {}),
		fx.Provide(func(cfg config.Component) telemetry.TelemetryCollector { return telemetry.NewCollector(cfg.Object()) }),
		fx.Supply(&agent.Params{}),
		fx.Invoke(func(_ agent.Component) {}),
		tagger.Module,
		fx.Supply(tagger.NewTaggerParams()),
		MockBundle,
	))

	require.NotNil(t, cfg.Object())
}
