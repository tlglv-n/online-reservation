package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"reservation-system/internal/domain/recruiter"
	"reservation-system/pkg/store"
	"strings"
)

type RecruiterRepository struct {
	db *sqlx.DB
}

func NewRecruiterRepository(db *sqlx.DB) *RecruiterRepository {
	return &RecruiterRepository{
		db: db,
	}
}

func (r *RecruiterRepository) List(ctx context.Context) (dest []recruiter.Entity, err error) {
	query := `
		SELECT id, full_name, email, phone
		FROM recruiters
		ORDER BY id`

	err = r.db.SelectContext(ctx, &dest, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list candidates: %w", err)
	}

	return
}

func (r *RecruiterRepository) Add(ctx context.Context, data recruiter.Entity) (id string, err error) {
	query := `
		INSERT INTO recruiters (full_name, email, phone)
		VALUES ($1, $2, $3)
		RETURNING id`

	args := []any{data.FullName, data.Email, data.Phone}

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to add candidate: %w", err)
	}

	return
}

func (r *RecruiterRepository) Get(ctx context.Context, id string) (dest recruiter.Entity, err error) {
	query := `
		SELECT id, full_name, email, phone
		FROM recruiters
		WHERE id = $1`

	args := []any{id}

	err = r.db.GetContext(ctx, &dest, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dest, store.ErrorNotFound
		}
		return dest, fmt.Errorf("failed to get candidate with id %s: %w", id, err)
	}

	return
}

func (r *RecruiterRepository) Update(ctx context.Context, id string, data recruiter.Entity) (err error) {
	sets, args := r.prepareArgs(data)
	if len(args) == 0 {
		return errors.New("no fields to update")
	}

	args = append(args, id)
	sets = append(sets, "updated_at = CURRENT_TIMESTAMP")

	setClause := strings.Join(sets, ", ")
	argPosition := len(args)

	query := fmt.Sprintf("UPDATE recruiters SET %s WHERE id = $%d RETURNING id", setClause, argPosition)

	var returnedID string
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&returnedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return store.ErrorNotFound
		}
		return fmt.Errorf("failed to update candidate with id %s: %w", id, err)
	}

	return
}

func (r *RecruiterRepository) prepareArgs(data recruiter.Entity) (sets []string, args []any) {
	if data.Email != nil {
		args = append(args, data.Email)
		sets = append(sets, fmt.Sprintf("email = $%d", len(args)))
	}

	if data.FullName != nil {
		args = append(args, data.FullName)
		sets = append(sets, fmt.Sprintf("full_name = $%d", len(args)))
	}

	if data.Phone != nil {
		args = append(args, data.Phone)
		sets = append(sets, fmt.Sprintf("phone = $%d", len(args)))
	}

	return
}

func (r *RecruiterRepository) Delete(ctx context.Context, id string) (err error) {
	query := `
		DELETE FROM recruiters
		WHERE id = $1
		RETURNING id`

	args := []any{id}

	var returnedID string
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&returnedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return store.ErrorNotFound
		}
		return fmt.Errorf("failed to delete candidate with id %s: %w", id, err)
	}

	return
}
