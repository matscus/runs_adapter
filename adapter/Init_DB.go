package adapter

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

var (
	DB     *sqlx.DB
	scheme = `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	CREATE EXTENSION IF NOT EXISTS "pg_stat_statements";
	CREATE SCHEMA IF NOT EXISTS tests;
	CREATE TABLE IF NOT EXISTS tests.tSpaces (id UUID, name VARCHAR(45) UNIQUE,login VARCHAR(45),password VARCHAR(45), PRIMARY KEY (id));
	CREATE INDEX IF NOT EXISTS idx_tSpaces_name ON tests.tSpaces(name);
	INSERT INTO tests.tSpaces (id,name)values(uuid_generate_v4(), 'detmir') ON CONFLICT DO NOTHING;
	CREATE TABLE IF NOT EXISTS tests.tProjects (id UUID, name VARCHAR(45) UNIQUE, space_id UUID REFERENCES tests.tSpaces(id) ON DELETE CASCADE, PRIMARY KEY (id));
	CREATE INDEX IF NOT EXISTS idx_tProjects_name ON tests.tProjects(name);
	CREATE INDEX IF NOT EXISTS idx_tProjects_space_id ON tests.tProjects(space_id);
	INSERT INTO tests.tprojects (id,name,space_id)VALUES(uuid_generate_v4(), 'global',(SELECT id FROM tests.tspaces WHERE name = 'detmir')) ON CONFLICT DO NOTHING;
	CREATE TABLE IF NOT EXISTS tests.tReleases (id UUID, name VARCHAR(45), release_date_time TIMESTAMP,project_id UUID REFERENCES tests.tProjects(id) ON DELETE CASCADE, PRIMARY KEY (id));
	CREATE INDEX IF NOT EXISTS idx_tReleases_name ON tests.tReleases(name);
	CREATE INDEX IF NOT EXISTS idx_tReleases_project_id ON tests.tReleases(project_id);
	CREATE TABLE IF NOT EXISTS tests.tTestTypes (id UUID, name VARCHAR(80) UNIQUE NOT NULL,project_id UUID REFERENCES tests.tProjects(id) ON DELETE CASCADE, PRIMARY KEY (id));
	CREATE INDEX IF NOT EXISTS idx_tTestTypes_name ON tests.tTestTypes(name);
	CREATE INDEX IF NOT EXISTS idx_tTestTypes_project_id ON tests.tTestTypes(project_id); 
	INSERT INTO tests.ttesttypes(id, "name", project_id)VALUES(uuid_generate_v4(), 'Объемное тестирование', (select id from tests.tprojects where name='global'))ON CONFLICT DO NOTHING;
	INSERT INTO tests.ttesttypes(id, "name", project_id)VALUES(uuid_generate_v4(), 'Подтверждение максимальной производительности', (select id from tests.tprojects where name='global'))ON CONFLICT DO NOTHING;
	INSERT INTO tests.ttesttypes(id, "name", project_id)VALUES(uuid_generate_v4(), 'Тест надежности', (select id from tests.tprojects where name='global'))ON CONFLICT DO NOTHING;
	INSERT INTO tests.ttesttypes(id, "name", project_id)VALUES(uuid_generate_v4(), 'Стресс тест', (select id from tests.tprojects where name='global'))ON CONFLICT DO NOTHING;
	INSERT INTO tests.ttesttypes(id, "name", project_id)VALUES(uuid_generate_v4(), 'Тестирование интеграционного взаимодействия', (select id from tests.tprojects where name='global'))ON CONFLICT DO NOTHING;
	INSERT INTO tests.ttesttypes(id, "name", project_id)VALUES(uuid_generate_v4(), 'Тестирование внештатных ситуаций', (select id from tests.tprojects where name='global'))ON CONFLICT DO NOTHING;
	INSERT INTO tests.ttesttypes(id, "name", project_id)VALUES(uuid_generate_v4(), 'Объемное тестирование', (select id from tests.tprojects where name='global'))ON CONFLICT DO NOTHING;
	CREATE TABLE IF NOT EXISTS tests.tVersions (id UUID, name VARCHAR(45), version_date_time TIMESTAMP,release_id UUID REFERENCES tests.tReleases(id) ON DELETE CASCADE,  PRIMARY KEY (id));
	CREATE INDEX IF NOT EXISTS idx_tVersions_name ON tests.tVersions(name);
	CREATE INDEX IF NOT EXISTS idx_tVersions_release_id ON tests.tVersions(release_id); 
	CREATE TABLE IF NOT EXISTS tests.tRuns (id UUID, run_id INT UNIQUE NOT NULL, space_id UUID REFERENCES tests.tSpaces(id) ON DELETE CASCADE,project_id UUID REFERENCES tests.tProjects(id) ON DELETE CASCADE, release_id UUID REFERENCES tests.tReleases(id) ON DELETE CASCADE, test_type_id UUID REFERENCES tests.tTestTypes(id) ON DELETE CASCADE, version_id UUID REFERENCES tests.tVersions(id) ON DELETE CASCADE, project_page_id INT, release_page_id INT,version_page_id INT, testtype_page_id INT, run_page_id INT, start_time TIMESTAMP NOT NULL, end_time TIMESTAMP,data JSONB, dashboards_uids text[], PRIMARY KEY (id));
	CREATE INDEX IF NOT EXISTS idx_tRuns_run_id ON tests.tRuns(run_id);
	CREATE TABLE IF NOT EXISTS tests.tProfiles (id UUID,  space_id UUID REFERENCES tests.tSpaces(id) ON DELETE CASCADE, project_id UUID REFERENCES tests.tProjects(id) ON DELETE CASCADE, release_id UUID REFERENCES tests.tReleases(id) ON DELETE CASCADE, version_id UUID REFERENCES tests.tVersions(id) ON DELETE CASCADE, test_type_id UUID REFERENCES tests.tTestTypes(id) ON DELETE CASCADE, scenario_name VARCHAR(80),tps NUMERIC,sla  NUMERIC,rump_up_time INT,rump_up_steps_count INT,test_duration INT,replicas INT, cpu NUMERIC,  memory INT)
	CREATE TABLE IF NOT EXISTS tests.tLoadType (id UUID, name VARCHAR(45) UNIQUE);
	INSERT INTO tests.tLoadType (id,name)values(uuid_generate_v4(), 'steps') ON CONFLICT DO NOTHING;
	INSERT INTO tests.tLoadType (id,name)values(uuid_generate_v4(), 'lenear') ON CONFLICT DO NOTHING;
`
)

func InitDB(connStr string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("InitDB must exec panic recover ", err)
		}
	}()
	DB, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		return errors.New("Database connection error " + err.Error())
	}
	DB.DB.SetMaxOpenConns(5)
	DB.DB.SetMaxIdleConns(1)
	go func() {
		for {
			err := DB.Ping()
			if err != nil {
				log.Error("Database ping error ", err)
			}
			time.Sleep(10 * time.Second)
		}
	}()
	DB.MustExec(scheme)
	return nil
}
