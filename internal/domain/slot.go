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

type Slot struct {
	ID        uuid.UUID `json:"id" db:"id"`
	WeekID    uuid.UUID `json:"weekId" db:"week_id"`
	DayOfWeek DayOfWeek `json:"dayOfWeek" db:"day_of_week"`
	Activity  string    `json:"activity" db:"activity"`
	StartTime int       `json:"startTime" db:"start_time"`
	EndTime   int       `json:"endTime" db:"end_time"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
