package adapter

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Run struct {
	ID        int       `json:"id" db:"id"`
	StartTime time.Time `json:"starttime" db:"starttime"`
	EndTime   time.Time `json:"endtime" db:"endtime"`
	Data      Data      `json:"data" db:"data"`
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

func (r Run) New() (err error) {
	jsonb, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}
	_, err = DB.Exec("INSERT INTO tRuns (id, starttime, endtime, data) VALUES($1, $2, $3, $4)", r.ID, r.StartTime, r.EndTime, string(jsonb))
	return err
}

func GetAllRuns() (result []Run, err error) {
	return result, DB.Select(&result, "SELECT * FROM tRuns")
}

func GetLastsRuns(count int) (result []Run, err error) {
	return result, DB.Select(&result, "SELECT * FROM tRuns order by id desc limit $1", count)
}
func GetRangeRuns(start int, end int) (result []Run, err error) {
	return result, DB.Select(&result, "SELECT * FROM tRuns")
}

func GetTableHTML(runs []Run) (string, error) {
	return getTable(runs)
}
