package repo

// User is a minimal domain type for harness tests.
type User struct {
	ID   int
	Name string
}

// Store is a trivial data accessor.
type Store struct{}

// Get returns a user by id (stub).
func (s *Store) Get(id int) (*User, error) {
	return &User{ID: id, Name: "stub"}, nil
}
