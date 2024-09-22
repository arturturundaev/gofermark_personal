package errors

import "errors"

var ErrExists = errors.New("exists")
var ErrConflict = errors.New("conflicct")
var ErrNotFound = errors.New("not found")
var ErrInternal = errors.New("internal error")
var ErrTimeout = errors.New("timeout")
var ErrInvalid = errors.New("invalid")
var ErrNotEnoughtMoney = errors.New("not enought money")
