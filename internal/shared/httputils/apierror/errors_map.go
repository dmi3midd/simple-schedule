package apierror

import (
	"errors"

	"github.com/dmi3midd/simple-schedule/internal/service"
)

var ErrorMap = map[error]func(err error) error{
	service.ErrWeekNotFound: func(err error) error {
		return NewNotFoundError(err, "Week not found")
	},
	service.ErrTagNotFound: func(err error) error {
		return NewNotFoundError(err, "Tag not found")
	},
	service.ErrSlotNotFound: func(err error) error {
		return NewNotFoundError(err, "Slot not found")
	},
	service.ErrTagAlreadyExists: func(err error) error {
		return NewConflictError(err, "Tag already exists")
	},
	service.ErrSlotOverlap: func(err error) error {
		return NewConflictError(err, "Slot overlaps with an existing slot in the schedule")
	},
	service.ErrInvalidTimeRange: func(err error) error {
		return NewBadRequestError(err, "Invalid time range")
	},
	service.ErrInvalidDayOfWeek: func(err error) error {
		return NewBadRequestError(err, "Invalid day of week")
	},
	service.ErrMaxWeeksReached: func(err error) error {
		return NewBadRequestError(err, "Maximum number of weeks reached")
	},
	service.ErrMaxSlotsPerDayReached: func(err error) error {
		return NewBadRequestError(err, "Maximum slots per day limit reached")
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
