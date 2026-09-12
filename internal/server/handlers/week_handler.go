package handlers

import (
	"net/http"

	"github.com/dmi3midd/simple-schedule/internal/domain"
	"github.com/dmi3midd/simple-schedule/internal/service"
	"github.com/dmi3midd/simple-schedule/internal/shared/httputils"
	"github.com/go-playground/validator/v10"
)

type CreateWeekRequest struct {
	Title string `json:"title" validate:"required,min=1,max=255" example:"Spring Semester 2026"`
}

type UpdateWeekRequest struct {
	Title string `json:"title" validate:"required,min=1,max=255" example:"Autumn Semester 2026"`
}

type WeekHandler struct {
	weekService service.WeekService
	val         *validator.Validate
}

func NewWeekHandler(weekService service.WeekService, val *validator.Validate) *WeekHandler {
	return &WeekHandler{
		weekService: weekService,
		val:         val,
	}
}

// GetAll retrieves all weeks.
// @Summary      Get all weeks
// @Description  Retrieve a list of all schedule weeks
// @Tags         weeks
// @Produce      json
// @Success      200  {object}  WeeksResponse
// @Failure      500  {object}  apierror.UserError
// @Router       /api/weeks [get]
func (h *WeekHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	weeks, err := h.weekService.GetAllWeeks(r.Context())
	if err != nil {
		return err
	}
	if weeks == nil {
		weeks = []domain.Week{}
	}
	return httputils.WriteJSON(w, http.StatusOK, WeeksResponse{Weeks: weeks})
}

// GetById retrieves a week by ID.
// @Summary      Get week by ID
// @Description  Retrieve a single week by its unique identifier
// @Tags         weeks
// @Produce      json
// @Param        id   path      string  true  "Week ID" example("d40c6c2b-e48f-4cb1-80a5-f8c5b6b801a2")
// @Success      200  {object}  WeekResponse
// @Failure      400  {object}  apierror.UserError
// @Failure      404  {object}  apierror.UserError
// @Failure      500  {object}  apierror.UserError
// @Router       /api/weeks/{id} [get]
func (h *WeekHandler) GetById(w http.ResponseWriter, r *http.Request) error {
	id, err := httputils.ParseUUID(r, "id")
	if err != nil {
		return err
	}

	week, err := h.weekService.GetWeek(r.Context(), id)
	if err != nil {
		return err
	}
	return httputils.WriteJSON(w, http.StatusOK, WeekResponse{Week: *week})
}

// GetSchedule retrieves the full schedule for a week.
// @Summary      Get week schedule
// @Description  Retrieve the complete schedule for a week including all slots grouped by days
// @Tags         weeks
// @Produce      json
// @Param        id   path      string  true  "Week ID" example("d40c6c2b-e48f-4cb1-80a5-f8c5b6b801a2")
// @Success      200  {object}  WeekScheduleResponse
// @Failure      400  {object}  apierror.UserError
// @Failure      404  {object}  apierror.UserError
// @Failure      500  {object}  apierror.UserError
// @Router       /api/weeks/{id}/schedule [get]
func (h *WeekHandler) GetSchedule(w http.ResponseWriter, r *http.Request) error {
	id, err := httputils.ParseUUID(r, "id")
	if err != nil {
		return err
	}

	schedule, err := h.weekService.GetWeekSchedule(r.Context(), id)
	if err != nil {
		return err
	}
	return httputils.WriteJSON(w, http.StatusOK, WeekScheduleResponse{Schedule: *schedule})
}

// Create creates a new week.
// @Summary      Create week
// @Description  Create a new schedule week with title
// @Tags         weeks
// @Accept       json
// @Produce      json
// @Param        request  body      CreateWeekRequest  true  "Week creation payload"
// @Success      201      {object}  IDResponse
// @Failure      400      {object}  apierror.UserError
// @Failure      500      {object}  apierror.UserError
// @Router       /api/weeks [post]
func (h *WeekHandler) Create(w http.ResponseWriter, r *http.Request) error {
	body, err := httputils.BindAndValidate[CreateWeekRequest](r, h.val)
	if err != nil {
		return err
	}

	id, err := h.weekService.CreateWeek(r.Context(), body.Title)
	if err != nil {
		return err
	}
	return httputils.WriteJSON(w, http.StatusCreated, IDResponse{ID: id})
}

// Update updates an existing week by ID.
// @Summary      Update week
// @Description  Update the title of an existing week
// @Tags         weeks
// @Accept       json
// @Produce      json
// @Param        id       path      string             true  "Week ID" example("d40c6c2b-e48f-4cb1-80a5-f8c5b6b801a2")
// @Param        request  body      UpdateWeekRequest  true  "Week update payload"
// @Success      200      {object}  IDResponse
// @Failure      400      {object}  apierror.UserError
// @Failure      404      {object}  apierror.UserError
// @Failure      500      {object}  apierror.UserError
// @Router       /api/weeks/{id} [put]
func (h *WeekHandler) Update(w http.ResponseWriter, r *http.Request) error {
	id, err := httputils.ParseUUID(r, "id")
	if err != nil {
		return err
	}

	body, err := httputils.BindAndValidate[UpdateWeekRequest](r, h.val)
	if err != nil {
		return err
	}

	updatedId, err := h.weekService.UpdateWeek(r.Context(), id, body.Title)
	if err != nil {
		return err
	}
	return httputils.WriteJSON(w, http.StatusOK, IDResponse{ID: updatedId})
}

// Delete deletes a week by ID.
// @Summary      Delete week
// @Description  Delete a week and cascade delete all its slots
// @Tags         weeks
// @Param        id   path      string  true  "Week ID" example("d40c6c2b-e48f-4cb1-80a5-f8c5b6b801a2")
// @Success      204  "No Content"
// @Failure      400  {object}  apierror.UserError
// @Failure      404  {object}  apierror.UserError
// @Failure      500  {object}  apierror.UserError
// @Router       /api/weeks/{id} [delete]
func (h *WeekHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	id, err := httputils.ParseUUID(r, "id")
	if err != nil {
		return err
	}

	if err := h.weekService.DeleteWeek(r.Context(), id); err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}
