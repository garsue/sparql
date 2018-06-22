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
	// DataType is the parameter type.
	DataType Serializable
	// LanguageTag is the parameter language tag.
	LanguageTag string
}

// Serializable is a serializer to replace parameters in a query.
type Serializable interface {
	Serialize() string
}

// URI is URI string.
type URI string

// Serialize returns this URI string surrounded with `<` and `>`.
func (s URI) Serialize() string {
	return "<" + string(s) + ">"
}

// DataType
type DataType struct {
	Prefix string
	Name   string
}

// Serialize returns the string formatted as "{prefix}:{name}".
func (d DataType) Serialize() string {
	return d.Prefix + ":" + d.Name
}

func (p Param) Serialize() string {
	if p.LanguageTag != "" {
		s := strings.Replace(fmt.Sprintf("%v", p.Value), `"""`, `\"\"\"`, -1)
		return fmt.Sprintf(`"""%v"""@%s`, s, p.LanguageTag)
	}
	if p.DataType != nil {
		s := strings.Replace(fmt.Sprintf("%v", p.Value), `"""`, `\"\"\"`, -1)
		return fmt.Sprintf(`"""%v"""^^%s`, s, p.DataType.Serialize())
	}

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
		return `"""` + strings.Replace(string(v), `"""`, `\"\"\"`, -1) + `"""`
	case string:
		return `"""` + strings.Replace(v, `"""`, `\"\"\"`, -1) + `"""`
	case time.Time:
		return v.Format(dateTimeFormat)
	case Serializable:
		return v.Serialize()
	default:
		return `"""` + strings.Replace(fmt.Sprintf("%v", v), `"""`, `\"\"\"`, -1) + `"""`
	}
}
