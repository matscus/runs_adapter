package adapter

import (
	"database/sql"

	"github.com/google/uuid"
)

type LoadType struct {
	ID   uuid.UUID `json:"id,omitempty" db:"id"`
	Name string    `json:"name" binding:"required" db:"scenario_name"`
}

func (t LoadType) Create() (sql.Result, error) {
	return DB.NamedExec(`INSERT INTO tests.tLoadType (id,name)`, t)
}

func (t LoadType) Update() (sql.Result, error) {
	return DB.Exec(`UPDATE tests.tLoadType SET name=$1 WHERE id=$2`, t.Name, t.ID)
}

func (t LoadType) Delete() (sql.Result, error) {
	return DB.Exec(`DELETE FROM tests.tLoadType WHERE id=$1`, t.ID)
}

func GetLoadTypes() (result []LoadType, err error) {
	err = DB.Select(&result, "SELECT id, name from tests.tLoadType")
	if err == nil && result == nil {
		return nil, sql.ErrNoRows
	}
	return result, err
}

func GetLoadTypeByID(id uuid.UUID) (result LoadType, err error) {
	return result, DB.Get(&result, "SELECT id, name WHERE id=$1", id)
}
