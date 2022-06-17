package adapter

import (
	"database/sql"

	"github.com/google/uuid"
)

type LoadType struct {
	ID   uuid.UUID `json:"id,omitempty" db:"id"`
	Name string    `json:"name" binding:"required" db:"name"`
}

func (t LoadType) Create() (sql.Result, error) {
	return DB.NamedExec(`INSERT INTO tests.tLoadTypes (id,name)`, t)
}

func (t LoadType) Update() (sql.Result, error) {
	return DB.Exec(`UPDATE tests.tLoadTypes SET name=$1 WHERE id=$2`, t.Name, t.ID)
}

func (t LoadType) Delete() (sql.Result, error) {
	return DB.Exec(`DELETE FROM tests.tLoadTypes WHERE id=$1`, t.ID)
}

func GetLoadTypes() (result []LoadType, err error) {
	err = DB.Select(&result, "SELECT id, name FROM tests.tLoadTypes")
	if err == nil && result == nil {
		return nil, sql.ErrNoRows
	}
	return result, err
}

func GetDefaultLoadTypeID() (id uuid.UUID, err error) {
	return id, DB.Get(&id, "SELECT id FROM tests.tLoadTypes WHERE name=$1", "steps")
}
func GetLoadTypeID(name string) (id uuid.UUID, err error) {
	return id, DB.Get(&id, "SELECT id FROM tests.tLoadTypes WHERE name=$1", name)
}
func GetLoadTypeByID(id uuid.UUID) (result LoadType, err error) {
	return result, DB.Get(&result, "SELECT id, name FROM tests.tLoadTypes WHERE id=$1", id)
}
