// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package universal_testing

import (
	"github.com/DataDog/datadog-agent/test/universal-testing/client"
	"github.com/DataDog/datadog-agent/test/universal-testing/infra"
	"github.com/DataDog/datadog-agent/test/universal-testing/infra/localvmparams"
)

func LocalVMDef(options ...localvmparams.Option) InfraProvider[VMEnv] {
	return NewLocalProvider(
		func(vmManager *infra.LocalVMManager) (*VMEnv, error) {
			vm, err := infra.NewLocalVM(options...)
			if err != nil {
				return nil, err
			}
			vmManager.AddVM(vm)

			return &VMEnv{
				VM: client.NewSSHVM(vm.Name(), vm.OSType()),
			}, nil
		},
	)
}
