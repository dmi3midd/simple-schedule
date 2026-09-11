package service

import "errors"

var (
	ErrSlotOverlap      = errors.New("slot overlaps with an existing slot in the schedule")
	ErrInvalidTimeRange = errors.New("invalid time range: start time must be >= 0, end time <= 86400, and start time < end time")
	ErrInvalidDayOfWeek = errors.New("invalid day of week")
	ErrWeekNotFound     = errors.New("week not found")
	ErrSlotNotFound     = errors.New("slot not found")
	ErrTagNotFound      = errors.New("tag not found")
	ErrTagAlreadyExists = errors.New("tag already exists")
)
