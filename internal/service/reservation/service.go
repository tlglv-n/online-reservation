package reservation

import (
	"reservation-system/internal/domain/candidate"
	"reservation-system/internal/domain/recruiter"
)

type Configuration func(s *Service) error

// Service is an implementation of the Service
type Service struct {
	candidateRepository candidate.Repository
	recruiterRepository recruiter.Repository
}

// New takes a variable amount of Configuration functions and returns a new Service
// Each Configuration will be called in the order they are passed in
func New(configs ...Configuration) (s *Service, err error) {
	// Add the service
	s = &Service{}

	// Apply all Configurations passed in
	for _, cfg := range configs {
		// Pass the service into the configuration function
		if err = cfg(s); err != nil {
			return
		}
	}
	return
}

func WithCandidateRepository(candidateRepository candidate.Repository) Configuration {
	// return a function that matches the Configuration alias,
	// You need to return this so that the parent function can take in all the needed parameters
	return func(s *Service) error {
		s.candidateRepository = candidateRepository
		return nil
	}
}

func WithRecruiterRepository(recruiterRepository recruiter.Repository) Configuration {
	// return a function that matches the Configuration alias,
	// You need to return this so that the parent function can take in all the needed parameters
	return func(s *Service) error {
		s.recruiterRepository = recruiterRepository
		return nil
	}
}
