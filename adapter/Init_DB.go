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
	CREATE TABLE IF NOT EXISTS tests.tSpaces (id UUID UNIQUE NOT NULL, name VARCHAR(45) UNIQUE NOT NULL, PRIMARY KEY (id));
	CREATE INDEX IF NOT EXISTS idx_tSpaces_name ON tests.tSpaces(name);
	INSERT INTO tests.tSpaces (id,name)values(uuid_generate_v4(), 'detmir') ON CONFLICT DO NOTHING;
	CREATE TABLE IF NOT EXISTS tests.tProjects (id UUID UNIQUE NOT NULL, name VARCHAR(45) UNIQUE NOT NULL, space_id UUID REFERENCES tests.tSpaces(id) ON DELETE CASCADE, PRIMARY KEY (id));
	CREATE INDEX IF NOT EXISTS idx_tProjects_name ON tests.tProjects(name);
	CREATE INDEX IF NOT EXISTS idx_tProjects_space_id ON tests.tProjects(space_id);
	INSERT INTO tests.tprojects (id,name,space_id)VALUES(uuid_generate_v4(), 'global',(SELECT id FROM tests.tspaces WHERE name = 'detmir')) ON CONFLICT DO NOTHING;
	CREATE TABLE IF NOT EXISTS tests.tReleases (id UUID NOT NULL, name VARCHAR(45), release_date_time TIMESTAMP,project_id UUID REFERENCES tests.tProjects(id)  ON DELETE CASCADE,UNIQUE (name,project_id), PRIMARY KEY (id));
	CREATE INDEX IF NOT EXISTS idx_tReleases_name ON tests.tReleases(name);
	CREATE INDEX IF NOT EXISTS idx_tReleases_project_id ON tests.tReleases(project_id);
	INSERT INTO tests.treleases(id, "name", release_date_time, project_id)VALUES(uuid_generate_v4(), 'global', now(), (select id from tests.tProjects where space_id=(select id from tests.tSpaces where name='detmir')))ON CONFLICT DO NOTHING;
	CREATE TABLE IF NOT EXISTS tests.tTestTypes (id UUID UNIQUE NOT NULL, name VARCHAR(80) UNIQUE NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_tTestTypes_name ON tests.tTestTypes(name);
	INSERT INTO tests.tTestTypes(id, name) VALUES (uuid_generate_v4(), 'Поиск максимальной производительности')ON CONFLICT DO NOTHING;
	INSERT INTO tests.tTestTypes(id, name) VALUES (uuid_generate_v4(), 'Подтверждение максимальной производительности')ON CONFLICT DO NOTHING;
	INSERT INTO tests.tTestTypes(id, name) VALUES (uuid_generate_v4(), 'Тестирование надежности')ON CONFLICT DO NOTHING;
	INSERT INTO tests.tTestTypes(id, name) VALUES (uuid_generate_v4(), 'Стресс тестирование')ON CONFLICT DO NOTHING;
	INSERT INTO tests.tTestTypes(id, name) VALUES (uuid_generate_v4(), 'Тестирование интеграционного взаимодействия')ON CONFLICT DO NOTHING;
	INSERT INTO tests.tTestTypes(id, name) VALUES (uuid_generate_v4(), 'Тестирование внештатных ситуаций')ON CONFLICT DO NOTHING;
	INSERT INTO tests.tTestTypes(id, name) VALUES (uuid_generate_v4(), 'Объемное тестирование')ON CONFLICT DO NOTHING;
	INSERT INTO tests.tTestTypes(id, name) VALUES (uuid_generate_v4(), 'Тестирование масштабируемости')ON CONFLICT DO NOTHING;
	INSERT INTO tests.tTestTypes(id, name) VALUES (uuid_generate_v4(), 'Узконаправленные кейсы')ON CONFLICT DO NOTHING;
	CREATE TABLE IF NOT EXISTS tests.tLoadTypes (id UUID UNIQUE NOT NULL, name VARCHAR(45) UNIQUE);
	INSERT INTO tests.tLoadTypes (id,name)values(uuid_generate_v4(),'steps') ON CONFLICT DO NOTHING;
	INSERT INTO tests.tLoadTypes (id,name)values(uuid_generate_v4(),'linear') ON CONFLICT DO NOTHING;
	CREATE TABLE IF NOT EXISTS tests.tRuns (id UUID, run_id INT NOT NULL, space_id UUID NOT NULL REFERENCES  tests.tSpaces(id) ON DELETE CASCADE, project_id UUID NOT NULL REFERENCES tests.tProjects(id) ON DELETE CASCADE, release_id UUID NOT NULL REFERENCES tests.tReleases(id) ON DELETE CASCADE, test_type_id UUID NOT NULL REFERENCES tests.tTestTypes(id) ON DELETE CASCADE,load_type_id UUID NOT NULL REFERENCES tests.tLoadTypes(id) ON DELETE CASCADE, start_time TIMESTAMP NOT NULL, end_time TIMESTAMP, data JSONB, PRIMARY KEY (id));
	CREATE INDEX IF NOT EXISTS idx_tRuns_run_id ON tests.tRuns(run_id);
	CREATE TABLE IF NOT EXISTS tests.tProfiles (id UUID,  space_id UUID REFERENCES tests.tSpaces(id) ON DELETE CASCADE, project_id UUID REFERENCES tests.tProjects(id) ON DELETE CASCADE, release_id UUID REFERENCES tests.tReleases(id) ON DELETE CASCADE, test_type_id UUID REFERENCES tests.tTestTypes(id) ON DELETE CASCADE, load_type_id UUID REFERENCES tests.tLoadTypes(id) ON DELETE CASCADE, scenario_name VARCHAR(80),tps NUMERIC,sla NUMERIC,rump_up_time INT,rump_up_steps_count INT,test_duration INT,replicas INT, cpu NUMERIC,  memory INT)
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
