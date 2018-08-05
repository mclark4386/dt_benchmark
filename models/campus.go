package models

import (
	"encoding/json"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type Campus struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Name      string    `json:"name" db:"name"`
	Admins    Users     `many_to_many:"campus_admins"`
}

// String is not required by pop and may be deleted
func (c Campus) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Campus is not required by pop and may be deleted
type Campuses []Campus

// String is not required by pop and may be deleted
func (c Campuses) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (c *Campus) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: c.Name, Name: "Name"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (c *Campus) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (c *Campus) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

//Should remove excess and add new admins
func (c *Campus) UpdateAdmins(tx *pop.Connection, admins []string) {
	//Remove existing that have been removed from the list
	args := []interface{}{c.ID}
	for _, admin := range admins {
		args = append(args, admin)
	}

	var sql string
	if len(admins) > 0 {
		sql = "delete from campus_admins where campus_id = ? and user_id not in (?" + strings.Repeat(",?", len(admins)-1) + ")"
	} else {
		sql = "delete from campus_admins where campus_id = ?"
	}
	log.Printf(sql)
	del_count, del_err := tx.RawQuery(sql, args...).ExecWithCount()
	log.Printf("Delete count: %d err: %v\n", del_count, del_err)

	//Get existing IDs
	current_ids := []string{}

	if err := tx.RawQuery("select user_id from campus_admins where campus_id = ?", c.ID).All(&current_ids); err != nil {
		log.Printf("Couldn't get existing ids: %v\n", err)
	}

	//Add new links
	for _, id := range admins {
		if i := sort.SearchStrings(current_ids, id); i >= len(current_ids) || current_ids[i] != id {
			rid, uerr := uuid.FromString(id)
			if uerr != nil {
				log.Printf("failed to convert str id(%s) to uuid, err: %v\n", id, uerr)
				continue
			}
			ca := CampusAdmin{}
			ca.CampusID = c.ID
			ca.UserID = rid
			verrs, err := tx.ValidateAndCreate(&ca)
			if err != nil {
				log.Printf("Error creating ca for %s <-> %s with err: %v and verrs: %v\n", c.ID, id, err, verrs)
			}
			if verrs.HasAny() {
				log.Printf("Error creating ca for %s <-> %s with verrs: %v\n", c.ID, id, verrs)
			}
		}
	}
}
