//go:build linux
// +build linux

// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package probe

import (
	json "encoding/json"
	serializers "github.com/DataDog/datadog-agent/pkg/security/serializers"
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

func easyjsonF8f9ddd1DecodeGithubComDataDogDatadogAgentPkgSecurityProbe(in *jlexer.Lexer, out *EventLostWrite) {
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
		case "map":
			out.Name = string(in.String())
		case "per_event":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				out.Lost = make(map[string]uint64)
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v1 uint64
					v1 = uint64(in.Uint64())
					(out.Lost)[key] = v1
					in.WantComma()
				}
				in.Delim('}')
			}
		case "date":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Timestamp).UnmarshalJSON(data))
			}
		case "service":
			out.Service = string(in.String())
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
func easyjsonF8f9ddd1EncodeGithubComDataDogDatadogAgentPkgSecurityProbe(out *jwriter.Writer, in EventLostWrite) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"map\":"
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"per_event\":"
		out.RawString(prefix)
		if in.Lost == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v2First := true
			for v2Name, v2Value := range in.Lost {
				if v2First {
					v2First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v2Name))
				out.RawByte(':')
				out.Uint64(uint64(v2Value))
			}
			out.RawByte('}')
		}
	}
	{
		const prefix string = ",\"date\":"
		out.RawString(prefix)
		out.Raw((in.Timestamp).MarshalJSON())
	}
	{
		const prefix string = ",\"service\":"
		out.RawString(prefix)
		out.String(string(in.Service))
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EventLostWrite) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF8f9ddd1EncodeGithubComDataDogDatadogAgentPkgSecurityProbe(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EventLostWrite) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF8f9ddd1DecodeGithubComDataDogDatadogAgentPkgSecurityProbe(l, v)
}
func easyjsonF8f9ddd1DecodeGithubComDataDogDatadogAgentPkgSecurityProbe1(in *jlexer.Lexer, out *EventLostRead) {
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
		case "map":
			out.Name = string(in.String())
		case "lost":
			out.Lost = float64(in.Float64())
		case "date":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Timestamp).UnmarshalJSON(data))
			}
		case "service":
			out.Service = string(in.String())
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
func easyjsonF8f9ddd1EncodeGithubComDataDogDatadogAgentPkgSecurityProbe1(out *jwriter.Writer, in EventLostRead) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"map\":"
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"lost\":"
		out.RawString(prefix)
		out.Float64(float64(in.Lost))
	}
	{
		const prefix string = ",\"date\":"
		out.RawString(prefix)
		out.Raw((in.Timestamp).MarshalJSON())
	}
	{
		const prefix string = ",\"service\":"
		out.RawString(prefix)
		out.String(string(in.Service))
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EventLostRead) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF8f9ddd1EncodeGithubComDataDogDatadogAgentPkgSecurityProbe1(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EventLostRead) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF8f9ddd1DecodeGithubComDataDogDatadogAgentPkgSecurityProbe1(l, v)
}
func easyjsonF8f9ddd1DecodeGithubComDataDogDatadogAgentPkgSecurityProbe2(in *jlexer.Lexer, out *AbnormalEvent) {
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
		case "triggering_event":
			if in.IsNull() {
				in.Skip()
				out.Event = nil
			} else {
				if out.Event == nil {
					out.Event = new(serializers.EventSerializer)
				}
				(*out.Event).UnmarshalEasyJSON(in)
			}
		case "error":
			out.Error = string(in.String())
		case "date":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Timestamp).UnmarshalJSON(data))
			}
		case "service":
			out.Service = string(in.String())
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
func easyjsonF8f9ddd1EncodeGithubComDataDogDatadogAgentPkgSecurityProbe2(out *jwriter.Writer, in AbnormalEvent) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"triggering_event\":"
		out.RawString(prefix[1:])
		if in.Event == nil {
			out.RawString("null")
		} else {
			(*in.Event).MarshalEasyJSON(out)
		}
	}
	{
		const prefix string = ",\"error\":"
		out.RawString(prefix)
		out.String(string(in.Error))
	}
	{
		const prefix string = ",\"date\":"
		out.RawString(prefix)
		out.Raw((in.Timestamp).MarshalJSON())
	}
	{
		const prefix string = ",\"service\":"
		out.RawString(prefix)
		out.String(string(in.Service))
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v AbnormalEvent) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF8f9ddd1EncodeGithubComDataDogDatadogAgentPkgSecurityProbe2(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *AbnormalEvent) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF8f9ddd1DecodeGithubComDataDogDatadogAgentPkgSecurityProbe2(l, v)
}