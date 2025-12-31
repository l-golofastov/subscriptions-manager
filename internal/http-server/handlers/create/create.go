package create

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/l-golofastov/subscriptions-manager/internal/domain"
	"github.com/l-golofastov/subscriptions-manager/internal/http-server/handlers"
	"github.com/l-golofastov/subscriptions-manager/internal/http-server/lib"
	"github.com/l-golofastov/subscriptions-manager/internal/http-server/middleware"
)

// @Summary Create subscription
// @Description Create new subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param input body domain.CreateSubscriptionInput true "Create subscription"
// @Success 201 {object} domain.Subscription
// @Failure 400 {object} lib.ErrorResponse
// @Failure 500 {object} lib.ErrorResponse
// @Router /subscriptions [post]
func NewCreateHandler(log *slog.Logger, repo handlers.SubscriptionRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.create.NewCreateHandler"

		ctx := r.Context()

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetRequestID(ctx)),
		)

		var in domain.CreateSubscriptionInput
		err := json.NewDecoder(r.Body).Decode(&in)
		if err != nil {
			lib.RespondWithError(w, http.StatusBadRequest, "invalid subscription input")
			return
		}

		err = validateCreateInput(in)
		if err != nil {
			lib.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		sub, err := repo.CreateSubscription(ctx, in)
		if err != nil {
			log.Error("error creating subscription", "error", err)
			lib.RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		lib.RespondWithJSON(w, http.StatusCreated, sub)
	}
}

func validateCreateInput(in domain.CreateSubscriptionInput) error {
	if strings.TrimSpace(in.ServiceName) == "" {
		return fmt.Errorf("service name is required")
	}

	if in.Price < 0 {
		return fmt.Errorf("price must be positive")
	}

	return nil
}
