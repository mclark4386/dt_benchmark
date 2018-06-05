package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type Resource struct {
	ID            uuid.UUID     `json:"id" db:"id"`
	CreatedAt     time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at" db:"updated_at"`
	Title         string        `json:"title" db:"title"`
	Url           string        `json:"url" db:"url"`
	Teams         Teams         `many_to_many:"team_resources"`
	BenchmarkItem BenchmarkItem `many_to_many:"benchmark_item_resources"`
}

// String is not required by pop and may be deleted
func (r Resource) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// Resources is not required by pop and may be deleted
type Resources []Resource

// String is not required by pop and may be deleted
func (r Resources) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *Resource) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: r.Title, Name: "Title"},
		&validators.StringIsPresent{Field: r.Url, Name: "Url"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (r *Resource) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (r *Resource) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
