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
	DataType string
	// LanguageTag is the parameter language tag.
	LanguageTag string
}

func (p Param) Serialize() string {
	if p.LanguageTag != "" {
		s := strings.Replace(fmt.Sprintf("%v", p.Value), `"""`, `\"\"\"`, -1)
		return fmt.Sprintf(`"""%v"""@%s`, s, p.LanguageTag)
	}
	if p.DataType != "" {
		s := strings.Replace(fmt.Sprintf("%v", p.Value), `"""`, `\"\"\"`, -1)
		return fmt.Sprintf(`"""%v"""^^%s`, s, p.DataType)
	}

	switch v := p.Value.(type) {
	// TODO deal all built-in types
	case int64:
		return strconv.FormatInt(v, 10)
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
	default:
		return `"""` + strings.Replace(fmt.Sprintf("%v", v), `"""`, `\"\"\"`, -1) + `"""`
	}
}
