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

type CreateActivityRequest struct {
	Title string `json:"title" validate:"required,min=1,max=255"`
}

type UpdateActivityRequest struct {
	Title string `json:"title" validate:"required,min=1,max=255"`
}

type ActivityHandler struct {
	service  service.ActivityService
	validate *validator.Validate
}

func NewActivityHandler(service service.ActivityService, validate *validator.Validate) *ActivityHandler {
	return &ActivityHandler{
		service:  service,
		validate: validate,
	}
}

func (h *ActivityHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /activities", apierror.ErrorHandler(h.Create))
	mux.HandleFunc("GET /activities", apierror.ErrorHandler(h.List))
	mux.HandleFunc("GET /activities/{id}", apierror.ErrorHandler(h.Get))
	mux.HandleFunc("PUT /activities/{id}", apierror.ErrorHandler(h.Update))
	mux.HandleFunc("DELETE /activities/{id}", apierror.ErrorHandler(h.Delete))
	mux.HandleFunc("DELETE /activities", apierror.ErrorHandler(h.DeleteAll))
}

func (h *ActivityHandler) Create(w http.ResponseWriter, r *http.Request) error {
	req, err := httputils.BindAndValidate[CreateActivityRequest](r, h.validate)
	if err != nil {
		return err
	}

	id, err := h.service.Create(r.Context(), req.Title)
	if err != nil {
		return err
	}

	return httputils.WriteJSON(w, http.StatusCreated, struct{ ID string }{ID: id})
}

func (h *ActivityHandler) Get(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	if id == "" {
		return apierror.NewBadRequestError(errors.New("missing id"), "ID parameter is required")
	}

	activity, err := h.service.Get(r.Context(), id)
	if err != nil {
		return err
	}

	return httputils.WriteJSON(w, http.StatusOK, activity)
}

func (h *ActivityHandler) List(w http.ResponseWriter, r *http.Request) error {
	activities := h.service.List(r.Context())
	if activities == nil {
		activities = []domain.Activity{}
	}

	return httputils.WriteJSON(w, http.StatusOK, activities)
}

func (h *ActivityHandler) Update(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	if id == "" {
		return apierror.NewBadRequestError(errors.New("missing id"), "ID parameter is required")
	}

	req, err := httputils.BindAndValidate[UpdateActivityRequest](r, h.validate)
	if err != nil {
		return err
	}

	id, err = h.service.Update(r.Context(), id, req.Title)
	if err != nil {
		return err
	}

	return httputils.WriteJSON(w, http.StatusOK, struct{ ID string }{ID: id})
}

func (h *ActivityHandler) Delete(w http.ResponseWriter, r *http.Request) error {
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

func (h *ActivityHandler) DeleteAll(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.DeleteAll(r.Context()); err != nil {
		return err
	}

	return httputils.WriteJSON(w, http.StatusNoContent, nil)
}
