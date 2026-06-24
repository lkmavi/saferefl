// Package main demonstrates CopyFields — copying matching exported fields
// between two structs with different types. Common in DTO↔entity mapping,
// API response projection, and version migration patterns.
package main

import (
	"fmt"

	"github.com/lkmavi/saferefl"
)

// UserEntity is a full database entity.
type UserEntity struct {
	ID        int
	Name      string
	Email     string
	Password  string // sensitive — not in DTO
	Score     float64
	Active    bool
	CreatedAt string
}

// UserDTO is the public API representation — a subset of UserEntity.
type UserDTO struct {
	ID     int
	Name   string
	Email  string
	Score  float64
	Active bool
	// Password is intentionally absent
}

// UserUpdateRequest is what the API client sends — only mutable fields.
type UserUpdateRequest struct {
	Name  string
	Email string
	Score float64
}

func main() {
	entity := &UserEntity{
		ID:        1,
		Name:      "Alice",
		Email:     "alice@example.com",
		Password:  "s3cr3t",
		Score:     9.5,
		Active:    true,
		CreatedAt: "2024-01-15",
	}

	// --- Entity → DTO: copy only matching fields ---
	// Password and CreatedAt are absent from UserDTO and are silently skipped.
	fmt.Println("=== Entity → DTO ===")
	dto := &UserDTO{}
	if err := saferefl.CopyFields(entity, dto); err != nil {
		panic(err)
	}
	fmt.Printf("  DTO: ID=%d Name=%q Email=%q Score=%.1f Active=%v\n",
		dto.ID, dto.Name, dto.Email, dto.Score, dto.Active)
	fmt.Printf("  Password not copied (absent in DTO): %q\n", "")

	// --- Request → Entity: apply only the fields the client may change ---
	// ID, Active, Password, CreatedAt are not in the request and are untouched.
	fmt.Println("\n=== UpdateRequest → Entity ===")
	req := &UserUpdateRequest{Name: "Alice B.", Email: "alice.b@example.com", Score: 9.8}
	if err := saferefl.CopyFields(req, entity); err != nil {
		panic(err)
	}
	fmt.Printf("  Entity after update: Name=%q Email=%q Score=%.1f\n",
		entity.Name, entity.Email, entity.Score)
	fmt.Printf("  Untouched: ID=%d Active=%v Password=%q\n",
		entity.ID, entity.Active, entity.Password)

	// --- Same-type copy: clone an entity ---
	fmt.Println("\n=== Same-type clone ===")
	clone := &UserEntity{}
	_ = saferefl.CopyFields(entity, clone)
	fmt.Printf("  Clone: %+v\n", *clone)

	// --- Convertible types: int32 → int64, float32 → float64 ---
	fmt.Println("\n=== Convertible-type copy (int32 → int64) ===")
	type Narrow struct {
		Count int32
		Rate  float32
	}
	type Wide struct {
		Count int64
		Rate  float64
		Extra string // no matching field in Narrow — left at zero
	}
	narrow := &Narrow{Count: 42, Rate: 1.5}
	wide := &Wide{}
	_ = saferefl.CopyFields(narrow, wide)
	fmt.Printf("  Wide: Count=%d Rate=%.1f Extra=%q\n", wide.Count, wide.Rate, wide.Extra)

	// --- Incompatible types are silently skipped ---
	fmt.Println("\n=== Incompatible types skipped ===")
	type SrcBad struct{ Tag []string }
	type DstBad struct{ Tag int } // []string → int: incompatible, silently skipped
	sbad := &SrcBad{Tag: []string{"go", "fast"}}
	dbad := &DstBad{Tag: 99}
	_ = saferefl.CopyFields(sbad, dbad)
	fmt.Printf("  DstBad.Tag unchanged: %d\n", dbad.Tag)

	// --- Promoted fields from embedded structs are included ---
	fmt.Println("\n=== Promoted fields from embedding ===")
	type Base struct{ Name string }
	type SrcEmbed struct {
		Base
		Score int
	}
	type DstFlat struct {
		Name  string // receives promoted field from SrcEmbed.Base.Name
		Score int
	}
	src := &SrcEmbed{Base: Base{Name: "Carol"}, Score: 88}
	dst := &DstFlat{}
	_ = saferefl.CopyFields(src, dst)
	fmt.Printf("  DstFlat: Name=%q Score=%d\n", dst.Name, dst.Score)
}
