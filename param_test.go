package sparql

import (
	"testing"
	"time"
)

func TestLiteral_Serialize(t *testing.T) {
	type fields struct {
		Value       string
		DataType    IRIRef
		LanguageTag string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "1",
			fields: fields{
				Value: "1",
			},
			want: `"""1"""`,
		},
		{
			name: "with datatype",
			fields: fields{
				Value:    "1",
				DataType: IRI("foo"),
			},
			want: `"""1"""^^<foo>`,
		},
		{
			name: "with language tag",
			fields: fields{
				Value:       "1",
				LanguageTag: "foo",
			},
			want: `"""1"""@foo`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Literal{
				Value:       tt.fields.Value,
				DataType:    tt.fields.DataType,
				LanguageTag: tt.fields.LanguageTag,
			}
			if got := l.Serialize(); got != tt.want {
				t.Errorf("Literal.Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParam_Serialize(t *testing.T) {
	type fields struct {
		Ordinal int
		Value   interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "int",
			fields: fields{
				Value: int(1),
			},
			want: "1",
		},
		{
			name: "int8",
			fields: fields{
				Value: int8(1),
			},
			want: "1",
		},
		{
			name: "int16",
			fields: fields{
				Value: int16(1),
			},
			want: "1",
		},
		{
			name: "int32",
			fields: fields{
				Value: int32(1),
			},
			want: "1",
		},
		{
			name: "int64",
			fields: fields{
				Value: int64(1),
			},
			want: "1",
		},
		{
			name: "uint",
			fields: fields{
				Value: uint(1),
			},
			want: "1",
		},
		{
			name: "uint8",
			fields: fields{
				Value: uint8(1),
			},
			want: "1",
		},
		{
			name: "uint16",
			fields: fields{
				Value: uint16(1),
			},
			want: "1",
		},
		{
			name: "uint32",
			fields: fields{
				Value: uint32(1),
			},
			want: "1",
		},
		{
			name: "uint64",
			fields: fields{
				Value: uint64(1),
			},
			want: "1",
		},
		{
			name: "float32",
			fields: fields{
				Value: float32(1),
			},
			want: "1e+00",
		},
		{
			name: "float64",
			fields: fields{
				Value: float64(1),
			},
			want: "1e+00",
		},
		{
			name: "bool",
			fields: fields{
				Value: true,
			},
			want: "true",
		},
		{
			name: "bytes",
			fields: fields{
				Value: []byte("hello"),
			},
			want: `"""hello"""`,
		},
		{
			name: "string",
			fields: fields{
				Value: "hello",
			},
			want: `"""hello"""`,
		},
		{
			name: "time",
			fields: fields{
				Value: time.Date(2018, time.September, 21, 12, 8, 10, 20, time.UTC),
			},
			want: `"2018-09-21T12:08:10Z"^^xsd:dateTime`,
		},
		{
			name: "IRI",
			fields: fields{
				Value: IRI("foo"),
			},
			want: `<foo>`,
		},
		{
			name: "Serializable",
			fields: fields{
				Value: Literal{},
			},
			want: `""""""`,
		},
		{
			name: "default",
			fields: fields{
				Value: complex(1, 1),
			},
			want: `"""(1+1i)"""`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Param{
				Ordinal: tt.fields.Ordinal,
				Value:   tt.fields.Value,
			}
			if got := p.Serialize(); got != tt.want {
				t.Errorf("Param.Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}
