//go:build windows
// +build windows

// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package serializers

import (
	json "encoding/json"
	utils "github.com/DataDog/datadog-agent/pkg/security/utils"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson6b24c4ebDecodeGithubComDataDogDatadogAgentPkgSecuritySerializers(in *jlexer.Lexer, out *ProcessContextSerializer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	out.ProcessSerializer = new(ProcessSerializer)
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "parent":
			if in.IsNull() {
				in.Skip()
				out.Parent = nil
			} else {
				if out.Parent == nil {
					out.Parent = new(ProcessSerializer)
				}
				(*out.Parent).UnmarshalEasyJSON(in)
			}
		case "ancestors":
			if in.IsNull() {
				in.Skip()
				out.Ancestors = nil
			} else {
				in.Delim('[')
				if out.Ancestors == nil {
					if !in.IsDelim(']') {
						out.Ancestors = make([]*ProcessSerializer, 0, 8)
					} else {
						out.Ancestors = []*ProcessSerializer{}
					}
				} else {
					out.Ancestors = (out.Ancestors)[:0]
				}
				for !in.IsDelim(']') {
					var v1 *ProcessSerializer
					if in.IsNull() {
						in.Skip()
						v1 = nil
					} else {
						if v1 == nil {
							v1 = new(ProcessSerializer)
						}
						(*v1).UnmarshalEasyJSON(in)
					}
					out.Ancestors = append(out.Ancestors, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "pid":
			out.Pid = uint32(in.Uint32())
		case "ppid":
			if in.IsNull() {
				in.Skip()
				out.PPid = nil
			} else {
				if out.PPid == nil {
					out.PPid = new(uint32)
				}
				*out.PPid = uint32(in.Uint32())
			}
		case "exec_time":
			if in.IsNull() {
				in.Skip()
				out.ExecTime = nil
			} else {
				if out.ExecTime == nil {
					out.ExecTime = new(utils.EasyjsonTime)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.ExecTime).UnmarshalJSON(data))
				}
			}
		case "exit_time":
			if in.IsNull() {
				in.Skip()
				out.ExitTime = nil
			} else {
				if out.ExitTime == nil {
					out.ExitTime = new(utils.EasyjsonTime)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.ExitTime).UnmarshalJSON(data))
				}
			}
		case "executable":
			if in.IsNull() {
				in.Skip()
				out.Executable = nil
			} else {
				if out.Executable == nil {
					out.Executable = new(FileSerializer)
				}
				(*out.Executable).UnmarshalEasyJSON(in)
			}
		case "container":
			if in.IsNull() {
				in.Skip()
				out.Container = nil
			} else {
				if out.Container == nil {
					out.Container = new(ContainerContextSerializer)
				}
				(*out.Container).UnmarshalEasyJSON(in)
			}
		case "cmdline":
			out.CmdLine = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6b24c4ebEncodeGithubComDataDogDatadogAgentPkgSecuritySerializers(out *jwriter.Writer, in ProcessContextSerializer) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Parent != nil {
		const prefix string = ",\"parent\":"
		first = false
		out.RawString(prefix[1:])
		(*in.Parent).MarshalEasyJSON(out)
	}
	if len(in.Ancestors) != 0 {
		const prefix string = ",\"ancestors\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v2, v3 := range in.Ancestors {
				if v2 > 0 {
					out.RawByte(',')
				}
				if v3 == nil {
					out.RawString("null")
				} else {
					(*v3).MarshalEasyJSON(out)
				}
			}
			out.RawByte(']')
		}
	}
	if in.Pid != 0 {
		const prefix string = ",\"pid\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint32(uint32(in.Pid))
	}
	if in.PPid != nil {
		const prefix string = ",\"ppid\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint32(uint32(*in.PPid))
	}
	if in.ExecTime != nil {
		const prefix string = ",\"exec_time\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.ExecTime).MarshalEasyJSON(out)
	}
	if in.ExitTime != nil {
		const prefix string = ",\"exit_time\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.ExitTime).MarshalEasyJSON(out)
	}
	if in.Executable != nil {
		const prefix string = ",\"executable\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.Executable).MarshalEasyJSON(out)
	}
	if in.Container != nil {
		const prefix string = ",\"container\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.Container).MarshalEasyJSON(out)
	}
	if in.CmdLine != "" {
		const prefix string = ",\"cmdline\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.CmdLine))
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ProcessContextSerializer) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6b24c4ebEncodeGithubComDataDogDatadogAgentPkgSecuritySerializers(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ProcessContextSerializer) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6b24c4ebDecodeGithubComDataDogDatadogAgentPkgSecuritySerializers(l, v)
}
func easyjson6b24c4ebDecodeGithubComDataDogDatadogAgentPkgSecuritySerializers1(in *jlexer.Lexer, out *NetworkContextSerializer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "device":
			if in.IsNull() {
				in.Skip()
				out.Device = nil
			} else {
				if out.Device == nil {
					out.Device = new(NetworkDeviceSerializer)
				}
				(*out.Device).UnmarshalEasyJSON(in)
			}
		case "l3_protocol":
			out.L3Protocol = string(in.String())
		case "l4_protocol":
			out.L4Protocol = string(in.String())
		case "source":
			(out.Source).UnmarshalEasyJSON(in)
		case "destination":
			(out.Destination).UnmarshalEasyJSON(in)
		case "size":
			out.Size = uint32(in.Uint32())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6b24c4ebEncodeGithubComDataDogDatadogAgentPkgSecuritySerializers1(out *jwriter.Writer, in NetworkContextSerializer) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Device != nil {
		const prefix string = ",\"device\":"
		first = false
		out.RawString(prefix[1:])
		(*in.Device).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"l3_protocol\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.L3Protocol))
	}
	{
		const prefix string = ",\"l4_protocol\":"
		out.RawString(prefix)
		out.String(string(in.L4Protocol))
	}
	{
		const prefix string = ",\"source\":"
		out.RawString(prefix)
		(in.Source).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"destination\":"
		out.RawString(prefix)
		(in.Destination).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"size\":"
		out.RawString(prefix)
		out.Uint32(uint32(in.Size))
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v NetworkContextSerializer) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6b24c4ebEncodeGithubComDataDogDatadogAgentPkgSecuritySerializers1(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *NetworkContextSerializer) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6b24c4ebDecodeGithubComDataDogDatadogAgentPkgSecuritySerializers1(l, v)
}
func easyjson6b24c4ebDecodeGithubComDataDogDatadogAgentPkgSecuritySerializers2(in *jlexer.Lexer, out *MatchedRuleSerializer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = string(in.String())
		case "version":
			out.Version = string(in.String())
		case "tags":
			if in.IsNull() {
				in.Skip()
				out.Tags = nil
			} else {
				in.Delim('[')
				if out.Tags == nil {
					if !in.IsDelim(']') {
						out.Tags = make([]string, 0, 4)
					} else {
						out.Tags = []string{}
					}
				} else {
					out.Tags = (out.Tags)[:0]
				}
				for !in.IsDelim(']') {
					var v4 string
					v4 = string(in.String())
					out.Tags = append(out.Tags, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "policy_name":
			out.PolicyName = string(in.String())
		case "policy_version":
			out.PolicyVersion = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6b24c4ebEncodeGithubComDataDogDatadogAgentPkgSecuritySerializers2(out *jwriter.Writer, in MatchedRuleSerializer) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != "" {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	if in.Version != "" {
		const prefix string = ",\"version\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Version))
	}
	if len(in.Tags) != 0 {
		const prefix string = ",\"tags\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v5, v6 := range in.Tags {
				if v5 > 0 {
					out.RawByte(',')
				}
				out.String(string(v6))
			}
			out.RawByte(']')
		}
	}
	if in.PolicyName != "" {
		const prefix string = ",\"policy_name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.PolicyName))
	}
	if in.PolicyVersion != "" {
		const prefix string = ",\"policy_version\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.PolicyVersion))
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v MatchedRuleSerializer) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6b24c4ebEncodeGithubComDataDogDatadogAgentPkgSecuritySerializers2(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *MatchedRuleSerializer) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6b24c4ebDecodeGithubComDataDogDatadogAgentPkgSecuritySerializers2(l, v)
}
func easyjson6b24c4ebDecodeGithubComDataDogDatadogAgentPkgSecuritySerializers3(in *jlexer.Lexer, out *IPPortSerializer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "ip":
			out.IP = string(in.String())
		case "port":
			out.Port = uint16(in.Uint16())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6b24c4ebEncodeGithubComDataDogDatadogAgentPkgSecuritySerializers3(out *jwriter.Writer, in IPPortSerializer) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"ip\":"
		out.RawString(prefix[1:])
		out.String(string(in.IP))
	}
	{
		const prefix string = ",\"port\":"
		out.RawString(prefix)
		out.Uint16(uint16(in.Port))
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v IPPortSerializer) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6b24c4ebEncodeGithubComDataDogDatadogAgentPkgSecuritySerializers3(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *IPPortSerializer) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6b24c4ebDecodeGithubComDataDogDatadogAgentPkgSecuritySerializers3(l, v)
}
func easyjson6b24c4ebDecodeGithubComDataDogDatadogAgentPkgSecuritySerializers4(in *jlexer.Lexer, out *IPPortFamilySerializer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "family":
			out.Family = string(in.String())
		case "ip":
			out.IP = string(in.String())
		case "port":
			out.Port = uint16(in.Uint16())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6b24c4ebEncodeGithubComDataDogDatadogAgentPkgSecuritySerializers4(out *jwriter.Writer, in IPPortFamilySerializer) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"family\":"
		out.RawString(prefix[1:])
		out.String(string(in.Family))
	}
	{
		const prefix string = ",\"ip\":"
		out.RawString(prefix)
		out.String(string(in.IP))
	}
	{
		const prefix string = ",\"port\":"
		out.RawString(prefix)
		out.Uint16(uint16(in.Port))
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v IPPortFamilySerializer) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6b24c4ebEncodeGithubComDataDogDatadogAgentPkgSecuritySerializers4(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *IPPortFamilySerializer) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6b24c4ebDecodeGithubComDataDogDatadogAgentPkgSecuritySerializers4(l, v)
}
func easyjson6b24c4ebDecodeGithubComDataDogDatadogAgentPkgSecuritySerializers5(in *jlexer.Lexer, out *ExitEventSerializer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "cause":
			out.Cause = string(in.String())
		case "code":
			out.Code = uint32(in.Uint32())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6b24c4ebEncodeGithubComDataDogDatadogAgentPkgSecuritySerializers5(out *jwriter.Writer, in ExitEventSerializer) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"cause\":"
		out.RawString(prefix[1:])
		out.String(string(in.Cause))
	}
	{
		const prefix string = ",\"code\":"
		out.RawString(prefix)
		out.Uint32(uint32(in.Code))
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ExitEventSerializer) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6b24c4ebEncodeGithubComDataDogDatadogAgentPkgSecuritySerializers5(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ExitEventSerializer) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6b24c4ebDecodeGithubComDataDogDatadogAgentPkgSecuritySerializers5(l, v)
}
func easyjson6b24c4ebDecodeGithubComDataDogDatadogAgentPkgSecuritySerializers6(in *jlexer.Lexer, out *EventContextSerializer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "name":
			out.Name = string(in.String())
		case "category":
			out.Category = string(in.String())
		case "outcome":
			out.Outcome = string(in.String())
		case "async":
			out.Async = bool(in.Bool())
		case "matched_rules":
			if in.IsNull() {
				in.Skip()
				out.MatchedRules = nil
			} else {
				in.Delim('[')
				if out.MatchedRules == nil {
					if !in.IsDelim(']') {
						out.MatchedRules = make([]MatchedRuleSerializer, 0, 0)
					} else {
						out.MatchedRules = []MatchedRuleSerializer{}
					}
				} else {
					out.MatchedRules = (out.MatchedRules)[:0]
				}
				for !in.IsDelim(']') {
					var v7 MatchedRuleSerializer
					(v7).UnmarshalEasyJSON(in)
					out.MatchedRules = append(out.MatchedRules, v7)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6b24c4ebEncodeGithubComDataDogDatadogAgentPkgSecuritySerializers6(out *jwriter.Writer, in EventContextSerializer) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Name != "" {
		const prefix string = ",\"name\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	if in.Category != "" {
		const prefix string = ",\"category\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Category))
	}
	if in.Outcome != "" {
		const prefix string = ",\"outcome\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Outcome))
	}
	if in.Async {
		const prefix string = ",\"async\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.Async))
	}
	if len(in.MatchedRules) != 0 {
		const prefix string = ",\"matched_rules\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v8, v9 := range in.MatchedRules {
				if v8 > 0 {
					out.RawByte(',')
				}
				(v9).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EventContextSerializer) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6b24c4ebEncodeGithubComDataDogDatadogAgentPkgSecuritySerializers6(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EventContextSerializer) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6b24c4ebDecodeGithubComDataDogDatadogAgentPkgSecuritySerializers6(l, v)
}
func easyjson6b24c4ebDecodeGithubComDataDogDatadogAgentPkgSecuritySerializers7(in *jlexer.Lexer, out *DNSQuestionSerializer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "class":
			out.Class = string(in.String())
		case "type":
			out.Type = string(in.String())
		case "name":
			out.Name = string(in.String())
		case "size":
			out.Size = uint16(in.Uint16())
		case "count":
			out.Count = uint16(in.Uint16())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6b24c4ebEncodeGithubComDataDogDatadogAgentPkgSecuritySerializers7(out *jwriter.Writer, in DNSQuestionSerializer) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"class\":"
		out.RawString(prefix[1:])
		out.String(string(in.Class))
	}
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix)
		out.String(string(in.Type))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"size\":"
		out.RawString(prefix)
		out.Uint16(uint16(in.Size))
	}
	{
		const prefix string = ",\"count\":"
		out.RawString(prefix)
		out.Uint16(uint16(in.Count))
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v DNSQuestionSerializer) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6b24c4ebEncodeGithubComDataDogDatadogAgentPkgSecuritySerializers7(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *DNSQuestionSerializer) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6b24c4ebDecodeGithubComDataDogDatadogAgentPkgSecuritySerializers7(l, v)
}
func easyjson6b24c4ebDecodeGithubComDataDogDatadogAgentPkgSecuritySerializers8(in *jlexer.Lexer, out *DNSEventSerializer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = uint16(in.Uint16())
		case "question":
			(out.Question).UnmarshalEasyJSON(in)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6b24c4ebEncodeGithubComDataDogDatadogAgentPkgSecuritySerializers8(out *jwriter.Writer, in DNSEventSerializer) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint16(uint16(in.ID))
	}
	{
		const prefix string = ",\"question\":"
		out.RawString(prefix)
		(in.Question).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v DNSEventSerializer) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6b24c4ebEncodeGithubComDataDogDatadogAgentPkgSecuritySerializers8(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *DNSEventSerializer) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6b24c4ebDecodeGithubComDataDogDatadogAgentPkgSecuritySerializers8(l, v)
}
func easyjson6b24c4ebDecodeGithubComDataDogDatadogAgentPkgSecuritySerializers9(in *jlexer.Lexer, out *DDContextSerializer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "span_id":
			out.SpanID = uint64(in.Uint64())
		case "trace_id":
			out.TraceID = uint64(in.Uint64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6b24c4ebEncodeGithubComDataDogDatadogAgentPkgSecuritySerializers9(out *jwriter.Writer, in DDContextSerializer) {
	out.RawByte('{')
	first := true
	_ = first
	if in.SpanID != 0 {
		const prefix string = ",\"span_id\":"
		first = false
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.SpanID))
	}
	if in.TraceID != 0 {
		const prefix string = ",\"trace_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint64(uint64(in.TraceID))
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v DDContextSerializer) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6b24c4ebEncodeGithubComDataDogDatadogAgentPkgSecuritySerializers9(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *DDContextSerializer) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6b24c4ebDecodeGithubComDataDogDatadogAgentPkgSecuritySerializers9(l, v)
}
func easyjson6b24c4ebDecodeGithubComDataDogDatadogAgentPkgSecuritySerializers10(in *jlexer.Lexer, out *ContainerContextSerializer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = string(in.String())
		case "created_at":
			if in.IsNull() {
				in.Skip()
				out.CreatedAt = nil
			} else {
				if out.CreatedAt == nil {
					out.CreatedAt = new(utils.EasyjsonTime)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.CreatedAt).UnmarshalJSON(data))
				}
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6b24c4ebEncodeGithubComDataDogDatadogAgentPkgSecuritySerializers10(out *jwriter.Writer, in ContainerContextSerializer) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != "" {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	if in.CreatedAt != nil {
		const prefix string = ",\"created_at\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.CreatedAt).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ContainerContextSerializer) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6b24c4ebEncodeGithubComDataDogDatadogAgentPkgSecuritySerializers10(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ContainerContextSerializer) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6b24c4ebDecodeGithubComDataDogDatadogAgentPkgSecuritySerializers10(l, v)
}
func easyjson6b24c4ebDecodeGithubComDataDogDatadogAgentPkgSecuritySerializers11(in *jlexer.Lexer, out *BaseEventSerializer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	out.FileEventSerializer = new(FileEventSerializer)
	out.DNSEventSerializer = new(DNSEventSerializer)
	out.NetworkContextSerializer = new(NetworkContextSerializer)
	out.ExitEventSerializer = new(ExitEventSerializer)
	out.ProcessContextSerializer = new(ProcessContextSerializer)
	out.DDContextSerializer = new(DDContextSerializer)
	out.ContainerContextSerializer = new(ContainerContextSerializer)
	out.SecurityProfileContextSerializer = new(SecurityProfileContextSerializer)
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "evt":
			(out.EventContextSerializer).UnmarshalEasyJSON(in)
		case "date":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Date).UnmarshalJSON(data))
			}
		case "file":
			if in.IsNull() {
				in.Skip()
				out.FileEventSerializer = nil
			} else {
				if out.FileEventSerializer == nil {
					out.FileEventSerializer = new(FileEventSerializer)
				}
				(*out.FileEventSerializer).UnmarshalEasyJSON(in)
			}
		case "dns":
			if in.IsNull() {
				in.Skip()
				out.DNSEventSerializer = nil
			} else {
				if out.DNSEventSerializer == nil {
					out.DNSEventSerializer = new(DNSEventSerializer)
				}
				(*out.DNSEventSerializer).UnmarshalEasyJSON(in)
			}
		case "network":
			if in.IsNull() {
				in.Skip()
				out.NetworkContextSerializer = nil
			} else {
				if out.NetworkContextSerializer == nil {
					out.NetworkContextSerializer = new(NetworkContextSerializer)
				}
				(*out.NetworkContextSerializer).UnmarshalEasyJSON(in)
			}
		case "exit":
			if in.IsNull() {
				in.Skip()
				out.ExitEventSerializer = nil
			} else {
				if out.ExitEventSerializer == nil {
					out.ExitEventSerializer = new(ExitEventSerializer)
				}
				(*out.ExitEventSerializer).UnmarshalEasyJSON(in)
			}
		case "process":
			if in.IsNull() {
				in.Skip()
				out.ProcessContextSerializer = nil
			} else {
				if out.ProcessContextSerializer == nil {
					out.ProcessContextSerializer = new(ProcessContextSerializer)
				}
				(*out.ProcessContextSerializer).UnmarshalEasyJSON(in)
			}
		case "dd":
			if in.IsNull() {
				in.Skip()
				out.DDContextSerializer = nil
			} else {
				if out.DDContextSerializer == nil {
					out.DDContextSerializer = new(DDContextSerializer)
				}
				(*out.DDContextSerializer).UnmarshalEasyJSON(in)
			}
		case "container":
			if in.IsNull() {
				in.Skip()
				out.ContainerContextSerializer = nil
			} else {
				if out.ContainerContextSerializer == nil {
					out.ContainerContextSerializer = new(ContainerContextSerializer)
				}
				(*out.ContainerContextSerializer).UnmarshalEasyJSON(in)
			}
		case "security_profile":
			if in.IsNull() {
				in.Skip()
				out.SecurityProfileContextSerializer = nil
			} else {
				if out.SecurityProfileContextSerializer == nil {
					out.SecurityProfileContextSerializer = new(SecurityProfileContextSerializer)
				}
				easyjson6b24c4ebDecodeGithubComDataDogDatadogAgentPkgSecuritySerializers12(in, out.SecurityProfileContextSerializer)
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6b24c4ebEncodeGithubComDataDogDatadogAgentPkgSecuritySerializers11(out *jwriter.Writer, in BaseEventSerializer) {
	out.RawByte('{')
	first := true
	_ = first
	if true {
		const prefix string = ",\"evt\":"
		first = false
		out.RawString(prefix[1:])
		(in.EventContextSerializer).MarshalEasyJSON(out)
	}
	if true {
		const prefix string = ",\"date\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.Date).MarshalEasyJSON(out)
	}
	if in.FileEventSerializer != nil {
		const prefix string = ",\"file\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.FileEventSerializer).MarshalEasyJSON(out)
	}
	if in.DNSEventSerializer != nil {
		const prefix string = ",\"dns\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.DNSEventSerializer).MarshalEasyJSON(out)
	}
	if in.NetworkContextSerializer != nil {
		const prefix string = ",\"network\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.NetworkContextSerializer).MarshalEasyJSON(out)
	}
	if in.ExitEventSerializer != nil {
		const prefix string = ",\"exit\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.ExitEventSerializer).MarshalEasyJSON(out)
	}
	if in.ProcessContextSerializer != nil {
		const prefix string = ",\"process\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.ProcessContextSerializer).MarshalEasyJSON(out)
	}
	if in.DDContextSerializer != nil {
		const prefix string = ",\"dd\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.DDContextSerializer).MarshalEasyJSON(out)
	}
	if in.ContainerContextSerializer != nil {
		const prefix string = ",\"container\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.ContainerContextSerializer).MarshalEasyJSON(out)
	}
	if in.SecurityProfileContextSerializer != nil {
		const prefix string = ",\"security_profile\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		easyjson6b24c4ebEncodeGithubComDataDogDatadogAgentPkgSecuritySerializers12(out, *in.SecurityProfileContextSerializer)
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v BaseEventSerializer) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6b24c4ebEncodeGithubComDataDogDatadogAgentPkgSecuritySerializers11(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *BaseEventSerializer) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6b24c4ebDecodeGithubComDataDogDatadogAgentPkgSecuritySerializers11(l, v)
}
func easyjson6b24c4ebDecodeGithubComDataDogDatadogAgentPkgSecuritySerializers12(in *jlexer.Lexer, out *SecurityProfileContextSerializer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "name":
			out.Name = string(in.String())
		case "status":
			out.Status = string(in.String())
		case "version":
			out.Version = string(in.String())
		case "tags":
			if in.IsNull() {
				in.Skip()
				out.Tags = nil
			} else {
				in.Delim('[')
				if out.Tags == nil {
					if !in.IsDelim(']') {
						out.Tags = make([]string, 0, 4)
					} else {
						out.Tags = []string{}
					}
				} else {
					out.Tags = (out.Tags)[:0]
				}
				for !in.IsDelim(']') {
					var v10 string
					v10 = string(in.String())
					out.Tags = append(out.Tags, v10)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6b24c4ebEncodeGithubComDataDogDatadogAgentPkgSecuritySerializers12(out *jwriter.Writer, in SecurityProfileContextSerializer) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix)
		out.String(string(in.Status))
	}
	{
		const prefix string = ",\"version\":"
		out.RawString(prefix)
		out.String(string(in.Version))
	}
	{
		const prefix string = ",\"tags\":"
		out.RawString(prefix)
		if in.Tags == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v11, v12 := range in.Tags {
				if v11 > 0 {
					out.RawByte(',')
				}
				out.String(string(v12))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}