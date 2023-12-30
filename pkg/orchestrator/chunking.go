// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package orchestrator

import "github.com/DataDog/datadog-agent/pkg/process/util"

// ChunkOrchestratorPayloadsBySizeAndWeight chunks orchestrator payloads by max allowed size and max allowed weight of a chunk
func ChunkOrchestratorPayloadsBySizeAndWeight(orchestratorPayloads []interface{}, orchestratorYaml []interface{}, weightForOrchestratorPayload func([]interface{}, int) int, maxChunkSize, maxChunkWeight int) [][]interface{} {
	if len(orchestratorPayloads) == 0 {
		return make([][]interface{}, 0)
	}

	chunker := &util.ChunkAllocator[[]interface{}, interface{}]{
		AppendToChunk: func(chunk *[]interface{}, payloads []interface{}) {
			*chunk = append(*chunk, payloads...)
		},
	}

	list := &util.PayloadList[interface{}]{
		Items: orchestratorPayloads,
		WeightAt: func(i int) int {
			return weightForOrchestratorPayload(orchestratorYaml, i)
		},
	}

	util.ChunkPayloadsBySizeAndWeight[[]interface{}, interface{}](list, chunker, maxChunkSize, maxChunkWeight)

	return *chunker.GetChunks()
}
