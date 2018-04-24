package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                   uuid.UUID `json:"id" db:"id"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
	Email                string    `json:"email" db:"email"`
	PasswordHash         string    `json:"-" db:"password_hash"`
	Password             string    `json:"-" db:"-"`
	PasswordConfirmation string    `json:"-" db:"-"`
	IsSuperAdmin         bool      `json:"super_admin" db:"super_admin"`
	TeamsIAdmin          Teams     `many_to_many:"team_admins"`
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

func (u *User) PrepFields() (*validate.Errors, error) {
	if len(strings.TrimSpace(u.Email)) > 0 {
		u.Email = strings.ToLower(u.Email)
	}
	if len(strings.TrimSpace(u.Password)) > 0 {
		ph, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return validate.NewErrors(), errors.WithStack(err)
		}
		u.PasswordHash = string(ph)
	}
	return nil, nil
}

func (u *User) Update(tx *pop.Connection) (*validate.Errors, error) {
	if verrs, err := u.PrepFields(); (verrs != nil && verrs.HasAny()) || err != nil {
		return verrs, err
	}
	return tx.ValidateAndSave(u)
}

func (u *User) Create(tx *pop.Connection) (*validate.Errors, error) {
	if verrs, err := u.PrepFields(); (verrs != nil && verrs.HasAny()) || err != nil {
		return verrs, err
	}
	return tx.ValidateAndCreate(u)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Email, Name: "Email"},
		&validators.StringIsPresent{Field: u.PasswordHash, Name: "PasswordHash"},
		//check email is unique
		&validators.FuncValidator{
			Field:   u.Email,
			Name:    "Email",
			Message: "%s is already taken",
			Fn: func() bool {
				var b bool
				q := tx.Where("email = ?", u.Email)
				if u.ID != uuid.Nil {
					q = q.Where("id != ?", u.ID)
				}
				b, err = q.Exists(u)
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
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Password, Name: "Password"},
		&validators.StringsMatch{Name: "Password", Field: u.Password, Field2: u.PasswordConfirmation, Message: "Password does not match confirmation"},
	), err
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
