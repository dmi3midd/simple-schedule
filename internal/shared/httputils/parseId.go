package httputils

import (
	"fmt"
	"net/http"

	"github.com/dmi3midd/simple-schedule/internal/shared/httputils/apierror"
	"github.com/google/uuid"
)

func ParseUUID(r *http.Request, param string) (uuid.UUID, error) {
	val := r.PathValue(param)
	id, err := uuid.Parse(val)
	if err != nil {
		return uuid.Nil, apierror.NewBadRequestError(err, fmt.Sprintf("invalid %s UUID parameter", param))
	}
	return id, nil
}
