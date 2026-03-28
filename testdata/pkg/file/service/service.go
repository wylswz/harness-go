package service

import "github.com/wylswz/harness-go/testdata/pkg/file/repo"

// Service coordinates file operations.
type Service struct {
	R *repo.Store
}

// Open opens a file via the store.
func (s *Service) Open(path string) (*repo.File, error) {
	return s.R.Open(path)
}
