// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package workloadmeta

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/DataDog/datadog-agent/pkg/util/pointer"
)

func container1(testTime time.Time) Container {
	return Container{
		EntityID: EntityID{
			Kind: KindContainer,
			ID:   "foo1",
		},
		EntityMeta: EntityMeta{
			Name:      "foo1-name",
			Namespace: "",
		},
		Ports: []ContainerPort{
			{
				Name:     "port1",
				Port:     42000,
				Protocol: "tcp",
			},
			{
				Port:     42001,
				Protocol: "udp",
			},
			{
				Port: 42002,
			},
		},
		Volumes: []ContainerVolume{
			{
				DockerName:  "vol1",
				Destination: "/host/path",
			},
		},
		Networks: []ContainerNetwork{
			{
				NetworkMode:   "bridge",
				IPv4Addresses: []string{"172.17.0.2"},
			},
		},
		State: ContainerState{
			Running:    true,
			CreatedAt:  testTime,
			StartedAt:  testTime,
			FinishedAt: time.Time{},
		},
		CollectorTags: []string{"tag1", "tag2"},
	}
}

func container2(testTime time.Time) Container { //nolint:revive // TODO fix revive unused-parameter
	return Container{
		EntityID: EntityID{
			Kind: KindContainer,
			ID:   "foo1",
		},
		EntityMeta: EntityMeta{
			Name:      "foo1-name",
			Namespace: "",
		},
		Ports: []ContainerPort{
			{
				Port:     42000,
				Protocol: "tcp",
			},
			{
				Port:     42001,
				Protocol: "udp",
			},
			{
				Port:     42002,
				Protocol: "tcp",
			},
			{
				Port: 42003,
			},
		},
		Volumes: []ContainerVolume{
			{
				DockerName:  "vol1",
				Destination: "/host/path",
			},
		},
		Networks: []ContainerNetwork{
			{
				NetworkMode:   "bridge",
				IPv4Addresses: []string{"172.17.0.2"},
			},
		},
		State: ContainerState{
			CreatedAt:  time.Time{},
			StartedAt:  time.Time{},
			FinishedAt: time.Time{},
			ExitCode:   pointer.Ptr(uint32(100)),
		},
		CollectorTags: []string{"tag3"},
	}
}

func TestMerge(t *testing.T) {
	testTime := time.Now()

	expectedContainer := Container{
		EntityID: EntityID{
			Kind: KindContainer,
			ID:   "foo1",
		},
		EntityMeta: EntityMeta{
			Name:      "foo1-name",
			Namespace: "",
		},
		Volumes: []ContainerVolume{
			{
				DockerName:  "vol1",
				Destination: "/host/path",
			},
		},
		Networks: []ContainerNetwork{
			{
				NetworkMode:   "bridge",
				IPv4Addresses: []string{"172.17.0.2"},
			},
		},
		State: ContainerState{
			Running:    true,
			CreatedAt:  testTime,
			StartedAt:  testTime,
			FinishedAt: time.Time{},
			ExitCode:   pointer.Ptr(uint32(100)),
		},
	}

	expectedPorts := []ContainerPort{
		{
			Name:     "port1",
			Port:     42000,
			Protocol: "tcp",
		},
		{
			Port:     42001,
			Protocol: "udp",
		},
		{
			Port: 42002,
		},
		{
			Port:     42002,
			Protocol: "tcp",
		},
		{
			Port: 42003,
		},
	}

	expectedTags := []string{"tag1", "tag2", "tag3"}

	// Test merging both ways
	fromSource1 := container1(testTime)
	fromSource2 := container2(testTime)
	err := merge(&fromSource1, &fromSource2)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedPorts, fromSource1.Ports)
	assert.ElementsMatch(t, expectedTags, fromSource1.CollectorTags)
	fromSource1.Ports = nil
	fromSource1.CollectorTags = nil
	assert.Equal(t, expectedContainer, fromSource1)

	fromSource1 = container1(testTime)
	fromSource2 = container2(testTime)
	err = merge(&fromSource2, &fromSource1)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedPorts, fromSource2.Ports)
	assert.ElementsMatch(t, expectedTags, fromSource2.CollectorTags)
	fromSource2.Ports = nil
	fromSource2.CollectorTags = nil
	assert.Equal(t, expectedContainer, fromSource2)

	// Test merging nil slice in src/dst
	fromSource1 = container1(testTime)
	fromSource2 = container2(testTime)
	fromSource2.Ports = nil
	err = merge(&fromSource1, &fromSource2)
	assert.NoError(t, err)
	assert.ElementsMatch(t, container1(testTime).Ports, fromSource1.Ports)

	fromSource1 = container1(testTime)
	fromSource2 = container2(testTime)
	fromSource2.Ports = nil
	err = merge(&fromSource2, &fromSource1)
	assert.NoError(t, err)
	assert.ElementsMatch(t, container1(testTime).Ports, fromSource2.Ports)
}

func TestMergeECSTask(t *testing.T) {
	task1 := ECSTask{
		EntityID: EntityID{
			Kind: KindECSTask,
			ID:   "foo1",
		},
		EntityMeta: EntityMeta{
			Name: "foo1-name",
		},
		Family: "foo1-family",
		Containers: []OrchestratorContainer{
			{
				ID:   "foo1-c1-id",
				Name: "foo1-c1-name",
			},
		},
		Tags: ECSTaskTags{
			"tag1": "value1",
		},
	}

	task2 := ECSTask{
		EntityID: EntityID{
			Kind: KindECSTask,
			ID:   "foo1",
		},
		EntityMeta: EntityMeta{
			Name: "foo1-name",
		},
		ServiceName: "foo1-service-name",
		Containers: []OrchestratorContainer{
			{
				ID:   "foo1-c1-id",
				Name: "foo1-c1-name",
			},
		},
	}

	expected := ECSTask{
		EntityID: EntityID{
			Kind: KindECSTask,
			ID:   "foo1",
		},
		EntityMeta: EntityMeta{
			Name: "foo1-name",
		},
		Family:      "foo1-family",
		ServiceName: "foo1-service-name",
		Containers: []OrchestratorContainer{
			{
				ID:   "foo1-c1-id",
				Name: "foo1-c1-name",
			},
		},
		Tags: ECSTaskTags{
			"tag1": "value1",
		},
	}

	err := merge(&task1, &task2)
	assert.NoError(t, err)
	assert.Equal(t, expected, task1)

	err = merge(&task2, &task1)
	assert.NoError(t, err)
	assert.Equal(t, expected, task2)

	// Tags should be overwritten by incoming task
	task2.Tags = ECSTaskTags{
		"tag2": "value2",
	}
	expected.Tags = ECSTaskTags{
		"tag2": "value2",
	}
	err = merge(&task1, &task2)
	assert.NoError(t, err)
	assert.Equal(t, expected, task1)

	err = merge(&task2, &task1)
	assert.NoError(t, err)
	assert.Equal(t, expected, task2)
}
