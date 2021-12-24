package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"runs_adapter/adapter"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Spaces(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodGet:
		all := c.Query("all")
		if all != "" {
			res, err := adapter.GetAllSpaces()
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
			res, err := adapter.GetSpaceByID(uuid)
			if err != nil {
				CheckSQLError(c, err)
				return
			}
			c.JSON(200, res)
			return
		}
		space := c.Query("space")
		if space != "" {
			res, err := adapter.GetSpace(space)
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
		space := adapter.Space{}
		err = json.Unmarshal(body, &space)
		if err != nil {
			c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
			return
		}
		space.ID = uuid.New()
		_, err = space.Create()
		if err != nil {
			CheckSQLError(c, err)
			return
		}
		c.JSON(200, gin.H{"Status": "ok", "Message": "Space created", "ID": space.ID.String()})
		return
	case http.MethodPut:
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
			return
		}
		space := adapter.Space{}
		err = json.Unmarshal(body, &space)
		if err != nil {
			c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
			return
		}
		_, err = space.Update()
		if err != nil {
			CheckSQLError(c, err)
			return
		}
		c.JSON(200, gin.H{"Status": "ok", "Message": "Space updates"})
		return
	case http.MethodDelete:
		id := c.Query("id")
		if id == "" {
			c.JSON(400, gin.H{"Status": "error", "Message": "param id is empty"})
			return
		}
		uuid := uuid.MustParse(id)
		_, err := adapter.Space{ID: uuid}.Delete()
		if err != nil {
			CheckSQLError(c, err)
			return
		}
		c.JSON(200, gin.H{"Status": "ok", "Message": "Space deleted"})
	}
}