package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/l-golofastov/subscriptions-manager/internal/domain"
)

type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, in domain.CreateSubscriptionInput) (*domain.Subscription, error)
	GetSubscriptionByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	ListSubscriptions(ctx context.Context) ([]domain.Subscription, error)
	UpdateSubscription(ctx context.Context, id uuid.UUID, in domain.UpdateSubscriptionInput) (*domain.Subscription, error)
	DeleteSubscription(ctx context.Context, id uuid.UUID) error
	SumSubscriptionsPrices(ctx context.Context, in domain.SumSubscriptionsFilter) (int, error)
}
