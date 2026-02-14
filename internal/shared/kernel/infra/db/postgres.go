package db

import (
	"strings"

	"github.com/Jose-Salazar-27/go-university-server/internal/shared/kernel/domain"
	"github.com/lib/pq"
)

func ExchangePGError(e *pq.Error) error {
	switch {
	case IsUniqueConstraintViolation(e):
		return domain.ErrConflict
	default:
		return e
	}
}

func IsPgError(e error) (bool, *pq.Error) {
	pq_error, ok := e.(*pq.Error)
	if !ok {
		return ok, nil
	}
	return ok, pq_error
}

// Integrity Constraint Violation
func IsUniqueConstraintViolation(err *pq.Error) bool {
	return strings.Contains(err.Code.Name(), "unique")
}
