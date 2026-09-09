package apierror

import (
	"errors"

	"github.com/dmi3midd/simple-schedule/internal/service"
)

var ErrorMap = map[error]func(err error) error{
	service.ErrActivityNotFound: func(err error) error {
		return NewNotFoundError(err, "Activity not found")
	},
	service.ErrSlotNotFound: func(err error) error {
		return NewNotFoundError(err, "Slot not found")
	},
}

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
