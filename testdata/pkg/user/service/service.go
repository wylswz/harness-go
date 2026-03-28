package service

import (
	fileRepo "github.com/wylswz/harness-go/testdata/pkg/file/repo"
	server "github.com/wylswz/harness-go/testdata/pkg/server"
	"github.com/wylswz/harness-go/testdata/pkg/user/repo"
)

// Service coordinates user operations.
type Service struct {
	R        *repo.Store
	FileRepo *fileRepo.Store
	Server   *server.Server
}

// User loads a user from the store.
func (s *Service) User(id int) (*repo.User, error) {
	return s.R.Get(id)
}
