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
	CREATE TABLE IF NOT EXISTS tRuns (
		id SERIAL PRIMARY key,
		starttime TIMESTAMP NOT NULL,
		endtime  TIMESTAMP NOT NULL,
		data jsonb NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idxstart ON tRuns (starttime);
	CREATE INDEX IF NOT EXISTS idxend ON tRuns (endtime);
	CREATE INDEX IF NOT EXISTS idxgin ON tRuns USING GIN (data);
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
