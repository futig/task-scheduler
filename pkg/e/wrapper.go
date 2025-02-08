package e

import "fmt"

func WrapErr(err error, message string) error {
	return fmt.Errorf("%s:%w", message, err)
}

func WrapErrIfNotNil(err error, message string) error {
	if err != nil {
		return fmt.Errorf("%s:%w", message, err)
	}
	return nil
}
