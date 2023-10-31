// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package universal_testing

import (
	"context"
	"testing"
)

type InfraProvider[Env any] interface {
	ProvisionInfraAndInitializeEnv(t *testing.T, ctx context.Context, name string, failOnMissing bool) (*Env, error)
	DeleteInfra(ctx context.Context, name string) error
}
