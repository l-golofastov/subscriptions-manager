package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Subscription represents a subscription entity
type Subscription struct {
	ID          uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ServiceName string    `json:"service_name" db:"service_name" example:"Netflix"`
	Price       int       `json:"price" db:"price" example:"499"`

	UserID    uuid.UUID  `json:"user_id" db:"user_id" example:"111e8400-e29b-41d4-a716-446655440000"`
	StartDate MonthYear  `json:"start_date" db:"start_date" example:"07-2025"`
	EndDate   *MonthYear `json:"end_date" db:"end_date" example:"12-2025"`

	CreatedAt time.Time `json:"created_at" db:"created_at" example:"2025-01-01T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" example:"2025-01-01T12:00:00Z"`
}

// CreateSubscriptionInput input payload
type CreateSubscriptionInput struct {
	// @Schema(required=true)
	ServiceName string `json:"service_name" example:"Netflix"`

	// @Schema(required=true)
	Price int `json:"price" example:"499"`

	// @Schema(required=true)
	UserID uuid.UUID `json:"user_id" example:"111e8400-e29b-41d4-a716-446655440000"`

	// @Schema(required=true)
	StartDate MonthYear `json:"start_date" example:"07-2025"`

	EndDate *MonthYear `json:"end_date,omitempty" example:"12-2025"`
}

// UpdateSubscriptionInput update payload
type UpdateSubscriptionInput struct {
	ServiceName *string     `json:"service_name,omitempty" example:"Spotify"`
	Price       *int        `json:"price,omitempty" example:"299"`
	StartDate   *MonthYear  `json:"start_date,omitempty" example:"08-2025"`
	EndDate     **MonthYear `json:"end_date,omitempty" example:"11-2025"`
}

// SumSubscriptionsFilter sum filter
type SumSubscriptionsFilter struct {
	// @Schema(required=true)
	ServiceName string `json:"service_name" example:"Netflix"`

	// @Schema(required=true)
	UserID uuid.UUID `json:"user_id" example:"111e8400-e29b-41d4-a716-446655440000"`

	// @Schema(required=true)
	From MonthYear `json:"from" example:"01-2025"`

	// @Schema(required=true)
	To MonthYear `json:"to" example:"12-2025"`
}

// MonthYear represents month-year date (MM-YYYY)
type MonthYear time.Time

func (my *MonthYear) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse("01-2006", s)
	if err != nil {
		return err
	}
	*my = MonthYear(t)
	return nil
}

func (my *MonthYear) MarshalJSON() ([]byte, error) {
	t := time.Time(*my)
	s := fmt.Sprintf("\"%02d-%d\"", t.Month(), t.Year())
	return []byte(s), nil
}

func (my *MonthYear) MonthYearPtrToTimePtr() *time.Time {
	if my == nil {
		return nil
	}
	t := time.Time(*my)
	return &t
}
