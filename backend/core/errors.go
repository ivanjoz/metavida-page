package core

import (
	"errors"
	"fmt"
)

func Err(msgs ...any) error {
	return errors.New(fmt.Sprint(msgs...))
}
