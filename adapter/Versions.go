package adapter

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Version struct {
	ID              uuid.UUID `json:"id" db:"id"`
	Name            string    `json:"name" binding:"required" db:"name"`
	VersionDateTime time.Time `json:"version_date_time" db:"version_date_time"`
	ReleaseID       uuid.UUID `json:"release_id" binding:"required" db:"release_id"`
}

func (v Version) Create() (sql.Result, error) {
	return DB.NamedExec(`INSERT INTO tests.tVersions(id,name,version_date_time,release_id) VALUES(:id,:name,:version_date_time,:release_id)`, v)
}

func (v Version) Update() (sql.Result, error) {
	return DB.Exec(`UPDATE tests.tVersions SET name=$1, version_date_time=$2,release_id=$3 WHERE id=$4`, v.Name, v.VersionDateTime, v.ReleaseID, v.ID)
}

func (v Version) Delete() (sql.Result, error) {
	return DB.Exec(`DELETE FROM tests.tVersions WHERE id=$1`, v.ID)
}

func GetAllVersions() (result []Version, err error) {
	err = DB.Select(&result, "SELECT * FROM tests.tVersions")
	if err == nil && result == nil {
		return nil, sql.ErrNoRows
	}
	return result, err
}

func GetVersion(space string, project string, release string, version string) (result Version, err error) {
	return result, DB.Get(&result, "SELECT v.* FROM tests.tVersions AS v INNER JOIN tests.tReleases AS r ON v.release_id = r.id INNER JOIN tests.tProjects AS p ON r.project_id = p.id INNER JOIN tests.tSpaces AS s ON p.space_id = s.id WHERE  s.name=$1 AND r.name =$2 AND p.name=$3 AND  v.name=$4", space, project, release, version)
}
func GetVersions(space string, project string, release string) (result []Version, err error) {
	err = DB.Select(&result, "SELECT v.* FROM tests.tVersions AS v INNER JOIN tests.tReleases AS r ON v.release_id = r.id INNER JOIN tests.tProjects AS p ON r.project_id = p.id INNER JOIN tests.tSpaces AS s ON p.space_id = s.id WHERE  s.name=$1 AND p.name=$2 AND r.name=$3 ORDER BY version_date_time DESC", space, project, release)
	if err == nil && result == nil {
		return nil, sql.ErrNoRows
	}
	return result, err
}

func GetVersionID(releaseID uuid.UUID, name string) (id uuid.UUID, err error) {
	return id, DB.Get(&id, "SELECT id FROM tests.tVersions WHERE release_id=$1 and name=$2", releaseID, name)
}

func GetVersionByID(id uuid.UUID) (result Version, err error) {
	return result, DB.Get(&result, "SELECT * FROM tests.tVersions WHERE id=$1", id)
}
