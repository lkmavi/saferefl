package spike

// User is the benchmark subject — covers int64, string, float64, bool field types.
type User struct {
	ID     int64
	Name   string
	Email  string
	Score  float64
	Active bool
}

func newUser() *User {
	return &User{
		ID:     42,
		Name:   "Alice",
		Email:  "alice@example.com",
		Score:  9.5,
		Active: true,
	}
}
