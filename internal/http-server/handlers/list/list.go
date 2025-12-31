package list

import (
	"log/slog"
	"net/http"

	"github.com/l-golofastov/subscriptions-manager/internal/http-server/handlers"
	"github.com/l-golofastov/subscriptions-manager/internal/http-server/lib"
	"github.com/l-golofastov/subscriptions-manager/internal/http-server/middleware"
)

// @Summary List subscriptions
// @Description Get all subscriptions
// @Tags subscriptions
// @Produce json
// @Success 200 {array} domain.Subscription
// @Failure 500 {object} lib.ErrorResponse
// @Router /subscriptions [get]
func NewListHandler(log *slog.Logger, repo handlers.SubscriptionRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.list.NewListHandler"

		ctx := r.Context()

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetRequestID(ctx)),
		)

		subs, err := repo.ListSubscriptions(ctx)
		if err != nil {
			log.Error("error getting subscriptions", "error", err)
			lib.RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		lib.RespondWithJSON(w, http.StatusOK, subs)
	}
}
