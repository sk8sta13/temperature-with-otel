package entity

import "errors"

var (
	ErrInvalidZipCode    = errors.New("Invalid zipcode")
	ErrCanNotFindZipcode = errors.New("Can not find zipcode")
	ErrInternalServer    = errors.New("Internal server error")
	ErrZipCodeRequired   = errors.New("Zipcode required")
	ErrNotFound          = errors.New("Not Found")
)
