// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
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

func easyjsonC8d74561DecodeGithubComAlSergeyTechMain2021DbForumInternalAppModels(in *jlexer.Lexer, out *ForumParams) {
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
		case "limit":
			out.Limit = int(in.Int())
		case "since":
			out.Since = string(in.String())
		case "desc":
			out.Desc = bool(in.Bool())
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
func easyjsonC8d74561EncodeGithubComAlSergeyTechMain2021DbForumInternalAppModels(out *jwriter.Writer, in ForumParams) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"limit\":"
		out.RawString(prefix[1:])
		out.Int(int(in.Limit))
	}
	{
		const prefix string = ",\"since\":"
		out.RawString(prefix)
		out.String(string(in.Since))
	}
	{
		const prefix string = ",\"desc\":"
		out.RawString(prefix)
		out.Bool(bool(in.Desc))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ForumParams) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC8d74561EncodeGithubComAlSergeyTechMain2021DbForumInternalAppModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ForumParams) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC8d74561EncodeGithubComAlSergeyTechMain2021DbForumInternalAppModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ForumParams) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC8d74561DecodeGithubComAlSergeyTechMain2021DbForumInternalAppModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ForumParams) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC8d74561DecodeGithubComAlSergeyTechMain2021DbForumInternalAppModels(l, v)
}
func easyjsonC8d74561DecodeGithubComAlSergeyTechMain2021DbForumInternalAppModels1(in *jlexer.Lexer, out *Forum) {
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
		case "title":
			out.Title = string(in.String())
		case "user":
			out.User = string(in.String())
		case "slug":
			out.Slug = string(in.String())
		case "posts":
			out.Posts = int(in.Int())
		case "threads":
			out.Threads = int(in.Int())
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
func easyjsonC8d74561EncodeGithubComAlSergeyTechMain2021DbForumInternalAppModels1(out *jwriter.Writer, in Forum) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix[1:])
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"user\":"
		out.RawString(prefix)
		out.String(string(in.User))
	}
	{
		const prefix string = ",\"slug\":"
		out.RawString(prefix)
		out.String(string(in.Slug))
	}
	{
		const prefix string = ",\"posts\":"
		out.RawString(prefix)
		out.Int(int(in.Posts))
	}
	{
		const prefix string = ",\"threads\":"
		out.RawString(prefix)
		out.Int(int(in.Threads))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Forum) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC8d74561EncodeGithubComAlSergeyTechMain2021DbForumInternalAppModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Forum) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC8d74561EncodeGithubComAlSergeyTechMain2021DbForumInternalAppModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Forum) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC8d74561DecodeGithubComAlSergeyTechMain2021DbForumInternalAppModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Forum) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC8d74561DecodeGithubComAlSergeyTechMain2021DbForumInternalAppModels1(l, v)
}