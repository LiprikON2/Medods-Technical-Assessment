package chi

import (
	"fmt"
	"net/http"
	"strings"

	auth "github.com/medods-technical-assessment"
	internalchi "github.com/medods-technical-assessment/internal/chi"
)

func Authorization(jwtService auth.JWTService) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			accessToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			if accessToken == "" {
				internalchi.ForbiddenErrorHandler(w, fmt.Errorf("error verifying Authorization header: Authorization header must be provided"))
				return

			}

			err := jwtService.VerifyAccessToken(accessToken)
			if err != nil {
				internalchi.ForbiddenErrorHandler(w, fmt.Errorf("error verifying Authorization header: %w", err))
				return
			}
			next.ServeHTTP(w, r)
		})
	}

}
