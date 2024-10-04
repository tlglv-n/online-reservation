package memory

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"reservation-system/internal/domain/candidate"
	"sync"
)

type CandidateRepository struct {
	db map[string]candidate.Entity
	sync.RWMutex
}

func NewCandidateRepository() *CandidateRepository {
	return &CandidateRepository{
		db: make(map[string]candidate.Entity),
	}
}

func (r *CandidateRepository) List(ctx context.Context) (dest []candidate.Entity, err error) {
	r.RLock()
	defer r.RUnlock()

	dest = make([]candidate.Entity, 0, len(r.db))
	for _, data := range r.db {
		dest = append(dest, data)
	}

	return
}

func (r *CandidateRepository) Add(ctx context.Context, data candidate.Entity) (dest string, err error) {
	r.Lock()
	defer r.Unlock()

	id := r.generateID()
	data.ID = id
	r.db[id] = data

	return id, nil
}

func (r *CandidateRepository) Get(ctx context.Context, id string) (dest candidate.Entity, err error) {
	r.RLock()
	defer r.RUnlock()

	dest, ok := r.db[id]
	if !ok {
		err = sql.ErrNoRows
		return
	}

	return
}

func (r *CandidateRepository) Update(ctx context.Context, id string, data candidate.Entity) (err error) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.db[id]; !ok {
		return sql.ErrNoRows
	}
	r.db[id] = data

	return
}

func (r *CandidateRepository) Delete(ctx context.Context, id string) (err error) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.db[id]; !ok {
		return sql.ErrNoRows
	}
	delete(r.db, id)

	return
}

func (r *CandidateRepository) generateID() string {
	return uuid.New().String()
}
