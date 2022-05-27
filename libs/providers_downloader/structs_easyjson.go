// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package providers_downloader

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

func easyjson6a975c40DecodeOttplayFossServerEpgLibsProvidersDownloader(in *jlexer.Lexer, out *ProviderData) {
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
		case "meta":
			easyjson6a975c40DecodeOttplayFossServerEpgLibsProvidersDownloader1(in, &out.Meta)
		case "data":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				out.Data = make(map[uint32][]string)
				for !in.IsDelim('}') {
					key := uint32(in.Uint32Str())
					in.WantColon()
					var v1 []string
					if in.IsNull() {
						in.Skip()
						v1 = nil
					} else {
						in.Delim('[')
						if v1 == nil {
							if !in.IsDelim(']') {
								v1 = make([]string, 0, 4)
							} else {
								v1 = []string{}
							}
						} else {
							v1 = (v1)[:0]
						}
						for !in.IsDelim(']') {
							var v2 string
							v2 = string(in.String())
							v1 = append(v1, v2)
							in.WantComma()
						}
						in.Delim(']')
					}
					(out.Data)[key] = v1
					in.WantComma()
				}
				in.Delim('}')
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
func easyjson6a975c40EncodeOttplayFossServerEpgLibsProvidersDownloader(out *jwriter.Writer, in ProviderData) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"meta\":"
		out.RawString(prefix[1:])
		easyjson6a975c40EncodeOttplayFossServerEpgLibsProvidersDownloader1(out, in.Meta)
	}
	{
		const prefix string = ",\"data\":"
		out.RawString(prefix)
		if in.Data == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v3First := true
			for v3Name, v3Value := range in.Data {
				if v3First {
					v3First = false
				} else {
					out.RawByte(',')
				}
				out.Uint32Str(uint32(v3Name))
				out.RawByte(':')
				if v3Value == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
					out.RawString("null")
				} else {
					out.RawByte('[')
					for v4, v5 := range v3Value {
						if v4 > 0 {
							out.RawByte(',')
						}
						out.String(string(v5))
					}
					out.RawByte(']')
				}
			}
			out.RawByte('}')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ProviderData) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6a975c40EncodeOttplayFossServerEpgLibsProvidersDownloader(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ProviderData) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6a975c40EncodeOttplayFossServerEpgLibsProvidersDownloader(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ProviderData) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6a975c40DecodeOttplayFossServerEpgLibsProvidersDownloader(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ProviderData) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6a975c40DecodeOttplayFossServerEpgLibsProvidersDownloader(l, v)
}
func easyjson6a975c40DecodeOttplayFossServerEpgLibsProvidersDownloader1(in *jlexer.Lexer, out *ProvMeta) {
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
			out.Id = string(in.String())
		case "url-hashes":
			if in.IsNull() {
				in.Skip()
				out.Urls = nil
			} else {
				in.Delim('[')
				if out.Urls == nil {
					if !in.IsDelim(']') {
						out.Urls = make([]uint32, 0, 16)
					} else {
						out.Urls = []uint32{}
					}
				} else {
					out.Urls = (out.Urls)[:0]
				}
				for !in.IsDelim(']') {
					var v6 uint32
					v6 = uint32(in.Uint32())
					out.Urls = append(out.Urls, v6)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "last-upd":
			out.LastUpd = uint64(in.Uint64())
		case "last-epg":
			out.LastEpg = uint64(in.Uint64())
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
func easyjson6a975c40EncodeOttplayFossServerEpgLibsProvidersDownloader1(out *jwriter.Writer, in ProvMeta) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.String(string(in.Id))
	}
	{
		const prefix string = ",\"url-hashes\":"
		out.RawString(prefix)
		if in.Urls == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v7, v8 := range in.Urls {
				if v7 > 0 {
					out.RawByte(',')
				}
				out.Uint32(uint32(v8))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"last-upd\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.LastUpd))
	}
	{
		const prefix string = ",\"last-epg\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.LastEpg))
	}
	out.RawByte('}')
}
