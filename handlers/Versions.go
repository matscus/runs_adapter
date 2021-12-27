package handlers

import (
	"net/http"
	"runs_adapter/adapter"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Versions(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodGet:
		all := c.Query("all")
		if all != "" {
			res, err := adapter.GetAllVersions()
			if err != nil {
				c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
				return
			}
			c.JSON(200, res)
			return
		}
		id := c.Query("id")
		if id != "" {
			uuid := uuid.MustParse(id)
			res, err := adapter.GetVersionByID(uuid)
			if err != nil {
				c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
				return
			}
			c.JSON(200, res)
			return
		}
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
		if projectName == "" {
			c.JSON(400, gin.H{"Status": "error", "Message": "param release is empty"})
			return
		}
		versionName := c.Query("version")
		if projectName == "" {
			c.JSON(400, gin.H{"Status": "error", "Message": "param release is empty"})
			return
		}
		res, err := adapter.GetVersion(spaceName, projectName, releaseName, versionName)
		if err != nil {
			c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
			return
		}
		c.JSON(200, res)
		return
	case http.MethodPost:
		version := adapter.Version{}
		err := c.BindJSON(&version)
		if err != nil {
			c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
			return
		}
		version.ID = uuid.New()
		_, err = version.Create()
		if err != nil {
			c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
			return
		}
		c.JSON(200, gin.H{"Status": "ok", "Message": "Version created", "ID": version.ID.String()})
		return
	case http.MethodPut:
		version := adapter.Version{}
		err := c.BindJSON(&version)
		if err != nil {
			c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
			return
		}
		_, err = version.Update()
		if err != nil {
			c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
			return
		}
		c.JSON(200, gin.H{"Status": "ok", "Message": "Version updates"})
		return
	case http.MethodDelete:
		id := c.Query("id")
		if id == "" {
			c.JSON(400, gin.H{"Status": "error", "Message": "param id is empty"})
			return
		}
		uuid := uuid.MustParse(id)
		_, err := adapter.Version{ID: uuid}.Delete()
		if err != nil {
			c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
			return
		}
		c.JSON(200, gin.H{"Status": "ok", "Message": "Version deleted"})
	}
}
