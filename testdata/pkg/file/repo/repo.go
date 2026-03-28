package repo

// File is a minimal file handle for harness tests.
type File struct {
	Path string
}

// Store opens files (stub).
type Store struct{}

// Open returns a file for path (stub).
func (s *Store) Open(path string) (*File, error) {
	return &File{Path: path}, nil
}
