package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"runs_adapter/adapter"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
			c.JSON(200, res)
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
			c.JSON(200, res)
			return
		}
	case http.MethodPost:
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
			return
		}
		run := adapter.Run{}
		err = json.Unmarshal(body, &run)
		if err != nil {
			c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
			return
		}
		run.ID = uuid.New()
		_, err = run.Create()
		if err != nil {
			CheckSQLError(c, err)
			return
		}
		c.JSON(200, gin.H{"Status": "ok", "Message": "Run created", "ID": run.ID.String()})
		return
	case http.MethodPut:
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
			return
		}
		run := adapter.Run{}
		err = json.Unmarshal(body, &run)
		if err != nil {
			c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
			return
		}
		run.ID = uuid.New()
		_, err = run.Update()
		if err != nil {
			CheckSQLError(c, err)
			return
		}
		c.JSON(200, gin.H{"Status": "ok", "Message": "Run created", "ID": run.ID.String()})
		return
	case http.MethodDelete:
		id := c.Query("id")
		if id == "" {
			c.JSON(400, gin.H{"Status": "error", "Message": "param id is empty"})
			return
		}
		uuid := uuid.MustParse(id)
		_, err := adapter.Run{ID: uuid}.Delete()
		if err != nil {
			CheckSQLError(c, err)
			return
		}
		c.JSON(200, gin.H{"Status": "ok", "Message": "Run deleted"})
	}
}

func LastRunID(c *gin.Context) {
	spaceName := c.Query("space")
	if spaceName == "" {
		c.JSON(400, gin.H{"Status": "error", "Message": "param space is empty"})
		return
	}
	projectName := c.Query("project")
	if projectName == "" {
		c.JSON(400, gin.H{"Status": "error", "Message": "param project is empty"})
		return
	}
	res, err := adapter.GetLastRunID(spaceName, projectName)
	if err != nil {
		CheckSQLError(c, err)
		return
	}
	c.JSON(200, gin.H{"RunID": res})
}

func SpaceRuns(c *gin.Context) {
	spaceName := c.Query("space")
	if spaceName == "" {
		c.JSON(400, gin.H{"Status": "error", "Message": "param space is empty"})
		return
	}
	res, err := adapter.GetAllRunsBySpace(spaceName)
	if err != nil {
		CheckSQLError(c, err)
		return
	}
	c.JSON(200, res)
}

func ProjectRuns(c *gin.Context) {
	spaceName := c.Query("space")
	if spaceName == "" {
		c.JSON(400, gin.H{"Status": "error", "Message": "param space is empty"})
		return
	}
	projectName := c.Query("project")
	if projectName == "" {
		c.JSON(400, gin.H{"Status": "error", "Message": "param project is empty"})
		return
	}
	res, err := adapter.GetAllRunsByProject(spaceName, projectName)
	if err != nil {
		CheckSQLError(c, err)
		return
	}
	c.JSON(200, res)
}

func ReleaseRuns(c *gin.Context) {
	spaceName := c.Query("space")
	if spaceName == "" {
		c.JSON(400, gin.H{"Status": "error", "Message": "param space is empty"})
		return
	}
	projectName := c.Query("project")
	if projectName == "" {
		c.JSON(400, gin.H{"Status": "error", "Message": "param project is empty"})
		return
	}
	releaseName := c.Query("release")
	if releaseName == "" {
		c.JSON(400, gin.H{"Status": "error", "Message": "param release is empty"})
		return
	}
	res, err := adapter.GetAllRunsByRelease(spaceName, projectName, releaseName)
	if err != nil {
		CheckSQLError(c, err)
		return
	}
	c.JSON(200, res)
}

func VersionRuns(c *gin.Context) {
	spaceName := c.Query("space")
	if spaceName == "" {
		c.JSON(400, gin.H{"Status": "error", "Message": "param space is empty"})
		return
	}
	projectName := c.Query("project")
	if projectName == "" {
		c.JSON(400, gin.H{"Status": "error", "Message": "param project is empty"})
		return
	}
	releaseName := c.Query("release")
	if releaseName == "" {
		c.JSON(400, gin.H{"Status": "error", "Message": "param release is empty"})
		return
	}
	versionName := c.Query("version")
	if versionName == "" {
		c.JSON(400, gin.H{"Status": "error", "Message": "param version is empty"})
		return
	}
	res, err := adapter.GetAllRunsByVersion(spaceName, projectName, releaseName, versionName)
	if err != nil {
		CheckSQLError(c, err)
		return
	}
	c.JSON(200, res)
}

func TestTypeRuns(c *gin.Context) {
	spaceName := c.Query("space")
	if spaceName == "" {
		c.JSON(400, gin.H{"Status": "error", "Message": "param space is empty"})
		return
	}
	projectName := c.Query("project")
	if projectName == "" {
		c.JSON(400, gin.H{"Status": "error", "Message": "param project is empty"})
		return
	}
	releaseName := c.Query("release")
	if releaseName == "" {
		c.JSON(400, gin.H{"Status": "error", "Message": "param release is empty"})
		return
	}
	versionName := c.Query("version")
	if versionName == "" {
		c.JSON(400, gin.H{"Status": "error", "Message": "param version is empty"})
		return
	}
	testtypeName := c.Query("testtype")
	if testtypeName == "" {
		c.JSON(400, gin.H{"Status": "error", "Message": "param testtype is empty"})
		return
	}
	res, err := adapter.GetAllRunsByTestType(spaceName, projectName, releaseName, versionName, testtypeName)
	if err != nil {
		CheckSQLError(c, err)
		return
	}
	c.JSON(200, res)
}
