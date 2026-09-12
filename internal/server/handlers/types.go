package handlers

import (
	"github.com/dmi3midd/simple-schedule/internal/domain"
	"github.com/google/uuid"
)

// IDResponse represents a response containing a resource ID.
type IDResponse struct {
	ID uuid.UUID `json:"id" example:"97ec8552-9562-4d79-ab4c-1ca959663629"`
}

// TagResponse represents a single tag response.
type TagResponse struct {
	Tag domain.Tag `json:"tag"`
}

// TagsResponse represents a list of tags response.
type TagsResponse struct {
	Tags []domain.Tag `json:"tags"`
}

// WeekResponse represents a single week response.
type WeekResponse struct {
	Week domain.Week `json:"week"`
}

// WeeksResponse represents a list of weeks response.
type WeeksResponse struct {
	Weeks []domain.Week `json:"weeks"`
}

// WeekScheduleResponse represents a full week schedule response.
type WeekScheduleResponse struct {
	Schedule domain.WeekSchedule `json:"schedule"`
}

// SlotResponse represents a single slot with its tag response.
type SlotResponse struct {
	Slot domain.SlotWithTag `json:"slot"`
}

// SlotsResponse represents a list of slots with tags response.
type SlotsResponse struct {
	Slots []domain.SlotWithTag `json:"slots"`
}
