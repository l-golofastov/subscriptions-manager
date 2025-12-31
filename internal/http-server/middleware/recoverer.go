package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/l-golofastov/subscriptions-manager/internal/http-server/lib"
)

func NewRecovererMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {

				// логируем панику
				log.Printf(
					"panic recovered: %v\n%s",
					err,
					debug.Stack(),
				)

				lib.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}
