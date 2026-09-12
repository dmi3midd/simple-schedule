package handlers

import (
	"net/http"

	"github.com/dmi3midd/simple-schedule/internal/domain"
	"github.com/dmi3midd/simple-schedule/internal/service"
	"github.com/dmi3midd/simple-schedule/internal/shared/httputils"
	"github.com/go-playground/validator/v10"
)

type CreateTagRequest struct {
	Title    string `json:"title" validate:"required,min=1,max=255" example:"Lecture"`
	HexColor string `json:"hexColor" validate:"required,hexcolor" example:"#FF5733"`
}

type UpdateTagRequest struct {
	Title    string `json:"title" validate:"required,min=1,max=255" example:"Practice"`
	HexColor string `json:"hexColor" validate:"required,hexcolor" example:"#33FF57"`
}

type TagHandler struct {
	tagService service.TagService
	val        *validator.Validate
}

func NewTagHandler(tagService service.TagService, val *validator.Validate) *TagHandler {
	return &TagHandler{
		tagService: tagService,
		val:        val,
	}
}

// GetAll retrieves all tags.
// @Summary      Get all tags
// @Description  Retrieve a list of all tags
// @Tags         tags
// @Produce      json
// @Success      200  {object}  TagsResponse
// @Failure      500  {object}  apierror.UserError
// @Router       /api/tags [get]
func (h *TagHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	tags, err := h.tagService.GetAll(r.Context())
	if err != nil {
		return err
	}
	if tags == nil {
		tags = []domain.Tag{}
	}
	return httputils.WriteJSON(w, http.StatusOK, TagsResponse{Tags: tags})
}

// GetById retrieves a tag by ID.
// @Summary      Get tag by ID
// @Description  Retrieve a single tag by its unique identifier
// @Tags         tags
// @Produce      json
// @Param        id   path      string  true  "Tag ID" example("d40c6c2b-e48f-4cb1-80a5-f8c5b6b801a2")
// @Success      200  {object}  TagResponse
// @Failure      400  {object}  apierror.UserError
// @Failure      404  {object}  apierror.UserError
// @Failure      500  {object}  apierror.UserError
// @Router       /api/tags/{id} [get]
func (h *TagHandler) GetById(w http.ResponseWriter, r *http.Request) error {
	id, err := httputils.ParseUUID(r, "id")
	if err != nil {
		return err
	}

	tag, err := h.tagService.GetById(r.Context(), id)
	if err != nil {
		return err
	}
	return httputils.WriteJSON(w, http.StatusOK, TagResponse{Tag: *tag})
}

// Create creates a new tag.
// @Summary      Create tag
// @Description  Create a new tag with title and hex color
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        request  body      CreateTagRequest  true  "Tag creation payload"
// @Success      201      {object}  IDResponse
// @Failure      400      {object}  apierror.UserError
// @Failure      409      {object}  apierror.UserError
// @Failure      500      {object}  apierror.UserError
// @Router       /api/tags [post]
func (h *TagHandler) Create(w http.ResponseWriter, r *http.Request) error {
	body, err := httputils.BindAndValidate[CreateTagRequest](r, h.val)
	if err != nil {
		return err
	}

	id, err := h.tagService.Create(r.Context(), body.Title, body.HexColor)
	if err != nil {
		return err
	}
	return httputils.WriteJSON(w, http.StatusCreated, IDResponse{ID: id})
}

// Update updates an existing tag by ID.
// @Summary      Update tag
// @Description  Update title and hex color of an existing tag
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        id       path      string            true  "Tag ID" example("d40c6c2b-e48f-4cb1-80a5-f8c5b6b801a2")
// @Param        request  body      UpdateTagRequest  true  "Tag update payload"
// @Success      200      {object}  IDResponse
// @Failure      400      {object}  apierror.UserError
// @Failure      404      {object}  apierror.UserError
// @Failure      409      {object}  apierror.UserError
// @Failure      500      {object}  apierror.UserError
// @Router       /api/tags/{id} [put]
func (h *TagHandler) Update(w http.ResponseWriter, r *http.Request) error {
	id, err := httputils.ParseUUID(r, "id")
	if err != nil {
		return err
	}

	body, err := httputils.BindAndValidate[UpdateTagRequest](r, h.val)
	if err != nil {
		return err
	}

	updatedId, err := h.tagService.Update(r.Context(), id, body.Title, body.HexColor)
	if err != nil {
		return err
	}
	return httputils.WriteJSON(w, http.StatusOK, IDResponse{ID: updatedId})
}

// Delete deletes a tag by ID.
// @Summary      Delete tag
// @Description  Delete a tag by its ID
// @Tags         tags
// @Param        id   path      string  true  "Tag ID" example("d40c6c2b-e48f-4cb1-80a5-f8c5b6b801a2")
// @Success      204  "No Content"
// @Failure      400  {object}  apierror.UserError
// @Failure      404  {object}  apierror.UserError
// @Failure      500  {object}  apierror.UserError
// @Router       /api/tags/{id} [delete]
func (h *TagHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	id, err := httputils.ParseUUID(r, "id")
	if err != nil {
		return err
	}

	if err := h.tagService.Delete(r.Context(), id); err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}
