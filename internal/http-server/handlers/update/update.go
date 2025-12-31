package update

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/l-golofastov/subscriptions-manager/internal/domain"
	"github.com/l-golofastov/subscriptions-manager/internal/http-server/handlers"
	"github.com/l-golofastov/subscriptions-manager/internal/http-server/lib"
	"github.com/l-golofastov/subscriptions-manager/internal/http-server/middleware"
	"github.com/l-golofastov/subscriptions-manager/internal/repository"
)

// @Summary Update subscription
// @Description Update subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param input body domain.UpdateSubscriptionInput true "Update subscription"
// @Success 200 {object} domain.Subscription
// @Failure 400 {object} lib.ErrorResponse
// @Failure 404 {object} lib.ErrorResponse
// @Failure 500 {object} lib.ErrorResponse
// @Router /subscriptions/{id} [patch]
func NewUpdateHandler(log *slog.Logger, repo handlers.SubscriptionRepository, id uuid.UUID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.update.NewUpdateHandler"

		ctx := r.Context()

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetRequestID(ctx)),
		)

		var in domain.UpdateSubscriptionInput
		err := json.NewDecoder(r.Body).Decode(&in)
		if err != nil {
			lib.RespondWithError(w, http.StatusBadRequest, "invalid update subscription input")
			return
		}

		sub, err := repo.UpdateSubscription(ctx, id, in)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				lib.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
			log.Error("error updating subscription", "error", err)
			lib.RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		lib.RespondWithJSON(w, http.StatusOK, sub)
	}
}
