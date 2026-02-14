package httpx

import (
	"net/http"

	"github.com/Jose-Salazar-27/go-university-server/internal/shared/kernel/domain"
)

func GetStatusCode(err *domain.AppError) int {
	switch err.Code {

	case domain.CodeConflict:
		return http.StatusConflict

	case domain.CodeInvalidInput:
		return http.StatusBadRequest

	default:
		return http.StatusInternalServerError
	}
}
