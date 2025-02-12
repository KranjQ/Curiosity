package utils

import "errors"

var ErrUsernameExists = errors.New("Username is already taken")
var ErrUsernameNotExists = errors.New("No user with that username")
