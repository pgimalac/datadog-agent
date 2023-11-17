// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2022-present Datadog, Inc.
// Code generated - DO NOT EDIT.

//go:build unix && ebpfless

package model

import (
	"github.com/DataDog/datadog-agent/pkg/security/secl/compiler/eval"
	"reflect"
)

func (m *Model) GetIterator(field eval.Field) (eval.Iterator, error) {
	switch field {
	case "process.ancestors":
		return &ProcessAncestorsIterator{}, nil
	}
	return nil, &eval.ErrIteratorNotSupported{Field: field}
}
func (m *Model) GetEventTypes() []eval.EventType {
	return []eval.EventType{
		eval.EventType("exec"),
		eval.EventType("exit"),
		eval.EventType("open"),
	}
}
func (m *Model) GetEvaluator(field eval.Field, regID eval.RegisterID) (eval.Evaluator, error) {
	switch field {
	case "container.created_at":
		return &eval.IntEvaluator{
			EvalFnc: func(ctx *eval.Context) int {
				ev := ctx.Event.(*Event)
				return int(ev.FieldHandlers.ResolveContainerCreatedAt(ev, ev.BaseEvent.ContainerContext))
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "container.id":
		return &eval.StringEvaluator{
			EvalFnc: func(ctx *eval.Context) string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveContainerID(ev, ev.BaseEvent.ContainerContext)
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "container.tags":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveContainerTags(ev, ev.BaseEvent.ContainerContext)
			},
			Field:  field,
			Weight: 9999 * eval.HandlerWeight,
		}, nil
	case "event.timestamp":
		return &eval.IntEvaluator{
			EvalFnc: func(ctx *eval.Context) int {
				ev := ctx.Event.(*Event)
				return int(ev.FieldHandlers.ResolveEventTimestamp(ev, &ev.BaseEvent))
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "exec.args":
		return &eval.StringEvaluator{
			EvalFnc: func(ctx *eval.Context) string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveProcessArgs(ev, ev.Exec.Process)
			},
			Field:  field,
			Weight: 100 * eval.HandlerWeight,
		}, nil
	case "exec.args_flags":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveProcessArgsFlags(ev, ev.Exec.Process)
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "exec.args_options":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveProcessArgsOptions(ev, ev.Exec.Process)
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "exec.args_truncated":
		return &eval.BoolEvaluator{
			EvalFnc: func(ctx *eval.Context) bool {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveProcessArgsTruncated(ev, ev.Exec.Process)
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "exec.argv":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveProcessArgv(ev, ev.Exec.Process)
			},
			Field:  field,
			Weight: 100 * eval.HandlerWeight,
		}, nil
	case "exec.argv0":
		return &eval.StringEvaluator{
			EvalFnc: func(ctx *eval.Context) string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveProcessArgv0(ev, ev.Exec.Process)
			},
			Field:  field,
			Weight: 100 * eval.HandlerWeight,
		}, nil
	case "exec.comm":
		return &eval.StringEvaluator{
			EvalFnc: func(ctx *eval.Context) string {
				ev := ctx.Event.(*Event)
				return ev.Exec.Process.Comm
			},
			Field:  field,
			Weight: eval.FunctionWeight,
		}, nil
	case "exec.container.id":
		return &eval.StringEvaluator{
			EvalFnc: func(ctx *eval.Context) string {
				ev := ctx.Event.(*Event)
				return ev.Exec.Process.ContainerID
			},
			Field:  field,
			Weight: eval.FunctionWeight,
		}, nil
	case "exec.created_at":
		return &eval.IntEvaluator{
			EvalFnc: func(ctx *eval.Context) int {
				ev := ctx.Event.(*Event)
				return int(ev.FieldHandlers.ResolveProcessCreatedAt(ev, ev.Exec.Process))
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "exec.envp":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveProcessEnvp(ev, ev.Exec.Process)
			},
			Field:  field,
			Weight: 100 * eval.HandlerWeight,
		}, nil
	case "exec.envs":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveProcessEnvs(ev, ev.Exec.Process)
			},
			Field:  field,
			Weight: 100 * eval.HandlerWeight,
		}, nil
	case "exec.envs_truncated":
		return &eval.BoolEvaluator{
			EvalFnc: func(ctx *eval.Context) bool {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveProcessEnvsTruncated(ev, ev.Exec.Process)
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "exec.file.name":
		return &eval.StringEvaluator{
			EvalFnc: func(ctx *eval.Context) string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveFileBasename(ev, &ev.Exec.Process.FileEvent)
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "exec.file.name.length":
		return &eval.IntEvaluator{
			EvalFnc: func(ctx *eval.Context) int {
				ev := ctx.Event.(*Event)
				return len(ev.FieldHandlers.ResolveFileBasename(ev, &ev.Exec.Process.FileEvent))
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "exec.file.path":
		return &eval.StringEvaluator{
			EvalFnc: func(ctx *eval.Context) string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveFilePath(ev, &ev.Exec.Process.FileEvent)
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "exec.file.path.length":
		return &eval.IntEvaluator{
			EvalFnc: func(ctx *eval.Context) int {
				ev := ctx.Event.(*Event)
				return len(ev.FieldHandlers.ResolveFilePath(ev, &ev.Exec.Process.FileEvent))
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "exec.pid":
		return &eval.IntEvaluator{
			EvalFnc: func(ctx *eval.Context) int {
				ev := ctx.Event.(*Event)
				return int(ev.Exec.Process.PIDContext.Pid)
			},
			Field:  field,
			Weight: eval.FunctionWeight,
		}, nil
	case "exec.ppid":
		return &eval.IntEvaluator{
			EvalFnc: func(ctx *eval.Context) int {
				ev := ctx.Event.(*Event)
				return int(ev.Exec.Process.PPid)
			},
			Field:  field,
			Weight: eval.FunctionWeight,
		}, nil
	case "exit.args":
		return &eval.StringEvaluator{
			EvalFnc: func(ctx *eval.Context) string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveProcessArgs(ev, ev.Exit.Process)
			},
			Field:  field,
			Weight: 100 * eval.HandlerWeight,
		}, nil
	case "exit.args_flags":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveProcessArgsFlags(ev, ev.Exit.Process)
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "exit.args_options":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveProcessArgsOptions(ev, ev.Exit.Process)
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "exit.args_truncated":
		return &eval.BoolEvaluator{
			EvalFnc: func(ctx *eval.Context) bool {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveProcessArgsTruncated(ev, ev.Exit.Process)
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "exit.argv":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveProcessArgv(ev, ev.Exit.Process)
			},
			Field:  field,
			Weight: 100 * eval.HandlerWeight,
		}, nil
	case "exit.argv0":
		return &eval.StringEvaluator{
			EvalFnc: func(ctx *eval.Context) string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveProcessArgv0(ev, ev.Exit.Process)
			},
			Field:  field,
			Weight: 100 * eval.HandlerWeight,
		}, nil
	case "exit.cause":
		return &eval.IntEvaluator{
			EvalFnc: func(ctx *eval.Context) int {
				ev := ctx.Event.(*Event)
				return int(ev.Exit.Cause)
			},
			Field:  field,
			Weight: eval.FunctionWeight,
		}, nil
	case "exit.code":
		return &eval.IntEvaluator{
			EvalFnc: func(ctx *eval.Context) int {
				ev := ctx.Event.(*Event)
				return int(ev.Exit.Code)
			},
			Field:  field,
			Weight: eval.FunctionWeight,
		}, nil
	case "exit.comm":
		return &eval.StringEvaluator{
			EvalFnc: func(ctx *eval.Context) string {
				ev := ctx.Event.(*Event)
				return ev.Exit.Process.Comm
			},
			Field:  field,
			Weight: eval.FunctionWeight,
		}, nil
	case "exit.container.id":
		return &eval.StringEvaluator{
			EvalFnc: func(ctx *eval.Context) string {
				ev := ctx.Event.(*Event)
				return ev.Exit.Process.ContainerID
			},
			Field:  field,
			Weight: eval.FunctionWeight,
		}, nil
	case "exit.created_at":
		return &eval.IntEvaluator{
			EvalFnc: func(ctx *eval.Context) int {
				ev := ctx.Event.(*Event)
				return int(ev.FieldHandlers.ResolveProcessCreatedAt(ev, ev.Exit.Process))
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "exit.envp":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveProcessEnvp(ev, ev.Exit.Process)
			},
			Field:  field,
			Weight: 100 * eval.HandlerWeight,
		}, nil
	case "exit.envs":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveProcessEnvs(ev, ev.Exit.Process)
			},
			Field:  field,
			Weight: 100 * eval.HandlerWeight,
		}, nil
	case "exit.envs_truncated":
		return &eval.BoolEvaluator{
			EvalFnc: func(ctx *eval.Context) bool {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveProcessEnvsTruncated(ev, ev.Exit.Process)
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "exit.file.name":
		return &eval.StringEvaluator{
			EvalFnc: func(ctx *eval.Context) string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveFileBasename(ev, &ev.Exit.Process.FileEvent)
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "exit.file.name.length":
		return &eval.IntEvaluator{
			EvalFnc: func(ctx *eval.Context) int {
				ev := ctx.Event.(*Event)
				return len(ev.FieldHandlers.ResolveFileBasename(ev, &ev.Exit.Process.FileEvent))
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "exit.file.path":
		return &eval.StringEvaluator{
			EvalFnc: func(ctx *eval.Context) string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveFilePath(ev, &ev.Exit.Process.FileEvent)
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "exit.file.path.length":
		return &eval.IntEvaluator{
			EvalFnc: func(ctx *eval.Context) int {
				ev := ctx.Event.(*Event)
				return len(ev.FieldHandlers.ResolveFilePath(ev, &ev.Exit.Process.FileEvent))
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "exit.pid":
		return &eval.IntEvaluator{
			EvalFnc: func(ctx *eval.Context) int {
				ev := ctx.Event.(*Event)
				return int(ev.Exit.Process.PIDContext.Pid)
			},
			Field:  field,
			Weight: eval.FunctionWeight,
		}, nil
	case "exit.ppid":
		return &eval.IntEvaluator{
			EvalFnc: func(ctx *eval.Context) int {
				ev := ctx.Event.(*Event)
				return int(ev.Exit.Process.PPid)
			},
			Field:  field,
			Weight: eval.FunctionWeight,
		}, nil
	case "open.file.destination.mode":
		return &eval.IntEvaluator{
			EvalFnc: func(ctx *eval.Context) int {
				ev := ctx.Event.(*Event)
				return int(ev.Open.Mode)
			},
			Field:  field,
			Weight: eval.FunctionWeight,
		}, nil
	case "open.file.name":
		return &eval.StringEvaluator{
			EvalFnc: func(ctx *eval.Context) string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveFileBasename(ev, &ev.Open.File)
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "open.file.name.length":
		return &eval.IntEvaluator{
			EvalFnc: func(ctx *eval.Context) int {
				ev := ctx.Event.(*Event)
				return len(ev.FieldHandlers.ResolveFileBasename(ev, &ev.Open.File))
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "open.file.path":
		return &eval.StringEvaluator{
			EvalFnc: func(ctx *eval.Context) string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveFilePath(ev, &ev.Open.File)
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "open.file.path.length":
		return &eval.IntEvaluator{
			EvalFnc: func(ctx *eval.Context) int {
				ev := ctx.Event.(*Event)
				return len(ev.FieldHandlers.ResolveFilePath(ev, &ev.Open.File))
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "open.flags":
		return &eval.IntEvaluator{
			EvalFnc: func(ctx *eval.Context) int {
				ev := ctx.Event.(*Event)
				return int(ev.Open.Flags)
			},
			Field:  field,
			Weight: eval.FunctionWeight,
		}, nil
	case "process.ancestors.args":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				if result, ok := ctx.StringCache[field]; ok {
					return result
				}
				var results []string
				iterator := &ProcessAncestorsIterator{}
				value := iterator.Front(ctx)
				for value != nil {
					element := (*ProcessCacheEntry)(value)
					result := ev.FieldHandlers.ResolveProcessArgs(ev, &element.ProcessContext.Process)
					results = append(results, result)
					value = iterator.Next()
				}
				ctx.StringCache[field] = results
				return results
			}, Field: field,
			Weight: 100 * eval.IteratorWeight,
		}, nil
	case "process.ancestors.args_flags":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				if result, ok := ctx.StringCache[field]; ok {
					return result
				}
				var results []string
				iterator := &ProcessAncestorsIterator{}
				value := iterator.Front(ctx)
				for value != nil {
					element := (*ProcessCacheEntry)(value)
					result := ev.FieldHandlers.ResolveProcessArgsFlags(ev, &element.ProcessContext.Process)
					results = append(results, result...)
					value = iterator.Next()
				}
				ctx.StringCache[field] = results
				return results
			}, Field: field,
			Weight: eval.IteratorWeight,
		}, nil
	case "process.ancestors.args_options":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				if result, ok := ctx.StringCache[field]; ok {
					return result
				}
				var results []string
				iterator := &ProcessAncestorsIterator{}
				value := iterator.Front(ctx)
				for value != nil {
					element := (*ProcessCacheEntry)(value)
					result := ev.FieldHandlers.ResolveProcessArgsOptions(ev, &element.ProcessContext.Process)
					results = append(results, result...)
					value = iterator.Next()
				}
				ctx.StringCache[field] = results
				return results
			}, Field: field,
			Weight: eval.IteratorWeight,
		}, nil
	case "process.ancestors.args_truncated":
		return &eval.BoolArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []bool {
				ev := ctx.Event.(*Event)
				if result, ok := ctx.BoolCache[field]; ok {
					return result
				}
				var results []bool
				iterator := &ProcessAncestorsIterator{}
				value := iterator.Front(ctx)
				for value != nil {
					element := (*ProcessCacheEntry)(value)
					result := ev.FieldHandlers.ResolveProcessArgsTruncated(ev, &element.ProcessContext.Process)
					results = append(results, result)
					value = iterator.Next()
				}
				ctx.BoolCache[field] = results
				return results
			}, Field: field,
			Weight: eval.IteratorWeight,
		}, nil
	case "process.ancestors.argv":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				if result, ok := ctx.StringCache[field]; ok {
					return result
				}
				var results []string
				iterator := &ProcessAncestorsIterator{}
				value := iterator.Front(ctx)
				for value != nil {
					element := (*ProcessCacheEntry)(value)
					result := ev.FieldHandlers.ResolveProcessArgv(ev, &element.ProcessContext.Process)
					results = append(results, result...)
					value = iterator.Next()
				}
				ctx.StringCache[field] = results
				return results
			}, Field: field,
			Weight: 100 * eval.IteratorWeight,
		}, nil
	case "process.ancestors.argv0":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				if result, ok := ctx.StringCache[field]; ok {
					return result
				}
				var results []string
				iterator := &ProcessAncestorsIterator{}
				value := iterator.Front(ctx)
				for value != nil {
					element := (*ProcessCacheEntry)(value)
					result := ev.FieldHandlers.ResolveProcessArgv0(ev, &element.ProcessContext.Process)
					results = append(results, result)
					value = iterator.Next()
				}
				ctx.StringCache[field] = results
				return results
			}, Field: field,
			Weight: 100 * eval.IteratorWeight,
		}, nil
	case "process.ancestors.comm":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				if result, ok := ctx.StringCache[field]; ok {
					return result
				}
				var results []string
				iterator := &ProcessAncestorsIterator{}
				value := iterator.Front(ctx)
				for value != nil {
					element := (*ProcessCacheEntry)(value)
					result := element.ProcessContext.Process.Comm
					results = append(results, result)
					value = iterator.Next()
				}
				ctx.StringCache[field] = results
				return results
			}, Field: field,
			Weight: eval.IteratorWeight,
		}, nil
	case "process.ancestors.container.id":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				if result, ok := ctx.StringCache[field]; ok {
					return result
				}
				var results []string
				iterator := &ProcessAncestorsIterator{}
				value := iterator.Front(ctx)
				for value != nil {
					element := (*ProcessCacheEntry)(value)
					result := element.ProcessContext.Process.ContainerID
					results = append(results, result)
					value = iterator.Next()
				}
				ctx.StringCache[field] = results
				return results
			}, Field: field,
			Weight: eval.IteratorWeight,
		}, nil
	case "process.ancestors.created_at":
		return &eval.IntArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []int {
				ev := ctx.Event.(*Event)
				if result, ok := ctx.IntCache[field]; ok {
					return result
				}
				var results []int
				iterator := &ProcessAncestorsIterator{}
				value := iterator.Front(ctx)
				for value != nil {
					element := (*ProcessCacheEntry)(value)
					result := int(ev.FieldHandlers.ResolveProcessCreatedAt(ev, &element.ProcessContext.Process))
					results = append(results, result)
					value = iterator.Next()
				}
				ctx.IntCache[field] = results
				return results
			}, Field: field,
			Weight: eval.IteratorWeight,
		}, nil
	case "process.ancestors.envp":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				if result, ok := ctx.StringCache[field]; ok {
					return result
				}
				var results []string
				iterator := &ProcessAncestorsIterator{}
				value := iterator.Front(ctx)
				for value != nil {
					element := (*ProcessCacheEntry)(value)
					result := ev.FieldHandlers.ResolveProcessEnvp(ev, &element.ProcessContext.Process)
					results = append(results, result...)
					value = iterator.Next()
				}
				ctx.StringCache[field] = results
				return results
			}, Field: field,
			Weight: 100 * eval.IteratorWeight,
		}, nil
	case "process.ancestors.envs":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				if result, ok := ctx.StringCache[field]; ok {
					return result
				}
				var results []string
				iterator := &ProcessAncestorsIterator{}
				value := iterator.Front(ctx)
				for value != nil {
					element := (*ProcessCacheEntry)(value)
					result := ev.FieldHandlers.ResolveProcessEnvs(ev, &element.ProcessContext.Process)
					results = append(results, result...)
					value = iterator.Next()
				}
				ctx.StringCache[field] = results
				return results
			}, Field: field,
			Weight: 100 * eval.IteratorWeight,
		}, nil
	case "process.ancestors.envs_truncated":
		return &eval.BoolArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []bool {
				ev := ctx.Event.(*Event)
				if result, ok := ctx.BoolCache[field]; ok {
					return result
				}
				var results []bool
				iterator := &ProcessAncestorsIterator{}
				value := iterator.Front(ctx)
				for value != nil {
					element := (*ProcessCacheEntry)(value)
					result := ev.FieldHandlers.ResolveProcessEnvsTruncated(ev, &element.ProcessContext.Process)
					results = append(results, result)
					value = iterator.Next()
				}
				ctx.BoolCache[field] = results
				return results
			}, Field: field,
			Weight: eval.IteratorWeight,
		}, nil
	case "process.ancestors.file.name":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				if result, ok := ctx.StringCache[field]; ok {
					return result
				}
				var results []string
				iterator := &ProcessAncestorsIterator{}
				value := iterator.Front(ctx)
				for value != nil {
					element := (*ProcessCacheEntry)(value)
					result := ev.FieldHandlers.ResolveFileBasename(ev, &element.ProcessContext.Process.FileEvent)
					results = append(results, result)
					value = iterator.Next()
				}
				ctx.StringCache[field] = results
				return results
			}, Field: field,
			Weight: eval.IteratorWeight,
		}, nil
	case "process.ancestors.file.name.length":
		return &eval.IntArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []int {
				ev := ctx.Event.(*Event)
				if result, ok := ctx.IntCache[field]; ok {
					return result
				}
				var results []int
				iterator := &ProcessAncestorsIterator{}
				value := iterator.Front(ctx)
				for value != nil {
					element := (*ProcessCacheEntry)(value)
					result := len(ev.FieldHandlers.ResolveFileBasename(ev, &element.ProcessContext.Process.FileEvent))
					results = append(results, result)
					value = iterator.Next()
				}
				ctx.IntCache[field] = results
				return results
			}, Field: field,
			Weight: eval.IteratorWeight,
		}, nil
	case "process.ancestors.file.path":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				if result, ok := ctx.StringCache[field]; ok {
					return result
				}
				var results []string
				iterator := &ProcessAncestorsIterator{}
				value := iterator.Front(ctx)
				for value != nil {
					element := (*ProcessCacheEntry)(value)
					result := ev.FieldHandlers.ResolveFilePath(ev, &element.ProcessContext.Process.FileEvent)
					results = append(results, result)
					value = iterator.Next()
				}
				ctx.StringCache[field] = results
				return results
			}, Field: field,
			Weight: eval.IteratorWeight,
		}, nil
	case "process.ancestors.file.path.length":
		return &eval.IntArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []int {
				ev := ctx.Event.(*Event)
				if result, ok := ctx.IntCache[field]; ok {
					return result
				}
				var results []int
				iterator := &ProcessAncestorsIterator{}
				value := iterator.Front(ctx)
				for value != nil {
					element := (*ProcessCacheEntry)(value)
					result := len(ev.FieldHandlers.ResolveFilePath(ev, &element.ProcessContext.Process.FileEvent))
					results = append(results, result)
					value = iterator.Next()
				}
				ctx.IntCache[field] = results
				return results
			}, Field: field,
			Weight: eval.IteratorWeight,
		}, nil
	case "process.ancestors.pid":
		return &eval.IntArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []int {
				if result, ok := ctx.IntCache[field]; ok {
					return result
				}
				var results []int
				iterator := &ProcessAncestorsIterator{}
				value := iterator.Front(ctx)
				for value != nil {
					element := (*ProcessCacheEntry)(value)
					result := int(element.ProcessContext.Process.PIDContext.Pid)
					results = append(results, result)
					value = iterator.Next()
				}
				ctx.IntCache[field] = results
				return results
			}, Field: field,
			Weight: eval.IteratorWeight,
		}, nil
	case "process.ancestors.ppid":
		return &eval.IntArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []int {
				if result, ok := ctx.IntCache[field]; ok {
					return result
				}
				var results []int
				iterator := &ProcessAncestorsIterator{}
				value := iterator.Front(ctx)
				for value != nil {
					element := (*ProcessCacheEntry)(value)
					result := int(element.ProcessContext.Process.PPid)
					results = append(results, result)
					value = iterator.Next()
				}
				ctx.IntCache[field] = results
				return results
			}, Field: field,
			Weight: eval.IteratorWeight,
		}, nil
	case "process.args":
		return &eval.StringEvaluator{
			EvalFnc: func(ctx *eval.Context) string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveProcessArgs(ev, &ev.BaseEvent.ProcessContext.Process)
			},
			Field:  field,
			Weight: 100 * eval.HandlerWeight,
		}, nil
	case "process.args_flags":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveProcessArgsFlags(ev, &ev.BaseEvent.ProcessContext.Process)
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "process.args_options":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveProcessArgsOptions(ev, &ev.BaseEvent.ProcessContext.Process)
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "process.args_truncated":
		return &eval.BoolEvaluator{
			EvalFnc: func(ctx *eval.Context) bool {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveProcessArgsTruncated(ev, &ev.BaseEvent.ProcessContext.Process)
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "process.argv":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveProcessArgv(ev, &ev.BaseEvent.ProcessContext.Process)
			},
			Field:  field,
			Weight: 100 * eval.HandlerWeight,
		}, nil
	case "process.argv0":
		return &eval.StringEvaluator{
			EvalFnc: func(ctx *eval.Context) string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveProcessArgv0(ev, &ev.BaseEvent.ProcessContext.Process)
			},
			Field:  field,
			Weight: 100 * eval.HandlerWeight,
		}, nil
	case "process.comm":
		return &eval.StringEvaluator{
			EvalFnc: func(ctx *eval.Context) string {
				ev := ctx.Event.(*Event)
				return ev.BaseEvent.ProcessContext.Process.Comm
			},
			Field:  field,
			Weight: eval.FunctionWeight,
		}, nil
	case "process.container.id":
		return &eval.StringEvaluator{
			EvalFnc: func(ctx *eval.Context) string {
				ev := ctx.Event.(*Event)
				return ev.BaseEvent.ProcessContext.Process.ContainerID
			},
			Field:  field,
			Weight: eval.FunctionWeight,
		}, nil
	case "process.created_at":
		return &eval.IntEvaluator{
			EvalFnc: func(ctx *eval.Context) int {
				ev := ctx.Event.(*Event)
				return int(ev.FieldHandlers.ResolveProcessCreatedAt(ev, &ev.BaseEvent.ProcessContext.Process))
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "process.envp":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveProcessEnvp(ev, &ev.BaseEvent.ProcessContext.Process)
			},
			Field:  field,
			Weight: 100 * eval.HandlerWeight,
		}, nil
	case "process.envs":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveProcessEnvs(ev, &ev.BaseEvent.ProcessContext.Process)
			},
			Field:  field,
			Weight: 100 * eval.HandlerWeight,
		}, nil
	case "process.envs_truncated":
		return &eval.BoolEvaluator{
			EvalFnc: func(ctx *eval.Context) bool {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveProcessEnvsTruncated(ev, &ev.BaseEvent.ProcessContext.Process)
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "process.file.name":
		return &eval.StringEvaluator{
			EvalFnc: func(ctx *eval.Context) string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveFileBasename(ev, &ev.BaseEvent.ProcessContext.Process.FileEvent)
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "process.file.name.length":
		return &eval.IntEvaluator{
			EvalFnc: func(ctx *eval.Context) int {
				ev := ctx.Event.(*Event)
				return len(ev.FieldHandlers.ResolveFileBasename(ev, &ev.BaseEvent.ProcessContext.Process.FileEvent))
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "process.file.path":
		return &eval.StringEvaluator{
			EvalFnc: func(ctx *eval.Context) string {
				ev := ctx.Event.(*Event)
				return ev.FieldHandlers.ResolveFilePath(ev, &ev.BaseEvent.ProcessContext.Process.FileEvent)
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "process.file.path.length":
		return &eval.IntEvaluator{
			EvalFnc: func(ctx *eval.Context) int {
				ev := ctx.Event.(*Event)
				return len(ev.FieldHandlers.ResolveFilePath(ev, &ev.BaseEvent.ProcessContext.Process.FileEvent))
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "process.parent.args":
		return &eval.StringEvaluator{
			EvalFnc: func(ctx *eval.Context) string {
				ev := ctx.Event.(*Event)
				if !ev.BaseEvent.ProcessContext.HasParent() {
					return ""
				}
				return ev.FieldHandlers.ResolveProcessArgs(ev, ev.BaseEvent.ProcessContext.Parent)
			},
			Field:  field,
			Weight: 100 * eval.HandlerWeight,
		}, nil
	case "process.parent.args_flags":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				if !ev.BaseEvent.ProcessContext.HasParent() {
					return []string{}
				}
				return ev.FieldHandlers.ResolveProcessArgsFlags(ev, ev.BaseEvent.ProcessContext.Parent)
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "process.parent.args_options":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				if !ev.BaseEvent.ProcessContext.HasParent() {
					return []string{}
				}
				return ev.FieldHandlers.ResolveProcessArgsOptions(ev, ev.BaseEvent.ProcessContext.Parent)
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "process.parent.args_truncated":
		return &eval.BoolEvaluator{
			EvalFnc: func(ctx *eval.Context) bool {
				ev := ctx.Event.(*Event)
				if !ev.BaseEvent.ProcessContext.HasParent() {
					return false
				}
				return ev.FieldHandlers.ResolveProcessArgsTruncated(ev, ev.BaseEvent.ProcessContext.Parent)
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "process.parent.argv":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				if !ev.BaseEvent.ProcessContext.HasParent() {
					return []string{}
				}
				return ev.FieldHandlers.ResolveProcessArgv(ev, ev.BaseEvent.ProcessContext.Parent)
			},
			Field:  field,
			Weight: 100 * eval.HandlerWeight,
		}, nil
	case "process.parent.argv0":
		return &eval.StringEvaluator{
			EvalFnc: func(ctx *eval.Context) string {
				ev := ctx.Event.(*Event)
				if !ev.BaseEvent.ProcessContext.HasParent() {
					return ""
				}
				return ev.FieldHandlers.ResolveProcessArgv0(ev, ev.BaseEvent.ProcessContext.Parent)
			},
			Field:  field,
			Weight: 100 * eval.HandlerWeight,
		}, nil
	case "process.parent.comm":
		return &eval.StringEvaluator{
			EvalFnc: func(ctx *eval.Context) string {
				ev := ctx.Event.(*Event)
				if !ev.BaseEvent.ProcessContext.HasParent() {
					return ""
				}
				return ev.BaseEvent.ProcessContext.Parent.Comm
			},
			Field:  field,
			Weight: eval.FunctionWeight,
		}, nil
	case "process.parent.container.id":
		return &eval.StringEvaluator{
			EvalFnc: func(ctx *eval.Context) string {
				ev := ctx.Event.(*Event)
				if !ev.BaseEvent.ProcessContext.HasParent() {
					return ""
				}
				return ev.BaseEvent.ProcessContext.Parent.ContainerID
			},
			Field:  field,
			Weight: eval.FunctionWeight,
		}, nil
	case "process.parent.created_at":
		return &eval.IntEvaluator{
			EvalFnc: func(ctx *eval.Context) int {
				ev := ctx.Event.(*Event)
				if !ev.BaseEvent.ProcessContext.HasParent() {
					return 0
				}
				return int(ev.FieldHandlers.ResolveProcessCreatedAt(ev, ev.BaseEvent.ProcessContext.Parent))
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "process.parent.envp":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				if !ev.BaseEvent.ProcessContext.HasParent() {
					return []string{}
				}
				return ev.FieldHandlers.ResolveProcessEnvp(ev, ev.BaseEvent.ProcessContext.Parent)
			},
			Field:  field,
			Weight: 100 * eval.HandlerWeight,
		}, nil
	case "process.parent.envs":
		return &eval.StringArrayEvaluator{
			EvalFnc: func(ctx *eval.Context) []string {
				ev := ctx.Event.(*Event)
				if !ev.BaseEvent.ProcessContext.HasParent() {
					return []string{}
				}
				return ev.FieldHandlers.ResolveProcessEnvs(ev, ev.BaseEvent.ProcessContext.Parent)
			},
			Field:  field,
			Weight: 100 * eval.HandlerWeight,
		}, nil
	case "process.parent.envs_truncated":
		return &eval.BoolEvaluator{
			EvalFnc: func(ctx *eval.Context) bool {
				ev := ctx.Event.(*Event)
				if !ev.BaseEvent.ProcessContext.HasParent() {
					return false
				}
				return ev.FieldHandlers.ResolveProcessEnvsTruncated(ev, ev.BaseEvent.ProcessContext.Parent)
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "process.parent.file.name":
		return &eval.StringEvaluator{
			EvalFnc: func(ctx *eval.Context) string {
				ev := ctx.Event.(*Event)
				if !ev.BaseEvent.ProcessContext.HasParent() {
					return ""
				}
				return ev.FieldHandlers.ResolveFileBasename(ev, &ev.BaseEvent.ProcessContext.Parent.FileEvent)
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "process.parent.file.name.length":
		return &eval.IntEvaluator{
			EvalFnc: func(ctx *eval.Context) int {
				ev := ctx.Event.(*Event)
				return len(ev.FieldHandlers.ResolveFileBasename(ev, &ev.BaseEvent.ProcessContext.Parent.FileEvent))
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "process.parent.file.path":
		return &eval.StringEvaluator{
			EvalFnc: func(ctx *eval.Context) string {
				ev := ctx.Event.(*Event)
				if !ev.BaseEvent.ProcessContext.HasParent() {
					return ""
				}
				return ev.FieldHandlers.ResolveFilePath(ev, &ev.BaseEvent.ProcessContext.Parent.FileEvent)
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "process.parent.file.path.length":
		return &eval.IntEvaluator{
			EvalFnc: func(ctx *eval.Context) int {
				ev := ctx.Event.(*Event)
				return len(ev.FieldHandlers.ResolveFilePath(ev, &ev.BaseEvent.ProcessContext.Parent.FileEvent))
			},
			Field:  field,
			Weight: eval.HandlerWeight,
		}, nil
	case "process.parent.pid":
		return &eval.IntEvaluator{
			EvalFnc: func(ctx *eval.Context) int {
				ev := ctx.Event.(*Event)
				if !ev.BaseEvent.ProcessContext.HasParent() {
					return 0
				}
				return int(ev.BaseEvent.ProcessContext.Parent.PIDContext.Pid)
			},
			Field:  field,
			Weight: eval.FunctionWeight,
		}, nil
	case "process.parent.ppid":
		return &eval.IntEvaluator{
			EvalFnc: func(ctx *eval.Context) int {
				ev := ctx.Event.(*Event)
				if !ev.BaseEvent.ProcessContext.HasParent() {
					return 0
				}
				return int(ev.BaseEvent.ProcessContext.Parent.PPid)
			},
			Field:  field,
			Weight: eval.FunctionWeight,
		}, nil
	case "process.pid":
		return &eval.IntEvaluator{
			EvalFnc: func(ctx *eval.Context) int {
				ev := ctx.Event.(*Event)
				return int(ev.BaseEvent.ProcessContext.Process.PIDContext.Pid)
			},
			Field:  field,
			Weight: eval.FunctionWeight,
		}, nil
	case "process.ppid":
		return &eval.IntEvaluator{
			EvalFnc: func(ctx *eval.Context) int {
				ev := ctx.Event.(*Event)
				return int(ev.BaseEvent.ProcessContext.Process.PPid)
			},
			Field:  field,
			Weight: eval.FunctionWeight,
		}, nil
	}
	return nil, &eval.ErrFieldNotFound{Field: field}
}
func (ev *Event) GetFields() []eval.Field {
	return []eval.Field{
		"container.created_at",
		"container.id",
		"container.tags",
		"event.timestamp",
		"exec.args",
		"exec.args_flags",
		"exec.args_options",
		"exec.args_truncated",
		"exec.argv",
		"exec.argv0",
		"exec.comm",
		"exec.container.id",
		"exec.created_at",
		"exec.envp",
		"exec.envs",
		"exec.envs_truncated",
		"exec.file.name",
		"exec.file.name.length",
		"exec.file.path",
		"exec.file.path.length",
		"exec.pid",
		"exec.ppid",
		"exit.args",
		"exit.args_flags",
		"exit.args_options",
		"exit.args_truncated",
		"exit.argv",
		"exit.argv0",
		"exit.cause",
		"exit.code",
		"exit.comm",
		"exit.container.id",
		"exit.created_at",
		"exit.envp",
		"exit.envs",
		"exit.envs_truncated",
		"exit.file.name",
		"exit.file.name.length",
		"exit.file.path",
		"exit.file.path.length",
		"exit.pid",
		"exit.ppid",
		"open.file.destination.mode",
		"open.file.name",
		"open.file.name.length",
		"open.file.path",
		"open.file.path.length",
		"open.flags",
		"process.ancestors.args",
		"process.ancestors.args_flags",
		"process.ancestors.args_options",
		"process.ancestors.args_truncated",
		"process.ancestors.argv",
		"process.ancestors.argv0",
		"process.ancestors.comm",
		"process.ancestors.container.id",
		"process.ancestors.created_at",
		"process.ancestors.envp",
		"process.ancestors.envs",
		"process.ancestors.envs_truncated",
		"process.ancestors.file.name",
		"process.ancestors.file.name.length",
		"process.ancestors.file.path",
		"process.ancestors.file.path.length",
		"process.ancestors.pid",
		"process.ancestors.ppid",
		"process.args",
		"process.args_flags",
		"process.args_options",
		"process.args_truncated",
		"process.argv",
		"process.argv0",
		"process.comm",
		"process.container.id",
		"process.created_at",
		"process.envp",
		"process.envs",
		"process.envs_truncated",
		"process.file.name",
		"process.file.name.length",
		"process.file.path",
		"process.file.path.length",
		"process.parent.args",
		"process.parent.args_flags",
		"process.parent.args_options",
		"process.parent.args_truncated",
		"process.parent.argv",
		"process.parent.argv0",
		"process.parent.comm",
		"process.parent.container.id",
		"process.parent.created_at",
		"process.parent.envp",
		"process.parent.envs",
		"process.parent.envs_truncated",
		"process.parent.file.name",
		"process.parent.file.name.length",
		"process.parent.file.path",
		"process.parent.file.path.length",
		"process.parent.pid",
		"process.parent.ppid",
		"process.pid",
		"process.ppid",
	}
}
func (ev *Event) GetFieldValue(field eval.Field) (interface{}, error) {
	switch field {
	case "container.created_at":
		return int(ev.FieldHandlers.ResolveContainerCreatedAt(ev, ev.BaseEvent.ContainerContext)), nil
	case "container.id":
		return ev.FieldHandlers.ResolveContainerID(ev, ev.BaseEvent.ContainerContext), nil
	case "container.tags":
		return ev.FieldHandlers.ResolveContainerTags(ev, ev.BaseEvent.ContainerContext), nil
	case "event.timestamp":
		return int(ev.FieldHandlers.ResolveEventTimestamp(ev, &ev.BaseEvent)), nil
	case "exec.args":
		return ev.FieldHandlers.ResolveProcessArgs(ev, ev.Exec.Process), nil
	case "exec.args_flags":
		return ev.FieldHandlers.ResolveProcessArgsFlags(ev, ev.Exec.Process), nil
	case "exec.args_options":
		return ev.FieldHandlers.ResolveProcessArgsOptions(ev, ev.Exec.Process), nil
	case "exec.args_truncated":
		return ev.FieldHandlers.ResolveProcessArgsTruncated(ev, ev.Exec.Process), nil
	case "exec.argv":
		return ev.FieldHandlers.ResolveProcessArgv(ev, ev.Exec.Process), nil
	case "exec.argv0":
		return ev.FieldHandlers.ResolveProcessArgv0(ev, ev.Exec.Process), nil
	case "exec.comm":
		return ev.Exec.Process.Comm, nil
	case "exec.container.id":
		return ev.Exec.Process.ContainerID, nil
	case "exec.created_at":
		return int(ev.FieldHandlers.ResolveProcessCreatedAt(ev, ev.Exec.Process)), nil
	case "exec.envp":
		return ev.FieldHandlers.ResolveProcessEnvp(ev, ev.Exec.Process), nil
	case "exec.envs":
		return ev.FieldHandlers.ResolveProcessEnvs(ev, ev.Exec.Process), nil
	case "exec.envs_truncated":
		return ev.FieldHandlers.ResolveProcessEnvsTruncated(ev, ev.Exec.Process), nil
	case "exec.file.name":
		return ev.FieldHandlers.ResolveFileBasename(ev, &ev.Exec.Process.FileEvent), nil
	case "exec.file.name.length":
		return ev.FieldHandlers.ResolveFileBasename(ev, &ev.Exec.Process.FileEvent), nil
	case "exec.file.path":
		return ev.FieldHandlers.ResolveFilePath(ev, &ev.Exec.Process.FileEvent), nil
	case "exec.file.path.length":
		return ev.FieldHandlers.ResolveFilePath(ev, &ev.Exec.Process.FileEvent), nil
	case "exec.pid":
		return int(ev.Exec.Process.PIDContext.Pid), nil
	case "exec.ppid":
		return int(ev.Exec.Process.PPid), nil
	case "exit.args":
		return ev.FieldHandlers.ResolveProcessArgs(ev, ev.Exit.Process), nil
	case "exit.args_flags":
		return ev.FieldHandlers.ResolveProcessArgsFlags(ev, ev.Exit.Process), nil
	case "exit.args_options":
		return ev.FieldHandlers.ResolveProcessArgsOptions(ev, ev.Exit.Process), nil
	case "exit.args_truncated":
		return ev.FieldHandlers.ResolveProcessArgsTruncated(ev, ev.Exit.Process), nil
	case "exit.argv":
		return ev.FieldHandlers.ResolveProcessArgv(ev, ev.Exit.Process), nil
	case "exit.argv0":
		return ev.FieldHandlers.ResolveProcessArgv0(ev, ev.Exit.Process), nil
	case "exit.cause":
		return int(ev.Exit.Cause), nil
	case "exit.code":
		return int(ev.Exit.Code), nil
	case "exit.comm":
		return ev.Exit.Process.Comm, nil
	case "exit.container.id":
		return ev.Exit.Process.ContainerID, nil
	case "exit.created_at":
		return int(ev.FieldHandlers.ResolveProcessCreatedAt(ev, ev.Exit.Process)), nil
	case "exit.envp":
		return ev.FieldHandlers.ResolveProcessEnvp(ev, ev.Exit.Process), nil
	case "exit.envs":
		return ev.FieldHandlers.ResolveProcessEnvs(ev, ev.Exit.Process), nil
	case "exit.envs_truncated":
		return ev.FieldHandlers.ResolveProcessEnvsTruncated(ev, ev.Exit.Process), nil
	case "exit.file.name":
		return ev.FieldHandlers.ResolveFileBasename(ev, &ev.Exit.Process.FileEvent), nil
	case "exit.file.name.length":
		return ev.FieldHandlers.ResolveFileBasename(ev, &ev.Exit.Process.FileEvent), nil
	case "exit.file.path":
		return ev.FieldHandlers.ResolveFilePath(ev, &ev.Exit.Process.FileEvent), nil
	case "exit.file.path.length":
		return ev.FieldHandlers.ResolveFilePath(ev, &ev.Exit.Process.FileEvent), nil
	case "exit.pid":
		return int(ev.Exit.Process.PIDContext.Pid), nil
	case "exit.ppid":
		return int(ev.Exit.Process.PPid), nil
	case "open.file.destination.mode":
		return int(ev.Open.Mode), nil
	case "open.file.name":
		return ev.FieldHandlers.ResolveFileBasename(ev, &ev.Open.File), nil
	case "open.file.name.length":
		return ev.FieldHandlers.ResolveFileBasename(ev, &ev.Open.File), nil
	case "open.file.path":
		return ev.FieldHandlers.ResolveFilePath(ev, &ev.Open.File), nil
	case "open.file.path.length":
		return ev.FieldHandlers.ResolveFilePath(ev, &ev.Open.File), nil
	case "open.flags":
		return int(ev.Open.Flags), nil
	case "process.ancestors.args":
		var values []string
		ctx := eval.NewContext(ev)
		iterator := &ProcessAncestorsIterator{}
		ptr := iterator.Front(ctx)
		for ptr != nil {
			element := (*ProcessCacheEntry)(ptr)
			result := ev.FieldHandlers.ResolveProcessArgs(ev, &element.ProcessContext.Process)
			values = append(values, result)
			ptr = iterator.Next()
		}
		return values, nil
	case "process.ancestors.args_flags":
		var values []string
		ctx := eval.NewContext(ev)
		iterator := &ProcessAncestorsIterator{}
		ptr := iterator.Front(ctx)
		for ptr != nil {
			element := (*ProcessCacheEntry)(ptr)
			result := ev.FieldHandlers.ResolveProcessArgsFlags(ev, &element.ProcessContext.Process)
			values = append(values, result...)
			ptr = iterator.Next()
		}
		return values, nil
	case "process.ancestors.args_options":
		var values []string
		ctx := eval.NewContext(ev)
		iterator := &ProcessAncestorsIterator{}
		ptr := iterator.Front(ctx)
		for ptr != nil {
			element := (*ProcessCacheEntry)(ptr)
			result := ev.FieldHandlers.ResolveProcessArgsOptions(ev, &element.ProcessContext.Process)
			values = append(values, result...)
			ptr = iterator.Next()
		}
		return values, nil
	case "process.ancestors.args_truncated":
		var values []bool
		ctx := eval.NewContext(ev)
		iterator := &ProcessAncestorsIterator{}
		ptr := iterator.Front(ctx)
		for ptr != nil {
			element := (*ProcessCacheEntry)(ptr)
			result := ev.FieldHandlers.ResolveProcessArgsTruncated(ev, &element.ProcessContext.Process)
			values = append(values, result)
			ptr = iterator.Next()
		}
		return values, nil
	case "process.ancestors.argv":
		var values []string
		ctx := eval.NewContext(ev)
		iterator := &ProcessAncestorsIterator{}
		ptr := iterator.Front(ctx)
		for ptr != nil {
			element := (*ProcessCacheEntry)(ptr)
			result := ev.FieldHandlers.ResolveProcessArgv(ev, &element.ProcessContext.Process)
			values = append(values, result...)
			ptr = iterator.Next()
		}
		return values, nil
	case "process.ancestors.argv0":
		var values []string
		ctx := eval.NewContext(ev)
		iterator := &ProcessAncestorsIterator{}
		ptr := iterator.Front(ctx)
		for ptr != nil {
			element := (*ProcessCacheEntry)(ptr)
			result := ev.FieldHandlers.ResolveProcessArgv0(ev, &element.ProcessContext.Process)
			values = append(values, result)
			ptr = iterator.Next()
		}
		return values, nil
	case "process.ancestors.comm":
		var values []string
		ctx := eval.NewContext(ev)
		iterator := &ProcessAncestorsIterator{}
		ptr := iterator.Front(ctx)
		for ptr != nil {
			element := (*ProcessCacheEntry)(ptr)
			result := element.ProcessContext.Process.Comm
			values = append(values, result)
			ptr = iterator.Next()
		}
		return values, nil
	case "process.ancestors.container.id":
		var values []string
		ctx := eval.NewContext(ev)
		iterator := &ProcessAncestorsIterator{}
		ptr := iterator.Front(ctx)
		for ptr != nil {
			element := (*ProcessCacheEntry)(ptr)
			result := element.ProcessContext.Process.ContainerID
			values = append(values, result)
			ptr = iterator.Next()
		}
		return values, nil
	case "process.ancestors.created_at":
		var values []int
		ctx := eval.NewContext(ev)
		iterator := &ProcessAncestorsIterator{}
		ptr := iterator.Front(ctx)
		for ptr != nil {
			element := (*ProcessCacheEntry)(ptr)
			result := int(ev.FieldHandlers.ResolveProcessCreatedAt(ev, &element.ProcessContext.Process))
			values = append(values, result)
			ptr = iterator.Next()
		}
		return values, nil
	case "process.ancestors.envp":
		var values []string
		ctx := eval.NewContext(ev)
		iterator := &ProcessAncestorsIterator{}
		ptr := iterator.Front(ctx)
		for ptr != nil {
			element := (*ProcessCacheEntry)(ptr)
			result := ev.FieldHandlers.ResolveProcessEnvp(ev, &element.ProcessContext.Process)
			values = append(values, result...)
			ptr = iterator.Next()
		}
		return values, nil
	case "process.ancestors.envs":
		var values []string
		ctx := eval.NewContext(ev)
		iterator := &ProcessAncestorsIterator{}
		ptr := iterator.Front(ctx)
		for ptr != nil {
			element := (*ProcessCacheEntry)(ptr)
			result := ev.FieldHandlers.ResolveProcessEnvs(ev, &element.ProcessContext.Process)
			values = append(values, result...)
			ptr = iterator.Next()
		}
		return values, nil
	case "process.ancestors.envs_truncated":
		var values []bool
		ctx := eval.NewContext(ev)
		iterator := &ProcessAncestorsIterator{}
		ptr := iterator.Front(ctx)
		for ptr != nil {
			element := (*ProcessCacheEntry)(ptr)
			result := ev.FieldHandlers.ResolveProcessEnvsTruncated(ev, &element.ProcessContext.Process)
			values = append(values, result)
			ptr = iterator.Next()
		}
		return values, nil
	case "process.ancestors.file.name":
		var values []string
		ctx := eval.NewContext(ev)
		iterator := &ProcessAncestorsIterator{}
		ptr := iterator.Front(ctx)
		for ptr != nil {
			element := (*ProcessCacheEntry)(ptr)
			result := ev.FieldHandlers.ResolveFileBasename(ev, &element.ProcessContext.Process.FileEvent)
			values = append(values, result)
			ptr = iterator.Next()
		}
		return values, nil
	case "process.ancestors.file.name.length":
		var values []int
		ctx := eval.NewContext(ev)
		iterator := &ProcessAncestorsIterator{}
		ptr := iterator.Front(ctx)
		for ptr != nil {
			element := (*ProcessCacheEntry)(ptr)
			result := len(ev.FieldHandlers.ResolveFileBasename(ev, &element.ProcessContext.Process.FileEvent))
			values = append(values, result)
			ptr = iterator.Next()
		}
		return values, nil
	case "process.ancestors.file.path":
		var values []string
		ctx := eval.NewContext(ev)
		iterator := &ProcessAncestorsIterator{}
		ptr := iterator.Front(ctx)
		for ptr != nil {
			element := (*ProcessCacheEntry)(ptr)
			result := ev.FieldHandlers.ResolveFilePath(ev, &element.ProcessContext.Process.FileEvent)
			values = append(values, result)
			ptr = iterator.Next()
		}
		return values, nil
	case "process.ancestors.file.path.length":
		var values []int
		ctx := eval.NewContext(ev)
		iterator := &ProcessAncestorsIterator{}
		ptr := iterator.Front(ctx)
		for ptr != nil {
			element := (*ProcessCacheEntry)(ptr)
			result := len(ev.FieldHandlers.ResolveFilePath(ev, &element.ProcessContext.Process.FileEvent))
			values = append(values, result)
			ptr = iterator.Next()
		}
		return values, nil
	case "process.ancestors.pid":
		var values []int
		ctx := eval.NewContext(ev)
		iterator := &ProcessAncestorsIterator{}
		ptr := iterator.Front(ctx)
		for ptr != nil {
			element := (*ProcessCacheEntry)(ptr)
			result := int(element.ProcessContext.Process.PIDContext.Pid)
			values = append(values, result)
			ptr = iterator.Next()
		}
		return values, nil
	case "process.ancestors.ppid":
		var values []int
		ctx := eval.NewContext(ev)
		iterator := &ProcessAncestorsIterator{}
		ptr := iterator.Front(ctx)
		for ptr != nil {
			element := (*ProcessCacheEntry)(ptr)
			result := int(element.ProcessContext.Process.PPid)
			values = append(values, result)
			ptr = iterator.Next()
		}
		return values, nil
	case "process.args":
		return ev.FieldHandlers.ResolveProcessArgs(ev, &ev.BaseEvent.ProcessContext.Process), nil
	case "process.args_flags":
		return ev.FieldHandlers.ResolveProcessArgsFlags(ev, &ev.BaseEvent.ProcessContext.Process), nil
	case "process.args_options":
		return ev.FieldHandlers.ResolveProcessArgsOptions(ev, &ev.BaseEvent.ProcessContext.Process), nil
	case "process.args_truncated":
		return ev.FieldHandlers.ResolveProcessArgsTruncated(ev, &ev.BaseEvent.ProcessContext.Process), nil
	case "process.argv":
		return ev.FieldHandlers.ResolveProcessArgv(ev, &ev.BaseEvent.ProcessContext.Process), nil
	case "process.argv0":
		return ev.FieldHandlers.ResolveProcessArgv0(ev, &ev.BaseEvent.ProcessContext.Process), nil
	case "process.comm":
		return ev.BaseEvent.ProcessContext.Process.Comm, nil
	case "process.container.id":
		return ev.BaseEvent.ProcessContext.Process.ContainerID, nil
	case "process.created_at":
		return int(ev.FieldHandlers.ResolveProcessCreatedAt(ev, &ev.BaseEvent.ProcessContext.Process)), nil
	case "process.envp":
		return ev.FieldHandlers.ResolveProcessEnvp(ev, &ev.BaseEvent.ProcessContext.Process), nil
	case "process.envs":
		return ev.FieldHandlers.ResolveProcessEnvs(ev, &ev.BaseEvent.ProcessContext.Process), nil
	case "process.envs_truncated":
		return ev.FieldHandlers.ResolveProcessEnvsTruncated(ev, &ev.BaseEvent.ProcessContext.Process), nil
	case "process.file.name":
		return ev.FieldHandlers.ResolveFileBasename(ev, &ev.BaseEvent.ProcessContext.Process.FileEvent), nil
	case "process.file.name.length":
		return ev.FieldHandlers.ResolveFileBasename(ev, &ev.BaseEvent.ProcessContext.Process.FileEvent), nil
	case "process.file.path":
		return ev.FieldHandlers.ResolveFilePath(ev, &ev.BaseEvent.ProcessContext.Process.FileEvent), nil
	case "process.file.path.length":
		return ev.FieldHandlers.ResolveFilePath(ev, &ev.BaseEvent.ProcessContext.Process.FileEvent), nil
	case "process.parent.args":
		if !ev.BaseEvent.ProcessContext.HasParent() {
			return "", &eval.ErrNotSupported{Field: field}
		}
		return ev.FieldHandlers.ResolveProcessArgs(ev, ev.BaseEvent.ProcessContext.Parent), nil
	case "process.parent.args_flags":
		if !ev.BaseEvent.ProcessContext.HasParent() {
			return []string{}, &eval.ErrNotSupported{Field: field}
		}
		return ev.FieldHandlers.ResolveProcessArgsFlags(ev, ev.BaseEvent.ProcessContext.Parent), nil
	case "process.parent.args_options":
		if !ev.BaseEvent.ProcessContext.HasParent() {
			return []string{}, &eval.ErrNotSupported{Field: field}
		}
		return ev.FieldHandlers.ResolveProcessArgsOptions(ev, ev.BaseEvent.ProcessContext.Parent), nil
	case "process.parent.args_truncated":
		if !ev.BaseEvent.ProcessContext.HasParent() {
			return false, &eval.ErrNotSupported{Field: field}
		}
		return ev.FieldHandlers.ResolveProcessArgsTruncated(ev, ev.BaseEvent.ProcessContext.Parent), nil
	case "process.parent.argv":
		if !ev.BaseEvent.ProcessContext.HasParent() {
			return []string{}, &eval.ErrNotSupported{Field: field}
		}
		return ev.FieldHandlers.ResolveProcessArgv(ev, ev.BaseEvent.ProcessContext.Parent), nil
	case "process.parent.argv0":
		if !ev.BaseEvent.ProcessContext.HasParent() {
			return "", &eval.ErrNotSupported{Field: field}
		}
		return ev.FieldHandlers.ResolveProcessArgv0(ev, ev.BaseEvent.ProcessContext.Parent), nil
	case "process.parent.comm":
		if !ev.BaseEvent.ProcessContext.HasParent() {
			return "", &eval.ErrNotSupported{Field: field}
		}
		return ev.BaseEvent.ProcessContext.Parent.Comm, nil
	case "process.parent.container.id":
		if !ev.BaseEvent.ProcessContext.HasParent() {
			return "", &eval.ErrNotSupported{Field: field}
		}
		return ev.BaseEvent.ProcessContext.Parent.ContainerID, nil
	case "process.parent.created_at":
		if !ev.BaseEvent.ProcessContext.HasParent() {
			return 0, &eval.ErrNotSupported{Field: field}
		}
		return int(ev.FieldHandlers.ResolveProcessCreatedAt(ev, ev.BaseEvent.ProcessContext.Parent)), nil
	case "process.parent.envp":
		if !ev.BaseEvent.ProcessContext.HasParent() {
			return []string{}, &eval.ErrNotSupported{Field: field}
		}
		return ev.FieldHandlers.ResolveProcessEnvp(ev, ev.BaseEvent.ProcessContext.Parent), nil
	case "process.parent.envs":
		if !ev.BaseEvent.ProcessContext.HasParent() {
			return []string{}, &eval.ErrNotSupported{Field: field}
		}
		return ev.FieldHandlers.ResolveProcessEnvs(ev, ev.BaseEvent.ProcessContext.Parent), nil
	case "process.parent.envs_truncated":
		if !ev.BaseEvent.ProcessContext.HasParent() {
			return false, &eval.ErrNotSupported{Field: field}
		}
		return ev.FieldHandlers.ResolveProcessEnvsTruncated(ev, ev.BaseEvent.ProcessContext.Parent), nil
	case "process.parent.file.name":
		if !ev.BaseEvent.ProcessContext.HasParent() {
			return "", &eval.ErrNotSupported{Field: field}
		}
		return ev.FieldHandlers.ResolveFileBasename(ev, &ev.BaseEvent.ProcessContext.Parent.FileEvent), nil
	case "process.parent.file.name.length":
		return ev.FieldHandlers.ResolveFileBasename(ev, &ev.BaseEvent.ProcessContext.Parent.FileEvent), nil
	case "process.parent.file.path":
		if !ev.BaseEvent.ProcessContext.HasParent() {
			return "", &eval.ErrNotSupported{Field: field}
		}
		return ev.FieldHandlers.ResolveFilePath(ev, &ev.BaseEvent.ProcessContext.Parent.FileEvent), nil
	case "process.parent.file.path.length":
		return ev.FieldHandlers.ResolveFilePath(ev, &ev.BaseEvent.ProcessContext.Parent.FileEvent), nil
	case "process.parent.pid":
		if !ev.BaseEvent.ProcessContext.HasParent() {
			return 0, &eval.ErrNotSupported{Field: field}
		}
		return int(ev.BaseEvent.ProcessContext.Parent.PIDContext.Pid), nil
	case "process.parent.ppid":
		if !ev.BaseEvent.ProcessContext.HasParent() {
			return 0, &eval.ErrNotSupported{Field: field}
		}
		return int(ev.BaseEvent.ProcessContext.Parent.PPid), nil
	case "process.pid":
		return int(ev.BaseEvent.ProcessContext.Process.PIDContext.Pid), nil
	case "process.ppid":
		return int(ev.BaseEvent.ProcessContext.Process.PPid), nil
	}
	return nil, &eval.ErrFieldNotFound{Field: field}
}
func (ev *Event) GetFieldEventType(field eval.Field) (eval.EventType, error) {
	switch field {
	case "container.created_at":
		return "*", nil
	case "container.id":
		return "*", nil
	case "container.tags":
		return "*", nil
	case "event.timestamp":
		return "*", nil
	case "exec.args":
		return "exec", nil
	case "exec.args_flags":
		return "exec", nil
	case "exec.args_options":
		return "exec", nil
	case "exec.args_truncated":
		return "exec", nil
	case "exec.argv":
		return "exec", nil
	case "exec.argv0":
		return "exec", nil
	case "exec.comm":
		return "exec", nil
	case "exec.container.id":
		return "exec", nil
	case "exec.created_at":
		return "exec", nil
	case "exec.envp":
		return "exec", nil
	case "exec.envs":
		return "exec", nil
	case "exec.envs_truncated":
		return "exec", nil
	case "exec.file.name":
		return "exec", nil
	case "exec.file.name.length":
		return "exec", nil
	case "exec.file.path":
		return "exec", nil
	case "exec.file.path.length":
		return "exec", nil
	case "exec.pid":
		return "exec", nil
	case "exec.ppid":
		return "exec", nil
	case "exit.args":
		return "exit", nil
	case "exit.args_flags":
		return "exit", nil
	case "exit.args_options":
		return "exit", nil
	case "exit.args_truncated":
		return "exit", nil
	case "exit.argv":
		return "exit", nil
	case "exit.argv0":
		return "exit", nil
	case "exit.cause":
		return "exit", nil
	case "exit.code":
		return "exit", nil
	case "exit.comm":
		return "exit", nil
	case "exit.container.id":
		return "exit", nil
	case "exit.created_at":
		return "exit", nil
	case "exit.envp":
		return "exit", nil
	case "exit.envs":
		return "exit", nil
	case "exit.envs_truncated":
		return "exit", nil
	case "exit.file.name":
		return "exit", nil
	case "exit.file.name.length":
		return "exit", nil
	case "exit.file.path":
		return "exit", nil
	case "exit.file.path.length":
		return "exit", nil
	case "exit.pid":
		return "exit", nil
	case "exit.ppid":
		return "exit", nil
	case "open.file.destination.mode":
		return "open", nil
	case "open.file.name":
		return "open", nil
	case "open.file.name.length":
		return "open", nil
	case "open.file.path":
		return "open", nil
	case "open.file.path.length":
		return "open", nil
	case "open.flags":
		return "open", nil
	case "process.ancestors.args":
		return "*", nil
	case "process.ancestors.args_flags":
		return "*", nil
	case "process.ancestors.args_options":
		return "*", nil
	case "process.ancestors.args_truncated":
		return "*", nil
	case "process.ancestors.argv":
		return "*", nil
	case "process.ancestors.argv0":
		return "*", nil
	case "process.ancestors.comm":
		return "*", nil
	case "process.ancestors.container.id":
		return "*", nil
	case "process.ancestors.created_at":
		return "*", nil
	case "process.ancestors.envp":
		return "*", nil
	case "process.ancestors.envs":
		return "*", nil
	case "process.ancestors.envs_truncated":
		return "*", nil
	case "process.ancestors.file.name":
		return "*", nil
	case "process.ancestors.file.name.length":
		return "*", nil
	case "process.ancestors.file.path":
		return "*", nil
	case "process.ancestors.file.path.length":
		return "*", nil
	case "process.ancestors.pid":
		return "*", nil
	case "process.ancestors.ppid":
		return "*", nil
	case "process.args":
		return "*", nil
	case "process.args_flags":
		return "*", nil
	case "process.args_options":
		return "*", nil
	case "process.args_truncated":
		return "*", nil
	case "process.argv":
		return "*", nil
	case "process.argv0":
		return "*", nil
	case "process.comm":
		return "*", nil
	case "process.container.id":
		return "*", nil
	case "process.created_at":
		return "*", nil
	case "process.envp":
		return "*", nil
	case "process.envs":
		return "*", nil
	case "process.envs_truncated":
		return "*", nil
	case "process.file.name":
		return "*", nil
	case "process.file.name.length":
		return "*", nil
	case "process.file.path":
		return "*", nil
	case "process.file.path.length":
		return "*", nil
	case "process.parent.args":
		return "*", nil
	case "process.parent.args_flags":
		return "*", nil
	case "process.parent.args_options":
		return "*", nil
	case "process.parent.args_truncated":
		return "*", nil
	case "process.parent.argv":
		return "*", nil
	case "process.parent.argv0":
		return "*", nil
	case "process.parent.comm":
		return "*", nil
	case "process.parent.container.id":
		return "*", nil
	case "process.parent.created_at":
		return "*", nil
	case "process.parent.envp":
		return "*", nil
	case "process.parent.envs":
		return "*", nil
	case "process.parent.envs_truncated":
		return "*", nil
	case "process.parent.file.name":
		return "*", nil
	case "process.parent.file.name.length":
		return "*", nil
	case "process.parent.file.path":
		return "*", nil
	case "process.parent.file.path.length":
		return "*", nil
	case "process.parent.pid":
		return "*", nil
	case "process.parent.ppid":
		return "*", nil
	case "process.pid":
		return "*", nil
	case "process.ppid":
		return "*", nil
	}
	return "", &eval.ErrFieldNotFound{Field: field}
}
func (ev *Event) GetFieldType(field eval.Field) (reflect.Kind, error) {
	switch field {
	case "container.created_at":
		return reflect.Int, nil
	case "container.id":
		return reflect.String, nil
	case "container.tags":
		return reflect.String, nil
	case "event.timestamp":
		return reflect.Int, nil
	case "exec.args":
		return reflect.String, nil
	case "exec.args_flags":
		return reflect.String, nil
	case "exec.args_options":
		return reflect.String, nil
	case "exec.args_truncated":
		return reflect.Bool, nil
	case "exec.argv":
		return reflect.String, nil
	case "exec.argv0":
		return reflect.String, nil
	case "exec.comm":
		return reflect.String, nil
	case "exec.container.id":
		return reflect.String, nil
	case "exec.created_at":
		return reflect.Int, nil
	case "exec.envp":
		return reflect.String, nil
	case "exec.envs":
		return reflect.String, nil
	case "exec.envs_truncated":
		return reflect.Bool, nil
	case "exec.file.name":
		return reflect.String, nil
	case "exec.file.name.length":
		return reflect.Int, nil
	case "exec.file.path":
		return reflect.String, nil
	case "exec.file.path.length":
		return reflect.Int, nil
	case "exec.pid":
		return reflect.Int, nil
	case "exec.ppid":
		return reflect.Int, nil
	case "exit.args":
		return reflect.String, nil
	case "exit.args_flags":
		return reflect.String, nil
	case "exit.args_options":
		return reflect.String, nil
	case "exit.args_truncated":
		return reflect.Bool, nil
	case "exit.argv":
		return reflect.String, nil
	case "exit.argv0":
		return reflect.String, nil
	case "exit.cause":
		return reflect.Int, nil
	case "exit.code":
		return reflect.Int, nil
	case "exit.comm":
		return reflect.String, nil
	case "exit.container.id":
		return reflect.String, nil
	case "exit.created_at":
		return reflect.Int, nil
	case "exit.envp":
		return reflect.String, nil
	case "exit.envs":
		return reflect.String, nil
	case "exit.envs_truncated":
		return reflect.Bool, nil
	case "exit.file.name":
		return reflect.String, nil
	case "exit.file.name.length":
		return reflect.Int, nil
	case "exit.file.path":
		return reflect.String, nil
	case "exit.file.path.length":
		return reflect.Int, nil
	case "exit.pid":
		return reflect.Int, nil
	case "exit.ppid":
		return reflect.Int, nil
	case "open.file.destination.mode":
		return reflect.Int, nil
	case "open.file.name":
		return reflect.String, nil
	case "open.file.name.length":
		return reflect.Int, nil
	case "open.file.path":
		return reflect.String, nil
	case "open.file.path.length":
		return reflect.Int, nil
	case "open.flags":
		return reflect.Int, nil
	case "process.ancestors.args":
		return reflect.String, nil
	case "process.ancestors.args_flags":
		return reflect.String, nil
	case "process.ancestors.args_options":
		return reflect.String, nil
	case "process.ancestors.args_truncated":
		return reflect.Bool, nil
	case "process.ancestors.argv":
		return reflect.String, nil
	case "process.ancestors.argv0":
		return reflect.String, nil
	case "process.ancestors.comm":
		return reflect.String, nil
	case "process.ancestors.container.id":
		return reflect.String, nil
	case "process.ancestors.created_at":
		return reflect.Int, nil
	case "process.ancestors.envp":
		return reflect.String, nil
	case "process.ancestors.envs":
		return reflect.String, nil
	case "process.ancestors.envs_truncated":
		return reflect.Bool, nil
	case "process.ancestors.file.name":
		return reflect.String, nil
	case "process.ancestors.file.name.length":
		return reflect.Int, nil
	case "process.ancestors.file.path":
		return reflect.String, nil
	case "process.ancestors.file.path.length":
		return reflect.Int, nil
	case "process.ancestors.pid":
		return reflect.Int, nil
	case "process.ancestors.ppid":
		return reflect.Int, nil
	case "process.args":
		return reflect.String, nil
	case "process.args_flags":
		return reflect.String, nil
	case "process.args_options":
		return reflect.String, nil
	case "process.args_truncated":
		return reflect.Bool, nil
	case "process.argv":
		return reflect.String, nil
	case "process.argv0":
		return reflect.String, nil
	case "process.comm":
		return reflect.String, nil
	case "process.container.id":
		return reflect.String, nil
	case "process.created_at":
		return reflect.Int, nil
	case "process.envp":
		return reflect.String, nil
	case "process.envs":
		return reflect.String, nil
	case "process.envs_truncated":
		return reflect.Bool, nil
	case "process.file.name":
		return reflect.String, nil
	case "process.file.name.length":
		return reflect.Int, nil
	case "process.file.path":
		return reflect.String, nil
	case "process.file.path.length":
		return reflect.Int, nil
	case "process.parent.args":
		return reflect.String, nil
	case "process.parent.args_flags":
		return reflect.String, nil
	case "process.parent.args_options":
		return reflect.String, nil
	case "process.parent.args_truncated":
		return reflect.Bool, nil
	case "process.parent.argv":
		return reflect.String, nil
	case "process.parent.argv0":
		return reflect.String, nil
	case "process.parent.comm":
		return reflect.String, nil
	case "process.parent.container.id":
		return reflect.String, nil
	case "process.parent.created_at":
		return reflect.Int, nil
	case "process.parent.envp":
		return reflect.String, nil
	case "process.parent.envs":
		return reflect.String, nil
	case "process.parent.envs_truncated":
		return reflect.Bool, nil
	case "process.parent.file.name":
		return reflect.String, nil
	case "process.parent.file.name.length":
		return reflect.Int, nil
	case "process.parent.file.path":
		return reflect.String, nil
	case "process.parent.file.path.length":
		return reflect.Int, nil
	case "process.parent.pid":
		return reflect.Int, nil
	case "process.parent.ppid":
		return reflect.Int, nil
	case "process.pid":
		return reflect.Int, nil
	case "process.ppid":
		return reflect.Int, nil
	}
	return reflect.Invalid, &eval.ErrFieldNotFound{Field: field}
}
func (ev *Event) SetFieldValue(field eval.Field, value interface{}) error {
	switch field {
	case "container.created_at":
		if ev.BaseEvent.ContainerContext == nil {
			ev.BaseEvent.ContainerContext = &ContainerContext{}
		}
		rv, ok := value.(int)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ContainerContext.CreatedAt"}
		}
		ev.BaseEvent.ContainerContext.CreatedAt = uint64(rv)
		return nil
	case "container.id":
		if ev.BaseEvent.ContainerContext == nil {
			ev.BaseEvent.ContainerContext = &ContainerContext{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ContainerContext.ID"}
		}
		ev.BaseEvent.ContainerContext.ID = rv
		return nil
	case "container.tags":
		if ev.BaseEvent.ContainerContext == nil {
			ev.BaseEvent.ContainerContext = &ContainerContext{}
		}
		switch rv := value.(type) {
		case string:
			ev.BaseEvent.ContainerContext.Tags = append(ev.BaseEvent.ContainerContext.Tags, rv)
		case []string:
			ev.BaseEvent.ContainerContext.Tags = append(ev.BaseEvent.ContainerContext.Tags, rv...)
		default:
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ContainerContext.Tags"}
		}
		return nil
	case "event.timestamp":
		rv, ok := value.(int)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.TimestampRaw"}
		}
		ev.BaseEvent.TimestampRaw = uint64(rv)
		return nil
	case "exec.args":
		if ev.Exec.Process == nil {
			ev.Exec.Process = &Process{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "Exec.Process.Args"}
		}
		ev.Exec.Process.Args = rv
		return nil
	case "exec.args_flags":
		if ev.Exec.Process == nil {
			ev.Exec.Process = &Process{}
		}
		switch rv := value.(type) {
		case string:
			ev.Exec.Process.Argv = append(ev.Exec.Process.Argv, rv)
		case []string:
			ev.Exec.Process.Argv = append(ev.Exec.Process.Argv, rv...)
		default:
			return &eval.ErrValueTypeMismatch{Field: "Exec.Process.Argv"}
		}
		return nil
	case "exec.args_options":
		if ev.Exec.Process == nil {
			ev.Exec.Process = &Process{}
		}
		switch rv := value.(type) {
		case string:
			ev.Exec.Process.Argv = append(ev.Exec.Process.Argv, rv)
		case []string:
			ev.Exec.Process.Argv = append(ev.Exec.Process.Argv, rv...)
		default:
			return &eval.ErrValueTypeMismatch{Field: "Exec.Process.Argv"}
		}
		return nil
	case "exec.args_truncated":
		if ev.Exec.Process == nil {
			ev.Exec.Process = &Process{}
		}
		rv, ok := value.(bool)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "Exec.Process.ArgsTruncated"}
		}
		ev.Exec.Process.ArgsTruncated = rv
		return nil
	case "exec.argv":
		if ev.Exec.Process == nil {
			ev.Exec.Process = &Process{}
		}
		switch rv := value.(type) {
		case string:
			ev.Exec.Process.Argv = append(ev.Exec.Process.Argv, rv)
		case []string:
			ev.Exec.Process.Argv = append(ev.Exec.Process.Argv, rv...)
		default:
			return &eval.ErrValueTypeMismatch{Field: "Exec.Process.Argv"}
		}
		return nil
	case "exec.argv0":
		if ev.Exec.Process == nil {
			ev.Exec.Process = &Process{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "Exec.Process.Argv0"}
		}
		ev.Exec.Process.Argv0 = rv
		return nil
	case "exec.comm":
		if ev.Exec.Process == nil {
			ev.Exec.Process = &Process{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "Exec.Process.Comm"}
		}
		ev.Exec.Process.Comm = rv
		return nil
	case "exec.container.id":
		if ev.Exec.Process == nil {
			ev.Exec.Process = &Process{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "Exec.Process.ContainerID"}
		}
		ev.Exec.Process.ContainerID = rv
		return nil
	case "exec.created_at":
		if ev.Exec.Process == nil {
			ev.Exec.Process = &Process{}
		}
		rv, ok := value.(int)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "Exec.Process.CreatedAt"}
		}
		ev.Exec.Process.CreatedAt = uint64(rv)
		return nil
	case "exec.envp":
		if ev.Exec.Process == nil {
			ev.Exec.Process = &Process{}
		}
		switch rv := value.(type) {
		case string:
			ev.Exec.Process.Envp = append(ev.Exec.Process.Envp, rv)
		case []string:
			ev.Exec.Process.Envp = append(ev.Exec.Process.Envp, rv...)
		default:
			return &eval.ErrValueTypeMismatch{Field: "Exec.Process.Envp"}
		}
		return nil
	case "exec.envs":
		if ev.Exec.Process == nil {
			ev.Exec.Process = &Process{}
		}
		switch rv := value.(type) {
		case string:
			ev.Exec.Process.Envs = append(ev.Exec.Process.Envs, rv)
		case []string:
			ev.Exec.Process.Envs = append(ev.Exec.Process.Envs, rv...)
		default:
			return &eval.ErrValueTypeMismatch{Field: "Exec.Process.Envs"}
		}
		return nil
	case "exec.envs_truncated":
		if ev.Exec.Process == nil {
			ev.Exec.Process = &Process{}
		}
		rv, ok := value.(bool)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "Exec.Process.EnvsTruncated"}
		}
		ev.Exec.Process.EnvsTruncated = rv
		return nil
	case "exec.file.name":
		if ev.Exec.Process == nil {
			ev.Exec.Process = &Process{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "Exec.Process.FileEvent.BasenameStr"}
		}
		ev.Exec.Process.FileEvent.BasenameStr = rv
		return nil
	case "exec.file.name.length":
		if ev.Exec.Process == nil {
			ev.Exec.Process = &Process{}
		}
		return &eval.ErrFieldReadOnly{Field: "exec.file.name.length"}
	case "exec.file.path":
		if ev.Exec.Process == nil {
			ev.Exec.Process = &Process{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "Exec.Process.FileEvent.PathnameStr"}
		}
		ev.Exec.Process.FileEvent.PathnameStr = rv
		return nil
	case "exec.file.path.length":
		if ev.Exec.Process == nil {
			ev.Exec.Process = &Process{}
		}
		return &eval.ErrFieldReadOnly{Field: "exec.file.path.length"}
	case "exec.pid":
		if ev.Exec.Process == nil {
			ev.Exec.Process = &Process{}
		}
		rv, ok := value.(int)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "Exec.Process.PIDContext.Pid"}
		}
		ev.Exec.Process.PIDContext.Pid = uint32(rv)
		return nil
	case "exec.ppid":
		if ev.Exec.Process == nil {
			ev.Exec.Process = &Process{}
		}
		rv, ok := value.(int)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "Exec.Process.PPid"}
		}
		ev.Exec.Process.PPid = uint32(rv)
		return nil
	case "exit.args":
		if ev.Exit.Process == nil {
			ev.Exit.Process = &Process{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "Exit.Process.Args"}
		}
		ev.Exit.Process.Args = rv
		return nil
	case "exit.args_flags":
		if ev.Exit.Process == nil {
			ev.Exit.Process = &Process{}
		}
		switch rv := value.(type) {
		case string:
			ev.Exit.Process.Argv = append(ev.Exit.Process.Argv, rv)
		case []string:
			ev.Exit.Process.Argv = append(ev.Exit.Process.Argv, rv...)
		default:
			return &eval.ErrValueTypeMismatch{Field: "Exit.Process.Argv"}
		}
		return nil
	case "exit.args_options":
		if ev.Exit.Process == nil {
			ev.Exit.Process = &Process{}
		}
		switch rv := value.(type) {
		case string:
			ev.Exit.Process.Argv = append(ev.Exit.Process.Argv, rv)
		case []string:
			ev.Exit.Process.Argv = append(ev.Exit.Process.Argv, rv...)
		default:
			return &eval.ErrValueTypeMismatch{Field: "Exit.Process.Argv"}
		}
		return nil
	case "exit.args_truncated":
		if ev.Exit.Process == nil {
			ev.Exit.Process = &Process{}
		}
		rv, ok := value.(bool)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "Exit.Process.ArgsTruncated"}
		}
		ev.Exit.Process.ArgsTruncated = rv
		return nil
	case "exit.argv":
		if ev.Exit.Process == nil {
			ev.Exit.Process = &Process{}
		}
		switch rv := value.(type) {
		case string:
			ev.Exit.Process.Argv = append(ev.Exit.Process.Argv, rv)
		case []string:
			ev.Exit.Process.Argv = append(ev.Exit.Process.Argv, rv...)
		default:
			return &eval.ErrValueTypeMismatch{Field: "Exit.Process.Argv"}
		}
		return nil
	case "exit.argv0":
		if ev.Exit.Process == nil {
			ev.Exit.Process = &Process{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "Exit.Process.Argv0"}
		}
		ev.Exit.Process.Argv0 = rv
		return nil
	case "exit.cause":
		rv, ok := value.(int)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "Exit.Cause"}
		}
		ev.Exit.Cause = uint32(rv)
		return nil
	case "exit.code":
		rv, ok := value.(int)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "Exit.Code"}
		}
		ev.Exit.Code = uint32(rv)
		return nil
	case "exit.comm":
		if ev.Exit.Process == nil {
			ev.Exit.Process = &Process{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "Exit.Process.Comm"}
		}
		ev.Exit.Process.Comm = rv
		return nil
	case "exit.container.id":
		if ev.Exit.Process == nil {
			ev.Exit.Process = &Process{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "Exit.Process.ContainerID"}
		}
		ev.Exit.Process.ContainerID = rv
		return nil
	case "exit.created_at":
		if ev.Exit.Process == nil {
			ev.Exit.Process = &Process{}
		}
		rv, ok := value.(int)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "Exit.Process.CreatedAt"}
		}
		ev.Exit.Process.CreatedAt = uint64(rv)
		return nil
	case "exit.envp":
		if ev.Exit.Process == nil {
			ev.Exit.Process = &Process{}
		}
		switch rv := value.(type) {
		case string:
			ev.Exit.Process.Envp = append(ev.Exit.Process.Envp, rv)
		case []string:
			ev.Exit.Process.Envp = append(ev.Exit.Process.Envp, rv...)
		default:
			return &eval.ErrValueTypeMismatch{Field: "Exit.Process.Envp"}
		}
		return nil
	case "exit.envs":
		if ev.Exit.Process == nil {
			ev.Exit.Process = &Process{}
		}
		switch rv := value.(type) {
		case string:
			ev.Exit.Process.Envs = append(ev.Exit.Process.Envs, rv)
		case []string:
			ev.Exit.Process.Envs = append(ev.Exit.Process.Envs, rv...)
		default:
			return &eval.ErrValueTypeMismatch{Field: "Exit.Process.Envs"}
		}
		return nil
	case "exit.envs_truncated":
		if ev.Exit.Process == nil {
			ev.Exit.Process = &Process{}
		}
		rv, ok := value.(bool)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "Exit.Process.EnvsTruncated"}
		}
		ev.Exit.Process.EnvsTruncated = rv
		return nil
	case "exit.file.name":
		if ev.Exit.Process == nil {
			ev.Exit.Process = &Process{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "Exit.Process.FileEvent.BasenameStr"}
		}
		ev.Exit.Process.FileEvent.BasenameStr = rv
		return nil
	case "exit.file.name.length":
		if ev.Exit.Process == nil {
			ev.Exit.Process = &Process{}
		}
		return &eval.ErrFieldReadOnly{Field: "exit.file.name.length"}
	case "exit.file.path":
		if ev.Exit.Process == nil {
			ev.Exit.Process = &Process{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "Exit.Process.FileEvent.PathnameStr"}
		}
		ev.Exit.Process.FileEvent.PathnameStr = rv
		return nil
	case "exit.file.path.length":
		if ev.Exit.Process == nil {
			ev.Exit.Process = &Process{}
		}
		return &eval.ErrFieldReadOnly{Field: "exit.file.path.length"}
	case "exit.pid":
		if ev.Exit.Process == nil {
			ev.Exit.Process = &Process{}
		}
		rv, ok := value.(int)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "Exit.Process.PIDContext.Pid"}
		}
		ev.Exit.Process.PIDContext.Pid = uint32(rv)
		return nil
	case "exit.ppid":
		if ev.Exit.Process == nil {
			ev.Exit.Process = &Process{}
		}
		rv, ok := value.(int)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "Exit.Process.PPid"}
		}
		ev.Exit.Process.PPid = uint32(rv)
		return nil
	case "open.file.destination.mode":
		rv, ok := value.(int)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "Open.Mode"}
		}
		ev.Open.Mode = uint32(rv)
		return nil
	case "open.file.name":
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "Open.File.BasenameStr"}
		}
		ev.Open.File.BasenameStr = rv
		return nil
	case "open.file.name.length":
		return &eval.ErrFieldReadOnly{Field: "open.file.name.length"}
	case "open.file.path":
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "Open.File.PathnameStr"}
		}
		ev.Open.File.PathnameStr = rv
		return nil
	case "open.file.path.length":
		return &eval.ErrFieldReadOnly{Field: "open.file.path.length"}
	case "open.flags":
		rv, ok := value.(int)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "Open.Flags"}
		}
		ev.Open.Flags = uint32(rv)
		return nil
	case "process.ancestors.args":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Ancestor == nil {
			ev.BaseEvent.ProcessContext.Ancestor = &ProcessCacheEntry{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Args"}
		}
		ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Args = rv
		return nil
	case "process.ancestors.args_flags":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Ancestor == nil {
			ev.BaseEvent.ProcessContext.Ancestor = &ProcessCacheEntry{}
		}
		switch rv := value.(type) {
		case string:
			ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Argv = append(ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Argv, rv)
		case []string:
			ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Argv = append(ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Argv, rv...)
		default:
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Argv"}
		}
		return nil
	case "process.ancestors.args_options":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Ancestor == nil {
			ev.BaseEvent.ProcessContext.Ancestor = &ProcessCacheEntry{}
		}
		switch rv := value.(type) {
		case string:
			ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Argv = append(ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Argv, rv)
		case []string:
			ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Argv = append(ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Argv, rv...)
		default:
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Argv"}
		}
		return nil
	case "process.ancestors.args_truncated":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Ancestor == nil {
			ev.BaseEvent.ProcessContext.Ancestor = &ProcessCacheEntry{}
		}
		rv, ok := value.(bool)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.ArgsTruncated"}
		}
		ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.ArgsTruncated = rv
		return nil
	case "process.ancestors.argv":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Ancestor == nil {
			ev.BaseEvent.ProcessContext.Ancestor = &ProcessCacheEntry{}
		}
		switch rv := value.(type) {
		case string:
			ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Argv = append(ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Argv, rv)
		case []string:
			ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Argv = append(ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Argv, rv...)
		default:
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Argv"}
		}
		return nil
	case "process.ancestors.argv0":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Ancestor == nil {
			ev.BaseEvent.ProcessContext.Ancestor = &ProcessCacheEntry{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Argv0"}
		}
		ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Argv0 = rv
		return nil
	case "process.ancestors.comm":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Ancestor == nil {
			ev.BaseEvent.ProcessContext.Ancestor = &ProcessCacheEntry{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Comm"}
		}
		ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Comm = rv
		return nil
	case "process.ancestors.container.id":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Ancestor == nil {
			ev.BaseEvent.ProcessContext.Ancestor = &ProcessCacheEntry{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.ContainerID"}
		}
		ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.ContainerID = rv
		return nil
	case "process.ancestors.created_at":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Ancestor == nil {
			ev.BaseEvent.ProcessContext.Ancestor = &ProcessCacheEntry{}
		}
		rv, ok := value.(int)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.CreatedAt"}
		}
		ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.CreatedAt = uint64(rv)
		return nil
	case "process.ancestors.envp":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Ancestor == nil {
			ev.BaseEvent.ProcessContext.Ancestor = &ProcessCacheEntry{}
		}
		switch rv := value.(type) {
		case string:
			ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Envp = append(ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Envp, rv)
		case []string:
			ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Envp = append(ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Envp, rv...)
		default:
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Envp"}
		}
		return nil
	case "process.ancestors.envs":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Ancestor == nil {
			ev.BaseEvent.ProcessContext.Ancestor = &ProcessCacheEntry{}
		}
		switch rv := value.(type) {
		case string:
			ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Envs = append(ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Envs, rv)
		case []string:
			ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Envs = append(ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Envs, rv...)
		default:
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.Envs"}
		}
		return nil
	case "process.ancestors.envs_truncated":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Ancestor == nil {
			ev.BaseEvent.ProcessContext.Ancestor = &ProcessCacheEntry{}
		}
		rv, ok := value.(bool)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.EnvsTruncated"}
		}
		ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.EnvsTruncated = rv
		return nil
	case "process.ancestors.file.name":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Ancestor == nil {
			ev.BaseEvent.ProcessContext.Ancestor = &ProcessCacheEntry{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.FileEvent.BasenameStr"}
		}
		ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.FileEvent.BasenameStr = rv
		return nil
	case "process.ancestors.file.name.length":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Ancestor == nil {
			ev.BaseEvent.ProcessContext.Ancestor = &ProcessCacheEntry{}
		}
		return &eval.ErrFieldReadOnly{Field: "process.ancestors.file.name.length"}
	case "process.ancestors.file.path":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Ancestor == nil {
			ev.BaseEvent.ProcessContext.Ancestor = &ProcessCacheEntry{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.FileEvent.PathnameStr"}
		}
		ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.FileEvent.PathnameStr = rv
		return nil
	case "process.ancestors.file.path.length":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Ancestor == nil {
			ev.BaseEvent.ProcessContext.Ancestor = &ProcessCacheEntry{}
		}
		return &eval.ErrFieldReadOnly{Field: "process.ancestors.file.path.length"}
	case "process.ancestors.pid":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Ancestor == nil {
			ev.BaseEvent.ProcessContext.Ancestor = &ProcessCacheEntry{}
		}
		rv, ok := value.(int)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.PIDContext.Pid"}
		}
		ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.PIDContext.Pid = uint32(rv)
		return nil
	case "process.ancestors.ppid":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Ancestor == nil {
			ev.BaseEvent.ProcessContext.Ancestor = &ProcessCacheEntry{}
		}
		rv, ok := value.(int)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.PPid"}
		}
		ev.BaseEvent.ProcessContext.Ancestor.ProcessContext.Process.PPid = uint32(rv)
		return nil
	case "process.args":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Process.Args"}
		}
		ev.BaseEvent.ProcessContext.Process.Args = rv
		return nil
	case "process.args_flags":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		switch rv := value.(type) {
		case string:
			ev.BaseEvent.ProcessContext.Process.Argv = append(ev.BaseEvent.ProcessContext.Process.Argv, rv)
		case []string:
			ev.BaseEvent.ProcessContext.Process.Argv = append(ev.BaseEvent.ProcessContext.Process.Argv, rv...)
		default:
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Process.Argv"}
		}
		return nil
	case "process.args_options":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		switch rv := value.(type) {
		case string:
			ev.BaseEvent.ProcessContext.Process.Argv = append(ev.BaseEvent.ProcessContext.Process.Argv, rv)
		case []string:
			ev.BaseEvent.ProcessContext.Process.Argv = append(ev.BaseEvent.ProcessContext.Process.Argv, rv...)
		default:
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Process.Argv"}
		}
		return nil
	case "process.args_truncated":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		rv, ok := value.(bool)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Process.ArgsTruncated"}
		}
		ev.BaseEvent.ProcessContext.Process.ArgsTruncated = rv
		return nil
	case "process.argv":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		switch rv := value.(type) {
		case string:
			ev.BaseEvent.ProcessContext.Process.Argv = append(ev.BaseEvent.ProcessContext.Process.Argv, rv)
		case []string:
			ev.BaseEvent.ProcessContext.Process.Argv = append(ev.BaseEvent.ProcessContext.Process.Argv, rv...)
		default:
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Process.Argv"}
		}
		return nil
	case "process.argv0":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Process.Argv0"}
		}
		ev.BaseEvent.ProcessContext.Process.Argv0 = rv
		return nil
	case "process.comm":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Process.Comm"}
		}
		ev.BaseEvent.ProcessContext.Process.Comm = rv
		return nil
	case "process.container.id":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Process.ContainerID"}
		}
		ev.BaseEvent.ProcessContext.Process.ContainerID = rv
		return nil
	case "process.created_at":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		rv, ok := value.(int)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Process.CreatedAt"}
		}
		ev.BaseEvent.ProcessContext.Process.CreatedAt = uint64(rv)
		return nil
	case "process.envp":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		switch rv := value.(type) {
		case string:
			ev.BaseEvent.ProcessContext.Process.Envp = append(ev.BaseEvent.ProcessContext.Process.Envp, rv)
		case []string:
			ev.BaseEvent.ProcessContext.Process.Envp = append(ev.BaseEvent.ProcessContext.Process.Envp, rv...)
		default:
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Process.Envp"}
		}
		return nil
	case "process.envs":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		switch rv := value.(type) {
		case string:
			ev.BaseEvent.ProcessContext.Process.Envs = append(ev.BaseEvent.ProcessContext.Process.Envs, rv)
		case []string:
			ev.BaseEvent.ProcessContext.Process.Envs = append(ev.BaseEvent.ProcessContext.Process.Envs, rv...)
		default:
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Process.Envs"}
		}
		return nil
	case "process.envs_truncated":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		rv, ok := value.(bool)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Process.EnvsTruncated"}
		}
		ev.BaseEvent.ProcessContext.Process.EnvsTruncated = rv
		return nil
	case "process.file.name":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Process.FileEvent.BasenameStr"}
		}
		ev.BaseEvent.ProcessContext.Process.FileEvent.BasenameStr = rv
		return nil
	case "process.file.name.length":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		return &eval.ErrFieldReadOnly{Field: "process.file.name.length"}
	case "process.file.path":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Process.FileEvent.PathnameStr"}
		}
		ev.BaseEvent.ProcessContext.Process.FileEvent.PathnameStr = rv
		return nil
	case "process.file.path.length":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		return &eval.ErrFieldReadOnly{Field: "process.file.path.length"}
	case "process.parent.args":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Parent == nil {
			ev.BaseEvent.ProcessContext.Parent = &Process{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Parent.Args"}
		}
		ev.BaseEvent.ProcessContext.Parent.Args = rv
		return nil
	case "process.parent.args_flags":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Parent == nil {
			ev.BaseEvent.ProcessContext.Parent = &Process{}
		}
		switch rv := value.(type) {
		case string:
			ev.BaseEvent.ProcessContext.Parent.Argv = append(ev.BaseEvent.ProcessContext.Parent.Argv, rv)
		case []string:
			ev.BaseEvent.ProcessContext.Parent.Argv = append(ev.BaseEvent.ProcessContext.Parent.Argv, rv...)
		default:
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Parent.Argv"}
		}
		return nil
	case "process.parent.args_options":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Parent == nil {
			ev.BaseEvent.ProcessContext.Parent = &Process{}
		}
		switch rv := value.(type) {
		case string:
			ev.BaseEvent.ProcessContext.Parent.Argv = append(ev.BaseEvent.ProcessContext.Parent.Argv, rv)
		case []string:
			ev.BaseEvent.ProcessContext.Parent.Argv = append(ev.BaseEvent.ProcessContext.Parent.Argv, rv...)
		default:
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Parent.Argv"}
		}
		return nil
	case "process.parent.args_truncated":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Parent == nil {
			ev.BaseEvent.ProcessContext.Parent = &Process{}
		}
		rv, ok := value.(bool)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Parent.ArgsTruncated"}
		}
		ev.BaseEvent.ProcessContext.Parent.ArgsTruncated = rv
		return nil
	case "process.parent.argv":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Parent == nil {
			ev.BaseEvent.ProcessContext.Parent = &Process{}
		}
		switch rv := value.(type) {
		case string:
			ev.BaseEvent.ProcessContext.Parent.Argv = append(ev.BaseEvent.ProcessContext.Parent.Argv, rv)
		case []string:
			ev.BaseEvent.ProcessContext.Parent.Argv = append(ev.BaseEvent.ProcessContext.Parent.Argv, rv...)
		default:
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Parent.Argv"}
		}
		return nil
	case "process.parent.argv0":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Parent == nil {
			ev.BaseEvent.ProcessContext.Parent = &Process{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Parent.Argv0"}
		}
		ev.BaseEvent.ProcessContext.Parent.Argv0 = rv
		return nil
	case "process.parent.comm":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Parent == nil {
			ev.BaseEvent.ProcessContext.Parent = &Process{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Parent.Comm"}
		}
		ev.BaseEvent.ProcessContext.Parent.Comm = rv
		return nil
	case "process.parent.container.id":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Parent == nil {
			ev.BaseEvent.ProcessContext.Parent = &Process{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Parent.ContainerID"}
		}
		ev.BaseEvent.ProcessContext.Parent.ContainerID = rv
		return nil
	case "process.parent.created_at":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Parent == nil {
			ev.BaseEvent.ProcessContext.Parent = &Process{}
		}
		rv, ok := value.(int)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Parent.CreatedAt"}
		}
		ev.BaseEvent.ProcessContext.Parent.CreatedAt = uint64(rv)
		return nil
	case "process.parent.envp":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Parent == nil {
			ev.BaseEvent.ProcessContext.Parent = &Process{}
		}
		switch rv := value.(type) {
		case string:
			ev.BaseEvent.ProcessContext.Parent.Envp = append(ev.BaseEvent.ProcessContext.Parent.Envp, rv)
		case []string:
			ev.BaseEvent.ProcessContext.Parent.Envp = append(ev.BaseEvent.ProcessContext.Parent.Envp, rv...)
		default:
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Parent.Envp"}
		}
		return nil
	case "process.parent.envs":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Parent == nil {
			ev.BaseEvent.ProcessContext.Parent = &Process{}
		}
		switch rv := value.(type) {
		case string:
			ev.BaseEvent.ProcessContext.Parent.Envs = append(ev.BaseEvent.ProcessContext.Parent.Envs, rv)
		case []string:
			ev.BaseEvent.ProcessContext.Parent.Envs = append(ev.BaseEvent.ProcessContext.Parent.Envs, rv...)
		default:
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Parent.Envs"}
		}
		return nil
	case "process.parent.envs_truncated":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Parent == nil {
			ev.BaseEvent.ProcessContext.Parent = &Process{}
		}
		rv, ok := value.(bool)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Parent.EnvsTruncated"}
		}
		ev.BaseEvent.ProcessContext.Parent.EnvsTruncated = rv
		return nil
	case "process.parent.file.name":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Parent == nil {
			ev.BaseEvent.ProcessContext.Parent = &Process{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Parent.FileEvent.BasenameStr"}
		}
		ev.BaseEvent.ProcessContext.Parent.FileEvent.BasenameStr = rv
		return nil
	case "process.parent.file.name.length":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Parent == nil {
			ev.BaseEvent.ProcessContext.Parent = &Process{}
		}
		return &eval.ErrFieldReadOnly{Field: "process.parent.file.name.length"}
	case "process.parent.file.path":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Parent == nil {
			ev.BaseEvent.ProcessContext.Parent = &Process{}
		}
		rv, ok := value.(string)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Parent.FileEvent.PathnameStr"}
		}
		ev.BaseEvent.ProcessContext.Parent.FileEvent.PathnameStr = rv
		return nil
	case "process.parent.file.path.length":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Parent == nil {
			ev.BaseEvent.ProcessContext.Parent = &Process{}
		}
		return &eval.ErrFieldReadOnly{Field: "process.parent.file.path.length"}
	case "process.parent.pid":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Parent == nil {
			ev.BaseEvent.ProcessContext.Parent = &Process{}
		}
		rv, ok := value.(int)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Parent.PIDContext.Pid"}
		}
		ev.BaseEvent.ProcessContext.Parent.PIDContext.Pid = uint32(rv)
		return nil
	case "process.parent.ppid":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		if ev.BaseEvent.ProcessContext.Parent == nil {
			ev.BaseEvent.ProcessContext.Parent = &Process{}
		}
		rv, ok := value.(int)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Parent.PPid"}
		}
		ev.BaseEvent.ProcessContext.Parent.PPid = uint32(rv)
		return nil
	case "process.pid":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		rv, ok := value.(int)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Process.PIDContext.Pid"}
		}
		ev.BaseEvent.ProcessContext.Process.PIDContext.Pid = uint32(rv)
		return nil
	case "process.ppid":
		if ev.BaseEvent.ProcessContext == nil {
			ev.BaseEvent.ProcessContext = &ProcessContext{}
		}
		rv, ok := value.(int)
		if !ok {
			return &eval.ErrValueTypeMismatch{Field: "BaseEvent.ProcessContext.Process.PPid"}
		}
		ev.BaseEvent.ProcessContext.Process.PPid = uint32(rv)
		return nil
	}
	return &eval.ErrFieldNotFound{Field: field}
}