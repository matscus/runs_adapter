package handlers

import (
	"net/http"
	"runs_adapter/adapter"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestTypes(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodGet:
		all := c.Query("all")
		if all != "" {
			res, err := adapter.GetAllTestTypes()
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
			res, err := adapter.GetTestTypeByID(uuid)
			if err != nil {
				CheckSQLError(c, err)
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
		testTypeName := c.Query("testtype")
		if projectName == "" {
			c.JSON(400, gin.H{"Status": "error", "Message": "param release is empty"})
			return
		}
		res, err := adapter.GetTestType(spaceName, testTypeName, projectName)
		if err != nil {
			CheckSQLError(c, err)
			return
		}
		c.JSON(200, res)
		return
	case http.MethodPost:
		testType := adapter.TestType{}
		err := c.BindJSON(&testType)
		if err != nil {
			c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
			return
		}
		testType.ID = uuid.New()
		_, err = testType.Create()
		if err != nil {
			CheckSQLError(c, err)
			return
		}
		c.JSON(200, gin.H{"Status": "ok", "Message": "Test type created", "ID": testType.ID.String()})
		return
	case http.MethodPut:
		testType := adapter.TestType{}
		err := c.BindJSON(&testType)
		if err != nil {
			c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
			return
		}
		_, err = testType.Update()
		if err != nil {
			CheckSQLError(c, err)
			return
		}
		c.JSON(200, gin.H{"Status": "ok", "Message": "Test type updates"})
		return
	case http.MethodDelete:
		id := c.Query("id")
		if id == "" {
			c.JSON(400, gin.H{"Status": "error", "Message": "param id is empty"})
			return
		}
		uuid := uuid.MustParse(id)
		_, err := adapter.TestType{ID: uuid}.Delete()
		if err != nil {
			CheckSQLError(c, err)
			return
		}
		c.JSON(200, gin.H{"Status": "ok", "Message": "Test type deleted"})
	}
}
