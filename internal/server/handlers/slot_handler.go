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
	ActivityID string `json:"activity_id" validate:"required"`
	Day        string `json:"day" validate:"required,min=1"`
	StartTime  string `json:"start_time" validate:"required"`
	Duration   string `json:"duration" validate:"required"`
}

type UpdateSlotRequest struct {
	ActivityID string `json:"activity_id" validate:"required"`
	Day        string `json:"day" validate:"required,min=1"`
	StartTime  string `json:"start_time" validate:"required"`
	Duration   string `json:"duration" validate:"required"`
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
}

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

	return httputils.WriteJSON(w, http.StatusCreated, map[string]string{"id": id})
}

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

	return httputils.WriteJSON(w, http.StatusOK, map[string]string{"id": id})
}

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
