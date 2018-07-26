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

type Team struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	PageSlug    string    `json:"page_slug" db:"page_slug"`
	Resources   Resources `many_to_many:"team_resources"`
	Admins      Users     `many_to_many:"team_admins"`
}

// String is not required by pop and may be deleted
func (t Team) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Teams is not required by pop and may be deleted
type Teams []Team

// String is not required by pop and may be deleted
func (t Teams) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Team) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	return validate.Validate(
		&validators.StringIsPresent{Field: t.Name, Name: "Name"},
		&validators.StringIsPresent{Field: t.Description, Name: "Description"},
		&validators.StringIsPresent{Field: t.PageSlug, Name: "PageSlug"},
		&validators.FuncValidator{
			Field:   t.Name,
			Name:    "Name",
			Message: "The name '%s' is already taken",
			Fn: func() bool {
				var b bool
				q := tx.Where("name = ?", t.Name)
				if t.ID != uuid.Nil {
					q = q.Where("id != ?", t.ID)
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

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *Team) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *Team) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

//Should remove excess and add new resources
func (t *Team) UpdateResources(tx *pop.Connection, resources []string) {
	//Remove existing that have been removed from the list
	args := []interface{}{t.ID}
	for _, resource := range resources {
		args = append(args, resource)
	}

	var sql string
	if len(resources) > 0 {
		sql = "delete from team_resources where team_id = ? and resource_id not in (?" + strings.Repeat(",?", len(resources)-1) + ")"
	} else {
		sql = "delete from team_resources where team_id = ?"
	}
	log.Printf(sql)
	del_count, del_err := tx.RawQuery(sql, args...).ExecWithCount()
	log.Printf("Delete count: %d err: %v\n", del_count, del_err)

	//Get existing IDs
	current_ids := []string{}

	if err := tx.RawQuery("select resource_id from team_resources where team_id = ?", t.ID).All(&current_ids); err != nil {
		log.Printf("Couldn't get existing ids: %v\n", err)
	}

	//Add new links
	for _, id := range resources {
		if i := sort.SearchStrings(current_ids, id); i >= len(current_ids) || current_ids[i] != id {
			rid, uerr := uuid.FromString(id)
			if uerr != nil {
				log.Printf("failed to convert str id(%s) to uuid, err: %v\n", id, uerr)
				continue
			}
			tr := TeamResource{}
			tr.TeamID = t.ID
			tr.ResourceID = rid
			verrs, err := tx.ValidateAndCreate(&tr)
			if err != nil {
				log.Printf("Error creating tr for %s <-> %s with err: %v and verrs: %v\n", t.ID, id, err, verrs)
			}
			if verrs.HasAny() {
				log.Printf("Error creating tr for %s <-> %s with verrs: %v\n", t.ID, id, verrs)
			}
		}
	}
}

// UpdateAdmins Should remove excess and add new admins
func (t *Team) UpdateAdmins(tx *pop.Connection, admins []string) {
	//Remove existing that have been removed from the list
	args := []interface{}{t.ID}
	for _, admin := range admins {
		args = append(args, admin)
	}

	var sql string
	if len(admins) > 0 {
		sql = "delete from team_admins where team_id = ? and user_id not in (?" + strings.Repeat(",?", len(admins)-1) + ")"
	} else {
		sql = "delete from team_admins where team_id = ?"
	}
	log.Printf(sql)
	del_count, del_err := tx.RawQuery(sql, args...).ExecWithCount()
	log.Printf("Delete count: %d err: %v\n", del_count, del_err)

	//Get existing IDs
	current_ids := []string{}

	if err := tx.RawQuery("select user_id from team_admins where team_id = ?", t.ID).All(&current_ids); err != nil {
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
			ta := TeamAdmin{}
			ta.TeamId = t.ID
			ta.UserId = rid
			verrs, err := tx.ValidateAndCreate(&ta)
			if err != nil {
				log.Printf("Error creating ta for %s <-> %s with err: %v and verrs: %v\n", t.ID, id, err, verrs)
			}
			if verrs.HasAny() {
				log.Printf("Error creating ta for %s <-> %s with verrs: %v\n", t.ID, id, verrs)
			}
		}
	}
}
