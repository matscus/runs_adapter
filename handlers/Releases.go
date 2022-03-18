package handlers

import (
	"fmt"
	"net/http"
	"runs_adapter/adapter"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Releases(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodGet:
		all := c.Query("all")
		if all != "" {
			res, err := adapter.GetAllReleases()
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
			res, err := adapter.GetReleaseByID(uuid)
			if err != nil {
				CheckSQLError(c, err)
				return
			}
			c.JSON(200, gin.H{"status": "ok", "data": res})
			return
		}
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
			fmt.Println("GetReleases")
			res, err := adapter.GetReleases(spaceName, projectName)
			if err != nil {
				CheckSQLError(c, err)
				return
			}
			c.JSON(200, gin.H{"status": "ok", "data": res})
			return
		}
		res, err := adapter.GetRelease(spaceName, projectName, releaseName)
		if err != nil {
			CheckSQLError(c, err)
			return
		}
		c.JSON(200, gin.H{"status": "ok", "data": res})
		return
	case http.MethodPost:
		release := adapter.Release{}
		err := c.BindJSON(&release)
		if err != nil {
			c.JSON(500, gin.H{"status": "error", "Message": err.Error()})
			return
		}
		release.ID = uuid.New()
		_, err = release.Create()
		if err != nil {
			CheckSQLError(c, err)
			return
		}
		c.JSON(200, gin.H{"status": "ok", "Message": "Release created", "ID": release.ID.String()})
		return
	case http.MethodPut:
		release := adapter.Release{}
		err := c.BindJSON(&release)
		if err != nil {
			c.JSON(500, gin.H{"status": "error", "Message": err.Error()})
			return
		}
		_, err = release.Update()
		if err != nil {
			CheckSQLError(c, err)
			return
		}
		c.JSON(200, gin.H{"status": "ok", "Message": "Release updated"})
		return
	case http.MethodDelete:
		id := c.Query("id")
		if id == "" {
			c.JSON(400, gin.H{"status": "error", "Message": "param id is empty"})
			return
		}
		uuid := uuid.MustParse(id)
		_, err := adapter.Release{ID: uuid}.Delete()
		if err != nil {
			CheckSQLError(c, err)
			return
		}
		c.JSON(200, gin.H{"status": "ok", "Message": "Release deleted"})
	}
}
