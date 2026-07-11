package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"wrappedweekly/backend/internal/domain"
)

type recapRepository struct {
	pool *pgxpool.Pool
}

func NewRecapRepository(pool *pgxpool.Pool) domain.RecapRepository {
	return &recapRepository{pool: pool}
}

func (r *recapRepository) Create(rec *domain.Recap) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	statsJSON, err := json.Marshal(rec.Stats)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO recaps (user_id, slug, week_start, week_end, stats, narrative)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`

	return r.pool.QueryRow(ctx, query, rec.UserID, rec.Slug, rec.WeekStart, rec.WeekEnd, statsJSON, rec.Narrative).
		Scan(&rec.ID, &rec.CreatedAt)
}

func scanRecap(row pgx.Row) (*domain.Recap, error) {
	rec := &domain.Recap{}
	var statsJSON []byte
	err := row.Scan(&rec.ID, &rec.UserID, &rec.Slug, &rec.WeekStart, &rec.WeekEnd, &statsJSON, &rec.Narrative, &rec.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(statsJSON, &rec.Stats); err != nil {
		return nil, err
	}
	return rec, nil
}

const recapColumns = `id, user_id, slug, week_start, week_end, stats, narrative, created_at`

func (r *recapRepository) FindByID(id string) (*domain.Recap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := r.pool.QueryRow(ctx, `SELECT `+recapColumns+` FROM recaps WHERE id = $1`, id)
	return scanRecap(row)
}

func (r *recapRepository) FindBySlug(slug string) (*domain.Recap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := r.pool.QueryRow(ctx, `SELECT `+recapColumns+` FROM recaps WHERE slug = $1`, slug)
	return scanRecap(row)
}

func (r *recapRepository) FindByUserAndWeek(userID string, weekStart time.Time) (*domain.Recap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := r.pool.QueryRow(ctx, `SELECT `+recapColumns+` FROM recaps WHERE user_id = $1 AND week_start = $2`, userID, weekStart)
	return scanRecap(row)
}

func (r *recapRepository) ListByUser(userID string) ([]*domain.Recap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.pool.Query(ctx, `SELECT `+recapColumns+` FROM recaps WHERE user_id = $1 ORDER BY week_start DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recaps []*domain.Recap
	for rows.Next() {
		rec, err := scanRecap(rows)
		if err != nil {
			return nil, err
		}
		recaps = append(recaps, rec)
	}
	return recaps, rows.Err()
}
