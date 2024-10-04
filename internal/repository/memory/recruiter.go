package memory

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"reservation-system/internal/domain/recruiter"
	"sync"
)

type RecruiterRepository struct {
	db map[string]recruiter.Entity
	sync.RWMutex
}

func NewRecruiterRepository() *RecruiterRepository {
	return &RecruiterRepository{
		db: make(map[string]recruiter.Entity),
	}
}

func (r *RecruiterRepository) List(ctx context.Context) (dest []recruiter.Entity, err error) {
	r.RLock()
	defer r.RUnlock()

	dest = make([]recruiter.Entity, 0, len(r.db))
	for _, data := range r.db {
		dest = append(dest, data)
	}

	return
}

func (r *RecruiterRepository) Add(ctx context.Context, data recruiter.Entity) (dest string, err error) {
	r.Lock()
	defer r.Unlock()

	id := r.generateID()
	data.ID = id
	r.db[id] = data

	return id, nil
}

func (r *RecruiterRepository) Get(ctx context.Context, id string) (dest recruiter.Entity, err error) {
	r.RLock()
	defer r.RUnlock()

	dest, ok := r.db[id]
	if !ok {
		err = sql.ErrNoRows
		return
	}

	return
}

func (r *RecruiterRepository) Update(ctx context.Context, id string, data recruiter.Entity) (err error) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.db[id]; !ok {
		return sql.ErrNoRows
	}
	r.db[id] = data

	return
}

func (r *RecruiterRepository) Delete(ctx context.Context, id string) (err error) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.db[id]; !ok {
		return sql.ErrNoRows
	}
	delete(r.db, id)

	return
}

func (r *RecruiterRepository) generateID() string {
	return uuid.New().String()
}
