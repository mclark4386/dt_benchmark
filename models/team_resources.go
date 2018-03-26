package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"

)

type TeamResource struct {
	ID         uuid.UUID `json:"id" db:"id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	TeamID     uuid.UUID `json:"team_id" db:"team_id"`
	ResourceID uuid.UUID `json:"resource_id" db:"resource_id"`
}

// String is not required by pop and may be deleted
func (t TeamResource) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// TeamResources is not required by pop and may be deleted
type TeamResources []TeamResource

// String is not required by pop and may be deleted
func (t TeamResources) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

func (t *TeamResource) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	return validate.Validate(

		&validators.UUIDIsPresent{Field: t.TeamID, Name:"TeamID"},
		&validators.UUIDIsPresent{Field: t.ResourceID, Name: "ResourceID"},
		&validators.FuncValidator{

			Fn: func() bool {
				var b bool
					q := tx.Where("team_id = ?", t.TeamID)
				if t.ResourceID != uuid.Nil {
					q = q.Where("resource_id != ?", t.ID)
				}
				b, err = q.Exists(t)
				if err != nil {
					return false
				}
				return !b
			},
		},
	), err
}


/*
// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *TeamResource) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}*/

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *TeamResource) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *TeamResource) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
