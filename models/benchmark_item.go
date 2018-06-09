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

type BenchmarkItem struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Name        string    `json:"name" db:"name"`
	BenchmarkID uuid.UUID `json:"benchmark_id" db:"benchmark_id"`
	Benchmark   Benchmark `belongs_to:"benchmark"`
	Resources   Resources `many_to_many:"benchmark_item_resources"`
}

// String is not required by pop and may be deleted
func (b BenchmarkItem) String() string {
	jb, _ := json.Marshal(b)
	return string(jb)
}

// BenchmarkItems is not required by pop and may be deleted
type BenchmarkItems []BenchmarkItem

// String is not required by pop and may be deleted
func (b BenchmarkItems) String() string {
	jb, _ := json.Marshal(b)
	return string(jb)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (b *BenchmarkItem) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: b.Name, Name: "Name"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (b *BenchmarkItem) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (b *BenchmarkItem) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

//Should remove excess and add new resources
func (t *BenchmarkItem) UpdateResources(tx *pop.Connection, resources []string) {
	//Remove existing that have been removed from the list
	args := []interface{}{t.ID}
	for _, resource := range resources {
		args = append(args, resource)
	}

	var sql string
	if len(resources) > 0 {
		sql = "delete from benchmark_item_resources where benchmark_item_id = ? and resource_id not in (?" + strings.Repeat(",?", len(resources)-1) + ")"
	} else {
		sql = "delete from benchmark_item_resources where benchmark_item_id = ?"
	}
	log.Printf(sql)
	del_count, del_err := tx.RawQuery(sql, args...).ExecWithCount()
	log.Printf("Delete count: %d err: %v\n", del_count, del_err)

	//Get existing IDs
	current_ids := []string{}

	if err := tx.RawQuery("select resource_id from benchmark_item_resources where benchmark_item_id = ?", t.ID).All(&current_ids); err != nil {
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
			tr := &BenchmarkItemResource{}
			tr.BenchmarkItemID = t.ID
			tr.ResourceID = rid
			err := tx.Create(tr)
			if err != nil {
				log.Printf("Error creating tr for %s <-> %s with err: %v\n", t.ID, id, err)
			}
		}
	}
}
