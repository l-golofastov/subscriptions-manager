package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/l-golofastov/subscriptions-manager/internal/config"
	"github.com/l-golofastov/subscriptions-manager/internal/domain"
	"github.com/l-golofastov/subscriptions-manager/internal/repository"
	_ "github.com/lib/pq"
)

type StoragePostgres struct {
	db *sqlx.DB
}

func NewStoragePostgres(cfg *config.Config) (*StoragePostgres, error) {
	const op = "repository.postgres.NewStoragePostgres"

	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DB,
	)

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)

	return &StoragePostgres{db: db}, nil
}

func (s *StoragePostgres) Close() error {
	return s.db.Close()
}

func (s *StoragePostgres) ListSubscriptions(ctx context.Context) ([]domain.Subscription, error) {
	const op = "repository.postgres.ListSubscriptions"

	subscriptions := make([]domain.Subscription, 0)

	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions
		ORDER BY created_at DESC;
	`

	err := s.db.SelectContext(ctx, &subscriptions, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return subscriptions, nil
}

func (s *StoragePostgres) GetSubscriptionByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	const op = "repository.postgres.GetSubscriptionByID"

	var subscription domain.Subscription

	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions
		WHERE id = $1;
	`

	err := s.db.GetContext(ctx, &subscription, query, id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &subscription, nil
}

func (s *StoragePostgres) CreateSubscription(ctx context.Context, in domain.CreateSubscriptionInput) (*domain.Subscription, error) {
	const op = "repository.postgres.CreateSubscription"

	var subscription domain.Subscription

	query := `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at;
	`

	startDate := in.StartDate.MonthYearPtrToTimePtr()
	endDate := in.EndDate.MonthYearPtrToTimePtr()

	err := s.db.QueryRowxContext(ctx, query, in.ServiceName, in.Price, in.UserID, startDate, endDate).StructScan(&subscription)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &subscription, nil
}

func (s *StoragePostgres) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	const op = "repository.postgres.DeleteSubscription"

	query := `
		DELETE FROM subscriptions
		WHERE id = $1;
	`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (s *StoragePostgres) UpdateSubscription(ctx context.Context, id uuid.UUID, in domain.UpdateSubscriptionInput) (*domain.Subscription, error) {
	const op = "repository.postgres.UpdateSubscription"

	sub, err := s.GetSubscriptionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if in.ServiceName != nil {
		sub.ServiceName = *in.ServiceName
	}

	if in.Price != nil {
		sub.Price = *in.Price
	}

	if in.StartDate != nil {
		sub.StartDate = *in.StartDate
	}

	if in.EndDate != nil {
		sub.EndDate = *in.EndDate
	}

	sub.UpdatedAt = time.Now()

	var updatedSubscription domain.Subscription

	query := `
		UPDATE subscriptions
		SET service_name = $1, price = $2, start_date = $3, end_date = $4, updated_at = $5
		WHERE id = $6
		RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at;
	`

	startDate := sub.StartDate.MonthYearPtrToTimePtr()
	endDate := sub.EndDate.MonthYearPtrToTimePtr()

	err = s.db.QueryRowxContext(ctx, query, sub.ServiceName, sub.Price, startDate, endDate, sub.UpdatedAt, id).StructScan(&updatedSubscription)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &updatedSubscription, nil
}

func (s *StoragePostgres) SumSubscriptionsPrices(ctx context.Context, in domain.SumSubscriptionsFilter) (int, error) {
	const op = "repository.postgres.SumSubscriptionsPrices"

	var total sql.NullInt64

	query := `
		SELECT COALESCE(SUM(price), 0)
		FROM subscriptions
		WHERE user_id = $1
		  AND service_name = $2
		  AND start_date <= $4
		  AND (end_date IS NULL OR end_date >= $3);
	`

	from := in.From.MonthYearPtrToTimePtr()
	to := in.To.MonthYearPtrToTimePtr()

	err := s.db.GetContext(ctx, &total, query, in.UserID, in.ServiceName, from, to)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return int(total.Int64), nil
}
