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
	StringValue string
	IntValue    int
	Label       string

	StringType bool
}
