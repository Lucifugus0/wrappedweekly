package usecase

import (
	"strings"
	"time"

	"wrappedweekly/backend/internal/domain"
	"wrappedweekly/backend/pkg/apperror"
)

type ActivityUsecase struct {
	activities domain.ActivityRepository
}

func NewActivityUsecase(activities domain.ActivityRepository) *ActivityUsecase {
	return &ActivityUsecase{activities: activities}
}

type ActivityInput struct {
	Category   string
	Value      float64
	Note       *string
	OccurredAt time.Time
}

func validateActivityInput(in ActivityInput) (domain.ActivityCategory, error) {
	category := domain.ActivityCategory(strings.ToLower(strings.TrimSpace(in.Category)))
	if !category.IsValid() {
		return "", apperror.BadRequest("category harus salah satu dari: workout, reading, coding, spending")
	}
	if in.Value < 0 {
		return "", apperror.BadRequest("value tidak boleh negatif")
	}
	if in.OccurredAt.IsZero() {
		return "", apperror.BadRequest("occurred_at wajib diisi")
	}
	return category, nil
}

func (u *ActivityUsecase) Create(userID string, in ActivityInput) (*domain.Activity, error) {
	category, err := validateActivityInput(in)
	if err != nil {
		return nil, err
	}

	activity := &domain.Activity{
		UserID:     userID,
		Category:   category,
		Value:      in.Value,
		Note:       in.Note,
		OccurredAt: in.OccurredAt,
	}
	if err := u.activities.Create(activity); err != nil {
		return nil, apperror.Internal("gagal menyimpan aktivitas")
	}
	return activity, nil
}

func (u *ActivityUsecase) Update(userID, id string, in ActivityInput) (*domain.Activity, error) {
	existing, err := u.getOwned(userID, id)
	if err != nil {
		return nil, err
	}

	category, err := validateActivityInput(in)
	if err != nil {
		return nil, err
	}

	existing.Category = category
	existing.Value = in.Value
	existing.Note = in.Note
	existing.OccurredAt = in.OccurredAt

	if err := u.activities.Update(existing); err != nil {
		return nil, apperror.Internal("gagal memperbarui aktivitas")
	}
	return existing, nil
}

func (u *ActivityUsecase) Delete(userID, id string) error {
	if _, err := u.getOwned(userID, id); err != nil {
		return err
	}
	if err := u.activities.Delete(id); err != nil {
		return apperror.Internal("gagal menghapus aktivitas")
	}
	return nil
}

func (u *ActivityUsecase) Get(userID, id string) (*domain.Activity, error) {
	return u.getOwned(userID, id)
}

func (u *ActivityUsecase) List(userID string, page, size int, filter domain.ActivityFilter) ([]*domain.Activity, int, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	if filter.Category != "" && !filter.Category.IsValid() {
		return nil, 0, apperror.BadRequest("category filter tidak valid")
	}
	activities, total, err := u.activities.ListByUser(userID, page, size, filter)
	if err != nil {
		return nil, 0, apperror.Internal("gagal mengambil daftar aktivitas")
	}
	return activities, total, nil
}

// getOwned fetches an activity and verifies it belongs to userID.
// Returns 404 (not 403) for both "not found" and "belongs to another user"
// so we don't leak whether a given ID exists to non-owners.
func (u *ActivityUsecase) getOwned(userID, id string) (*domain.Activity, error) {
	activity, err := u.activities.FindByID(id)
	if err != nil {
		return nil, apperror.Internal("gagal mengambil aktivitas")
	}
	if activity == nil || activity.UserID != userID {
		return nil, apperror.NotFound("aktivitas tidak ditemukan")
	}
	return activity, nil
}
