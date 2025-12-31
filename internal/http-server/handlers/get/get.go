package get

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/l-golofastov/subscriptions-manager/internal/http-server/handlers"
	"github.com/l-golofastov/subscriptions-manager/internal/http-server/lib"
	"github.com/l-golofastov/subscriptions-manager/internal/http-server/middleware"
	"github.com/l-golofastov/subscriptions-manager/internal/repository"
)

// @Summary Get subscription
// @Description Get subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} domain.Subscription
// @Failure 400 {object} lib.ErrorResponse
// @Failure 404 {object} lib.ErrorResponse
// @Failure 500 {object} lib.ErrorResponse
// @Router /subscriptions/{id} [get]
func NewGetHandler(log *slog.Logger, repo handlers.SubscriptionRepository, id uuid.UUID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.get.NewGetHandler"

		ctx := r.Context()

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetRequestID(ctx)),
		)

		sub, err := repo.GetSubscriptionByID(ctx, id)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				lib.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
			log.Error("error getting subscription", "error", err)
			lib.RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		lib.RespondWithJSON(w, http.StatusOK, sub)
	}
}
