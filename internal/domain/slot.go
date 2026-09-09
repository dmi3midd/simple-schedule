package domain

import "time"

type Slot struct {
	ID         string    `json:"id" db:"id"`
	ActivityID string    `json:"activity_id" db:"activity_id"`
	Day        string    `json:"day" db:"day"`
	StartTime  string    `json:"start_time" db:"start_time"`
	Duration   string    `json:"duration" db:"duration"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`
}
