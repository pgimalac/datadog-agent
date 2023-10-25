// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package payload defines the JSON payload we send to the events platform.
package payload

import (
	"encoding/json"
	"reflect"
	"strings"
)

// Device contains device details (device sending NetFlow flows)
type Device struct {
	Namespace string `json:"namespace"`
}

// Exporter contains NetFlow exporter details
type Exporter struct {
	IP string `json:"ip"`
}

// Endpoint contains source or destination endpoint details
type Endpoint struct {
	IP   string `json:"ip"`
	Port string `json:"port"` // Port number can be zero/positive or `*` (ephemeral port)
	Mac  string `json:"mac"`
	Mask string `json:"mask"`
}

// NextHop contains next hop details
type NextHop struct {
	IP string `json:"ip"`
}

// Interface contains interface details
type Interface struct {
	Index uint32 `json:"index"`
}

// ObservationPoint contains ingress or egress observation point
type ObservationPoint struct {
	Interface Interface `json:"interface"`
}

// AdditionalFields contains additional configured fields
type AdditionalFields = map[string]any

// FlowPayload contains network devices flows
type FlowPayload struct {
	FlushTimestamp   int64            `json:"flush_timestamp"`
	FlowType         string           `json:"type"`
	SamplingRate     uint64           `json:"sampling_rate"`
	Direction        string           `json:"direction"`
	Start            uint64           `json:"start"` // in seconds
	End              uint64           `json:"end"`   // in seconds
	Bytes            uint64           `json:"bytes"`
	Packets          uint64           `json:"packets"`
	EtherType        string           `json:"ether_type,omitempty"`
	IPProtocol       string           `json:"ip_protocol"`
	Device           Device           `json:"device"`
	Exporter         Exporter         `json:"exporter"`
	Source           Endpoint         `json:"source"`
	Destination      Endpoint         `json:"destination"`
	Ingress          ObservationPoint `json:"ingress"`
	Egress           ObservationPoint `json:"egress"`
	Host             string           `json:"host"`
	TCPFlags         []string         `json:"tcp_flags,omitempty"`
	NextHop          NextHop          `json:"next_hop,omitempty"`
	AdditionalFields AdditionalFields `json:"_"`
}

// MarshalJSON Custom marshaller that moves AdditionalFields to the root of the payload
func (p FlowPayload) MarshalWithAdditionalFields() ([]byte, error) {
	// Turn FlowPayload into a map
	b, err := json.Marshal(p)

	if err != nil {
		return nil, err
	}

	var m map[string]json.RawMessage
	err = json.Unmarshal(b, &m)

	if err != nil {
		return nil, err
	}

	// Move additional fields to the root of the payload
	for k, v := range p.AdditionalFields {
		b, err = json.Marshal(v)
		if err != nil {
			// We failed to marshall a custom field, continuing anyway TODO : log
			continue
		}
		m[k] = b
	}

	return json.Marshal(m)
}

// MarshalJSON Custom marshaller that moves AdditionalFields to the root of the payload
func (p FlowPayload) MarshalWithAdditionalFieldsLessMarshall() ([]byte, error) {
	marshaled, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	fields := make(map[string]any)
	err = json.Unmarshal(marshaled, &fields)
	if err != nil {
		return nil, err
	}

	delete(fields, "additional_fields")

	for k, v := range p.AdditionalFields {
		fields[k] = v
	}

	return json.Marshal(fields)
}

// MarshalJSON Custom marshaller that moves AdditionalFields to the root of the payload
func (p FlowPayload) MarshalWithAdditionalFieldsReflect() ([]byte, error) {
	v := reflect.ValueOf(p)
	typeOfV := v.Type()

	fields := make(map[string]any)

	for i := 0; i < v.NumField(); i++ {
		jsonTag := typeOfV.Field(i).Tag.Get("json")
		jsonTag = strings.TrimRight(jsonTag, ",omitempty")
		value := v.Field(i).Interface()
		if jsonTag != "additional_fields" && value != nil {
			fields[jsonTag] = v.Field(i).Interface()
		}
	}

	for k, v := range p.AdditionalFields {
		if _, ok := fields[k]; ok {
			continue
		}
		fields[k] = v
	}

	return json.Marshal(fields)
}

// MarshalJSON Custom marshaller that moves AdditionalFields to the root of the payload
func (p FlowPayload) MarshalJSON() ([]byte, error) {
	fields := map[string]any{
		"flush_timestamp": p.FlushTimestamp,
		"type":            p.FlowType,
		"sampling_rate":   p.SamplingRate,
		"direction":       p.Direction,
		"start":           p.Start,
		"end":             p.End,
		"bytes":           p.Bytes,
		"packets":         p.Packets,
		"ip_protocol":     p.IPProtocol,
		"device":          p.Device,
		"exporter":        p.Exporter,
		"source":          p.Source,
		"destination":     p.Destination,
		"ingress":         p.Ingress,
		"egress":          p.Egress,
		"host":            p.Host,
		"next_hop":        p.NextHop,
	}

	// omit empty
	if p.EtherType != "" {
		fields["ether_type"] = p.EtherType
	}

	// omit empty
	if p.TCPFlags != nil {
		fields["tcp_flags"] = p.TCPFlags
	}

	// Adding additional fields
	for k, v := range p.AdditionalFields {
		if _, ok := fields[k]; ok {
			// Do not override, this is handled in goflowlib/convert.go
			continue
		}
		fields[k] = v
	}

	return json.Marshal(fields)
}
