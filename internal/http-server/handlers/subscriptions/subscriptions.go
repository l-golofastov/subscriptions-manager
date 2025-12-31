package subscriptions

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/l-golofastov/subscriptions-manager/internal/http-server/handlers"
	"github.com/l-golofastov/subscriptions-manager/internal/http-server/handlers/create"
	del "github.com/l-golofastov/subscriptions-manager/internal/http-server/handlers/delete"
	"github.com/l-golofastov/subscriptions-manager/internal/http-server/handlers/get"
	"github.com/l-golofastov/subscriptions-manager/internal/http-server/handlers/list"
	"github.com/l-golofastov/subscriptions-manager/internal/http-server/handlers/update"
	"github.com/l-golofastov/subscriptions-manager/internal/http-server/lib"
)

func NewSubscriptionByIDHandler(log *slog.Logger, repo handlers.SubscriptionRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prefix := "/subscriptions/"

		if !strings.HasPrefix(r.URL.Path, prefix) {
			lib.RespondWithError(w, http.StatusNotFound, "invalid URL path")
			return
		}

		idStr := strings.TrimPrefix(r.URL.Path, prefix)
		if idStr == "" {
			lib.RespondWithError(w, http.StatusBadRequest, "empty path parameters")
			return
		}

		id, err := uuid.Parse(idStr)
		if err != nil {
			lib.RespondWithError(w, http.StatusBadRequest, "invalid id")
			return
		}

		log.Info(idStr)

		switch r.Method {
		case http.MethodGet:
			h := get.NewGetHandler(log, repo, id)
			h.ServeHTTP(w, r)
		case http.MethodDelete:
			h := del.NewDeleteHandler(log, repo, id)
			h.ServeHTTP(w, r)
		case http.MethodPatch:
			h := update.NewUpdateHandler(log, repo, id)
			h.ServeHTTP(w, r)
		default:
			log.Error("method not allowed")
			lib.RespondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func NewSubscriptionsHandler(log *slog.Logger, repo handlers.SubscriptionRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h := list.NewListHandler(log, repo)
			h.ServeHTTP(w, r)
		case http.MethodPost:
			h := create.NewCreateHandler(log, repo)
			h.ServeHTTP(w, r)
		default:
			lib.RespondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}
