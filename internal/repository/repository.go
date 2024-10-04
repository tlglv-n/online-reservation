package repository

import (
	"reservation-system/internal/domain/candidate"
	"reservation-system/internal/domain/recruiter"
	"reservation-system/internal/repository/memory"
	"reservation-system/pkg/store"
)

type Configuration func(r *Repository) error

type Repository struct {
	postgres store.SQLX
	//TODO: mongo?

	Recruiter recruiter.Repository
	Candidate candidate.Repository
}

func New(configs ...Configuration) (s *Repository, err error) {
	// Create the repository
	s = &Repository{}

	// Apply all Configurations passed in
	for _, cfg := range configs {
		// Pass the repository into the configuration function
		if err = cfg(s); err != nil {
			return
		}
	}

	return
}

func (r *Repository) Close() {
	if r.postgres.Client != nil {
		r.postgres.Client.Close()
	}
}

func WithMemoryStore() Configuration {
	return func(s *Repository) (err error) {
		s.Recruiter = memory.NewRecruiterRepository()
		s.Candidate = memory.NewCandidateRepository()

		return
	}
}
