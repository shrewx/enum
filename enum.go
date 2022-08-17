package enum

import "errors"

var InvalidTypeError = errors.New("invalid type error")

type Enum interface {
	Int() int
	String() string
	Label() string
	Type() string
	Values() []Enum
}

type EnumOffset interface {
	Offset() int
}

type EnumValue struct {
	Key         string
	StringValue *string
	IntValue    *int64
	FloatValue  *float64
	Label       string
}

func (e EnumValue) Type() interface{} {
	if e.StringValue != nil {
		return e.StringValue
	}
	if e.IntValue != nil {
		return e.IntValue
	}
	if e.FloatValue != nil {
		return e.FloatValue
	}

	return nil
}
