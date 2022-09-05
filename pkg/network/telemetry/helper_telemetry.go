// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package telemetry

import (
	"fmt"
	"hash/fnv"
	"syscall"
	"unsafe"

	"github.com/DataDog/datadog-agent/pkg/util/log"
	manager "github.com/DataDog/ebpf-manager"
	"github.com/cilium/ebpf"
)

const (
	maxErrno    = 64
	maxErrnoStr = "other"
)

// BPFTelemetry struct contains all the maps that
// are registered to have their telemetry collected.
type BPFTelemetry struct {
	MapErrMap    *ebpf.Map
	HelperErrMap *ebpf.Map
	maps         []string
	mapKeys      map[string]uint64
}

// Initialize a new BPFTelemetry object
func NewBPFTelemetry() *BPFTelemetry {
	b := new(BPFTelemetry)
	b.mapKeys = make(map[string]uint64)

	return b
}

// Register a ebpf map entry in the map error telemetry map, to have
// failing operation telemetry recorded.
func (b *BPFTelemetry) RegisterMaps(maps []string) error {
	b.maps = append(b.maps, maps...)
	return b.initializeMapErrTelemetryMap()
}

// Returns a map of error telemetry for each ebpf map
func (b *BPFTelemetry) GetMapsTelemetry() map[string]interface{} {
	var val MapErrTelemetry
	t := make(map[string]interface{})

	for m, k := range b.mapKeys {
		err := b.MapErrMap.Lookup(&k, &val)
		if err != nil {
			log.Debugf("failed to get telemetry for map:key %s:%d\n", m, k)
		}
		t[m] = getMapErrCount(&val)
	}

	return t
}

func getMapErrCount(v *MapErrTelemetry) map[string]uint32 {
	errCount := make(map[string]uint32)

	for i, count := range v.Count {
		if count == 0 {
			continue
		}

		if (i + 1) == maxErrno {
			errCount[maxErrnoStr] = count
		} else {
			errCount[syscall.Errno(i).Error()] = count
		}
	}

	return errCount
}

// This functions builds the keys used to index in the map error telemetry ebpf map
// for recording map operation failure telemetry. The keys are built using the map names
func BuildMapErrTelemetryKeys(mgr *manager.Manager) []manager.ConstantEditor {
	var keys []manager.ConstantEditor

	h := fnv.New64a()
	for _, m := range mgr.Maps {
		h.Write([]byte(m.Name))
		keys = append(keys, manager.ConstantEditor{
			Name:  m.Name + "_telemetry_key",
			Value: h.Sum64(),
		})
		h.Reset()
	}

	return keys
}

func (b *BPFTelemetry) initializeMapErrTelemetryMap() error {
	z := new(MapErrTelemetry)
	h := fnv.New64a()

	for _, m := range b.maps {
		h.Write([]byte(m))
		key := h.Sum64()
		err := b.MapErrMap.Put(unsafe.Pointer(&key), unsafe.Pointer(z))
		if err != nil {
			return fmt.Errorf("failed to initialize telemetry struct for map %s", m)
		}
		h.Reset()

		b.mapKeys[m] = key
	}

	return nil
}
