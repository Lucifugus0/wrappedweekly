package aiprovider

import (
	"fmt"
	"strings"

	"wrappedweekly/backend/internal/domain"
)

// MockProvider generates a deterministic, template-based narrative without
// calling any external LLM API. This is the default provider (AI_PROVIDER=mock)
// so the app runs end-to-end without an API key, as required by the study case.
type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (p *MockProvider) GenerateNarrative(stats domain.RecapStats, userName string) (string, error) {
	if stats.TotalActivities == 0 {
		return fmt.Sprintf(
			"Halo %s! Minggu ini (%s - %s) belum ada aktivitas yang tercatat. "+
				"Yuk mulai catat aktivitasmu minggu depan supaya recap-nya makin seru!",
			userName, stats.WeekStart.Format("02 Jan"), stats.WeekEnd.AddDate(0, 0, -1).Format("02 Jan"),
		), nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Halo %s! Ini rekap mingguanmu (%s - %s).\n\n",
		userName, stats.WeekStart.Format("02 Jan"), stats.WeekEnd.AddDate(0, 0, -1).Format("02 Jan"))

	fmt.Fprintf(&sb, "Kamu mencatat %d aktivitas dengan total nilai %.0f. ", stats.TotalActivities, stats.TotalValue)

	if stats.TopCategory != nil {
		fmt.Fprintf(&sb, "Kategori paling aktif minggu ini: %s. ", categoryLabel(*stats.TopCategory))
	}
	if stats.MostProductiveDay != nil {
		fmt.Fprintf(&sb, "Hari paling produktifmu adalah %s. ", *stats.MostProductiveDay)
	}

	switch {
	case stats.ChangeVsPrevWeekPct == nil:
		sb.WriteString("Belum ada data minggu sebelumnya untuk dibandingkan.")
	case *stats.ChangeVsPrevWeekPct > 0:
		fmt.Fprintf(&sb, "Aktivitasmu naik %.0f%% dibanding minggu lalu. Keren, terus pertahankan!", *stats.ChangeVsPrevWeekPct)
	case *stats.ChangeVsPrevWeekPct < 0:
		fmt.Fprintf(&sb, "Aktivitasmu turun %.0f%% dibanding minggu lalu. Yuk semangat lagi minggu depan!", -*stats.ChangeVsPrevWeekPct)
	default:
		sb.WriteString("Aktivitasmu stabil, sama seperti minggu lalu.")
	}

	return sb.String(), nil
}

func categoryLabel(c domain.ActivityCategory) string {
	labels := map[domain.ActivityCategory]string{
		domain.CategoryWorkout:  "olahraga",
		domain.CategoryReading:  "membaca",
		domain.CategoryCoding:   "ngoding",
		domain.CategorySpending: "pengeluaran",
	}
	if l, ok := labels[c]; ok {
		return l
	}
	return string(c)
}
