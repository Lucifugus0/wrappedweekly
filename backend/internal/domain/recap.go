package domain

import "time"

type CategoryTotal struct {
	Category ActivityCategory `json:"category"`
	Total    float64          `json:"total"`
	Count    int              `json:"count"`
}

type DayTotal struct {
	Date  string  `json:"date"`
	Total float64 `json:"total"`
}

type RecapStats struct {
	WeekStart          time.Time       `json:"week_start"`
	WeekEnd            time.Time       `json:"week_end"`
	TotalsByCategory    []CategoryTotal `json:"totals_by_category"`
	TopCategory         *ActivityCategory `json:"top_category"`
	MostProductiveDay   *string         `json:"most_productive_day"`
	DailyBreakdown      []DayTotal      `json:"daily_breakdown"`
	TotalActivities     int             `json:"total_activities"`
	TotalValue          float64         `json:"total_value"`
	PrevWeekTotalValue  float64         `json:"prev_week_total_value"`
	ChangeVsPrevWeekPct *float64        `json:"change_vs_prev_week_pct"`
}

type Recap struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	Slug      string     `json:"slug"`
	WeekStart time.Time  `json:"week_start"`
	WeekEnd   time.Time  `json:"week_end"`
	Stats     RecapStats `json:"stats"`
	Narrative string     `json:"narrative"`
	CreatedAt time.Time  `json:"created_at"`
}

type RecapRepository interface {
	Create(r *Recap) error
	FindByID(id string) (*Recap, error)
	FindBySlug(slug string) (*Recap, error)
	FindByUserAndWeek(userID string, weekStart time.Time) (*Recap, error)
	ListByUser(userID string) ([]*Recap, error)
}

type AIProvider interface {
	GenerateNarrative(stats RecapStats, userName string) (string, error)
}
