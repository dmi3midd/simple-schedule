package domain

import (
	"time"

	"github.com/google/uuid"
)

type DayOfWeek string

const (
	Monday    DayOfWeek = "monday"
	Tuesday   DayOfWeek = "tuesday"
	Wednesday DayOfWeek = "wednesday"
	Thursday  DayOfWeek = "thursday"
	Friday    DayOfWeek = "friday"
	Saturday  DayOfWeek = "saturday"
	Sunday    DayOfWeek = "sunday"
)

func (d DayOfWeek) IsValid() bool {
	switch d {
	case Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday:
		return true
	default:
		return false
	}
}

func (d DayOfWeek) Order() int {
	switch d {
	case Monday:
		return 1
	case Tuesday:
		return 2
	case Wednesday:
		return 3
	case Thursday:
		return 4
	case Friday:
		return 5
	case Saturday:
		return 6
	case Sunday:
		return 7
	default:
		return 0
	}
}

func AllDaysOfWeek() []DayOfWeek {
	return []DayOfWeek{
		Monday,
		Tuesday,
		Wednesday,
		Thursday,
		Friday,
		Saturday,
		Sunday,
	}
}


type Slot struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	WeekID    uuid.UUID  `json:"weekId" db:"week_id"`
	TagID     *uuid.UUID `json:"tagId" db:"tag_id"`
	DayOfWeek DayOfWeek  `json:"dayOfWeek" db:"day_of_week"`
	Activity  string     `json:"activity" db:"activity"`
	StartTime int        `json:"startTime" db:"start_time"`
	EndTime   int        `json:"endTime" db:"end_time"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt" db:"updated_at"`
}

type SlotWithTag struct {
	Slot
	Tag *Tag `json:"tag,omitempty"`
}

