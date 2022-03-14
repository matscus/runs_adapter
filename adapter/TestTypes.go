package adapter

import (
	"database/sql"

	"github.com/google/uuid"
)

type TestType struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" binding:"required" db:"name"`
	ProjectID uuid.UUID `json:"project_id" binding:"required" db:"project_id"`
}

func (t TestType) Create() (sql.Result, error) {
	return DB.NamedExec(`INSERT INTO tests.tTestTypes (id,name,project_id) VALUES(:id, :name,:project_id)`, t)
}

func (t TestType) Update() (sql.Result, error) {
	return DB.Exec(`UPDATE tests.tTestTypes SET name=$1,project_id=$2 WHERE id=$3`, t.Name, t.ProjectID, t.ID)
}

func (t TestType) Delete() (sql.Result, error) {
	return DB.Exec(`DELETE FROM tests.tTestTypes WHERE id=$1`, t.ID)
}

func GetAllTestTypes() (result []TestType, err error) {
	err = DB.Select(&result, "SELECT * FROM tests.tTestTypes")
	if err == nil && result == nil {
		return nil, sql.ErrNoRows
	}
	return result, err
}

func GetTestType(space string, project string, testType string) (result TestType, err error) {
	return result, DB.Get(&result, "SELECT t.* FROM tests.tTestTypes AS t INNER JOIN tests.tProjects AS p ON t.project_id = p.id  INNER JOIN tests.tSpaces AS s ON p.space_id = s.id WHERE  s.name=$1 AND p.name=$2 AND t.name =$3", space, project, testType)
}

func GetTestTypes(space string, project string) (result []TestType, err error) {
	err = DB.Select(&result, "SELECT t.* FROM tests.tTestTypes AS t INNER JOIN tests.tProjects AS p ON t.project_id = p.id  INNER JOIN tests.tSpaces AS s ON p.space_id = s.id WHERE  s.name=$1 AND p.name=$2", space, project)
	if err == nil && result == nil {
		return nil, sql.ErrNoRows
	}
	return result, err
}

func GetTestTypeID(projectID uuid.UUID, name string) (id uuid.UUID, err error) {
	return id, DB.Get(&id, "SELECT id FROM tests.tTestTypes WHERE project_id=$1 and name=$2", projectID, name)
}

func GetTestTypeByID(id uuid.UUID) (result TestType, err error) {
	return result, DB.Get(&result, "SELECT * FROM tests.tTestTypes WHERE id=$1", id)
}
