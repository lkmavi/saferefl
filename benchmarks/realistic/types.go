package realistic

// DTO mapping scenario (5 fields, heterogeneous types).

type UserSrc struct {
	ID     int
	Name   string
	Email  string
	Score  float64
	Active bool
}

type UserDst struct {
	ID     int
	Name   string
	Email  string
	Score  float64
	Active bool
}

func newUserSrc() UserSrc {
	return UserSrc{ID: 1, Name: "Alice", Email: "alice@example.com", Score: 9.5, Active: true}
}

// JSON-like decode scenario (10 fields, mixed types).

type Product struct {
	ID          int64
	SKU         string
	Title       string
	Description string
	Price       float64
	Stock       int
	Weight      float64
	Active      bool
	Category    string
	Tags        string
}

func newProduct() Product {
	return Product{
		ID:          42,
		SKU:         "SKU-001",
		Title:       "Widget",
		Description: "A fine widget",
		Price:       9.99,
		Stock:       100,
		Weight:      0.5,
		Active:      true,
		Category:    "Gadgets",
		Tags:        "new,sale",
	}
}

// DI resolve scenario: lookup by type, inject 3 dependencies.

type ServiceA struct{ Name string }
type ServiceB struct{ Name string }
type ServiceC struct{ Name string }

type AppServices struct {
	A *ServiceA
	B *ServiceB
	C *ServiceC
}

// ORM row scan scenario (10 columns → struct).

type ORMRow struct {
	Col1  int64
	Col2  string
	Col3  string
	Col4  float64
	Col5  int64
	Col6  bool
	Col7  string
	Col8  float64
	Col9  int64
	Col10 string
}

func newORMValues() []interface{} {
	return []interface{}{
		int64(1), "Alice", "alice@example.com", 9.5, int64(100),
		true, "admin", 0.75, int64(42), "active",
	}
}
