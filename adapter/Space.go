package adapter

import (
	"database/sql"

	"github.com/google/uuid"
)

type Space struct {
	ID   uuid.UUID `json:"id" db:"id"`
	Name string    `json:"name" binding:"required" db:"name"`
}

func (s Space) Create() (sql.Result, error) {
	return DB.NamedExec(`INSERT INTO tests.tSpaces (id,name) VALUES(:id,:name)`, s)
}

func (s Space) Update() (sql.Result, error) {
	return DB.Exec(`UPDATE tests.tSpaces SET name=$1, login=$2 ,password=$3 WHERE id=$4`, s.Name, s.ID)
}

func (s Space) Delete() (sql.Result, error) {
	return DB.Exec(`DELETE FROM tests.tSpaces WHERE id=$1`, s.ID)
}

func GetAllSpaces() (result []Space, err error) {
	err = DB.Select(&result, "SELECT * FROM tests.tSpaces")
	if err == nil && result == nil {
		return nil, sql.ErrNoRows
	}
	return result, err
}

func GetSpace(name string) (result Space, err error) {
	return result, DB.Get(&result, "SELECT * FROM tests.tSpaces WHERE name=$1", name)
}

func GetSpaceID(name string) (id uuid.UUID, err error) {
	return id, DB.Get(&id, "SELECT id FROM tests.tSpaces WHERE name=$1", name)
}

func GetSpaceByID(id uuid.UUID) (result Space, err error) {
	return result, DB.Get(&result, "SELECT * FROM tests.tSpaces WHERE id=$1", id)
}
