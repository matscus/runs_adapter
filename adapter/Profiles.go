package adapter

import (
	"database/sql"

	"github.com/google/uuid"
)

type Profile struct {
	ID               uuid.UUID `json:"id,omitempty" db:"id"`
	ScenarioName     string    `json:"scenario_name" binding:"required" db:"scenario_name"`
	SpaceID          uuid.UUID `json:"space_id,omitempty" db:"space_id"`
	ProjectID        uuid.UUID `json:"project_id,omitempty" db:"project_id"`
	ReleaseID        uuid.UUID `json:"release_id,omitempty" db:"release_id"`
	VersionID        uuid.UUID `json:"version_id,omitempty" db:"version_id"`
	TestTypeID       uuid.UUID `json:"test_type_id,omitempty" db:"test_type_id"`
	SpaceName        string    `json:"space_name,omitempty" binding:"required" db:"space_name"`
	ProjectName      string    `json:"project_name,omitempty" binding:"required" db:"project_name"`
	ReleaseName      string    `json:"release_name,omitempty" binding:"required" db:"release_name"`
	TestTypeName     string    `json:"testtype_name,omitempty" binding:"required" db:"testtype_name"`
	VersionName      string    `json:"version_name,omitempty" binding:"required" db:"version_name"`
	TPS              float64   `json:"tps,omitempty" binding:"required" db:"tps"`
	SLA              float64   `json:"sla,omitempty" binding:"required" db:"sla"`
	RumpUpTime       int       `json:"rump_up_time,omitempty" binding:"required" db:"rump_up_time"`
	RumpUpStepsCount int       `json:"rump_up_steps_count,omitempty" binding:"required" db:"rump_up_steps_count"`
	TestDuration     int       `json:"test_duration,omitempty" binding:"required" db:"test_duration"`
	Replicas         int       `json:"replicas,omitempty" binding:"required" db:"replicas"`
	CPU              float64   `json:"cpu,omitempty" binding:"required" db:"cpu"`
	Memory           int       `json:"memory,omitempty" binding:"required" db:"memory"`
}

func (p Profile) Create() (sql.Result, error) {
	return DB.NamedExec(`INSERT INTO tests.tProfiles (id,scenario_name,space_id,project_id,release_id,version_id,test_type_id,tps,sla,rump_up_time,rump_up_steps_count,test_duration,replicas,cpu,memory) VALUES(:id,:scenario_name,:space_id,:project_id,:release_id,:version_id,:test_type_id,:tps,:sla,:rump_up_time,:rump_up_steps_count,:test_duration,:replicas,:cpu,:memory)`, p)
}

func (p Profile) Update() (sql.Result, error) {
	return DB.Exec(`UPDATE tests.tProfiles SET scenario_name=$1,space_id=$2,project_id=$3,release_id=$4,version_id=$5,test_type_id=$6,tps=$7,sla=$8,rump_up_time=$9,rump_up_steps_count=$10,test_duration=$11,replicas=$12,cpu=$13,memory=$14 WHERE id=$15`, p.ScenarioName, p.SpaceID, p.ProjectID, p.ReleaseID, p.VersionID, p.TestTypeID, p.TPS, p.SLA, p.RumpUpTime, p.RumpUpStepsCount, p.TestDuration, p.Replicas, p.CPU, p.Memory, p.ID)
}

func (p Profile) Delete() (sql.Result, error) {
	return DB.Exec(`DELETE FROM tests.tProfiles WHERE id=$1`, p.ID)
}

func GetAllProfiles() (result []Profile, err error) {
	err = DB.Select(&result, "SELECT profile.id,profile.scenario_name,profile.space_id,s.name AS space_name,profile.project_id,p.name AS project_name, profile.release_id,r.name AS release_name,profile.version_id,v.name AS version_name,profile.test_type_id,t.name  AS testtype_name,profile.tps,profile.sla,profile.rump_up_time,profile.rump_up_steps_count,profile.test_duration,profile.replicas,profile.cpu,profile.memory from tests.tProfiles AS profile INNER JOIN tests.tSpaces AS s ON profile.space_id = s.id INNER JOIN tests.tProjects AS p ON profile.project_id = p.id INNER JOIN tests.tReleases AS r ON profile.release_id = r.id INNER JOIN tests.tVersions AS v ON profile.version_id = v.id INNER JOIN tests.tTestTypes AS t ON profile.test_type_id = t.id")
	if err == nil && result == nil {
		return nil, sql.ErrNoRows
	}
	return result, err
}

func GetProfile(space string, project string, release string, version string, testType string, scenario_name string) (result Profile, err error) {
	return result, DB.Get(&result, "SELECT profile.id,profile.scenario_name,profile.space_id,s.name AS space_name,profile.project_id,p.name AS project_name, profile.release_id,r.name AS release_name,profile.version_id,v.name AS version_name,profile.test_type_id,t.name  AS testtype_name,profile.tps,profile.sla,profile.rump_up_time,profile.rump_up_steps_count,profile.test_duration,profile.replicas,profile.cpu,profile.memory from tests.tProfiles AS profile INNER JOIN tests.tSpaces AS s ON profile.space_id = s.id INNER JOIN tests.tProjects AS p ON profile.project_id = p.id INNER JOIN tests.tReleases AS r ON profile.release_id = r.id INNER JOIN tests.tVersions AS v ON profile.version_id = v.id INNER JOIN tests.tTestTypes AS t ON profile.test_type_id = t.id WHERE s.name =$1 AND p.name=$2 AND r.name=$3 AND v.name =$4 AND t.name =$5 AND profile.scenario_name=$6", space, project, release, version, testType, scenario_name)
}

func GetProfiles(space string, project string, release string, version string, testType string) (result []Profile, err error) {
	err = DB.Select(&result, "SELECT profile.id,profile.scenario_name,profile.space_id,s.name AS space_name,profile.project_id,p.name AS project_name, profile.release_id,r.name AS release_name,profile.version_id,v.name AS version_name,profile.test_type_id,t.name  AS testtype_name,profile.tps,profile.sla,profile.rump_up_time,profile.rump_up_steps_count,profile.test_duration,profile.replicas,profile.cpu,profile.memory from tests.tProfiles AS profile INNER JOIN tests.tSpaces AS s ON profile.space_id = s.id INNER JOIN tests.tProjects AS p ON profile.project_id = p.id INNER JOIN tests.tReleases AS r ON profile.release_id = r.id INNER JOIN tests.tVersions AS v ON profile.version_id = v.id INNER JOIN tests.tTestTypes AS t ON profile.test_type_id = t.id WHERE s.name =$1 AND p.name=$2 AND r.name=$3 AND v.name =$4 AND t.name =$5 ", space, project, release, version, testType)
	if err == nil && result == nil {
		return nil, sql.ErrNoRows
	}
	return result, err
}

func GetProfileByID(id uuid.UUID) (result Profile, err error) {
	return result, DB.Get(&result, "SELECT profile.id,profile.scenario_name,profile.space_id,s.name AS space_name,profile.project_id,p.name AS project_name, profile.release_id,r.name AS release_name,profile.version_id,v.name AS version_name,profile.test_type_id,t.name  AS testtype_name,profile.tps,profile.sla,profile.rump_up_time,profile.rump_up_steps_count,profile.test_duration,profile.replicas,profile.cpu,profile.memory from tests.tProfiles AS profile INNER JOIN tests.tSpaces AS s ON profile.space_id = s.id INNER JOIN tests.tProjects AS p ON profile.project_id = p.id INNER JOIN tests.tReleases AS r ON profile.release_id = r.id INNER JOIN tests.tVersions AS v ON profile.version_id = v.id INNER JOIN tests.tTestTypes AS t ON profile.test_type_id = t.id WHERE profile.id=$1", id)
}
