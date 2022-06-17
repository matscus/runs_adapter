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
		res, err := adapter.GetAllTestTypes()
		if err != nil {
			CheckSQLError(c, err)
			return
		}
		c.JSON(200, gin.H{"status": "ok", "data": res})
		return
	case http.MethodPost:
		testType := adapter.TestType{}
		err := c.BindJSON(&testType)
		if err != nil {
			c.JSON(500, gin.H{"status": "error", "Message": err.Error()})
			return
		}
		testType.ID = uuid.New()
		_, err = testType.Create()
		if err != nil {
			CheckSQLError(c, err)
			return
		}
		c.JSON(200, gin.H{"status": "ok", "Message": "Test type created", "ID": testType.ID.String()})
		return
	case http.MethodPut:
		testType := adapter.TestType{}
		err := c.BindJSON(&testType)
		if err != nil {
			c.JSON(500, gin.H{"status": "error", "Message": err.Error()})
			return
		}
		_, err = testType.Update()
		if err != nil {
			CheckSQLError(c, err)
			return
		}
		c.JSON(200, gin.H{"status": "ok", "Message": "Test type updates"})
		return
	case http.MethodDelete:
		id := c.Query("id")
		if id == "" {
			c.JSON(400, gin.H{"status": "error", "Message": "param id is empty"})
			return
		}
		uuid := uuid.MustParse(id)
		_, err := adapter.TestType{ID: uuid}.Delete()
		if err != nil {
			CheckSQLError(c, err)
			return
		}
		c.JSON(200, gin.H{"status": "ok", "Message": "Test type deleted"})
	}
}
