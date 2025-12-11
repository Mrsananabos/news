package handlers

import (
	"fmt"
	"service/internal/apperrors"
)

const (
	minLimit  int64 = 1
	maxLimit  int64 = 100
	minOffset int64 = 0
)

func ValidatePaginationParams(limit, offset int64) error {
	if limit < minLimit {
		return apperrors.NewBadRequest(fmt.Sprintf("limit must be greater or equal %d", minLimit))
	}
	if limit > maxLimit {
		return apperrors.NewBadRequest(fmt.Sprintf("limit must be less or equal to %d", maxLimit))
	}

	if offset < minOffset {
		return apperrors.NewBadRequest("offset cannot be negative")
	}

	return nil
}
