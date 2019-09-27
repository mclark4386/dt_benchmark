package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const minPasswordLength int = 8

type User struct {
	ID                   uuid.UUID  `json:"id" db:"id"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
	Email                string     `json:"email" db:"email"`
	FirstName            string     `json:"first_name" db:"first_name"`
	LastName             string     `json:"last_name" db:"last_name"`
	Name                 string     `json:"name" db:"name" select:"(users.first_name||' '||users.last_name) as name" rw:"r"`
	PasswordHash         string     `json:"-" db:"password_hash"`
	Password             string     `json:"password" db:"-"`
	PasswordConfirmation string     `json:"password_confirm" db:"-"`
	IsSuperAdmin         bool       `json:"super_admin" db:"super_admin"`
	TeamsIAdmin          Teams      `many_to_many:"team_admins"`
	CampusesIAdmin       Campuses   `many_to_many:"campus_admins"`
	RefreshToken         nulls.UUID `json:"-" db:"refresh_token"`
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Create wraps up the pattern of encrypting the password and
// running validations. Useful when writing tests.
func (u *User) Create(tx *pop.Connection) (*validate.Errors, error) {
	return tx.ValidateAndCreate(u)
}

// Update handles the extra work possibly needed during user update,
// hashing password and making email lowercase for consistency
func (u *User) Update(tx *pop.Connection) (*validate.Errors, error) {
	return tx.Eager().ValidateAndUpdate(u)
}

func (u *User) sanitizeFields() *User {
	// force email to lowercase for better matching
	u.Email = strings.ToLower(u.Email)

	// wipe out password field after it's been hashed
	u.Password = ""
	u.PasswordConfirmation = ""
	return u
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Email, Name: "Email"},

		// check to see if the email address is already taken:
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
	return validate.NewErrors(), err
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (u *User) validatePassword() *validate.Errors {
	passwordLengthValidator := validators.StringLengthInRange{
		Field:   u.Password,
		Name:    "Password",
		Message: "Password is not long enough. Minimum of 8 characters required.",
		Min:     minPasswordLength,
		Max:     0,
	}
	passwordConfirmValidator := validators.StringsMatch{
		Field:   u.Password,
		Name:    "PasswordConfirm",
		Message: "Password does not match confirmation",
		Field2:  u.PasswordConfirmation,
	}
	verrs := validate.NewErrors()
	passwordLengthValidator.IsValid(verrs)
	passwordConfirmValidator.IsValid(verrs)
	// spew.Printf("validations at this point: %+v\n", verrs)
	return verrs
}

func encryptPassword(password string) (string, error) {
	ph, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(ph), err
}

func (u *User) BeforeSave(tx *pop.Connection) error {
	if len(u.Password) != 0 {
		verrs := u.validatePassword()
		if verrs.HasAny() {
			return errors.New(verrs.Error())
		}
		ph, err := encryptPassword(u.Password)
		if err != nil {
			return errors.WithStack(err)
		}
		u.PasswordHash = ph
	}

	u = u.sanitizeFields()
	return nil
}

// AfterCreate builds a ReferralCode for a user after it is created.
func (u *User) AfterCreate(tx *pop.Connection) error {
	tx.Reload(u)
	return nil
}
