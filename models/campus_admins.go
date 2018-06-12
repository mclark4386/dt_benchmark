package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

type CampusAdmin struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	CampusID  uuid.UUID `json:"campus_id" db:"campus_id"`
	Admin     User      `belongs_to:"user"`
	Campus    Campus    `belongs_to:"campus"`
}

// String is not required by pop and may be deleted
func (c CampusAdmin) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// CampusAdmins is not required by pop and may be deleted
type CampusAdmins []CampusAdmin

// String is not required by pop and may be deleted
func (c CampusAdmins) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (c *CampusAdmin) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (c *CampusAdmin) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (c *CampusAdmin) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
