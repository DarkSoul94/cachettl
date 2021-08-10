package cachettl

import "errors"

var (
	ErrObjNotFound = errors.New("no object found for this key")
	ErrObjExist    = errors.New("object with this key already exist")
	ErrObjNotValid = errors.New("object not valid")
	ErrKeyIsBlank  = errors.New("key is blank")
	ErrInvalidType = errors.New("invalid destination type")
)
