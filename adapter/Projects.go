package adapter

import (
	"database/sql"

	"github.com/google/uuid"
)

type Project struct {
	ID      uuid.UUID `json:"id,omitempty" db:"id"`
	Name    string    `json:"name" binding:"required" db:"name"`
	SpaceID uuid.UUID `json:"space_id" binding:"required" db:"space_id"`
}

func (p Project) Create() (sql.Result, error) {
	return DB.NamedExec(`INSERT INTO tests.tProjects (id,name,space_id) VALUES(:id, :name,:space_id)`, p)
}

func (p Project) Update() (sql.Result, error) {
	return DB.Exec(`UPDATE tests.tProjects SET name=$1, space_id=$2 WHERE id=$3`, p.Name, p.SpaceID, p.ID)
}

func (p Project) Delete() (sql.Result, error) {
	return DB.Exec(`DELETE FROM tests.tProjects WHERE id=$1`, p.ID)
}

func GetAllProjects() (result []Project, err error) {
	err = DB.Select(&result, "SELECT * FROM tests.tProjects")
	if err == nil && result == nil {
		return nil, sql.ErrNoRows
	}
	return result, err
}

func GetProject(space string, name string) (result Project, err error) {
	return result, DB.Get(&result, "SELECT p.* FROM tests.tProjects AS p INNER JOIN tests.tSpaces AS s ON p.space_id=s.id WHERE s.name=$1 AND p.name=$2", space, name)
}
func GetProjects(space string) (result []Project, err error) {
	err = DB.Select(&result, "SELECT p.* FROM tests.tProjects AS p INNER JOIN tests.tSpaces AS s ON p.space_id=s.id WHERE s.name=$1", space)
	if err == nil && result == nil {
		return nil, sql.ErrNoRows
	}
	return result, err
}

func GetProjectID(spaceID uuid.UUID, name string) (id uuid.UUID, err error) {
	return id, DB.Get(&id, "SELECT id FROM tests.tProjects WHERE space_id=$1 and name=$2", spaceID, name)
}

func GetProjectByID(id uuid.UUID) (result Project, err error) {
	return result, DB.Get(&result, "SELECT * FROM tests.tProjects WHERE id=$1", id)
}
