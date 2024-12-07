package models

import "errors"

var (
	JSONUnmarshalError = errors.New("unexpected content of the JSON message")
	KeyNotExistsError  = errors.New("key not exists")
	EmptyOrderIDError  = errors.New("order id is empty")
)
