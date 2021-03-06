package types

import (
	"reflect"
	"strings"
	"testing"

	"github.com/golang/protobuf/jsonpb"
)

// WrapperMessage implements protobuf.Message but is not a normal generated message type.
type WrapperMessage struct {
	JSON *JSONValue `protobuf:"bytes,1,opt,name=json,json=json" json:"json,omitempty"`
	UUID *UUIDValue `protobuf:"bytes,2,opt,name=uuid,json=uuid" json:"uuid,omitempty"`
}

func (m *WrapperMessage) Reset() {
	m.JSON = nil
}

func (m *WrapperMessage) String() string {
	return "null"
}

func (m *WrapperMessage) ProtoMessage() {
}

func TestSuccessfulUnmarshalTypes(t *testing.T) {
	unmarshaler := &jsonpb.Unmarshaler{}
	for in, expected := range map[string]WrapperMessage{
		`{}`: {JSON: nil, UUID: nil},
		// Can't unmarshal 'null' to nil like a WKT, only an invalid, empty state
		// which will be remarshalled to 'null'
		`{"json":null}`:                                      {JSON: &JSONValue{}},
		`{"uuid":null}`:                                      {UUID: &UUIDValue{}},
		`{"json":    {"key": "value"}}`:                      {JSON: &JSONValue{Value: `{"key": "value"}`}},
		`{"uuid":  "6ba7b810-9dad-11d1-80b4-00c04fd430c8" }`: {UUID: &UUIDValue{Value: `6ba7b810-9dad-11d1-80b4-00c04fd430c8`}},
		`{"uuid":  "6ba7b8109dad11d180b400c04fd430c8" }`:     {UUID: &UUIDValue{Value: `6ba7b8109dad11d180b400c04fd430c8`}},
	} {
		jv := &WrapperMessage{}
		err := unmarshaler.Unmarshal(strings.NewReader(in), jv)
		if err != nil {
			t.Error(err.Error())
		}
		if !reflect.DeepEqual(*jv, expected) {
			t.Errorf("Expected unmarshaled output '%+v' did not match actual output '%+v'",
				expected, *jv)
		}
	}
}

func TestBrokenUnmarshalTypes(t *testing.T) {
	unmarshaler := &jsonpb.Unmarshaler{}
	for in, expected := range map[string]string{
		// A couple cases to demo standard json unmarshaling handling
		`{"}`: "unexpected EOF",
		`{"uuid":"6ba7b810-9dad-11d1-80b4-00c04fd430c8}`:  "unexpected EOF",
		`{"json":[1,2,3,4,`:                               "unexpected EOF",
		`{"json":}`:                                       "invalid character '}' looking for beginning of value",
		`{"json":[1,2,3,4,]}`:                             "invalid character ']' looking for beginning of value",
		`{"json":{"top":{"something":1},2]}}`:             "invalid character '2' looking for beginning of object key string",
		`{"uuid":{"top":{"something":1}}}`:                "invalid uuid '{\"top\":{\"something\":1}}' does not match accepted format",
		`{"uuid":"6ba7b810-9dad-11d1-80b4-00c04fdX30c8"}`: "invalid uuid '6ba7b810-9dad-11d1-80b4-00c04fdX30c8' does not match accepted format",
		`{"uuid":6ba7b810-9dad-11d1-80b4-00c04fd430c8}`:   "invalid character 'b' after object key:value pair",
		`{"uuid":ba67b810-9dad-11d1-80b4-00c04fd430c8}`:   "invalid character 'b' looking for beginning of value",
	} {
		err := unmarshaler.Unmarshal(strings.NewReader(in), &WrapperMessage{})
		if err == nil || err.Error() != expected {
			if err == nil {
				t.Errorf("Expected error %q, but got no error", expected)
			} else {
				t.Errorf("Expected error %q, but got %q", expected, err.Error())
			}
		}
	}
}

func TestMarshalTypes(t *testing.T) {
	marshaler := &jsonpb.Marshaler{OrigName: true, EmitDefaults: true}
	for expected, in := range map[string]WrapperMessage{
		`{"json":null,"uuid":null}`:                                               {},
		`{"json":{"key": "value"},"uuid":"6ba7b810-9dad-11d1-80b4-00c04fd430c8"}`: {JSON: &JSONValue{Value: `{"key": "value"}`}, UUID: &UUIDValue{Value: `6ba7b810-9dad-11d1-80b4-00c04fd430c8`}},
	} {
		out, err := marshaler.MarshalToString(&in)
		if err != nil {
			t.Error(err.Error())
		}
		if string(out) != expected {
			t.Errorf("Expected marshaled output '%s' did not match actual output '%s'",
				expected, out)
		}
	}
}

func TestMarshalTypesOmitEmpty(t *testing.T) {
	marshaller := &jsonpb.Marshaler{OrigName: true}
	for expected, in := range map[string]WrapperMessage{
		`{}`:                                                                      {},
		`{"json":null}`:                                                           {JSON: &JSONValue{}},
		`{"uuid":null}`:                                                           {UUID: &UUIDValue{}},
		`{"json":{"key": "value"}}`:                                               {JSON: &JSONValue{Value: `{"key": "value"}`}},
		`{"uuid":"6ba7b810-9dad-11d1-80b4-00c04fd430c8"}`:                         {UUID: &UUIDValue{Value: `6ba7b810-9dad-11d1-80b4-00c04fd430c8`}},
		`{"json":{"key": "value"},"uuid":"6ba7b810-9dad-11d1-80b4-00c04fd430c8"}`: {JSON: &JSONValue{Value: `{"key": "value"}`}, UUID: &UUIDValue{Value: `6ba7b810-9dad-11d1-80b4-00c04fd430c8`}},
	} {
		out, err := marshaller.MarshalToString(&in)
		if err != nil {
			t.Error(err.Error())
		}
		if string(out) != expected {
			t.Errorf("Expected marshaled output '%s' did not match actual output '%s'",
				expected, out)
		}
	}
}
