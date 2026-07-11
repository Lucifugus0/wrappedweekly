package usecase

import (
	"time"

	"wrappedweekly/backend/internal/domain"
)

// WeekBounds returns the [start, end) bounds of the ISO week (Monday 00:00
// inclusive to next Monday 00:00 exclusive) that contains t, evaluated in loc.
//
// Using an explicit location is important: "which week" a timestamp belongs to
// depends on the timezone the user perceives their day in. All comparisons in
// this package treat the [start, end) range as half-open to avoid double-counting
// activity that lands exactly on a week boundary.
func WeekBounds(t time.Time, loc *time.Location) (start, end time.Time) {
	t = t.In(loc)
	// Go's Weekday: Sunday=0 ... Saturday=6. Convert to ISO (Monday=0 ... Sunday=6).
	isoWeekday := (int(t.Weekday()) + 6) % 7
	dayStart := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
	start = dayStart.AddDate(0, 0, -isoWeekday)
	end = start.AddDate(0, 0, 7)
	return start, end
}

// PrevWeekBounds returns the bounds of the week immediately before [start, end).
func PrevWeekBounds(start, end time.Time) (prevStart, prevEnd time.Time) {
	return start.AddDate(0, 0, -7), start
}

// AggregateWeek computes RecapStats for the activities that fall in [weekStart, weekEnd)
// and compares total value against prevWeekActivities (already scoped to the prior week).
//
// Both slices are expected to already be filtered to their respective week ranges
// by the caller (repository query) — this function does not re-filter, it only aggregates.
func AggregateWeek(weekStart, weekEnd time.Time, loc *time.Location, current, previous []*domain.Activity) domain.RecapStats {
	stats := domain.RecapStats{
		WeekStart: weekStart,
		WeekEnd:   weekEnd,
	}

	categoryTotals := map[domain.ActivityCategory]*domain.CategoryTotal{}
	dayTotals := map[string]float64{}
	dayOrder := make([]string, 0, 7)
	for i := 0; i < 7; i++ {
		d := weekStart.AddDate(0, 0, i).Format("2006-01-02")
		dayTotals[d] = 0
		dayOrder = append(dayOrder, d)
	}

	var totalValue float64
	for _, a := range current {
		totalValue += a.Value
		stats.TotalActivities++

		if _, ok := categoryTotals[a.Category]; !ok {
			categoryTotals[a.Category] = &domain.CategoryTotal{Category: a.Category}
		}
		categoryTotals[a.Category].Total += a.Value
		categoryTotals[a.Category].Count++

		dayKey := a.OccurredAt.In(loc).Format("2006-01-02")
		dayTotals[dayKey] += a.Value
	}
	stats.TotalValue = totalValue

	// Category totals: sorted by category name for deterministic output.
	for _, cat := range []domain.ActivityCategory{
		domain.CategoryWorkout, domain.CategoryReading, domain.CategoryCoding, domain.CategorySpending,
	} {
		if ct, ok := categoryTotals[cat]; ok {
			stats.TotalsByCategory = append(stats.TotalsByCategory, *ct)
		}
	}

	// Top category: highest total. Empty week -> nil (no crash, no fabricated winner).
	// Tie-break: first category in the fixed order above wins, so the result is
	// deterministic rather than depending on map iteration order.
	var topCategory *domain.ActivityCategory
	var topTotal float64
	for _, ct := range stats.TotalsByCategory {
		if topCategory == nil || ct.Total > topTotal {
			c := ct.Category
			topCategory = &c
			topTotal = ct.Total
		}
	}
	stats.TopCategory = topCategory

	// Daily breakdown in week order (Monday..Sunday), always 7 entries even if zero.
	for _, d := range dayOrder {
		stats.DailyBreakdown = append(stats.DailyBreakdown, domain.DayTotal{Date: d, Total: dayTotals[d]})
	}

	// Most productive day: highest total. Empty week -> nil.
	// Tie-break: earliest day in the week wins (dayOrder is already Monday..Sunday).
	var mostProductiveDay *string
	var maxDayTotal float64
	first := true
	for _, dt := range stats.DailyBreakdown {
		if dt.Total <= 0 {
			continue
		}
		if first || dt.Total > maxDayTotal {
			d := dt.Date
			mostProductiveDay = &d
			maxDayTotal = dt.Total
			first = false
		}
	}
	stats.MostProductiveDay = mostProductiveDay

	var prevTotal float64
	for _, a := range previous {
		prevTotal += a.Value
	}
	stats.PrevWeekTotalValue = prevTotal

	// Change vs previous week, as a percentage. Explicit rule for the zero-division
	// edge case (previous week had no activity at all):
	//   - current == 0 too  -> nil (no meaningful "change" from nothing to nothing)
	//   - current > 0       -> nil (percentage change is undefined/infinite when
	//                          going from zero; we surface the raw totals instead
	//                          of fabricating a number like "+Inf%" or "+100%")
	// Frontend/narrative should treat a nil pointer as "minggu pertama tercatat" /
	// "tidak ada data pembanding", not as 0%.
	if prevTotal > 0 {
		change := ((totalValue - prevTotal) / prevTotal) * 100
		stats.ChangeVsPrevWeekPct = &change
	}

	return stats
}
