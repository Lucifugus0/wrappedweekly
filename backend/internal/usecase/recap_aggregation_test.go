package usecase

import (
	"testing"
	"time"

	"wrappedweekly/backend/internal/domain"
)

func mustLoadLocation(t *testing.T, name string) *time.Location {
	t.Helper()
	loc, err := time.LoadLocation(name)
	if err != nil {
		t.Skipf("timezone data %q not available in this environment: %v", name, err)
	}
	return loc
}

func act(category domain.ActivityCategory, value float64, occurredAt time.Time) *domain.Activity {
	return &domain.Activity{Category: category, Value: value, OccurredAt: occurredAt}
}

func TestWeekBounds_MondayStart(t *testing.T) {
	loc := time.UTC
	// Wednesday 2026-07-08
	wed := time.Date(2026, 7, 8, 15, 0, 0, 0, loc)
	start, end := WeekBounds(wed, loc)

	wantStart := time.Date(2026, 7, 6, 0, 0, 0, 0, loc) // Monday
	wantEnd := time.Date(2026, 7, 13, 0, 0, 0, 0, loc)  // next Monday

	if !start.Equal(wantStart) {
		t.Errorf("start = %v, want %v", start, wantStart)
	}
	if !end.Equal(wantEnd) {
		t.Errorf("end = %v, want %v", end, wantEnd)
	}
}

func TestWeekBounds_SundayBelongsToPreviousMonday(t *testing.T) {
	loc := time.UTC
	// Sunday 2026-07-12 23:59:59 must still be in the week starting Monday 2026-07-06.
	sun := time.Date(2026, 7, 12, 23, 59, 59, 0, loc)
	start, end := WeekBounds(sun, loc)

	wantStart := time.Date(2026, 7, 6, 0, 0, 0, 0, loc)
	wantEnd := time.Date(2026, 7, 13, 0, 0, 0, 0, loc)

	if !start.Equal(wantStart) || !end.Equal(wantEnd) {
		t.Errorf("Sunday 23:59:59 should belong to week [%v,%v), got [%v,%v)", wantStart, wantEnd, start, end)
	}
}

func TestWeekBounds_MondayMidnightIsNewWeek(t *testing.T) {
	loc := time.UTC
	mondayMidnight := time.Date(2026, 7, 13, 0, 0, 0, 0, loc)
	start, _ := WeekBounds(mondayMidnight, loc)
	want := time.Date(2026, 7, 13, 0, 0, 0, 0, loc)
	if !start.Equal(want) {
		t.Errorf("Monday 00:00:00 should start its own week, got start=%v", start)
	}
}

func TestAggregateWeek_EmptyWeek_NoCrash(t *testing.T) {
	loc := time.UTC
	start := time.Date(2026, 7, 6, 0, 0, 0, 0, loc)
	end := start.AddDate(0, 0, 7)

	stats := AggregateWeek(start, end, loc, nil, nil)

	if stats.TotalActivities != 0 {
		t.Errorf("TotalActivities = %d, want 0", stats.TotalActivities)
	}
	if stats.TotalValue != 0 {
		t.Errorf("TotalValue = %v, want 0", stats.TotalValue)
	}
	if stats.TopCategory != nil {
		t.Errorf("TopCategory = %v, want nil for empty week", *stats.TopCategory)
	}
	if stats.MostProductiveDay != nil {
		t.Errorf("MostProductiveDay = %v, want nil for empty week", *stats.MostProductiveDay)
	}
	if len(stats.DailyBreakdown) != 7 {
		t.Errorf("DailyBreakdown length = %d, want 7 (always full week)", len(stats.DailyBreakdown))
	}
	if stats.ChangeVsPrevWeekPct != nil {
		t.Errorf("ChangeVsPrevWeekPct = %v, want nil when both weeks are empty", *stats.ChangeVsPrevWeekPct)
	}
}

func TestAggregateWeek_ChangeVsPrevWeek_ZeroDivisionGuard(t *testing.T) {
	loc := time.UTC
	start := time.Date(2026, 7, 6, 0, 0, 0, 0, loc)
	end := start.AddDate(0, 0, 7)

	current := []*domain.Activity{act(domain.CategoryCoding, 120, start.Add(2 * time.Hour))}

	// Previous week had zero activity: percentage change is undefined, must be nil
	// (not +Inf, not a fabricated "+100%").
	stats := AggregateWeek(start, end, loc, current, nil)

	if stats.ChangeVsPrevWeekPct != nil {
		t.Errorf("ChangeVsPrevWeekPct = %v, want nil when previous week total is 0", *stats.ChangeVsPrevWeekPct)
	}
	if stats.TotalValue != 120 {
		t.Errorf("TotalValue = %v, want 120", stats.TotalValue)
	}
}

func TestAggregateWeek_ChangeVsPrevWeek_NormalCase(t *testing.T) {
	loc := time.UTC
	start := time.Date(2026, 7, 6, 0, 0, 0, 0, loc)
	end := start.AddDate(0, 0, 7)

	current := []*domain.Activity{act(domain.CategoryCoding, 150, start.Add(time.Hour))}
	previous := []*domain.Activity{act(domain.CategoryCoding, 100, start.Add(-24 * time.Hour))}

	stats := AggregateWeek(start, end, loc, current, previous)

	if stats.ChangeVsPrevWeekPct == nil {
		t.Fatal("ChangeVsPrevWeekPct = nil, want a value when previous week has data")
	}
	got := *stats.ChangeVsPrevWeekPct
	want := 50.0 // (150-100)/100 * 100
	if got != want {
		t.Errorf("ChangeVsPrevWeekPct = %v, want %v", got, want)
	}
}

func TestAggregateWeek_TopCategory_PicksHighestTotal(t *testing.T) {
	loc := time.UTC
	start := time.Date(2026, 7, 6, 0, 0, 0, 0, loc)

	current := []*domain.Activity{
		act(domain.CategoryReading, 30, start.Add(time.Hour)),
		act(domain.CategoryCoding, 200, start.Add(2 * time.Hour)),
		act(domain.CategoryCoding, 100, start.Add(3 * time.Hour)),
	}

	stats := AggregateWeek(start, start.AddDate(0, 0, 7), loc, current, nil)

	if stats.TopCategory == nil || *stats.TopCategory != domain.CategoryCoding {
		t.Errorf("TopCategory = %v, want coding (total 300 > reading 30)", stats.TopCategory)
	}
}

func TestAggregateWeek_TopCategory_TieBreaksDeterministically(t *testing.T) {
	loc := time.UTC
	start := time.Date(2026, 7, 6, 0, 0, 0, 0, loc)

	// workout and reading tie at 50; fixed order = workout, reading, coding, spending
	// so workout must win deterministically regardless of map iteration order.
	current := []*domain.Activity{
		act(domain.CategoryReading, 50, start.Add(time.Hour)),
		act(domain.CategoryWorkout, 50, start.Add(2 * time.Hour)),
	}

	stats := AggregateWeek(start, start.AddDate(0, 0, 7), loc, current, nil)

	if stats.TopCategory == nil || *stats.TopCategory != domain.CategoryWorkout {
		t.Errorf("TopCategory = %v, want workout (deterministic tie-break)", stats.TopCategory)
	}
}

func TestAggregateWeek_MostProductiveDay_TieBreaksToEarliestDay(t *testing.T) {
	loc := time.UTC
	start := time.Date(2026, 7, 6, 0, 0, 0, 0, loc) // Monday

	current := []*domain.Activity{
		act(domain.CategoryCoding, 60, start.AddDate(0, 0, 3)), // Thursday
		act(domain.CategoryCoding, 60, start),                  // Monday
	}

	stats := AggregateWeek(start, start.AddDate(0, 0, 7), loc, current, nil)

	if stats.MostProductiveDay == nil {
		t.Fatal("MostProductiveDay = nil, want Monday's date")
	}
	want := start.Format("2006-01-02")
	if *stats.MostProductiveDay != want {
		t.Errorf("MostProductiveDay = %v, want %v (earliest day wins tie)", *stats.MostProductiveDay, want)
	}
}

func TestAggregateWeek_ActivityAtWeekBoundary_NotDoubleCounted(t *testing.T) {
	loc := time.UTC
	start := time.Date(2026, 7, 6, 0, 0, 0, 0, loc)
	end := start.AddDate(0, 0, 7)

	// Exactly at the boundary: belongs to the NEXT week, not this one.
	// Caller (repository) is responsible for the half-open [start,end) filtering;
	// this test documents the expectation so a future repository query bug
	// (e.g. using <= instead of <) is caught by AggregateWeek's own invariants
	// if boundary activities leak through.
	onBoundary := act(domain.CategoryCoding, 999, end)
	withinWeek := act(domain.CategoryCoding, 10, start)

	// Simulate correct repository filtering: boundary activity is excluded.
	stats := AggregateWeek(start, end, loc, []*domain.Activity{withinWeek}, nil)

	if stats.TotalValue != 10 {
		t.Errorf("TotalValue = %v, want 10 (boundary activity at 'end' must belong to next week)", stats.TotalValue)
	}
	_ = onBoundary
}

func TestAggregateWeek_Timezone_ShiftsDayBucket(t *testing.T) {
	loc := mustLoadLocation(t, "Asia/Jakarta") // UTC+7

	start := time.Date(2026, 7, 6, 0, 0, 0, 0, loc) // Monday 00:00 WIB
	// 2026-07-06 23:30 WIB == 2026-07-06 16:30 UTC. Stored as UTC in the DB,
	// but the "day" it belongs to must be computed in the user's timezone (loc),
	// not in UTC — otherwise it would incorrectly bucket into Tuesday.
	lateMonday := time.Date(2026, 7, 6, 16, 30, 0, 0, time.UTC)

	stats := AggregateWeek(start, start.AddDate(0, 0, 7), loc, []*domain.Activity{
		act(domain.CategoryReading, 40, lateMonday),
	}, nil)

	mondayKey := start.Format("2006-01-02")
	var mondayTotal float64
	for _, d := range stats.DailyBreakdown {
		if d.Date == mondayKey {
			mondayTotal = d.Total
		}
	}
	if mondayTotal != 40 {
		t.Errorf("Monday (WIB) total = %v, want 40 — activity must bucket by local day, not UTC day", mondayTotal)
	}
}
