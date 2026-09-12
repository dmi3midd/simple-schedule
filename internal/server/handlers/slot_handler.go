package handlers

import (
	"net/http"

	"github.com/dmi3midd/simple-schedule/internal/domain"
	"github.com/dmi3midd/simple-schedule/internal/service"
	"github.com/dmi3midd/simple-schedule/internal/shared/httputils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CreateSlotRequest struct {
	WeekID    uuid.UUID        `json:"weekId" validate:"required" example:"d40c6c2b-e48f-4cb1-80a5-f8c5b6b801a2"`
	TagID     *uuid.UUID       `json:"tagId,omitempty" example:"97ec8552-9562-4d79-ab4c-1ca959663629"`
	DayOfWeek domain.DayOfWeek `json:"dayOfWeek" validate:"required,oneof=monday tuesday wednesday thursday friday saturday sunday" example:"monday"`
	Activity  string           `json:"activity" validate:"required,min=1,max=255" example:"Higher Mathematics"`
	StartTime int              `json:"startTime" validate:"min=0,max=86400" example:"36000"`
	EndTime   int              `json:"endTime" validate:"min=0,max=86400,gtfield=StartTime" example:"39600"`
}

type UpdateSlotRequest struct {
	WeekID    uuid.UUID        `json:"weekId" validate:"required" example:"d40c6c2b-e48f-4cb1-80a5-f8c5b6b801a2"`
	TagID     *uuid.UUID       `json:"tagId,omitempty" example:"97ec8552-9562-4d79-ab4c-1ca959663629"`
	DayOfWeek domain.DayOfWeek `json:"dayOfWeek" validate:"required,oneof=monday tuesday wednesday thursday friday saturday sunday" example:"monday"`
	Activity  string           `json:"activity" validate:"required,min=1,max=255" example:"Physics Lecture"`
	StartTime int              `json:"startTime" validate:"min=0,max=86400" example:"36000"`
	EndTime   int              `json:"endTime" validate:"min=0,max=86400,gtfield=StartTime" example:"39600"`
}

type SlotHandler struct {
	slotService service.SlotService
	val         *validator.Validate
}

func NewSlotHandler(slotService service.SlotService, val *validator.Validate) *SlotHandler {
	return &SlotHandler{
		slotService: slotService,
		val:         val,
	}
}

// GetById retrieves a slot by ID with its tag.
// @Summary      Get slot by ID
// @Description  Retrieve a single schedule slot with its associated tag by its unique identifier
// @Tags         slots
// @Produce      json
// @Param        id   path      string  true  "Slot ID" example("d40c6c2b-e48f-4cb1-80a5-f8c5b6b801a2")
// @Success      200  {object}  SlotResponse
// @Failure      400  {object}  apierror.UserError
// @Failure      404  {object}  apierror.UserError
// @Failure      500  {object}  apierror.UserError
// @Router       /api/slots/{id} [get]
func (h *SlotHandler) GetById(w http.ResponseWriter, r *http.Request) error {
	id, err := httputils.ParseUUID(r, "id")
	if err != nil {
		return err
	}

	slot, err := h.slotService.GetSlot(r.Context(), id)
	if err != nil {
		return err
	}
	return httputils.WriteJSON(w, http.StatusOK, SlotResponse{Slot: *slot})
}

// GetByWeekId retrieves all slots for a given week.
// @Summary      Get slots by week ID
// @Description  Retrieve all schedule slots with tags for a given week
// @Tags         slots
// @Produce      json
// @Param        id   path      string  true  "Week ID" example("d40c6c2b-e48f-4cb1-80a5-f8c5b6b801a2")
// @Success      200  {object}  SlotsResponse
// @Failure      400  {object}  apierror.UserError
// @Failure      404  {object}  apierror.UserError
// @Failure      500  {object}  apierror.UserError
// @Router       /api/weeks/{id}/slots [get]
func (h *SlotHandler) GetByWeekId(w http.ResponseWriter, r *http.Request) error {
	weekId, err := httputils.ParseUUID(r, "id")
	if err != nil {
		return err
	}

	slots, err := h.slotService.GetSlotsByWeek(r.Context(), weekId)
	if err != nil {
		return err
	}
	if slots == nil {
		slots = []domain.SlotWithTag{}
	}
	return httputils.WriteJSON(w, http.StatusOK, SlotsResponse{Slots: slots})
}

// Create creates a new slot.
// @Summary      Create slot
// @Description  Create a new schedule slot with time range and day validation
// @Tags         slots
// @Accept       json
// @Produce      json
// @Param        request  body      CreateSlotRequest  true  "Slot creation payload"
// @Success      201      {object}  IDResponse
// @Failure      400      {object}  apierror.UserError
// @Failure      404      {object}  apierror.UserError
// @Failure      409      {object}  apierror.UserError
// @Failure      500      {object}  apierror.UserError
// @Router       /api/slots [post]
func (h *SlotHandler) Create(w http.ResponseWriter, r *http.Request) error {
	body, err := httputils.BindAndValidate[CreateSlotRequest](r, h.val)
	if err != nil {
		return err
	}

	input := service.CreateSlotInput{
		WeekID:    body.WeekID,
		TagID:     body.TagID,
		DayOfWeek: body.DayOfWeek,
		Activity:  body.Activity,
		StartTime: body.StartTime,
		EndTime:   body.EndTime,
	}

	id, err := h.slotService.CreateSlot(r.Context(), input)
	if err != nil {
		return err
	}
	return httputils.WriteJSON(w, http.StatusCreated, IDResponse{ID: id})
}

// Update updates an existing slot by ID.
// @Summary      Update slot
// @Description  Update details of an existing slot with collision validation
// @Tags         slots
// @Accept       json
// @Produce      json
// @Param        id       path      string             true  "Slot ID" example("d40c6c2b-e48f-4cb1-80a5-f8c5b6b801a2")
// @Param        request  body      UpdateSlotRequest  true  "Slot update payload"
// @Success      200      {object}  IDResponse
// @Failure      400      {object}  apierror.UserError
// @Failure      404      {object}  apierror.UserError
// @Failure      409      {object}  apierror.UserError
// @Failure      500      {object}  apierror.UserError
// @Router       /api/slots/{id} [put]
func (h *SlotHandler) Update(w http.ResponseWriter, r *http.Request) error {
	id, err := httputils.ParseUUID(r, "id")
	if err != nil {
		return err
	}

	body, err := httputils.BindAndValidate[UpdateSlotRequest](r, h.val)
	if err != nil {
		return err
	}

	input := service.UpdateSlotInput{
		ID:        id,
		WeekID:    body.WeekID,
		TagID:     body.TagID,
		DayOfWeek: body.DayOfWeek,
		Activity:  body.Activity,
		StartTime: body.StartTime,
		EndTime:   body.EndTime,
	}

	if err := h.slotService.UpdateSlot(r.Context(), input); err != nil {
		return err
	}
	return httputils.WriteJSON(w, http.StatusOK, IDResponse{ID: id})
}

// Delete deletes a slot by ID.
// @Summary      Delete slot
// @Description  Delete a schedule slot by its unique identifier
// @Tags         slots
// @Param        id   path      string  true  "Slot ID" example("d40c6c2b-e48f-4cb1-80a5-f8c5b6b801a2")
// @Success      204  "No Content"
// @Failure      400  {object}  apierror.UserError
// @Failure      404  {object}  apierror.UserError
// @Failure      500  {object}  apierror.UserError
// @Router       /api/slots/{id} [delete]
func (h *SlotHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	id, err := httputils.ParseUUID(r, "id")
	if err != nil {
		return err
	}

	if err := h.slotService.DeleteSlot(r.Context(), id); err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}
