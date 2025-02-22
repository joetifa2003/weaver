package vm

import "errors"

var (
	ErrInvalidNumberOfArguments = errors.New("missing argument for function")
	ErrInvalidArgType           = errors.New("invalid argument type")
)
