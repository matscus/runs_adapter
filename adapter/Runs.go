package adapter

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Run struct {
	ID           uuid.UUID `json:"id" db:"id"`
	RunID        int       `json:"run_id" binding:"required" db:"run_id"`
	SpaceID      uuid.UUID `json:"space_id" db:"space_id"`
	ProjectID    uuid.UUID `json:"project_id" db:"project_id"`
	ReleaseID    uuid.UUID `json:"release_id" db:"release_id"`
	TestTypeID   uuid.UUID `json:"test_type_id" db:"test_type_id"`
	VersionID    uuid.UUID `json:"version_id" db:"version_id"`
	SpaceName    string    `json:"space_name" binding:"required" db:"space_name"`
	ProjectName  string    `json:"project_name" binding:"required" db:"project_name"`
	ReleaseName  string    `json:"release_name" binding:"required" db:"release_name"`
	VersionName  string    `json:"version_name" binding:"required" db:"version_name"`
	TestTypeName string    `json:"testtype_name" binding:"required" db:"testtype_name"`
	StartTime    time.Time `json:"start_time,omitempty" binding:"required" db:"start_time"`
	EndTime      time.Time `json:"end_time,omitempty" db:"end_time"`
	Data         Data      `json:"data" binding:"required" db:"data"`
}

type Data struct {
	Project     string     `json:"project,omitempty"`
	Grafanalink string     `json:"grafanalink,omitempty"`
	Description string     `json:"description,omitempty"`
	Status      string     `json:"status,omitempty"`
	Scenarios   []Scenario `json:"scenarios,omitempty"`
}

type Scenario struct {
	Name     string `json:"name,omitempty"`
	TPS      int    `json:"tps,omitempty"`
	SLA      int    `json:"sla,omitempty"`
	Duration int    `json:"duration,omitempty"`
}

func (r Run) Create() (sql.Result, error) {
	return DB.NamedExec(`INSERT INTO tests.tRuns (id,run_id,space_id,project_id,release_id,test_type_id,version_id,start_time,end_time,data) VALUES (:id,:run_id,:space_id,:project_id,:release_id,:test_type_id,:version_id,:start_time,:end_time,:data)`, r)
}

func (r Run) Update() (sql.Result, error) {
	return DB.Exec(`UPDATE tests.tRuns SET run_id=$1,space_id=$2,project_id=$3,release_id=$4,version_id=$5,test_type_id=$6,start_time=$7,end_time=$8,data=$9 WHERE id=$10`, r.RunID, r.SpaceID, r.ProjectID, r.ReleaseID, r.VersionID, r.TestTypeID, r.StartTime, r.EndTime, r.Data, r.ID)
}

func (r Run) Delete() (sql.Result, error) {
	return DB.Exec(`DELETE FROM tests.tRuns WHERE id=$1`, r.ID)
}

func GetAllRuns() (result []Run, err error) {
	err = DB.Select(&result, "SELECT test.id,test.run_id,test.space_id,s.name AS space_name, test.project_id, p.name AS project_name, test.release_id, r.name AS release_name, test.version_id, v.name AS version_name, test.test_type_id, t.name  AS testtype_name, test.start_time, test.end_time, test.data FROM tests.tRuns AS test INNER JOIN tests.tSpaces AS s ON test.space_id = s.id INNER JOIN tests.tProjects AS p ON test.project_id = p.id INNER JOIN tests.tReleases AS r ON test.release_id = r.id INNER JOIN tests.tVersions AS v ON test.version_id = v.id INNER JOIN tests.tTestTypes AS t ON test.test_type_id = t.id")
	if err == nil && result == nil {
		return nil, sql.ErrNoRows
	}
	return result, err
}

func GetLastRunID() (result int, err error) {
	return result, DB.Get(&result, "SELECT run_id FROM tests.tRuns ORDER BY run_id DESC LIMIT 1")
}

func GetRuns(schema string, project string, limit int) (result []Run, err error) {
	err = DB.Select(&result, "SELECT test.id,test.run_id,test.space_id,s.name AS space_name, test.project_id, p.name AS project_name, test.release_id, r.name AS release_name, test.version_id, v.name AS version_name, test.test_type_id, t.name  AS testtype_name, test.start_time, test.end_time, test.data FROM tests.tRuns AS test INNER JOIN tests.tSpaces AS s ON test.space_id = s.id INNER JOIN tests.tProjects AS p ON test.project_id = p.id INNER JOIN tests.tReleases AS r ON test.release_id = r.id INNER JOIN tests.tVersions AS v ON test.version_id = v.id INNER JOIN tests.tTestTypes AS t ON test.test_type_id = t.id WHERE s.name=$1 AND p.name=$2 ORDER BY run_id DESC LIMIT $3", schema, project, limit)
	if err == nil && result == nil {
		return nil, sql.ErrNoRows
	}
	return result, err
}

func GetRunByID(id uuid.UUID) (result Run, err error) {
	return result, DB.Get(&result, "SELECT test.id,test.run_id,test.space_id,s.name AS space_name, test.project_id, p.name AS project_name, test.release_id, r.name AS release_name, test.version_id, v.name AS version_name, test.test_type_id, t.name  AS testtype_name, test.start_time, test.end_time, test.data FROM tests.tRuns AS test INNER JOIN tests.tSpaces AS s ON test.space_id = s.id INNER JOIN tests.tProjects AS p ON test.project_id = p.id INNER JOIN tests.tReleases AS r ON test.release_id = r.id INNER JOIN tests.tVersions AS v ON test.version_id = v.id INNER JOIN tests.tTestTypes AS t ON test.test_type_id = t.id WHERE test.id = $1", id)
}

func GetAllRunsBySpace(schema string) (result []Run, err error) {
	err = DB.Select(&result, "SELECT test.id,test.run_id,test.space_id,s.name AS space_name, test.project_id, p.name AS project_name, test.release_id, r.name AS release_name, test.version_id, v.name AS version_name, test.test_type_id, t.name  AS testtype_name, test.start_time, test.end_time, test.data FROM tests.tRuns AS test INNER JOIN tests.tSpaces AS s ON test.space_id = s.id INNER JOIN tests.tProjects AS p ON test.project_id = p.id INNER JOIN tests.tReleases AS r ON test.release_id = r.id INNER JOIN tests.tVersions AS v ON test.version_id = v.id INNER JOIN tests.tTestTypes AS t ON test.test_type_id = t.id WHERE s.name=$1", schema)
	if err == nil && result == nil {
		return nil, sql.ErrNoRows
	}
	return result, err
}

func GetAllRunsByProject(schema string, project string) (result []Run, err error) {
	err = DB.Select(&result, "SELECT test.id,test.run_id,test.space_id,s.name AS space_name, test.project_id, p.name AS project_name, test.release_id, r.name AS release_name, test.version_id, v.name AS version_name, test.test_type_id, t.name  AS testtype_name, test.start_time, test.end_time, test.data FROM tests.tRuns AS test INNER JOIN tests.tSpaces AS s ON test.space_id = s.id INNER JOIN tests.tProjects AS p ON test.project_id = p.id INNER JOIN tests.tReleases AS r ON test.release_id = r.id INNER JOIN tests.tVersions AS v ON test.version_id = v.id INNER JOIN tests.tTestTypes AS t ON test.test_type_id = t.id WHERE s.name=$1 AND p.name=$2", schema, project)
	if err == nil && result == nil {
		return nil, sql.ErrNoRows
	}
	return result, err
}

func GetAllRunsByRelease(schema string, project string, release string) (result []Run, err error) {
	err = DB.Select(&result, "SELECT test.id,test.run_id,test.space_id,s.name AS space_name, test.project_id, p.name AS project_name, test.release_id, r.name AS release_name, test.version_id, v.name AS version_name, test.test_type_id, t.name  AS testtype_name, test.start_time, test.end_time, test.data FROM tests.tRuns AS test INNER JOIN tests.tSpaces AS s ON test.space_id = s.id INNER JOIN tests.tProjects AS p ON test.project_id = p.id INNER JOIN tests.tReleases AS r ON test.release_id = r.id INNER JOIN tests.tVersions AS v ON test.version_id = v.id INNER JOIN tests.tTestTypes AS t ON test.test_type_id = t.id WHERE s.name=$1 AND p.name=$2 AND r.name=$3", schema, project, release)
	if err == nil && result == nil {
		return nil, sql.ErrNoRows
	}
	return result, err
}

func GetAllRunsByVersion(schema string, project string, release string, version string) (result []Run, err error) {
	err = DB.Select(&result, "SELECT test.id,test.run_id,test.space_id,s.name AS space_name, test.project_id, p.name AS project_name, test.release_id, r.name AS release_name, test.version_id, v.name AS version_name, test.test_type_id, t.name  AS testtype_name, test.start_time, test.end_time, test.data FROM tests.tRuns AS test INNER JOIN tests.tSpaces AS s ON test.space_id = s.id INNER JOIN tests.tProjects AS p ON test.project_id = p.id INNER JOIN tests.tReleases AS r ON test.release_id = r.id INNER JOIN tests.tVersions AS v ON test.version_id = v.id INNER JOIN tests.tTestTypes AS t ON test.test_type_id = t.id WHERE s.name=$1 AND p.name=$2 AND r.name=$3 AND v.name=$4", schema, project, release, version)
	if err == nil && result == nil {
		return nil, sql.ErrNoRows
	}
	return result, err
}

func GetAllRunsByTestType(schema string, project string, release string, version string, testType string) (result []Run, err error) {
	err = DB.Select(&result, "SELECT test.id,test.run_id,test.space_id,s.name AS space_name, test.project_id, p.name AS project_name, test.release_id, r.name AS release_name, test.version_id, v.name AS version_name, test.test_type_id, t.name  AS testtype_name, test.start_time, test.end_time, test.data FROM tests.tRuns AS test INNER JOIN tests.tSpaces AS s ON test.space_id = s.id INNER JOIN tests.tProjects AS p ON test.project_id = p.id INNER JOIN tests.tReleases AS r ON test.release_id = r.id INNER JOIN tests.tVersions AS v ON test.version_id = v.id INNER JOIN tests.tTestTypes AS t ON test.test_type_id = t.id WHERE s.name=$1 AND p.name=$2 AND r.name=$3 AND v.name=$4 AND t.name=$5", schema, project, release, version, testType)
	if err == nil && result == nil {
		return nil, sql.ErrNoRows
	}
	return result, err
}

func SetEndTime(id uuid.UUID, endTime time.Time) (sql.Result, error) {
	return DB.Exec(`UPDATE tests.tRuns SET end_time=$1 WHERE id=$2`, endTime, id)
}

func (a Data) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Data) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

func GetTableHTML(runs []Run) (string, error) {
	return getTable(runs)
}
