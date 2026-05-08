package exceptions

import "errors"

var (
	ErrNotFound           = errors.New("not found error")
	ErrNotExists          = errors.New("not exists error")
	ErrAlreadyDeleted     = errors.New("already deleted error")
	ErrAlreadyExists      = errors.New("already exists error")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNoTxObj            = errors.New("no transaction object in context")
	ErrTxType             = errors.New("unexpected tx type")
	ErrNoUserID           = errors.New("no user id in request context")
)
