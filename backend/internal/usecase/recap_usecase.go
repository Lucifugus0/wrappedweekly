package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"wrappedweekly/backend/internal/domain"
	"wrappedweekly/backend/pkg/apperror"
)

// AppTimezone is the fixed reference timezone used to determine week/day
// boundaries across the whole app. The study case flags timezone handling as
// an explicit edge case; we pick one consistent, documented rule (UTC) rather
// than trusting a per-request client timezone, so aggregation is always
// reproducible. See docs/backend_output.md "Asumsi Agregasi".
var AppTimezone = time.UTC

type RecapUsecase struct {
	recaps     domain.RecapRepository
	activities domain.ActivityRepository
	users      domain.UserRepository
	ai         domain.AIProvider
}

func NewRecapUsecase(recaps domain.RecapRepository, activities domain.ActivityRepository, users domain.UserRepository, ai domain.AIProvider) *RecapUsecase {
	return &RecapUsecase{recaps: recaps, activities: activities, users: users, ai: ai}
}

// GenerateForWeek generates (or regenerates, if refDate falls in an already
// recapped week and force=true) a recap for the ISO week containing refDate.
func (u *RecapUsecase) GenerateForWeek(userID string, refDate time.Time, force bool) (*domain.Recap, error) {
	weekStart, weekEnd := WeekBounds(refDate, AppTimezone)

	existing, err := u.recaps.FindByUserAndWeek(userID, weekStart)
	if err != nil {
		return nil, apperror.Internal("gagal memeriksa recap yang sudah ada")
	}
	if existing != nil && !force {
		return existing, nil
	}

	user, err := u.users.FindByID(userID)
	if err != nil || user == nil {
		return nil, apperror.NotFound("user tidak ditemukan")
	}

	current, err := u.activities.ListByUserInRange(userID, weekStart, weekEnd)
	if err != nil {
		return nil, apperror.Internal("gagal mengambil aktivitas minggu ini")
	}

	prevStart, prevEnd := PrevWeekBounds(weekStart, weekEnd)
	previous, err := u.activities.ListByUserInRange(userID, prevStart, prevEnd)
	if err != nil {
		return nil, apperror.Internal("gagal mengambil aktivitas minggu sebelumnya")
	}

	stats := AggregateWeek(weekStart, weekEnd, AppTimezone, current, previous)

	narrative, err := u.ai.GenerateNarrative(stats, user.Name)
	if err != nil {
		return nil, apperror.Internal("gagal membuat narasi recap")
	}

	slug, err := generateSlug()
	if err != nil {
		return nil, apperror.Internal("gagal membuat slug")
	}

	recap := &domain.Recap{
		UserID:    userID,
		Slug:      slug,
		WeekStart: weekStart,
		WeekEnd:   weekEnd,
		Stats:     stats,
		Narrative: narrative,
	}
	if err := u.recaps.Create(recap); err != nil {
		return nil, apperror.Internal("gagal menyimpan recap")
	}
	return recap, nil
}

func (u *RecapUsecase) List(userID string) ([]*domain.Recap, error) {
	recaps, err := u.recaps.ListByUser(userID)
	if err != nil {
		return nil, apperror.Internal("gagal mengambil daftar recap")
	}
	return recaps, nil
}

func (u *RecapUsecase) Get(userID, id string) (*domain.Recap, error) {
	recap, err := u.recaps.FindByID(id)
	if err != nil {
		return nil, apperror.Internal("gagal mengambil recap")
	}
	if recap == nil || recap.UserID != userID {
		return nil, apperror.NotFound("recap tidak ditemukan")
	}
	return recap, nil
}

// GetPublicBySlug intentionally does not check ownership — this is the
// public share endpoint accessible without auth.
func (u *RecapUsecase) GetPublicBySlug(slug string) (*domain.Recap, error) {
	recap, err := u.recaps.FindBySlug(slug)
	if err != nil {
		return nil, apperror.Internal("gagal mengambil recap")
	}
	if recap == nil {
		return nil, apperror.NotFound("recap tidak ditemukan")
	}
	return recap, nil
}

func generateSlug() (string, error) {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
