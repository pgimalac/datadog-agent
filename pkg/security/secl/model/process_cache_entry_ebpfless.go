// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build unix && ebpfless

// Package model holds model related files
package model

import (
	"time"
)

// SetAncestor sets the ancestor
func (pc *ProcessCacheEntry) SetAncestor(parent *ProcessCacheEntry) {
	if pc.Ancestor == parent {
		return
	}

	if pc.Ancestor != nil {
		pc.Ancestor.Release()
	}

	pc.Ancestor = parent
	pc.Parent = &parent.Process
	parent.Retain()
}

// Exit a process
func (pc *ProcessCacheEntry) Exit(exitTime time.Time) {
	pc.ExitTime = exitTime
}

func copyProcessContext(parent, child *ProcessCacheEntry) {
	// inherit the container ID from the parent if necessary. If a container is already running when system-probe
	// starts, the in-kernel process cache will have out of sync container ID values for the processes of that
	// container (the snapshot doesn't update the in-kernel cache with the container IDs). This can also happen if
	// the proc_cache LRU ejects an entry.
	// WARNING: this is why the user space cache should not be used to detect container breakouts. Dedicated
	// in-kernel probes will need to be added.
	if len(parent.ContainerID) > 0 && len(child.ContainerID) == 0 {
		child.ContainerID = parent.ContainerID
	}
}

// Exec replace a process
func (pc *ProcessCacheEntry) Exec(entry *ProcessCacheEntry) {
	entry.SetAncestor(pc)

	// use exec time as exit time
	pc.Exit(entry.ExecTime)

	// keep some context
	copyProcessContext(pc, entry)
}

// SetParentOfForkChild set the parent of a fork child
func (pc *ProcessCacheEntry) SetParentOfForkChild(parent *ProcessCacheEntry) {
	pc.SetAncestor(parent)
	if parent != nil {
		pc.ArgsEntry = parent.ArgsEntry
		pc.EnvsEntry = parent.EnvsEntry
	}
}

// Fork returns a copy of the current ProcessCacheEntry
func (pc *ProcessCacheEntry) Fork(childEntry *ProcessCacheEntry) {
	childEntry.PPid = pc.Pid
	childEntry.FileEvent = pc.FileEvent
	childEntry.ContainerID = pc.ContainerID
	childEntry.ExecTime = pc.ExecTime

	childEntry.SetParentOfForkChild(pc)
}

// Equals returns whether process cache entries share the same values for file and args/envs
func (pc *ProcessCacheEntry) Equals(entry *ProcessCacheEntry) bool {
	return (pc.FileEvent.Equals(&entry.FileEvent) &&
		pc.ArgsEntry.Equals(entry.ArgsEntry) &&
		pc.EnvsEntry.Equals(entry.EnvsEntry))
}

// NewPlaceholderProcessCacheEntry returns a new empty process cache entry for failed process resolutions
func NewPlaceholderProcessCacheEntry(pid uint32) *ProcessCacheEntry {
	return &ProcessCacheEntry{
		ProcessContext: ProcessContext{
			Process: Process{
				PIDContext: PIDContext{Pid: pid},
			},
		},
	}
}

var processContextZero = ProcessCacheEntry{}

// GetPlaceholderProcessCacheEntry returns an empty process cache entry for failed process resolutions
func GetPlaceholderProcessCacheEntry(pid uint32) *ProcessCacheEntry {
	processContextZero.Pid = pid
	return &processContextZero
}