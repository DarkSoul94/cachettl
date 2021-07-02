package cachettl

import "errors"

var (
	ErrObjNotFound = errors.New("No object found for this key")
	ErrObjExist    = errors.New("Object with this key already exist")
	ErrObjNotValid = errors.New("Object not valid")
	ErrKeyIsBlanc  = errors.New("Key is blank")
	ErrInvalidType = errors.New("Invalid type")
)
