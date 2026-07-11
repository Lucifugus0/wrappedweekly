package domain

import "time"

type ActivityCategory string

const (
	CategoryWorkout ActivityCategory = "workout"
	CategoryReading ActivityCategory = "reading"
	CategoryCoding  ActivityCategory = "coding"
	CategorySpending ActivityCategory = "spending"
)

func (c ActivityCategory) IsValid() bool {
	switch c {
	case CategoryWorkout, CategoryReading, CategoryCoding, CategorySpending:
		return true
	}
	return false
}

type Activity struct {
	ID         string           `json:"id"`
	UserID     string           `json:"user_id"`
	Category   ActivityCategory `json:"category"`
	Value      float64          `json:"value"`
	Note       *string          `json:"note"`
	OccurredAt time.Time        `json:"occurred_at"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

// ActivityFilter narrows ListByUser results. Zero values mean "no filter":
// empty Category matches all categories, zero-value From/To mean unbounded.
type ActivityFilter struct {
	Category ActivityCategory
	From     time.Time
	To       time.Time
}

type ActivityRepository interface {
	Create(a *Activity) error
	FindByID(id string) (*Activity, error)
	Update(a *Activity) error
	Delete(id string) error
	ListByUser(userID string, page, size int, filter ActivityFilter) ([]*Activity, int, error)
	ListByUserInRange(userID string, from, to time.Time) ([]*Activity, error)
}
