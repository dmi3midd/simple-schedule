package apierror

import (
	"errors"
)

var ErrorMap = map[error]func(err error) error{}

func MapError(err error) error {
	if err == nil {
		return nil
	}

	var apiErr APIError
	if errors.As(err, &apiErr) {
		return err
	}

	for serviceErr, mapFn := range ErrorMap {
		if errors.Is(err, serviceErr) {
			return mapFn(err)
		}
	}

	return NewInternalServerError(err)
}
