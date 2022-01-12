package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"runs_adapter/adapter"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Profiles(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodGet:
		all := c.Query("all")
		if all != "" {
			res, err := adapter.GetAllProfiles()
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
			res, err := adapter.GetProfileByID(uuid)
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
			c.JSON(400, gin.H{"status": "error", "Message": "param release is empty"})
			return
		}
		versionName := c.Query("version")
		if versionName == "" {
			if releaseName == "" {
				c.JSON(400, gin.H{"status": "error", "Message": "param version is empty"})
				return
			}
		}
		testTypeName := c.Query("testtype")
		if testTypeName == "" {
			if releaseName == "" {
				c.JSON(400, gin.H{"status": "error", "Message": "param testtype is empty"})
				return
			}
		}
		res, err := adapter.GetProfiles(spaceName, projectName, releaseName, versionName, testTypeName)
		if err != nil {
			CheckSQLError(c, err)
			return
		}
		c.JSON(200, gin.H{"status": "ok", "data": res})
		return

	case http.MethodPost:
		profile := adapter.Profile{}
		err := c.BindJSON(&profile)
		if err != nil {
			c.JSON(500, gin.H{"status": "error", "Message": err.Error()})
			return
		}
		new := false
		oldProfile, err := adapter.GetProfile(profile.SpaceName, profile.ProjectName, profile.ReleaseName, profile.VersionName, profile.TestTypeName, profile.Name)
		if err == sql.ErrNoRows {
			profile.ID = uuid.New()
			new = true
		} else {
			profile.ID = oldProfile.ID
		}
		if profile.SpaceID == nilUUID {
			profile.SpaceID, err = adapter.GetSpaceID(profile.SpaceName)
			if err != nil {
				if err == sql.ErrNoRows {
					profile.SpaceID = uuid.New()
					_, err = adapter.Space{ID: profile.SpaceID, Name: profile.SpaceName}.Create()
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
		if profile.ProjectID == nilUUID {
			profile.ProjectID, err = adapter.GetProjectID(profile.SpaceID, profile.ProjectName)
			if err != nil {
				if err == sql.ErrNoRows {
					profile.ProjectID = uuid.New()
					_, err = adapter.Project{ID: profile.ProjectID, Name: profile.ProjectName, SpaceID: profile.SpaceID}.Create()
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
		if profile.ReleaseID == nilUUID {
			profile.ReleaseID, err = adapter.GetReleaseID(profile.ProjectID, profile.ReleaseName)
			if err != nil {
				if err == sql.ErrNoRows {
					profile.ReleaseID = uuid.New()
					_, err = adapter.Release{ID: profile.ReleaseID, Name: profile.ReleaseName, ProjectID: profile.ProjectID}.Create()
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
		if profile.VersionID == nilUUID {
			profile.VersionID, err = adapter.GetVersionID(profile.ReleaseID, profile.VersionName)
			if err != nil {
				if err == sql.ErrNoRows {
					profile.VersionID = uuid.New()
					_, err = adapter.Version{ID: profile.VersionID, Name: profile.VersionName, ReleaseID: profile.ReleaseID}.Create()
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
		if profile.TestTypeID == nilUUID {
			profile.TestTypeID, err = adapter.GetTestTypeID(profile.ProjectID, profile.TestTypeName)
			if err != nil {
				if err == sql.ErrNoRows {
					profile.TestTypeID = uuid.New()
					_, err = adapter.TestType{ID: profile.TestTypeID, Name: profile.TestTypeName, ProjectID: profile.ProjectID}.Create()
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
		if new {
			_, err = profile.Create()
			if err != nil {
				CheckSQLError(c, err)
				return
			}
			c.JSON(200, gin.H{"status": "ok", "Message": "Profile created", "ID": profile.ID.String()})
			return
		}
		_, err = profile.Update()
		if err != nil {
			CheckSQLError(c, err)
			return
		}
		c.JSON(200, gin.H{"status": "ok", "Message": "Profile updated"})
		return

	case http.MethodPut:
		profile := adapter.Profile{}
		err := c.BindJSON(&profile)
		if err != nil {
			c.JSON(500, gin.H{"status": "error", "Message": err.Error()})
			return
		}

		fmt.Println(profile)
		_, err = profile.Update()
		if err != nil {
			CheckSQLError(c, err)
			return
		}
		c.JSON(200, gin.H{"status": "ok", "Message": "Profile updated", "ID": profile.ID.String()})
		return
	case http.MethodDelete:
		id := c.Query("id")
		if id == "" {
			c.JSON(400, gin.H{"status": "error", "Message": "param id is empty"})
			return
		}
		uuid := uuid.MustParse(id)
		_, err := adapter.Profile{ID: uuid}.Delete()
		if err != nil {
			CheckSQLError(c, err)
			return
		}
		c.JSON(200, gin.H{"status": "ok", "Message": "Profile deleted"})
	}
}
