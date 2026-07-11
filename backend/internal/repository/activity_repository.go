package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"wrappedweekly/backend/internal/domain"
)

type activityRepository struct {
	pool *pgxpool.Pool
}

func NewActivityRepository(pool *pgxpool.Pool) domain.ActivityRepository {
	return &activityRepository{pool: pool}
}

func (r *activityRepository) Create(a *domain.Activity) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO activities (user_id, category, value, note, occurred_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	return r.pool.QueryRow(ctx, query, a.UserID, a.Category, a.Value, a.Note, a.OccurredAt).
		Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

func (r *activityRepository) FindByID(id string) (*domain.Activity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, user_id, category, value::float8, note, occurred_at, created_at, updated_at
		FROM activities WHERE id = $1`

	a := &domain.Activity{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&a.ID, &a.UserID, &a.Category, &a.Value, &a.Note, &a.OccurredAt, &a.CreatedAt, &a.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (r *activityRepository) Update(a *domain.Activity) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE activities
		SET category = $1, value = $2, note = $3, occurred_at = $4, updated_at = now()
		WHERE id = $5
		RETURNING updated_at`

	return r.pool.QueryRow(ctx, query, a.Category, a.Value, a.Note, a.OccurredAt, a.ID).
		Scan(&a.UpdatedAt)
}

func (r *activityRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.pool.Exec(ctx, `DELETE FROM activities WHERE id = $1`, id)
	return err
}

func (r *activityRepository) ListByUser(userID string, page, size int, filter domain.ActivityFilter) ([]*domain.Activity, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	where := "WHERE user_id = $1"
	args := []interface{}{userID}

	if filter.Category != "" {
		args = append(args, filter.Category)
		where += fmt.Sprintf(" AND category = $%d", len(args))
	}
	if !filter.From.IsZero() {
		args = append(args, filter.From)
		where += fmt.Sprintf(" AND occurred_at >= $%d", len(args))
	}
	if !filter.To.IsZero() {
		args = append(args, filter.To)
		where += fmt.Sprintf(" AND occurred_at < $%d", len(args))
	}

	var total int
	countQuery := "SELECT count(*) FROM activities " + where
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	args = append(args, size, offset)
	listQuery := fmt.Sprintf(`
		SELECT id, user_id, category, value::float8, note, occurred_at, created_at, updated_at
		FROM activities
		%s
		ORDER BY occurred_at DESC
		LIMIT $%d OFFSET $%d`, where, len(args)-1, len(args))

	rows, err := r.pool.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var activities []*domain.Activity
	for rows.Next() {
		a := &domain.Activity{}
		if err := rows.Scan(&a.ID, &a.UserID, &a.Category, &a.Value, &a.Note, &a.OccurredAt, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, 0, err
		}
		activities = append(activities, a)
	}
	return activities, total, rows.Err()
}

func (r *activityRepository) ListByUserInRange(userID string, from, to time.Time) ([]*domain.Activity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, category, value::float8, note, occurred_at, created_at, updated_at
		FROM activities
		WHERE user_id = $1 AND occurred_at >= $2 AND occurred_at < $3
		ORDER BY occurred_at ASC`, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []*domain.Activity
	for rows.Next() {
		a := &domain.Activity{}
		if err := rows.Scan(&a.ID, &a.UserID, &a.Category, &a.Value, &a.Note, &a.OccurredAt, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		activities = append(activities, a)
	}
	return activities, rows.Err()
}
