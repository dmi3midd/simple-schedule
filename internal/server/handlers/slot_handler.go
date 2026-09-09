package handlers

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/dmi3midd/simple-schedule/internal/domain"
	"github.com/dmi3midd/simple-schedule/internal/service"
	"github.com/dmi3midd/simple-schedule/internal/shared/httputils"
	"github.com/dmi3midd/simple-schedule/internal/shared/httputils/apierror"
)

type CreateSlotRequest struct {
	ActivityID string `json:"activity_id" validate:"required" example:"d40c6c2b-e48f-4cb1-80a5-f8c5b6b801a2"`
	Day        string `json:"day" validate:"required,min=1" example:"Monday"`
	StartTime  string `json:"start_time" validate:"required" example:"09:00"`
	Duration   string `json:"duration" validate:"required" example:"1h30m"`
}

type UpdateSlotRequest struct {
	ActivityID string `json:"activity_id" validate:"required" example:"d40c6c2b-e48f-4cb1-80a5-f8c5b6b801a2"`
	Day        string `json:"day" validate:"required,min=1" example:"Tuesday"`
	StartTime  string `json:"start_time" validate:"required" example:"10:00"`
	Duration   string `json:"duration" validate:"required" example:"2h"`
}

type SlotHandler struct {
	service  service.SlotService
	validate *validator.Validate
}

func NewSlotHandler(service service.SlotService, validate *validator.Validate) *SlotHandler {
	return &SlotHandler{
		service:  service,
		validate: validate,
	}
}

func (h *SlotHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /slots", apierror.ErrorHandler(h.Create))
	mux.HandleFunc("GET /slots", apierror.ErrorHandler(h.List))
	mux.HandleFunc("GET /slots/{id}", apierror.ErrorHandler(h.Get))
	mux.HandleFunc("PUT /slots/{id}", apierror.ErrorHandler(h.Update))
	mux.HandleFunc("DELETE /slots/{id}", apierror.ErrorHandler(h.Delete))
	mux.HandleFunc("DELETE /slots", apierror.ErrorHandler(h.DeleteAll))
}

// Create creates a new slot.
// @Summary      Create a slot
// @Description  Create a new schedule slot linked to an activity
// @Tags         slots
// @Accept       json
// @Produce      json
// @Param        request  body      CreateSlotRequest  true  "Slot creation payload"
// @Success      201      {object}  IDResponse
// @Failure      400      {object}  apierror.UserError
// @Failure      500      {object}  apierror.UserError
// @Router       /slots [post]
func (h *SlotHandler) Create(w http.ResponseWriter, r *http.Request) error {
	req, err := httputils.BindAndValidate[CreateSlotRequest](r, h.validate)
	if err != nil {
		return err
	}

	id, err := h.service.Create(r.Context(), service.CreateSlotInput{
		ActivityID: req.ActivityID,
		Day:        req.Day,
		StartTime:  req.StartTime,
		Duration:   req.Duration,
	})
	if err != nil {
		return err
	}

	return httputils.WriteJSON(w, http.StatusCreated, IDResponse{ID: id})
}

// Get retrieves a slot by ID.
// @Summary      Get slot by ID
// @Description  Retrieve single schedule slot by its unique identifier
// @Tags         slots
// @Produce      json
// @Param        id   path      string  true  "Slot ID" example("f8a29b20-c23d-4299-8255-ec4319fb7914")
// @Success      200  {object}  domain.Slot
// @Failure      400  {object}  apierror.UserError
// @Failure      404  {object}  apierror.UserError
// @Failure      500  {object}  apierror.UserError
// @Router       /slots/{id} [get]
func (h *SlotHandler) Get(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	if id == "" {
		return apierror.NewBadRequestError(errors.New("missing id"), "ID parameter is required")
	}

	slot, err := h.service.Get(r.Context(), id)
	if err != nil {
		return err
	}

	return httputils.WriteJSON(w, http.StatusOK, slot)
}

// List retrieves all slots or filters by activity_id.
// @Summary      List slots
// @Description  Retrieve all schedule slots, optionally filtered by activity_id query parameter
// @Tags         slots
// @Produce      json
// @Param        activity_id  query     string  false  "Filter by Activity ID" example("d40c6c2b-e48f-4cb1-80a5-f8c5b6b801a2")
// @Success      200          {array}   domain.Slot
// @Failure      500          {object}  apierror.UserError
// @Router       /slots [get]
func (h *SlotHandler) List(w http.ResponseWriter, r *http.Request) error {
	activityID := r.URL.Query().Get("activity_id")

	var slots []domain.Slot
	if activityID != "" {
		slots = h.service.ListByActivityID(r.Context(), activityID)
	} else {
		slots = h.service.List(r.Context())
	}

	if slots == nil {
		slots = []domain.Slot{}
	}

	return httputils.WriteJSON(w, http.StatusOK, slots)
}

// Update updates an existing slot by ID.
// @Summary      Update slot
// @Description  Update details of an existing schedule slot
// @Tags         slots
// @Accept       json
// @Produce      json
// @Param        id       path      string             true  "Slot ID" example("f8a29b20-c23d-4299-8255-ec4319fb7914")
// @Param        request  body      UpdateSlotRequest  true  "Slot update payload"
// @Success      200      {object}  IDResponse
// @Failure      400      {object}  apierror.UserError
// @Failure      404      {object}  apierror.UserError
// @Failure      500      {object}  apierror.UserError
// @Router       /slots/{id} [put]
func (h *SlotHandler) Update(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	if id == "" {
		return apierror.NewBadRequestError(errors.New("missing id"), "ID parameter is required")
	}

	req, err := httputils.BindAndValidate[UpdateSlotRequest](r, h.validate)
	if err != nil {
		return err
	}

	id, err = h.service.Update(r.Context(), id, service.UpdateSlotInput{
		ActivityID: req.ActivityID,
		Day:        req.Day,
		StartTime:  req.StartTime,
		Duration:   req.Duration,
	})
	if err != nil {
		return err
	}

	return httputils.WriteJSON(w, http.StatusOK, IDResponse{ID: id})
}

// Delete deletes a slot by ID.
// @Summary      Delete slot
// @Description  Delete a schedule slot by its ID
// @Tags         slots
// @Param        id   path      string  true  "Slot ID" example("f8a29b20-c23d-4299-8255-ec4319fb7914")
// @Success      204  "No Content"
// @Failure      400  {object}  apierror.UserError
// @Failure      404  {object}  apierror.UserError
// @Failure      500  {object}  apierror.UserError
// @Router       /slots/{id} [delete]
func (h *SlotHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	if id == "" {
		return apierror.NewBadRequestError(errors.New("missing id"), "ID parameter is required")
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// DeleteAll deletes all slots.
// @Summary      Delete all slots
// @Description  Delete all schedule slots
// @Tags         slots
// @Success      204  "No Content"
// @Failure      500  {object}  apierror.UserError
// @Router       /slots [delete]
func (h *SlotHandler) DeleteAll(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.DeleteAll(r.Context()); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
