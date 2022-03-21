package config

import (
	"context"
	"regexp"
	"sync"
)

type Iterator interface {
	Next() bool
	Value() Value
	Label() string
}

type Value interface {
	Decode(interface{}) error
	DecodeWithCtx(context.Context, interface{}) error
	String() (string, error)
	Bytes() ([]byte, error)
	Bool() (bool, error)
	Float64() (float64, error)
	Int64() (int64, error)
	Uint64() (uint64, error)
	Interface() (interface{}, error)
	Ref() string
	File() string
	Lookup(path ...string) (Value, bool)
	List() (Iterator, error)
	StringList() ([]string, error)
	Struct() (Iterator, error)
	Kind() Kind
	Json() ([]byte, error)
}

type Loader interface {
	Type() string
	AllowDir() bool
	PathPattern() *regexp.Regexp
	Load(path string, files map[string][]byte) (Value, error)
	Clear()
}

type Kind uint8

const (
	UndefinedKind Kind = iota
	NullKind
	BoolKind
	StringKind
	BytesKind
	StructKind
	ListKind
	NumberKind
	DecimalKind
)

func (k Kind) String() string {
	switch k {
	default:
		fallthrough
	case UndefinedKind:
		return "Undefined"
	case NullKind:
		return "Null"
	case BoolKind:
		return "Bool"
	case StringKind:
		return "String"
	case BytesKind:
		return "Bytes"
	case StructKind:
		return "Struct"
	case ListKind:
		return "List"
	case NumberKind:
		return "Number"
	case DecimalKind:
		return "Decimal"
	}
}

type Decoder interface {
	Decode(Value) error
}

type CtxDecoder interface {
	Decode(context.Context, Value) error
}

type CtxConfigDecoder interface {
	DecodeConfig(context.Context, Value) error
}

var configInstance Value
var configRWLocker sync.RWMutex

func Lookup(fieldPath ...string) (Value, bool) {
	configRWLocker.RLock()
	defer configRWLocker.RUnlock()
	if configInstance == nil {
		return nil, false
	}
	return configInstance.Lookup(fieldPath...)
}

func Root() Value {
	configRWLocker.RLock()
	defer configRWLocker.RUnlock()
	return configInstance
}

func Clear() {
	configRWLocker.Lock()
	defer configRWLocker.Unlock()
	loader.Clear()
	loader = nil
	configInstance = nil
}
