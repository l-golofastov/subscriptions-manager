package sum

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/l-golofastov/subscriptions-manager/internal/domain"
	"github.com/l-golofastov/subscriptions-manager/internal/http-server/handlers"
	"github.com/l-golofastov/subscriptions-manager/internal/http-server/lib"
	"github.com/l-golofastov/subscriptions-manager/internal/http-server/middleware"
)

// SuccessSumResponse represents success summarizing subscriptions prices response with amount in body
type SuccessSumResponse struct {
	Amount int `json:"amount" example:"1000"`
}

// @Summary Sum subscriptions prices
// @Description Calculate total price of subscriptions
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param input body domain.SumSubscriptionsFilter true "Sum filter"
// @Success 200 {object} SuccessSumResponse
// @Failure 400 {object} lib.ErrorResponse
// @Failure 500 {object} lib.ErrorResponse
// @Router /subscriptions/sum [get]
func NewSumHandler(log *slog.Logger, repo handlers.SubscriptionRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.sum.NewSumHandler"

		ctx := r.Context()

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetRequestID(ctx)),
		)

		if r.Method != http.MethodGet {
			lib.RespondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var filter domain.SumSubscriptionsFilter
		err := json.NewDecoder(r.Body).Decode(&filter)
		if err != nil {
			lib.RespondWithError(w, http.StatusBadRequest, "invalid sum subscriptions prices filter")
			return
		}

		amount, err := repo.SumSubscriptionsPrices(ctx, filter)
		if err != nil {
			log.Error("error getting sum subscriptions prices", "error", err)
			lib.RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		lib.RespondWithJSON(w, http.StatusOK, SuccessSumResponse{Amount: amount})
	}
}
