package sparql

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const dateTimeFormat = `"2006-01-02T15:04:05Z07:00"^^xsd:dateTime`

// Param is a parameter to fill placeholders
type Param struct {
	// Ordinal position of the parameter starting from one and is always set.
	Ordinal int
	// Value is the parameter value.
	Value interface{}
}

// Serializable serialize data to embed to queries.
type Serializable interface {
	Serialize() string
}

// Literal http://www.w3.org/TR/2004/REC-rdf-concepts-20040210/#dfn-literal
type Literal struct {
	Value       interface{}
	DataType    IRIRef
	LanguageTag string
}

// Serialize returns the serialized literal string.
func (l Literal) Serialize() string {
	s := fmt.Sprint(l.Value)
	s = strings.Replace(s, `"""`, `\"\"\"`, -1)
	if l.LanguageTag != "" {
		return strings.Join([]string{`"""`, s, `"""@`, l.LanguageTag}, "")
	}
	if l.DataType != nil {
		return strings.Join([]string{`"""`, s, `"""^^`, l.DataType.Ref()}, "")
	}
	return strings.Join([]string{`"""`, s, `"""`}, "")
}

// Serialize returns the serialized as query parameter.
func (p Param) Serialize() string {
	switch v := p.Value.(type) {
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(uint64(v), 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'e', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'e', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case []byte:
		s := strings.Replace(string(v), `"""`, `\"\"\"`, -1)
		return strings.Join([]string{`"""`, s, `"""`}, "")
	case string:
		s := strings.Replace(v, `"""`, `\"\"\"`, -1)
		return strings.Join([]string{`"""`, s, `"""`}, "")
	case time.Time:
		return v.Format(dateTimeFormat)
	case IRIRef:
		return v.Ref()
	case Serializable:
		return v.Serialize()
	default:
		s := fmt.Sprint(v)
		s = strings.Replace(s, `"""`, `\"\"\"`, -1)
		return strings.Join([]string{`"""`, s, `"""`}, "")
	}
}
