package delete

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

// @Summary Delete subscription
// @Description Delete subscription by ID
// @Tags subscriptions
// @Param id path string true "Subscription ID"
// @Success 200 {object} lib.SuccessResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 404 {object} lib.ErrorResponse
// @Failure 500 {object} lib.ErrorResponse
// @Router /subscriptions/{id} [delete]
func NewDeleteHandler(log *slog.Logger, repo handlers.SubscriptionRepository, id uuid.UUID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.delete.NewDeleteHandler"

		ctx := r.Context()

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetRequestID(ctx)),
		)

		err := repo.DeleteSubscription(ctx, id)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				lib.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
			log.Error("error deleting subscription", "error", err)
			lib.RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		lib.RespondWithJSON(w, http.StatusOK, lib.NewSuccessResponse("success"))
	}
}
