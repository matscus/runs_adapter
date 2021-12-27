package adapter

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Release struct {
	ID              uuid.UUID `json:"id" db:"id"`
	Name            string    `json:"name" binding:"required" db:"name"`
	ReleaseDateTime time.Time `json:"release_date_time" db:"release_date_time"`
	ProjectID       uuid.UUID `json:"project_id" binding:"required" db:"project_id"`
}

func (r Release) Create() (sql.Result, error) {
	return DB.NamedExec(`INSERT INTO tests.tReleases (id,name,release_date_time,project_id) VALUES(:id, :name,:release_date_time,:project_id)`, r)
}

func (r Release) Update() (sql.Result, error) {
	return DB.Exec(`UPDATE tests.tReleases SET name=$1, release_date_time=$2,project_id=$3 WHERE id=$4`, r.Name, r.ReleaseDateTime, r.ProjectID, r.ID)
}

func (r Release) Delete() (sql.Result, error) {
	return DB.Exec(`DELETE FROM tests.tReleases WHERE id=$1`, r.ID)
}

func GetAllReleases() (result []Release, err error) {
	err = DB.Select(&result, "SELECT * FROM tests.tReleases")
	if err == nil && result == nil {
		return nil, sql.ErrNoRows
	}
	return result, err
}

func GetRelease(space string, project string, release string) (result Release, err error) {
	return result, DB.Get(&result, "SELECT r.* FROM tests.tReleases AS r  INNER JOIN tests.tProjects AS p ON r.project_id = p.id  INNER JOIN tests.tSpaces AS s ON p.space_id = s.id WHERE s.name=$1 AND p.name=$2 ABD r.name =$3", space, project, release)
}
func GetReleaseID(projectID uuid.UUID, release string) (result uuid.UUID, err error) {
	return result, DB.Get(&result, "SELECT r.id FROM tests.tReleases AS r  WHERE r.project_id=$1 AND r.name=$2", projectID, release)
}

func GetReleaseByID(id uuid.UUID) (result Release, err error) {
	return result, DB.Get(&result, "SELECT * FROM tests.tReleases WHERE id=$1", id)
}
