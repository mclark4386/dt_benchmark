package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
)

type BenchmarkItemResource struct {
	ID              uuid.UUID `json:"id" db:"id"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	BenchmarkItemID uuid.UUID `json:"benchmark_item_id" db:"benchmark_item_id"`
	ResourceID      uuid.UUID `json:"resource_id" db:"resource_id"`
}

// String is not required by pop and may be deleted
func (b BenchmarkItemResource) String() string {
	jb, _ := json.Marshal(b)
	return string(jb)
}

// BenchmarkItemResources is not required by pop and may be deleted
type BenchmarkItemResources []BenchmarkItemResource

// String is not required by pop and may be deleted
func (b BenchmarkItemResources) String() string {
	jb, _ := json.Marshal(b)
	return string(jb)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (b *BenchmarkItemResource) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (b *BenchmarkItemResource) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (b *BenchmarkItemResource) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
