package handlers

import (
	"database/sql"
	"net/http"
	"runs_adapter/adapter"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var nilUUID uuid.UUID

func Runs(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodGet:
		all := c.Query("all")
		if all != "" {
			res, err := adapter.GetAllRuns()
			if err != nil {
				CheckSQLError(c, err)
				return
			}
			c.JSON(200, gin.H{"status": "ok", "data": res})
			return
		}
		id := c.Query("id")
		if id != "" {
			uuid := uuid.MustParse(id)
			res, err := adapter.GetRunByID(uuid)
			if err != nil {
				CheckSQLError(c, err)
				return
			}
			c.JSON(200, gin.H{"status": "ok", "data": res})
			return
		}
	case http.MethodPost:
		run := adapter.Run{}
		err := c.BindJSON(&run)
		if err != nil {
			c.JSON(500, gin.H{"status": "error", "Message": err.Error()})
			return
		}
		run.ID = uuid.New()
		if run.SpaceID == nilUUID {
			run.SpaceID, err = adapter.GetSpaceID(run.SpaceName)
			if err != nil {
				if err == sql.ErrNoRows {
					run.SpaceID = uuid.New()
					_, err = adapter.Space{ID: run.SpaceID, Name: run.SpaceName}.Create()
					if err != nil {
						CheckSQLError(c, err)
						return
					}
				} else {
					CheckSQLError(c, err)
					return
				}
			}
		}
		if run.ProjectID == nilUUID {
			run.ProjectID, err = adapter.GetProjectID(run.SpaceID, run.ProjectName)
			if err != nil {
				if err == sql.ErrNoRows {
					run.ProjectID = uuid.New()
					_, err = adapter.Project{ID: run.ProjectID, Name: run.ProjectName, SpaceID: run.SpaceID}.Create()
					if err != nil {
						CheckSQLError(c, err)
						return
					}
				} else {
					CheckSQLError(c, err)
					return
				}
			}
		}
		if run.ReleaseID == nilUUID {
			run.ReleaseID, err = adapter.GetReleaseID(run.ProjectID, run.ReleaseName)
			if err != nil {
				if err == sql.ErrNoRows {
					run.ReleaseID = uuid.New()
					_, err = adapter.Release{ID: run.ReleaseID, Name: run.ReleaseName, ProjectID: run.ProjectID}.Create()
					if err != nil {
						CheckSQLError(c, err)
						return
					}
				} else {
					CheckSQLError(c, err)
					return
				}
			}
		}
		if run.VersionID == nilUUID {
			run.VersionID, err = adapter.GetVersionID(run.ReleaseID, run.VersionName)
			if err != nil {
				if err == sql.ErrNoRows {
					run.VersionID = uuid.New()
					_, err = adapter.Version{ID: run.VersionID, Name: run.VersionName, ReleaseID: run.ReleaseID}.Create()
					if err != nil {
						CheckSQLError(c, err)
						return
					}
				} else {
					CheckSQLError(c, err)
					return
				}
			}
		}
		if run.TestTypeID == nilUUID {
			run.TestTypeID, err = adapter.GetTestTypeID(run.ProjectID, run.TestTypeName)
			if err != nil {
				if err == sql.ErrNoRows {
					run.TestTypeID = uuid.New()
					_, err = adapter.TestType{ID: run.TestTypeID, Name: run.TestTypeName, ProjectID: run.ProjectID}.Create()
					if err != nil {
						CheckSQLError(c, err)
						return
					}
				} else {
					CheckSQLError(c, err)
					return
				}
			}
		}
		if run.RunID == 0 {
			c.JSON(400, gin.H{"status": "error", "Message": "Fields RunID must not be 0 "})
			return
		}
		_, err = run.Create()
		if err != nil {
			CheckSQLError(c, err)
			return
		}
		c.JSON(200, gin.H{"status": "ok", "Message": "Run created", "ID": run.ID.String()})
		return
	case http.MethodPut:
		run := adapter.Run{}
		err := c.BindJSON(&run)
		if err != nil {
			c.JSON(500, gin.H{"status": "error", "Message": err.Error()})
			return
		}
		_, err = run.Update()
		if err != nil {
			CheckSQLError(c, err)
			return
		}
		c.JSON(200, gin.H{"status": "ok", "Message": "Run updated", "ID": run.ID.String()})
		return
	case http.MethodDelete:
		id := c.Query("id")
		if id == "" {
			c.JSON(400, gin.H{"status": "error", "Message": "param id is empty"})
			return
		}
		uuid := uuid.MustParse(id)
		_, err := adapter.Run{ID: uuid}.Delete()
		if err != nil {
			CheckSQLError(c, err)
			return
		}
		c.JSON(200, gin.H{"status": "ok", "Message": "Run deleted"})
	}
}

func LastRunID(c *gin.Context) {
	spaceName := c.Query("space")
	if spaceName == "" {
		c.JSON(400, gin.H{"status": "error", "Message": "param space is empty"})
		return
	}
	projectName := c.Query("project")
	if projectName == "" {
		c.JSON(400, gin.H{"status": "error", "Message": "param project is empty"})
		return
	}
	res, err := adapter.GetLastRunID(spaceName, projectName)
	if err != nil {
		CheckSQLError(c, err)
		return
	}
	c.JSON(200, gin.H{"status": "ok", "data": gin.H{"RunID": res}})
}

func SpaceRuns(c *gin.Context) {
	spaceName := c.Query("space")
	if spaceName == "" {
		c.JSON(400, gin.H{"status": "error", "Message": "param space is empty"})
		return
	}
	res, err := adapter.GetAllRunsBySpace(spaceName)
	if err != nil {
		CheckSQLError(c, err)
		return
	}
	c.JSON(200, gin.H{"status": "ok", "data": res})
}

func ProjectRuns(c *gin.Context) {
	spaceName := c.Query("space")
	if spaceName == "" {
		c.JSON(400, gin.H{"status": "error", "Message": "param space is empty"})
		return
	}
	projectName := c.Query("project")
	if projectName == "" {
		c.JSON(400, gin.H{"status": "error", "Message": "param project is empty"})
		return
	}
	res, err := adapter.GetAllRunsByProject(spaceName, projectName)
	if err != nil {
		CheckSQLError(c, err)
		return
	}
	c.JSON(200, gin.H{"status": "ok", "data": res})
}

func ReleaseRuns(c *gin.Context) {
	spaceName := c.Query("space")
	if spaceName == "" {
		c.JSON(400, gin.H{"status": "error", "Message": "param space is empty"})
		return
	}
	projectName := c.Query("project")
	if projectName == "" {
		c.JSON(400, gin.H{"status": "error", "Message": "param project is empty"})
		return
	}
	releaseName := c.Query("release")
	if releaseName == "" {
		c.JSON(400, gin.H{"status": "error", "Message": "param release is empty"})
		return
	}
	res, err := adapter.GetAllRunsByRelease(spaceName, projectName, releaseName)
	if err != nil {
		CheckSQLError(c, err)
		return
	}
	c.JSON(200, gin.H{"status": "ok", "data": res})
}

func VersionRuns(c *gin.Context) {
	spaceName := c.Query("space")
	if spaceName == "" {
		c.JSON(400, gin.H{"status": "error", "Message": "param space is empty"})
		return
	}
	projectName := c.Query("project")
	if projectName == "" {
		c.JSON(400, gin.H{"status": "error", "Message": "param project is empty"})
		return
	}
	releaseName := c.Query("release")
	if releaseName == "" {
		c.JSON(400, gin.H{"status": "error", "Message": "param release is empty"})
		return
	}
	versionName := c.Query("version")
	if versionName == "" {
		c.JSON(400, gin.H{"status": "error", "Message": "param version is empty"})
		return
	}
	res, err := adapter.GetAllRunsByVersion(spaceName, projectName, releaseName, versionName)
	if err != nil {
		CheckSQLError(c, err)
		return
	}
	c.JSON(200, gin.H{"status": "ok", "data": res})
}

func TestTypeRuns(c *gin.Context) {
	spaceName := c.Query("space")
	if spaceName == "" {
		c.JSON(400, gin.H{"status": "error", "Message": "param space is empty"})
		return
	}
	projectName := c.Query("project")
	if projectName == "" {
		c.JSON(400, gin.H{"status": "error", "Message": "param project is empty"})
		return
	}
	releaseName := c.Query("release")
	if releaseName == "" {
		c.JSON(400, gin.H{"status": "error", "Message": "param release is empty"})
		return
	}
	versionName := c.Query("version")
	if versionName == "" {
		c.JSON(400, gin.H{"status": "error", "Message": "param version is empty"})
		return
	}
	testtypeName := c.Query("testtype")
	if testtypeName == "" {
		c.JSON(400, gin.H{"status": "error", "Message": "param testtype is empty"})
		return
	}
	res, err := adapter.GetAllRunsByTestType(spaceName, projectName, releaseName, versionName, testtypeName)
	if err != nil {
		CheckSQLError(c, err)
		return
	}
	c.JSON(200, gin.H{"status": "ok", "data": res})
}

func SetEndTime(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(400, gin.H{"status": "error", "Message": "param id is empty"})
		return
	}
	endtime := c.Query("endtime")
	if id == "" {
		c.JSON(400, gin.H{"status": "error", "Message": "param id is empty"})
		return
	}
	uuid := uuid.MustParse(id)
	t, err := strconv.ParseInt(endtime, 10, 64)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "Message": err})
		return
	}
	res, err := adapter.SetEndTime(uuid, time.Unix(t, 0))
	if err != nil {
		CheckSQLError(c, err)
		return
	}
	c.JSON(200, gin.H{"status": "ok", "data": res})
}
