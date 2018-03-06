package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/satori/go.uuid"
)

type TeamAdmin struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	UserId    uuid.UUID `json:"user_id" db:"user_id"` 
	TeamId    uuid.UUID `json:"team_id" db:"team_id"` 
	Admin     User  `belongs_to:"user"`
	Team      Team  `belongs_to:"team"`
}

// String is not required by pop and may be deleted
func (t TeamAdmin) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// TeamAdmins is not required by pop and may be deleted
type TeamAdmins []TeamAdmin

// String is not required by pop and may be deleted
func (t TeamAdmins) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *TeamAdmin) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *TeamAdmin) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *TeamAdmin) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
