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
	Title string `json:"title" validate:"required,min=1,max=255" example:"Gym Workout"`
}

type UpdateActivityRequest struct {
	Title string `json:"title" validate:"required,min=1,max=255" example:"Evening Run"`
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

// Create creates a new activity.
// @Summary      Create an activity
// @Description  Create a new activity with a given title
// @Tags         activities
// @Accept       json
// @Produce      json
// @Param        request  body      CreateActivityRequest  true  "Activity creation payload"
// @Success      201      {object}  IDResponse
// @Failure      400      {object}  apierror.UserError
// @Failure      500      {object}  apierror.UserError
// @Router       /activities [post]
func (h *ActivityHandler) Create(w http.ResponseWriter, r *http.Request) error {
	req, err := httputils.BindAndValidate[CreateActivityRequest](r, h.validate)
	if err != nil {
		return err
	}

	id, err := h.service.Create(r.Context(), req.Title)
	if err != nil {
		return err
	}

	return httputils.WriteJSON(w, http.StatusCreated, IDResponse{ID: id})
}

// Get retrieves an activity by ID.
// @Summary      Get activity by ID
// @Description  Retrieve single activity by its unique identifier
// @Tags         activities
// @Produce      json
// @Param        id   path      string  true  "Activity ID" example("d40c6c2b-e48f-4cb1-80a5-f8c5b6b801a2")
// @Success      200  {object}  domain.Activity
// @Failure      400  {object}  apierror.UserError
// @Failure      404  {object}  apierror.UserError
// @Failure      500  {object}  apierror.UserError
// @Router       /activities/{id} [get]
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

// List retrieves all activities.
// @Summary      List activities
// @Description  Retrieve all created activities
// @Tags         activities
// @Produce      json
// @Success      200  {array}   domain.Activity
// @Failure      500  {object}  apierror.UserError
// @Router       /activities [get]
func (h *ActivityHandler) List(w http.ResponseWriter, r *http.Request) error {
	activities := h.service.List(r.Context())
	if activities == nil {
		activities = []domain.Activity{}
	}

	return httputils.WriteJSON(w, http.StatusOK, activities)
}

// Update updates an existing activity by ID.
// @Summary      Update activity
// @Description  Update title of an activity by its ID
// @Tags         activities
// @Accept       json
// @Produce      json
// @Param        id       path      string                 true  "Activity ID" example("d40c6c2b-e48f-4cb1-80a5-f8c5b6b801a2")
// @Param        request  body      UpdateActivityRequest  true  "Activity update payload"
// @Success      200      {object}  IDResponse
// @Failure      400      {object}  apierror.UserError
// @Failure      404      {object}  apierror.UserError
// @Failure      500      {object}  apierror.UserError
// @Router       /activities/{id} [put]
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

	return httputils.WriteJSON(w, http.StatusOK, IDResponse{ID: id})
}

// Delete deletes an activity by ID.
// @Summary      Delete activity
// @Description  Delete an activity by its ID and cascade delete all associated slots
// @Tags         activities
// @Param        id   path      string  true  "Activity ID" example("d40c6c2b-e48f-4cb1-80a5-f8c5b6b801a2")
// @Success      204  "No Content"
// @Failure      400  {object}  apierror.UserError
// @Failure      404  {object}  apierror.UserError
// @Failure      500  {object}  apierror.UserError
// @Router       /activities/{id} [delete]
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

// DeleteAll deletes all activities.
// @Summary      Delete all activities
// @Description  Delete all activities and cascade delete all slots
// @Tags         activities
// @Success      204  "No Content"
// @Failure      500  {object}  apierror.UserError
// @Router       /activities [delete]
func (h *ActivityHandler) DeleteAll(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.DeleteAll(r.Context()); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
