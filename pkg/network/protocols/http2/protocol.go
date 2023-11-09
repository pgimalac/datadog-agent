// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build linux_bpf

package http2

import (
	"errors"
	"fmt"
	netebpf "github.com/DataDog/datadog-agent/pkg/network/ebpf"
	"strings"
	"time"
	"unsafe"

	"github.com/cilium/ebpf"
	"github.com/davecgh/go-spew/spew"

	manager "github.com/DataDog/ebpf-manager"

	ddebpf "github.com/DataDog/datadog-agent/pkg/ebpf"
	"github.com/DataDog/datadog-agent/pkg/network/config"
	"github.com/DataDog/datadog-agent/pkg/network/ebpf/probes"
	"github.com/DataDog/datadog-agent/pkg/network/protocols"
	"github.com/DataDog/datadog-agent/pkg/network/protocols/events"
	"github.com/DataDog/datadog-agent/pkg/network/protocols/http"
	"github.com/DataDog/datadog-agent/pkg/network/usm/buildmode"
	"github.com/DataDog/datadog-agent/pkg/network/usm/utils"
	"github.com/DataDog/datadog-agent/pkg/util/log"
)

type protocol struct {
	cfg *config.Config
	// TODO: Do we need to duplicate?
	telemetry *http.Telemetry
	// TODO: Do we need to duplicate?
	statkeeper     *http.StatKeeper
	mapCleaner     *ddebpf.MapCleaner
	eventsConsumer *events.Consumer
}

const (
	inFlightMap               = "http2_in_flight"
	dynamicTable              = "http2_dynamic_table"
	dynamicTableCounter       = "http2_dynamic_counter_table"
	http2IterationsTable      = "http2_iterations"
	staticTable               = "http2_static_table"
	firstFrameHandlerTailCall = "socket__http2_handle_first_frame"
	filterTailCall            = "socket__http2_filter"
	parserTailCall            = "socket__http2_frames_parser"
	eventStream               = "http2"
)

var Spec = &protocols.ProtocolSpec{
	Factory: newHTTP2Protocol,
	Maps: []*manager.Map{
		{
			Name: inFlightMap,
		},
		{
			Name: dynamicTable,
		},
		{
			Name: staticTable,
		},
		{
			Name: dynamicTableCounter,
		},
		{
			Name: http2IterationsTable,
		},
		{
			Name: "http2_headers_to_process",
		},
		{
			Name: "http2_frames_to_process",
		},
		{
			Name: "http2_stream_heap",
		},
		{
			Name: "http2_ctx_heap",
		},
	},
	TailCalls: []manager.TailCallRoute{
		{
			ProgArrayName: protocols.ProtocolDispatcherProgramsMap,
			Key:           uint32(protocols.ProgramHTTP2HandleFirstFrame),
			ProbeIdentificationPair: manager.ProbeIdentificationPair{
				EBPFFuncName: firstFrameHandlerTailCall,
			},
		},
		{
			ProgArrayName: protocols.ProtocolDispatcherProgramsMap,
			Key:           uint32(protocols.ProgramHTTP2FrameFilter),
			ProbeIdentificationPair: manager.ProbeIdentificationPair{
				EBPFFuncName: filterTailCall,
			},
		},
		{
			ProgArrayName: protocols.ProtocolDispatcherProgramsMap,
			Key:           uint32(protocols.ProgramHTTP2FrameParser),
			ProbeIdentificationPair: manager.ProbeIdentificationPair{
				EBPFFuncName: parserTailCall,
			},
		},
	},
}

func newHTTP2Protocol(cfg *config.Config) (protocols.Protocol, error) {
	if !cfg.EnableHTTP2Monitoring {
		return nil, nil
	}

	if !Supported() {
		return nil, fmt.Errorf("http2 feature not available on pre %s kernels", MinimumKernelVersion.String())
	}

	telemetry := http.NewTelemetry("http2")

	return &protocol{
		cfg:       cfg,
		telemetry: telemetry,
	}, nil
}

func (p *protocol) Name() string {
	return "HTTP2"
}

const (
	mapSizeValue = 1024
)

// ConfigureOptions add the necessary options for http2 monitoring to work,
// to be used by the manager. These are:
// - Set the `http2_in_flight` map size to the value of the `max_tracked_connection` configuration variable.
//
// We also configure the http2 event stream with the manager and its options.
func (p *protocol) ConfigureOptions(mgr *manager.Manager, opts *manager.Options) {
	opts.MapSpecEditors[inFlightMap] = manager.MapSpecEditor{
		MaxEntries: p.cfg.MaxUSMConcurrentRequests,
		EditorFlag: manager.EditMaxEntries,
	}

	opts.MapSpecEditors[dynamicTable] = manager.MapSpecEditor{
		MaxEntries: mapSizeValue,
		EditorFlag: manager.EditMaxEntries,
	}
	opts.MapSpecEditors[dynamicTableCounter] = manager.MapSpecEditor{
		MaxEntries: mapSizeValue,
		EditorFlag: manager.EditMaxEntries,
	}
	opts.MapSpecEditors[http2IterationsTable] = manager.MapSpecEditor{
		MaxEntries: mapSizeValue,
		EditorFlag: manager.EditMaxEntries,
	}

	utils.EnableOption(opts, "http2_monitoring_enabled")
	// Configure event stream
	events.Configure(eventStream, mgr, opts)
}

func (p *protocol) PreStart(mgr *manager.Manager) (err error) {
	p.eventsConsumer, err = events.NewConsumer(
		eventStream,
		mgr,
		p.processHTTP2,
	)
	if err != nil {
		return
	}

	p.statkeeper = http.NewStatkeeper(p.cfg, p.telemetry)
	p.eventsConsumer.Start()

	if err = p.createStaticTable(mgr); err != nil {
		return fmt.Errorf("error creating a static table for http2 monitoring: %w", err)
	}

	return
}

func (p *protocol) PostStart(mgr *manager.Manager) error {
	// Setup map cleaner after manager start.
	p.setupMapCleaner(mgr)
	p.setupMapCleaner2(mgr)

	return nil
}

func (p *protocol) Stop(_ *manager.Manager) {
	// mapCleaner handles nil pointer receivers
	p.mapCleaner.Stop()

	if p.eventsConsumer != nil {
		p.eventsConsumer.Stop()
	}

	if p.statkeeper != nil {
		p.statkeeper.Close()
	}
}

func (p *protocol) DumpMaps(output *strings.Builder, mapName string, currentMap *ebpf.Map) {
	if mapName == inFlightMap { // maps/http2_in_flight (BPF_MAP_TYPE_HASH), key ConnTuple, value httpTX
		output.WriteString("Map: '" + mapName + "', key: 'ConnTuple', value: 'httpTX'\n")
		iter := currentMap.Iterate()
		var key http2StreamKey
		var value EbpfTx
		for iter.Next(unsafe.Pointer(&key), unsafe.Pointer(&value)) {
			output.WriteString(spew.Sdump(key, value))
		}
	} else if mapName == dynamicTable {
		output.WriteString("Map: '" + mapName + "', key: 'ConnTuple', value: 'httpTX'\n")
		iter := currentMap.Iterate()
		var key http2DynamicTableIndex
		var value http2DynamicTableEntry
		for iter.Next(unsafe.Pointer(&key), unsafe.Pointer(&value)) {
			output.WriteString(spew.Sdump(key, value))
		}
	} else if mapName == "http2_terminated_conns" {
		output.WriteString("Map: '" + mapName + "', key: 'ConnTuple', value: 'httpTX'\n")
		var value netebpf.ConnTuple
		for {
			if err := currentMap.LookupAndDelete(nil, unsafe.Pointer(&value)); err != nil {
				return
			}

			output.WriteString(spew.Sdump(value))
		}
	}
}

func (p *protocol) processHTTP2(data []byte) {
	tx := (*EbpfTx)(unsafe.Pointer(&data[0]))
	p.telemetry.Count(tx)
	p.statkeeper.Process(tx)
}

func (p *protocol) setupMapCleaner(mgr *manager.Manager) {
	http2Map, _, err := mgr.GetMap(inFlightMap)
	if err != nil {
		log.Errorf("error getting %q map: %s", inFlightMap, err)
		return
	}
	mapCleaner, err := ddebpf.NewMapCleaner(http2Map, new(http2StreamKey), new(EbpfTx))
	if err != nil {
		log.Errorf("error creating map cleaner: %s", err)
		return
	}

	ttl := p.cfg.HTTPIdleConnectionTTL.Nanoseconds()
	mapCleaner.Clean(p.cfg.HTTPMapCleanerInterval, func(now int64, key, val interface{}) bool {
		http2Txn, ok := val.(*EbpfTx)
		if !ok {
			return false
		}

		if updated := int64(http2Txn.Response_last_seen); updated > 0 {
			return (now - updated) > ttl
		}

		started := int64(http2Txn.Request_started)
		return started > 0 && (now-started) > ttl
	})

	p.mapCleaner = mapCleaner
}

func (p *protocol) setupMapCleaner2(mgr *manager.Manager) {
	http2TerminatedConnectionsQueue, _, err := mgr.GetMap("http2_terminated_conns")
	if err != nil {
		log.Errorf("error getting http2_terminated_conns map: %s", err)
		return
	}
	dynamicTableMap, _, err := mgr.GetMap(dynamicTable)
	if err != nil {
		log.Errorf("error getting %q map: %s", dynamicTable, err)
		return
	}

	keys := make([]http2DynamicTableIndex, dynamicTableMap.MaxEntries())
	values := make([]http2DynamicTableEntry, dynamicTableMap.MaxEntries())
	keysForDeletion := make([]http2DynamicTableIndex, 0, dynamicTableMap.MaxEntries())
	go func() {
		for {
			// Poor people periodic
			time.Sleep(10 * time.Second)

			keys := keys[:0]
			next := uint32(0)
			// Pulling all keys from the map
			n, err := dynamicTableMap.BatchLookup(nil, &next, keys, values, nil)
			if n == 0 {
				continue
			}

			// While queue is not empty
			for {
				var value netebpf.ConnTuple
				// Pop from the queue
				if err := http2TerminatedConnectionsQueue.LookupAndDelete(nil, unsafe.Pointer(&value)); err != nil {
					log.Errorf("error lookup and delete: %s", err)
					break
				}

				// Search for matching entries in the dynamic table
				// matching entries are those that have the same connection tuple, regardless of the index value.
				for i := 0; i < n; i++ {
					if keys[i].Tup == value {
						// If found, add the index to the list of keys to be deleted
						keysForDeletion = append(keysForDeletion, keys[i])
						// Swap the last element with the current one so we can shrink the slice and avoid re-checking
						// the same (deleted) element again.
						keys[i] = keys[n-1]
						n--
						break
					}
				}
			}
			// Delete all matching entries
			if _, err = dynamicTableMap.BatchDelete(keysForDeletion, nil); err != nil {
				log.Errorf("error batch delete: %s", err)
				break
			}
		}
	}()
}

// GetStats returns a map of HTTP2 stats stored in the following format:
// [source, dest tuple, request path] -> RequestStats object
func (p *protocol) GetStats() *protocols.ProtocolStats {
	p.eventsConsumer.Sync()
	p.telemetry.Log()
	return &protocols.ProtocolStats{
		Type:  protocols.HTTP2,
		Stats: p.statkeeper.GetAndResetAllStats(),
	}
}

// The staticTableEntry represents an entry in the static table that contains an index in the table and a value.
// The value itself contains both the key and the corresponding value in the static table.
// For instance, index 2 in the static table has a value of method: GET, and index 3 has a value of method: POST.
// It is not possible to save the index by the key because we need to distinguish between the values attached to the key.
type staticTableEntry struct {
	Index uint64
	Value StaticTableEnumValue
}

var (
	staticTableEntries = []staticTableEntry{
		{
			Index: 2,
			Value: GetValue,
		},
		{
			Index: 3,
			Value: PostValue,
		},
		{
			Index: 4,
			Value: EmptyPathValue,
		},
		{
			Index: 5,
			Value: IndexPathValue,
		},
		{
			Index: 8,
			Value: K200Value,
		},
		{
			Index: 9,
			Value: K204Value,
		},
		{
			Index: 10,
			Value: K206Value,
		},
		{
			Index: 11,
			Value: K304Value,
		},
		{
			Index: 12,
			Value: K400Value,
		},
		{
			Index: 13,
			Value: K404Value,
		},
		{
			Index: 14,
			Value: K500Value,
		},
	}
)

// createStaticTable creates a static table for http2 monitor.
func (p *protocol) createStaticTable(mgr *manager.Manager) error {
	staticTable, _, _ := mgr.GetMap(probes.StaticTableMap)
	if staticTable == nil {
		return errors.New("http2 static table is null")
	}

	for _, entry := range staticTableEntries {
		err := staticTable.Put(unsafe.Pointer(&entry.Index), unsafe.Pointer(&entry.Value))

		if err != nil {
			return err
		}
	}
	return nil
}

// IsBuildModeSupported returns always true, as http2 module is supported by all modes.
func (*protocol) IsBuildModeSupported(buildmode.Type) bool {
	return true
}
