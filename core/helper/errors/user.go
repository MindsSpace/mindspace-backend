package errors

import "errors"

var (
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrPasswordWrong         = errors.New("password is not correct")
	ErrUserNotFound          = errors.New("user not found")
	ErrUserNoPicture         = errors.New("user don't have any picture")
)
