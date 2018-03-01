package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type Page struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	Slug      string    `json:"slug" db:"slug"`
}

// String is not required by pop and may be deleted
func (p Page) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Pages is not required by pop and may be deleted
type Pages []Page

// String is not required by pop and may be deleted
func (p Pages) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (p *Page) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	return validate.Validate(
		&validators.StringIsPresent{Field: p.Title, Name: "Title"},
		&validators.StringIsPresent{Field: p.Content, Name: "Content"},
		&validators.StringIsPresent{Field: p.Slug, Name: "Slug"},
		&validators.FuncValidator{
			Field:   p.Slug,
			Name:    "Slug",
			Message: "%s is already taken",
			Fn: func() bool {
				var b bool
				q := tx.Where("slug = ?", p.Slug)
				if p.ID != uuid.Nil {
					q = q.Where("id != ?", p.ID)
				}
				b, err = q.Exists(p)
				if err != nil {
					return false
				}
				return !b
			},
		},
	), err
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (p *Page) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (p *Page) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
