package sqlbuilder

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type literal interface {
	serializable
	Raw() interface{}
	IsNil() bool
}

type literalImpl struct {
	raw         interface{}
	placeholder bool
}

func toLiteral(v interface{}) literal {
	return &literalImpl{
		raw:         v,
		placeholder: true,
	}
}

func (l *literalImpl) serialize(bldr *builder) {
	if l.placeholder {
		bldr.AppendValue(l.raw)
	} else {
		bldr.Append(l.string())
	}
}

func (l *literalImpl) IsNil() bool {
	if l.raw == nil {
		return true
	}

	v := reflect.ValueOf(l.raw)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}

func (l *literalImpl) string() string {
	switch t := l.raw.(type) {
	case int:
		return strconv.FormatInt(int64(t), 10)
	case int8:
		return strconv.FormatInt(int64(t), 10)
	case int16:
		return strconv.FormatInt(int64(t), 10)
	case int32:
		return strconv.FormatInt(int64(t), 10)
	case int64:
		return strconv.FormatInt(t, 10)
	case uint:
		return strconv.FormatUint(uint64(t), 10)
	case uint8:
		return strconv.FormatUint(uint64(t), 10)
	case uint16:
		return strconv.FormatUint(uint64(t), 10)
	case uint32:
		return strconv.FormatUint(uint64(t), 10)
	case uint64:
		return strconv.FormatUint(t, 10)
	case float32:
		return strconv.FormatFloat(float64(t), 'f', 10, 32)
	case float64:
		return strconv.FormatFloat(t, 'f', 10, 64)
	case bool:
		return strconv.FormatBool(t)
	case string:
		return t
	case []byte:
		return string(t)
	case time.Time:
		return t.Format("2006-01-02 15:04:05")
	case fmt.Stringer:
		return t.String()
	case nil:
		return "NULL"
	default:
		return ""
	}
}

func (l *literalImpl) Raw() interface{} {
	return l.raw
}
