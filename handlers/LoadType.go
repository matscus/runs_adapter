package handlers

import (
	"net/http"
	"runs_adapter/adapter"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func LoadType(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodGet:
		id := c.Query("id")
		if id != "" {
			uuid := uuid.MustParse(id)
			res, err := adapter.GetLoadTypeByID(uuid)
			if err != nil {
				CheckSQLError(c, err)
				return
			}
			c.JSON(200, gin.H{"status": "ok", "data": res})
			return
		}
		res, err := adapter.GetLoadTypes()
		if err != nil {
			CheckSQLError(c, err)
			return
		}
		c.JSON(200, gin.H{"status": "ok", "data": res})
		return
	case http.MethodPost:
		loadType := adapter.LoadType{}
		err := c.BindJSON(&loadType)
		if err != nil {
			c.JSON(500, gin.H{"status": "error", "Message": err.Error()})
			return
		}
		loadType.ID = uuid.New()
		_, err = loadType.Create()
		if err != nil {
			CheckSQLError(c, err)
			return
		}
		c.JSON(200, gin.H{"status": "ok", "Message": "Space created", "ID": loadType.ID.String()})
		return
	case http.MethodPut:
		loadType := adapter.LoadType{}
		err := c.BindJSON(&loadType)
		if err != nil {
			c.JSON(500, gin.H{"status": "error", "Message": err.Error()})
			return
		}
		_, err = loadType.Update()
		if err != nil {
			CheckSQLError(c, err)
			return
		}
		c.JSON(200, gin.H{"status": "ok", "Message": "Space updates"})
		return
	case http.MethodDelete:
		id := c.Query("id")
		if id == "" {
			c.JSON(400, gin.H{"status": "error", "Message": "param id is empty"})
			return
		}
		uuid := uuid.MustParse(id)
		_, err := adapter.LoadType{ID: uuid}.Delete()
		if err != nil {
			CheckSQLError(c, err)
			return
		}
		c.JSON(200, gin.H{"status": "ok", "Message": "Space deleted"})
	}
}
