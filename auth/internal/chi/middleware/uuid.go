package chi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	chierrors "github.com/medods-technical-assessment/internal/chi"
)

func ValidateUUIDParam(paramName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			paramValue := chi.URLParam(r, paramName)
			parsedUUID, err := uuid.Parse(paramValue)
			if err != nil {
				chierrors.BadRequestErrorHandler(w, fmt.Errorf("invalid UUID format: %w", err))
				return
			}
			ctx := context.WithValue(r.Context(), paramName, parsedUUID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
